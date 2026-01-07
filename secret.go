package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

const (
	SecretTypeDockerConfigJSON = "kubernetes.io/dockerconfigjson"
	APIVersionV1               = "v1"
	KindSecret                 = "Secret"
	DataKeyDockerConfigJSON    = ".dockerconfigjson"
)

type DockerConfig struct {
	Auths map[string]AuthEntry `json:"auths"`
}

type AuthEntry struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Auth     string `json:"auth"`
}

type Secret struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   Metadata          `yaml:"metadata"`
	Type       string            `yaml:"type"`
	Data       map[string]string `yaml:"data"`
}

type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace,omitempty"`
}

func NewImagePullSecret(registry, user, pass, name, namespace string) (*Secret, error) {
	auth := base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", user, pass))

	cfg := DockerConfig{
		Auths: map[string]AuthEntry{
			registry: {
				Username: user,
				Password: pass,
				Auth:     auth,
			},
		},
	}

	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	dockerCfgB64 := base64.StdEncoding.EncodeToString(jsonBytes)

	secret := Secret{
		APIVersion: APIVersionV1,
		Kind:       KindSecret,
		Type:       SecretTypeDockerConfigJSON,
		Data: map[string]string{
			DataKeyDockerConfigJSON: dockerCfgB64,
		},
		Metadata: Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}

	return &secret, nil
}

// ToYAML converts the Secret struct to a YAML string with proper indentation.
func (s *Secret) ToYAML() (string, error) {
	var buf bytes.Buffer

	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	defer enc.Close()

	if err := enc.Encode(s); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// DecodeDockerConfig decodes the base64-encoded Docker config JSON from the secret data
func (s *Secret) DecodeDockerConfig() (string, error) {
	dockerCfgB64 := s.Data[DataKeyDockerConfigJSON]

	// Base64 → JSON
	dockerCfgJSON, err := base64.StdEncoding.DecodeString(dockerCfgB64)
	if err != nil {
		return "", err
	}

	secret := Secret{
		APIVersion: APIVersionV1,
		Kind:       KindSecret,
		Type:       SecretTypeDockerConfigJSON,
		Data: map[string]string{
			DataKeyDockerConfigJSON: string(dockerCfgJSON),
		},
		Metadata: Metadata{
			Name:      s.Metadata.Name,
			Namespace: s.Metadata.Namespace,
		},
	}

	return secret.ToYAML()
}
