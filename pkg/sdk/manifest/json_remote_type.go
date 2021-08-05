package manifest

import (
	hubpublicgraphql "capact.io/capact/pkg/hub/api/graphql/public"
	"capact.io/capact/pkg/sdk/apis/0.0.1/types"
	"encoding/json"
	"github.com/pkg/errors"
)


type RemoteTypeValidator struct {
	hub Hub
}

func NewRemoteTypeValidator(hub Hub) *RemoteTypeValidator {
	return &RemoteTypeValidator{
		hub: hub,
	}
}

func (v *RemoteTypeValidator) Do(metadata types.ManifestMetadata, jsonBytes []byte) (ValidationResult, error) {
	var entity types.Type
	err := json.Unmarshal(jsonBytes, &entity)
	if err != nil {
		return ValidationResult{}, errors.Wrap(err, "while unmarshalling JSON into Type type")
	}

	var typeRefsToCheck []hubpublicgraphql.TypeReference

	// Attributes
	for path, attr := range entity.Metadata.Attributes {
		typeRefsToCheck = append(typeRefsToCheck, hubpublicgraphql.TypeReference{
			Path:     path,
			Revision: attr.Revision,
		})
	}

	return checkTypeRevisionsExist(v.hub, typeRefsToCheck)
}

func (v *RemoteTypeValidator) Name() string {
	return "RemoteTypeValidator"
}
