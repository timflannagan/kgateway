package validator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
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

// Option is a functional option for New.
type Option func(*config)

// WithBinaryPath overrides the Envoy binary path.
func WithBinaryPath(p string) Option { return func(c *config) { c.binaryPath = p } }

// WithDockerImage overrides the Docker image used for validation.
func WithDockerImage(img string) Option { return func(c *config) { c.dockerImage = img } }

// config stores the validator configuration.
type config struct {
	binaryPath  string
	dockerImage string
}

// New chooses the best validator available.
func New(o ...Option) Validator {
	c := &config{}
	for _, opt := range o {
		opt(c)
	}
	// use defaults if not set by options
	binaryPath := c.binaryPath
	if binaryPath == "" {
		binaryPath = defaultEnvoyPath
	}
	dockerImage := c.dockerImage
	if dockerImage == "" {
		dockerImage = defaultEnvoyImage
	}
	// check if envoy is in the path
	if _, err := exec.LookPath(binaryPath); err == nil {
		return &binaryValidator{path: binaryPath}
	}
	// otherwise, fallback to docker
	return &dockerValidator{img: dockerImage}
}

// binaryValidator validates envoy using the binary.
type binaryValidator struct {
	path string
}

var _ Validator = &binaryValidator{}

func (b *binaryValidator) Validate(ctx context.Context, yaml string) error {
	cmd := exec.CommandContext(ctx, b.path, "--mode", "validate", "--config-yaml", yaml, "-l", "critical", "--log-format", "%v")
	cmd.Stdin = strings.NewReader(yaml)
	var e bytes.Buffer
	cmd.Stderr = &e
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidXDS, e.String())
	}
	return nil
}

type dockerValidator struct {
	img string
}

var _ Validator = &dockerValidator{}

func (d *dockerValidator) Validate(ctx context.Context, yaml string) error {
	cmd := exec.CommandContext(ctx,
		"docker", "run",
		"--rm",
		"-i",
		"--platform", "linux/amd64",
		d.img,
		"--mode",
		"validate",
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

	// TODO(tim): Just return first match from "error initializing configuration"?
	rawErr := strings.TrimSpace(stderr.String())
	rawErr = stripDockerWarn(rawErr)
	if _, ok := err.(*exec.ExitError); ok {
		if rawErr == "" {
			rawErr = err.Error()
		}
		return fmt.Errorf("%w: %s", ErrInvalidXDS, rawErr)
	}
	return fmt.Errorf("envoy validate invocation failed: %v", err)
}

// stripDockerWarn removes the platform-mismatch warning Docker prints on ARM hosts.
func stripDockerWarn(s string) string {
	lines := strings.Split(s, "\n")
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		if !strings.HasPrefix(line, "WARNING: The requested image's platform") {
			cleaned = append(cleaned, line)
		}
	}
	return strings.TrimSpace(strings.Join(cleaned, "\n"))
}
