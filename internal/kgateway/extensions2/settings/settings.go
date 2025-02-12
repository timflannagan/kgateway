package settings

import (
	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	EnableIstioIntegration bool
	EnableAutoMTLS         bool
	StsClusterName         string // TODO: Is this still relevant?
	StsUri                 string // TODO: Is this still relevant?
}

// BuildSettings returns a zero-valued Settings obj if error is encountered when parsing env
func BuildSettings() (*Settings, error) {
	settings := &Settings{}
	if err := envconfig.Process("KGW", settings); err != nil {
		return settings, err
	}
	return settings, nil
}
