package validator

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"strings"
	"sync"

	"github.com/avast/retry-go"
)

var (
	defaultEnvoyPath = "/usr/local/bin/envoy"
	// TODO(tim): avoid hardcoding the envoy image version in multiple places.
	defaultEnvoyImage = "quay.io/solo-io/envoy-gloo:1.34.1-patch3"
)

// ErrInvalidXDS is returned when Envoy rejects the supplied YAML.
var ErrInvalidXDS = errors.New("invalid xds configuration")

// Validator validates an Envoy bootstrap/partial YAML.
type Validator interface {
	// Validate validates the given YAML configuration. Returns an error
	// if the configuration is invalid.
	Validate(context.Context, string) error
}

// validationCache caches validation results to avoid redundant validations
type validationCache struct {
	cache map[string]error
	mu    sync.RWMutex
}

// cachedValidator wraps another validator with caching
type cachedValidator struct {
	underlying Validator
	cache      *validationCache
}

var _ Validator = &cachedValidator{}

// NewCachedValidator wraps any validator with caching to avoid redundant validations
func NewCachedValidator(underlying Validator) Validator {
	return &cachedValidator{
		underlying: underlying,
		cache: &validationCache{
			cache: make(map[string]error),
		},
	}
}

func (cv *cachedValidator) Validate(ctx context.Context, yaml string) error {
	return cv.cache.validate(ctx, yaml, cv.underlying.Validate)
}

// validate checks cache first, then delegates to actual validator if needed
func (vc *validationCache) validate(ctx context.Context, yaml string, actualValidator func(context.Context, string) error) error {
	// Create hash of the YAML content
	hash := sha256.Sum256([]byte(yaml))
	key := hex.EncodeToString(hash[:])

	// Check cache first
	vc.mu.RLock()
	if result, exists := vc.cache[key]; exists {
		vc.mu.RUnlock()
		return result
	}
	vc.mu.RUnlock()

	// Not in cache, perform actual validation
	result := actualValidator(ctx, yaml)

	// Store result in cache
	vc.mu.Lock()
	vc.cache[key] = result
	vc.mu.Unlock()

	return result
}

// binaryValidator validates envoy using the binary.
type binaryValidator struct {
	path string
}

func NewBinaryValidator(path string) Validator {
	if path == "" {
		path = defaultEnvoyPath
	}
	// TODO(tim): check if the binary exists.
	return &binaryValidator{path: path}
}

var _ Validator = &binaryValidator{}

func (b *binaryValidator) Validate(ctx context.Context, yaml string) error {
	cmd := exec.CommandContext(ctx, b.path, "--mode", "validate", "--config-yaml", yaml, "-l", "critical", "--log-format", "%v")
	cmd.Stdin = strings.NewReader(yaml)
	var e bytes.Buffer
	cmd.Stderr = &e
	if err := cmd.Run(); err != nil {
		rawErr := strings.TrimSpace(e.String())
		if _, ok := err.(*exec.ExitError); ok {
			if rawErr == "" {
				rawErr = err.Error()
			}
			return fmt.Errorf("%w: %s", ErrInvalidXDS, rawErr)
		}
		return fmt.Errorf("envoy validate invocation failed: %v", err)
	}
	return nil
}

// PersistentValidator is a validator that supports cleanup operations.
type PersistentValidator interface {
	Validator
	// Start explicitly starts the validator (e.g., container startup).
	Start(context.Context) error
	// Cleanup performs any necessary cleanup operations.
	Cleanup(context.Context) error
}

// persistentDockerValidator uses a long-running Docker container for validation
// to avoid the overhead of starting a new container for each validation.
type persistentDockerValidator struct {
	img           string
	containerName string
	started       bool
}

var _ PersistentValidator = &persistentDockerValidator{}

func NewPersistentDockerValidator(img string) PersistentValidator {
	if img == "" {
		img = defaultEnvoyImage
	}
	return &persistentDockerValidator{
		img:           img,
		containerName: "envoy-validator-persistent",
		started:       false,
	}
}

func (p *persistentDockerValidator) Start(ctx context.Context) error {
	// If it already exists, just start it
	if err := exec.CommandContext(ctx, "docker", "inspect", p.containerName).Run(); err == nil {
		if err := exec.CommandContext(ctx, "docker", "start", "-d", p.containerName).Run(); err != nil {
			return fmt.Errorf("failed to start existing container: %w", err)
		}
		p.started = true
		return nil
	}

	// Otherwise, create it with full config
	// Override entrypoint to keep container alive for validation commands
	if err := exec.CommandContext(
		ctx,
		"docker", "create",
		"--name", p.containerName,
		"--platform", "linux/amd64",
		"--memory", "512m",
		"--memory-swap", "512m", // no swap
		"--restart", "no",
		"--entrypoint", "sh",
		"--read-only",
		"--pids-limit", "256",
		"--network", "none",
		"--tmpfs", "/tmp:rw,exec",
		p.img,
		"-c", "trap 'exit 0' TERM; sleep 3600 & wait",
	).Run(); err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	startCmd := exec.CommandContext(ctx, "docker", "start", p.containerName)
	stderr := new(bytes.Buffer)
	startCmd.Stderr = stderr
	if err := startCmd.Run(); err != nil {
		rawErr := strings.TrimSpace(stderr.String())
		if rawErr != "" {
			return fmt.Errorf("failed to start container: %w: %s", err, rawErr)
		}
		return fmt.Errorf("failed to start container: %w", err)
	}
	p.started = true

	if err := waitRunningWithRetry(ctx, p.containerName); err != nil {
		return fmt.Errorf("failed to wait for container to be running: %w", err)
	}

	return nil
}

func waitRunningWithRetry(ctx context.Context, name string) error {
	return retry.Do(
		func() error {
			out, err := exec.CommandContext(ctx, "docker", "inspect", "-f", "{{.State.Running}}", name).Output()
			if err != nil {
				return err
			}
			if strings.TrimSpace(string(out)) != "true" {
				return fmt.Errorf("not running yet")
			}
			return nil
		},
		retry.Attempts(5),
	)
}

func (p *persistentDockerValidator) Validate(ctx context.Context, yaml string) error {
	cmd := exec.CommandContext(
		ctx, "docker", "exec", "-i",
		p.containerName,
		"/usr/local/bin/envoy",
		"--mode", "validate",
		"--config-yaml", yaml,
		"-l", "critical",
		"--log-format", "%v",
	)
	cmd.Stdin = strings.NewReader(yaml)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		return nil
	}

	rawErr := strings.TrimSpace(stderr.String())
	if _, ok := err.(*exec.ExitError); ok {
		// Extract just the envoy error message, ignoring Docker pull output
		if envoyErr := extractEnvoyError(rawErr); envoyErr != "" {
			return fmt.Errorf("%w: %s", ErrInvalidXDS, envoyErr)
		}
		if rawErr == "" {
			rawErr = err.Error()
		}
		return fmt.Errorf("%w: %s", ErrInvalidXDS, rawErr)
	}
	return fmt.Errorf("envoy validate invocation failed: %v", err)
}

// Cleanup stops and removes the persistent container
func (p *persistentDockerValidator) Cleanup(ctx context.Context) error {
	// Always try to clean up, regardless of started status
	// Use short timeout since sleep doesn't handle SIGTERM gracefully
	exec.CommandContext(ctx, "docker", "stop", "-t", "1", p.containerName).Run()
	exec.CommandContext(ctx, "docker", "rm", p.containerName).Run()
	p.started = false
	return nil
}

// extractEnvoyError extracts the actual Envoy validation error from stderr output,
// ignoring Docker pull progress and other noise that comes before the error.
func extractEnvoyError(stderr string) string {
	lines := strings.Split(stderr, "\n")
	// find the first line containing the Envoy error message. see:
	// https://github.com/envoyproxy/envoy/blob/d552b66f5d70ddd9e13c68c40f70729a45fb24e0/source/server/config_validation/server.cc#L75
	errorIndex := slices.IndexFunc(lines, func(line string) bool {
		return strings.Contains(strings.TrimSpace(line), "error initializing configuration")
	})
	if errorIndex == -1 {
		return ""
	}
	// extract all remaining lines that are relevant error context
	remainingLines := make([]string, 0, len(lines)-errorIndex)
	for i := errorIndex; i < len(lines); i++ {
		if trimmed := strings.TrimSpace(lines[i]); trimmed != "" {
			remainingLines = append(remainingLines, trimmed)
		}
	}
	return strings.Join(remainingLines, " ")
}
