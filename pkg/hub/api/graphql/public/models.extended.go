package graphql

import (
	"fmt"
	"strings"
)

type ManifestReference struct {
	Path     string `json:"path"`
	Revision string `json:"revision"`
}

func (r ManifestReference) GQLQueryName() (string, error) {
	parts := strings.Split(r.Path, ".")
	if len(parts) < 3 {
		return "", fmt.Errorf("path parts for %q cannot be less than 3", r.Path)
	}

	if parts[1] == "core" {
		return parts[2], nil
	}

	return parts[1], nil
}
