package manifest

import (
	"capact.io/capact/pkg/sdk/apis/0.0.1/types"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
)

type TypeValidator struct {}

func NewTypeValidator() *TypeValidator {
	return &TypeValidator{}
}

func (v *TypeValidator) Do(_ context.Context, _ types.ManifestMetadata, jsonBytes []byte) (ValidationResult, error) {
	var typeEntity types.Type
	err := json.Unmarshal(jsonBytes, &typeEntity)
	if err != nil {
		return ValidationResult{}, errors.Wrap(err, "while unmarshalling JSON into Type type")
	}

	jsonSchemaStr := typeEntity.Spec.JSONSchema.Value
	schemaLoader := gojsonschema.NewReferenceLoader("http://json-schema.org/draft-07/schema")
	manifestLoader := gojsonschema.NewStringLoader(jsonSchemaStr)

	jsonSchemaValidationResult, err := gojsonschema.Validate(schemaLoader, manifestLoader)
	if err != nil {
		return newValidationResult(err), nil
	}

	result := ValidationResult{}
	for _, err := range jsonSchemaValidationResult.Errors() {
		result.Errors = append(result.Errors, fmt.Errorf("%v", err.String()))
	}

	return result, nil
}

func (v *TypeValidator) Name() string {
	return "TypeValidator"
}
