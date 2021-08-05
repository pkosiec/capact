package manifest

import (
	hubpublicgraphql "capact.io/capact/pkg/hub/api/graphql/public"
	"capact.io/capact/pkg/sdk/apis/0.0.1/types"
	"encoding/json"
	"github.com/pkg/errors"
)


type RemoteInterfaceValidator struct {
	hub Hub
}

func NewRemoteInterfaceValidator(hub Hub) *RemoteInterfaceValidator {
	return &RemoteInterfaceValidator{
		hub: hub,
	}
}

func (v *RemoteInterfaceValidator) Do(metadata types.ManifestMetadata, jsonBytes []byte) (ValidationResult, error) {
	var entity types.Interface
	err := json.Unmarshal(jsonBytes, &entity)
	if err != nil {
		return ValidationResult{}, errors.Wrap(err, "while unmarshalling JSON into Interface type")
	}

	var typeRefsToCheck []hubpublicgraphql.TypeReference

	// Input Parameters
	for _, param := range entity.Spec.Input.Parameters.ParameterMap {
		if param.TypeRef == nil {
			continue
		}

		typeRefsToCheck = append(typeRefsToCheck, hubpublicgraphql.TypeReference{
			Path:     param.TypeRef.Path,
			Revision: param.TypeRef.Revision,
		})
	}

	// Input TypeInstances
	for _, ti := range entity.Spec.Input.TypeInstances {
		typeRefsToCheck = append(typeRefsToCheck, hubpublicgraphql.TypeReference{
			Path:     ti.TypeRef.Path,
			Revision: ti.TypeRef.Revision,
		})
	}

	// Output TypeInstances
	for _, ti := range entity.Spec.Output.TypeInstances {
		if ti.TypeRef == nil {
			continue
		}

		typeRefsToCheck = append(typeRefsToCheck, hubpublicgraphql.TypeReference{
			Path:     ti.TypeRef.Path,
			Revision: ti.TypeRef.Revision,
		})
	}

	return checkTypeRevisionsExist(v.hub, typeRefsToCheck)
}

func (v *RemoteInterfaceValidator) Name() string {
	return "RemoteInterfaceValidator"
}
