package main

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImagePullSecret(t *testing.T) {
	secretName := "mysecret"
	namespace := "default"
	registry := "docker.io"
	user := "admin"
	pass := "password123"

	secret, err := NewImagePullSecret(registry, user, pass, secretName, namespace)
	assert.NoError(t, err)
	assert.NotNil(t, secret)

	// Check metadata
	assert.Equal(t, secretName, secret.Metadata.Name)
	assert.Equal(t, namespace, secret.Metadata.Namespace)

	// Check type and kind
	assert.Equal(t, SecretTypeDockerConfigJSON, secret.Type)
	assert.Equal(t, KindSecret, secret.Kind)

	// Check Data key exists
	dockerCfgB64, ok := secret.Data[DataKeyDockerConfigJSON]
	assert.True(t, ok)
	assert.NotEmpty(t, dockerCfgB64)

	// Decode and verify content
	jsonBytes, err := base64.StdEncoding.DecodeString(dockerCfgB64)
	assert.NoError(t, err)

	var cfg DockerConfig
	err = json.Unmarshal(jsonBytes, &cfg)
	assert.NoError(t, err)

	entry, ok := cfg.Auths[registry]
	assert.True(t, ok)
	assert.Equal(t, user, entry.Username)
	assert.Equal(t, pass, entry.Password)

	expectedAuth := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
	assert.Equal(t, expectedAuth, entry.Auth)
}

func TestToYAML(t *testing.T) {
	secret, _ := NewImagePullSecret("docker.io", "user", "pass", "mysecret", "default")

	yamlStr, err := secret.ToYAML()
	assert.NoError(t, err)
	assert.Contains(t, yamlStr, "apiVersion: v1")
	assert.Contains(t, yamlStr, "kind: Secret")
	assert.Contains(t, yamlStr, "name: mysecret")
	assert.Contains(t, yamlStr, ".dockerconfigjson:")
}

func TestDecodeDockerConfig(t *testing.T) {
	registry := "docker.io"
	user := "user1"
	pass := "p@ssw0rd"
	name := "testsecret"
	namespace := "default"

	secret, _ := NewImagePullSecret(registry, user, pass, name, namespace)

	yamlStr, err := secret.DecodeDockerConfig()
	assert.NoError(t, err)
	assert.Contains(t, yamlStr, "apiVersion: v1")
	assert.Contains(t, yamlStr, "kind: Secret")
	assert.Contains(t, yamlStr, name)
	assert.Contains(t, yamlStr, DataKeyDockerConfigJSON)
	assert.Contains(t, yamlStr, user)
	assert.Contains(t, yamlStr, pass)
}
