package common

import (
	"capact.io/capact/internal/ptr"
	"capact.io/capact/pkg/sdk/apis/0.0.1/types"
)

// CreateManifestPath create a manifest path based on a manifest type and suffix.
func CreateManifestPath(manifestType types.ManifestKind, suffix string) string {
	suffixes := map[types.ManifestKind]string{
		types.AttributeManifestKind:      "attribute",
		types.TypeManifestKind:           "type",
		types.InterfaceManifestKind:      "interface",
		types.InterfaceGroupManifestKind: "interfaceGroup",
		types.ImplementationManifestKind: "implementation",
	}
	return "cap." + suffixes[manifestType] + "." + suffix
}

// AddRevisionToPath adds revision to manifest path.
func AddRevisionToPath(path string, revision string) string {
	return path + ":" + revision
}

// TODO: Consider moving them out of common

// GetDefaultAttributeMetadata creates a new Metadata object for Attribute and sets default values.
// TODO: Unfortunately generated Attribute type refer to InterfaceMetadata. That should be fixed in the `types` package.
func GetDefaultAttributeMetadata() types.InterfaceMetadata {
	return types.InterfaceMetadata{
		DocumentationURL: defaultURL(),
		SupportURL:       defaultURL(),
		IconURL:          defaultIconURL(),
		Maintainers:      defaultMaintainers(),
	}
}

// GetDefaultImplementationMetadata creates a new Metadata object for Implementation and sets default values.
func GetDefaultImplementationMetadata() types.ImplementationMetadata {
	return types.ImplementationMetadata{
		DocumentationURL: defaultURL(),
		SupportURL:       defaultURL(),
		IconURL:          defaultIconURL(),
		Maintainers:      defaultMaintainers(),
		License:          defaultLicense(),
	}
}

// GetDefaultInterfaceMetadata creates a new Metadata object for Interface and sets default values.
func GetDefaultInterfaceMetadata() types.InterfaceMetadata {
	return types.InterfaceMetadata{
		DocumentationURL: defaultURL(),
		SupportURL:       defaultURL(),
		IconURL:          defaultIconURL(),
		Maintainers:      defaultMaintainers(),
	}
}

func defaultURL() *string {
	return ptr.String("https://example.com")
}

func defaultIconURL() *string {
	return ptr.String("https://example.com/icon.png")
}

func defaultMaintainers() []types.Maintainer {
	return []types.Maintainer{
		{
			Email: "dev@example.com",
			Name:  ptr.String("Example Dev"),
			URL:   ptr.String("https://example.com"),
		},
	}
}

func defaultLicense() types.License {
	return types.License{
		Name: ApacheLicense,
	}
}
