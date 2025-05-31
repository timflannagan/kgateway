package validator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	envoyPath  = "/usr/local/bin/envoy"
	envoyImage = "quay.io/solo-io/envoy-gloo:1.34.1-patch1"
)

// ErrInvalidXDS is returned when Envoy rejects the supplied YAML.
var ErrInvalidXDS = errors.New("invalid xds configuration")

// Validator validates an Envoy bootstrap/partial YAML.
//
//	nil            → accepted
//	ErrInvalidXDS  → rejected by Envoy
//	other error    → validation mechanism failed
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
	// first try with default platform
	if err := d.run(ctx, yaml, platFlag()); err == nil || errors.Is(err, ErrInvalidXDS) {
		return err
	}
	// fallback for arm64 hosts to amd64 manifest
	if runtime.GOARCH == "arm64" {
		return d.run(ctx, yaml, "--platform=linux/amd64")
	}
	return fmt.Errorf("docker validator failed after retries")
}

func (d *dockerValidator) run(parent context.Context, yaml, flag string) error {
	args := []string{"docker", "run", "--rm", "-i"}
	if flag != "" {
		args = append(args, flag)
	}
	args = append(args, d.img, "--mode", "validate", "--config-yaml", yaml, "-l", "critical", "--log-format", "%v")
	cmd := exec.CommandContext(parent, args[0], args[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = strings.NewReader(yaml), nil, nil
	var buf bytes.Buffer
	cmd.Stderr = &buf
	cmd.Stdout = &buf
	err := cmd.Run()
	out := stripWarn(buf.String())
	if os.Getenv("ENVOY_VALIDATOR_DEBUG") != "" {
		fmt.Fprintln(os.Stderr, "docker validator output:\n"+out)
	}
	if err != nil {
		if strings.Contains(out, "error initializing configuration") || strings.Contains(out, "[error] envoy") {
			return fmt.Errorf("%w: %s", ErrInvalidXDS, trim(out))
		}
		return fmt.Errorf(trim(out))
	}
	return nil
}

func platFlag() string {
	if runtime.GOARCH == "arm64" {
		return "--platform=linux/amd64"
	}
	return ""
}

func trim(s string) string {
	s = strings.ReplaceAll(s, "''", "'")
	s = strings.TrimPrefix(s, "error initializing configuration '': ")
	s = strings.TrimPrefix(s, "Failed to parse request template: ")
	if i := strings.IndexByte(s, '\n'); i != -1 {
		s = s[:i]
	}
	return strings.Join(strings.Fields(s), " ")
}

func stripWarn(s string) string {
	var r []string
	for _, l := range strings.Split(s, "\n") {
		if !strings.HasPrefix(l, "WARNING: The requested image's platform") {
			r = append(r, l)
		}
	}
	return strings.Join(r, "\n")
}
