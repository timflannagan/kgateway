package validator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const (
	envoyPath  = "/usr/local/bin/envoy"
	envoyImage = "quay.io/solo-io/envoy-gloo:1.34.1-patch1"
)

// ErrInvalidXDS is returned when Envoy rejects the supplied YAML.
var ErrInvalidXDS = errors.New("invalid xds configuration")

// Validator validates an Envoy bootstrap/partial YAML.
type Validator interface {
	Validate(context.Context, string) error
}

// New chooses the best validator available.
func New() Validator {
	// check if envoy is in the path
	if _, err := exec.LookPath(envoyPath); err == nil {
		return &binaryValidator{path: envoyPath}
	}
	// otherwise, fallback to docker
	return &dockerValidator{img: envoyImage}
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

// stripDockerWarn drops the platform-mismatch line docker prints on ARM hosts.
func stripDockerWarn(s string) string {
	var cleaned []string
	// TODO(tim): fix the linter violation below.
	for _, l := range strings.Split(s, "\n") {
		if !strings.HasPrefix(l, "WARNING: The requested image's platform") {
			cleaned = append(cleaned, l)
		}
	}
	return strings.TrimSpace(strings.Join(cleaned, "\n"))
}
