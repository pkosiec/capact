package clusterpolicy

import (
	"capact.io/capact/pkg/sdk/apis/0.0.1/types"
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

const (
	CurrentAPIVersion        = "0.2.0"
	AnyInterfacePath  string = "cap.*"
)

type ClusterPolicy struct {
	APIVersion string    `json:"apiVersion"`
	Rules      RulesList `json:"rules"`
}

type RulesList []RulesForInterface

type RulesForInterface struct {
	// Interface refers to a given Interface manifest.
	Interface types.ManifestRef `json:"interface"`

	OneOf []Rule `json:"oneOf"`
}

type Rule struct {
	ImplementationConstraints ImplementationConstraints `json:"implementationConstraints,omitempty"`
	InjectTypeInstances       []TypeInstanceToInject    `json:"injectTypeInstances,omitempty"`
}

type ImplementationConstraints struct {
	// Requires refers a specific requirement path and optional revision.
	Requires *[]types.ManifestRef `json:"requires,omitempty"`

	// Attributes refers a specific Attribute by path and optional revision.
	Attributes *[]types.ManifestRef `json:"attributes,omitempty"`

	// Path refers a specific Implementation with exact path.
	Path *string `json:"path,omitempty"`
}

type TypeInstanceToInject struct {
	ID string `json:"id"`

	// TypeRef refers to a given Type.
	TypeRef types.ManifestRef `json:"typeRef"`
}

func (p ClusterPolicy) ToYAMLString() (string, error) {
	bytes, err := yaml.Marshal(&p)

	if err != nil {
		return "", errors.Wrap(err, "while marshalling policy to YAML bytes")
	}

	return string(bytes), nil
}
