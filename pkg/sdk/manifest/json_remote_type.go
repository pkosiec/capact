package manifest

import (
	hubpublicgraphql "capact.io/capact/pkg/hub/api/graphql/public"
	"capact.io/capact/pkg/sdk/apis/0.0.1/types"
	"context"
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

func (v *RemoteTypeValidator) Do(ctx context.Context, _ types.ManifestMetadata, jsonBytes []byte) (ValidationResult, error) {
	var entity types.Type
	err := json.Unmarshal(jsonBytes, &entity)
	if err != nil {
		return ValidationResult{}, errors.Wrap(err, "while unmarshalling JSON into Type type")
	}

	var manifestRefsToCheck []hubpublicgraphql.ManifestReference

	// Attributes
	for path, attr := range entity.Metadata.Attributes {
		manifestRefsToCheck = append(manifestRefsToCheck, hubpublicgraphql.ManifestReference{
			Path:     path,
			Revision: attr.Revision,
		})
	}

	return checkManifestRevisionsExist(ctx, v.hub, manifestRefsToCheck)
}

func (v *RemoteTypeValidator) Name() string {
	return "RemoteTypeValidator"
}
