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
		return fmt.Errorf("%w: %s", ErrInvalidXDS, trim(e.String()))
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
		d.img,
		"--mode",
		"validate",
		"--config-yaml", yaml,
		"-l", "critical",
		"--log-format", "%v",
	)
	var buf bytes.Buffer
	cmd.Stdin, cmd.Stdout, cmd.Stderr = strings.NewReader(yaml), &buf, &buf
	if err := cmd.Run(); err != nil {
		out := stripWarn(buf.String())
		if strings.Contains(out, "error initializing configuration") || strings.Contains(out, "[error] envoy") {
			return fmt.Errorf("%w: %s", ErrInvalidXDS, trim(out))
		}
		return errors.New(trim(out))
	}
	return nil
}

// trim removes some common prefixes from the output.
func trim(s string) string {
	s = strings.ReplaceAll(s, "''", "'")
	s = strings.TrimPrefix(s, "error initializing configuration '': ")
	s = strings.TrimPrefix(s, "Failed to parse request template: ")
	if i := strings.IndexByte(s, '\n'); i != -1 {
		s = s[:i]
	}
	return strings.Join(strings.Fields(s), " ")
}

// stripWarn removes the warning about the image platform from the output. This was
// causing cmd.Run to fail with an error.
func stripWarn(s string) string {
	var r []string
	for _, l := range strings.Split(s, "\n") {
		if !strings.HasPrefix(l, "WARNING: The requested image's platform") {
			r = append(r, l)
		}
	}
	return strings.Join(r, "\n")
}
