package utils

import (
	"slices"
	"strings"
)

// MaxHistory defines the maximum number of entries to keep in history lists
const MaxHistory = 100

type AppHistory struct {
	Names      []string
	Namespaces []string
	Registries []string
}

func NewAppHistory() *AppHistory {
	return &AppHistory{
		Names:      []string{},
		Namespaces: []string{},
		Registries: []string{},
	}
}

// addValue adds value to the list if not exists, removes overflow
func addValue(list []string, value string) []string {
	// Only add if not exists and it's not empty
	if strings.TrimSpace(value) == "" {
		return list
	}

	if slices.Contains(list, value) {
		return list
	}

	list = append(list, value)

	// remove Overflow
	if len(list) > MaxHistory {
		list = list[1:]
	}

	return list
}

func (h *AppHistory) SetNames(names []string) {
	h.Names = names
}

func (h *AppHistory) SetNamespaces(namespaces []string) {
	h.Namespaces = namespaces
}

func (h *AppHistory) SetRegistries(registries []string) {
	h.Registries = registries
}

// Add... methods add new entries to the history lists
func (h *AppHistory) AddSecretName(name string) {
	if IsK8sNameValid(name) {
		h.Names = addValue(h.Names, name)
	}
}

func (h *AppHistory) AddNamespace(namespace string) {
	// Only add valid k8s names
	if IsK8sNameValid(namespace) {
		h.Namespaces = addValue(h.Namespaces, namespace)
	}
}

func (h *AppHistory) AddRegistry(registry string) {
	h.Registries = addValue(h.Registries, registry)
}

// Sorted... methods return sorted copies of the history lists
func (h *AppHistory) SortedNames() []string {
	copyList := append([]string{}, h.Names...)
	slices.Sort(copyList)
	return copyList
}

func (h *AppHistory) SortedNamespaces() []string {
	copyList := append([]string{}, h.Namespaces...)
	slices.Sort(copyList)
	return copyList
}

func (h *AppHistory) SortedRegistries() []string {
	copyList := append([]string{}, h.Registries...)
	slices.Sort(copyList)
	return copyList
}

func (h *AppHistory) Clear() {
	h.Names = []string{}
	h.Namespaces = []string{}
	h.Registries = []string{}
}

func (h *AppHistory) IsEmpty() bool {
	return len(h.Names) == 0 && len(h.Namespaces) == 0 && len(h.Registries) == 0
}
