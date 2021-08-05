package manifest

import (
	hubpublicgraphql "capact.io/capact/pkg/hub/api/graphql/public"
	"capact.io/capact/pkg/sdk/apis/0.0.1/types"
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

func (v *RemoteImplementationValidator) Do(metadata types.ManifestMetadata, jsonBytes []byte) (ValidationResult, error) {
	var entity types.Implementation
	err := json.Unmarshal(jsonBytes, &entity)
	if err != nil {
		return ValidationResult{}, errors.Wrap(err, "while unmarshalling JSON into Implementation type")
	}

	var typeRefsToCheck []hubpublicgraphql.TypeReference

	// Attributes
	for path, attr := range entity.Metadata.Attributes {
		typeRefsToCheck = append(typeRefsToCheck, hubpublicgraphql.TypeReference{
			Path:     path,
			Revision: attr.Revision,
		})
	}

	// AdditionalInput - Parameters
	additionalInputParams := entity.Spec.AdditionalInput.Parameters
	if additionalInputParams != nil && additionalInputParams.TypeRef != nil {
		typeRefsToCheck = append(typeRefsToCheck, hubpublicgraphql.TypeReference{
			Path:     additionalInputParams.TypeRef.Path,
			Revision: additionalInputParams.TypeRef.Revision,
		})
	}

	// AdditionalInput - TypeInstances
	for _, ti := range entity.Spec.AdditionalInput.TypeInstances {
		typeRefsToCheck = append(typeRefsToCheck, hubpublicgraphql.TypeReference{
			Path:     ti.TypeRef.Path,
			Revision: ti.TypeRef.Revision,
		})
	}

	// AdditionalOutput
	for _, ti := range entity.Spec.AdditionalOutput.TypeInstances {
		if ti.TypeRef == nil {
			continue
		}

		typeRefsToCheck = append(typeRefsToCheck, hubpublicgraphql.TypeReference{
			Path:     ti.TypeRef.Path,
			Revision: ti.TypeRef.Revision,
		})
	}

	// Implements

	// Requires

	// Imports

	return checkTypeRevisionsExist(v.hub, typeRefsToCheck)
}

func (v *RemoteImplementationValidator) Name() string {
	return "RemoteImplementationValidator"
}
