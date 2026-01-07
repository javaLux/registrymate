package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPullSecretNameFormat(t *testing.T) {
	name := GeneratePullSecretName()

	// Regex for Kubernetes valid names
	k8sNameRegex := `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	assert.Regexp(t, k8sNameRegex, name, "Name must be Kubernetes valid name")

	// Regex for Docker-like Format: pullsecret-adjektiv-name-hex
	dockerLikeRegex := `^pullsecret-[a-z]+-[a-z]+-[a-f0-9]{4}$`
	assert.Regexp(t, dockerLikeRegex, name, "Name must be pullsecret-<adjektiv>-<name>-<hex> Format")
}

func TestPullSecretNameMultipleUnique(t *testing.T) {
	names := make(map[string]struct{})

	for i := 0; i < 100; i++ {
		name := GeneratePullSecretName()

		_, exists := names[name]
		assert.False(t, exists, "Name %s was already generated", name)

		names[name] = struct{}{}
	}
}

func TestPullSecretNameContainsAllowedChars(t *testing.T) {
	name := GeneratePullSecretName()

	// Check that only lowercase, digits and '-' are present
	for _, c := range name {
		assert.True(t, (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-', "Invalid char: %c", c)
	}
}

func TestPullSecretNameStartsWithPullSecret(t *testing.T) {
	name := GeneratePullSecretName()
	assert.Regexp(t, `^pullsecret-`, name, "Name have to start with 'pullsecret-'")
}

func TestIsK8sNameValid_Negative(t *testing.T) {
	invalidNames := []string{
		"a-very-long-name-that-exceeds-the-maximum-length-allowed-by-kubernetes-naming-conventions-which-is-two-hundred-and-fifty-three-characters-abcdefghijklmnopqrstuvwxyz-abcdefghijklmnopqrstuvwxyz-abcdefghijklmnopqrstuvwxyz-abcdefghijklmnopqrstuvwxyz-abcdefghij", // too long
		"-starts-dash",        // begins with dash
		"ends-dash-",          // ends with dash
		"contains_underscore", // _ not allowed
		"contains.dot",        // . not allowed
		"contains$money",      // $ not allowed
		"contains!bang",       // ! not allowed
		"contains@at",         // @ not allowed
		"contains space",      // spaces not allowed
		"contains+plus",       // + not allowed
		"contains=equals",     // = not allowed
		"UPPERCASE",           // uppercase letters not allowed
		"äöü",                 // umlauts not allowed
		"中文",                  // non-latin characters not allowed
	}

	for _, name := range invalidNames {
		assert.False(t, IsK8sNameValid(name), "Name should be invalid: %s", name)
	}
}

func TestIsK8sNameValid_Positive(t *testing.T) {
	validNames := []string{
		"a",                    // minimal
		"abc123",               // alphanumeric
		"my-pull-secret",       // with dashes
		"pullsecret-john-a3f2", // Docker-like Name
		"my-secret-123",        // normal
		"n1-n2-n3",             // multiple dashes
	}

	for _, name := range validNames {
		assert.True(t, IsK8sNameValid(name), "Name should be valid: %s", name)
	}
}
