package manifest

import (
	hubpublicgraphql "capact.io/capact/pkg/hub/api/graphql/public"
	"capact.io/capact/pkg/sdk/apis/0.0.1/types"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
)

type RemoteImplementationValidator struct {
	hub Hub
}

func NewRemoteImplementationValidator(hub Hub) *RemoteImplementationValidator {
	return &RemoteImplementationValidator{
		hub: hub,
	}
}

func (v *RemoteImplementationValidator) Do(ctx context.Context, _ types.ManifestMetadata, jsonBytes []byte) (ValidationResult, error) {
	var entity types.Implementation
	err := json.Unmarshal(jsonBytes, &entity)
	if err != nil {
		return ValidationResult{}, errors.Wrap(err, "while unmarshalling JSON into Implementation type")
	}

	var manifestRefsToCheck []hubpublicgraphql.ManifestReference

	// Attributes
	for path, attr := range entity.Metadata.Attributes {
		manifestRefsToCheck = append(manifestRefsToCheck, hubpublicgraphql.ManifestReference{
			Path:     path,
			Revision: attr.Revision,
		})
	}

	// AdditionalInput
	if entity.Spec.AdditionalInput != nil {

		// Parameters
		additionalInputParams := entity.Spec.AdditionalInput.Parameters
		if additionalInputParams != nil && additionalInputParams.TypeRef != nil {
			manifestRefsToCheck = append(manifestRefsToCheck, hubpublicgraphql.ManifestReference{
				Path:     additionalInputParams.TypeRef.Path,
				Revision: additionalInputParams.TypeRef.Revision,
			})
		}

		// TypeInstances
		for _, ti := range entity.Spec.AdditionalInput.TypeInstances {
			manifestRefsToCheck = append(manifestRefsToCheck, hubpublicgraphql.ManifestReference{
				Path:     ti.TypeRef.Path,
				Revision: ti.TypeRef.Revision,
			})
		}
	}

	// AdditionalOutput
	if entity.Spec.AdditionalOutput != nil {
		for _, ti := range entity.Spec.AdditionalOutput.TypeInstances {
			if ti.TypeRef == nil {
				continue
			}

			manifestRefsToCheck = append(manifestRefsToCheck, hubpublicgraphql.ManifestReference{
				Path:     ti.TypeRef.Path,
				Revision: ti.TypeRef.Revision,
			})
		}
	}

	// Implements

	// Requires

	// Imports

	return checkManifestRevisionsExist(ctx, v.hub, manifestRefsToCheck)
}

func (v *RemoteImplementationValidator) Name() string {
	return "RemoteImplementationValidator"
}
