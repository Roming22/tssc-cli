package integration

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redhat-appstudio/tssc-cli/pkg/config"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

// ImageRegistry represents the image registry integration coordinates. Works with
// different TSSC integrations managing image registry configuration.
type ImageRegistry struct {
	dockerConfig   string // registry credentials (JSON)
	dockerConfigRO string // registry read-only credentials (JSON)
	url            string // API endpoint
	token          string // API token
}

var _ Interface = &ImageRegistry{}

const dockerConfigEx = `{ "auths": { "registry.tld": { "auth": "username" } } }`

// PersistentFlags adds the persistent flags to the informed Cobra command.
func (i *ImageRegistry) PersistentFlags(cmd *cobra.Command) {
	p := cmd.PersistentFlags()

	p.StringVar(&i.dockerConfig, "dockerconfigjson", i.dockerConfig,
		fmt.Sprintf("JSON formatted registry credentials, e.g.: %q",
			dockerConfigEx))
	p.StringVar(
		&i.dockerConfigRO, "dockerconfigjsonreadonly", i.dockerConfigRO,
		fmt.Sprintf("JSON formatted read-only registry credentials, e.g.: %q",
			dockerConfigEx))
	p.StringVar(&i.url, "url", i.url, "Container registry API endpoint.")
	p.StringVar(&i.token, "token", i.token, "Container registry API token.")

	for _, f := range []string{"dockerconfigjson", "url"} {
		if err := cmd.MarkPersistentFlagRequired(f); err != nil {
			panic(err)
		}
	}
}

// SetArgument sets additional arguments to the integration.
func (i *ImageRegistry) SetArgument(_, _ string) error {
	return nil
}

// LoggerWith decorates the logger with the integration flags.
func (i *ImageRegistry) LoggerWith(logger *slog.Logger) *slog.Logger {
	return logger.With(
		"dockerconfigjson-len", len(i.dockerConfig),
		"dockerconfigjsonreadonly-len", len(i.dockerConfigRO),
		"url", i.url,
		"token-len", len(i.token),
	)
}

// Validate validates the integration configuration.
func (i *ImageRegistry) Validate() error {
	err := ValidateJSON("dockerconfigjson", i.dockerConfig)
	if err != nil {
		return err
	}

	if i.dockerConfigRO != "" {
		err = ValidateJSON("dockerconfigjsonreadonly", i.dockerConfigRO)
		if err != nil {
			return err
		}
	}

	return ValidateURL(i.url)
}

// Type returns the type of the integration.
func (i *ImageRegistry) Type() corev1.SecretType {
	return corev1.SecretTypeDockerConfigJson
}

// Data returns the integration data.
func (i *ImageRegistry) Data(
	_ context.Context,
	_ *config.Config,
) (map[string][]byte, error) {
	return map[string][]byte{
		".dockerconfigjson":         []byte(i.dockerConfig),
		".dockerconfigjsonreadonly": []byte(i.dockerConfigRO),
		"url":                       []byte(i.url),
		"token":                     []byte(i.token),
	}, nil
}

// NewContainerRegistry creates a new instance with the default URL.
func NewContainerRegistry(defaultURL string) *ImageRegistry {
	return &ImageRegistry{url: defaultURL}
}
