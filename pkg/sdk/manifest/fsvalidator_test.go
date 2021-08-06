package manifest_test

import (
	"context"
	"testing"

	"capact.io/capact/internal/cli/schema"
	"github.com/stretchr/testify/assert"

	"capact.io/capact/pkg/sdk/manifest"

	"github.com/stretchr/testify/require"
)

func TestFilesystemValidator_ValidateFile(t *testing.T) {
	// given
	validator := manifest.NewDefaultFilesystemValidator(&schema.LocalFileSystem{}, "../../../ocf-spec")

	tests := map[string]struct {
		manifestPath      string
		expectedErrorMsgs []string
	}{
		"Invalid Implementation": {
			manifestPath: "testdata/invalid-implementation.yaml",
			expectedErrorMsgs: []string{
				"spec: implements is required",
				"spec: appVersion is required",
			},
		},
		"Valid Implementation": {
			manifestPath:      "testdata/valid-implementation.yaml",
			expectedErrorMsgs: []string{},
		},
		"Invalid JSON Schema in Type": {
			manifestPath: "testdata/invalid-type1.yaml",
			expectedErrorMsgs: []string{
				"type: Must validate at least one schema (anyOf)",
				`type: type must be one of the following: "array", "boolean", "integer", "null", "number", "object", "string"`,
			},
		},
		"Invalid JSON for Type": {
			manifestPath: "testdata/invalid-type2.yaml",
			expectedErrorMsgs: []string{
				"while JSON schema validation: invalid character '}' looking for beginning of object key string",
			},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			// when
			result, err := validator.Do(context.Background(), tc.manifestPath)

			// then
			require.Nil(t, err, "failed to read file: %v", err)
			require.Len(t, result.Errors, len(tc.expectedErrorMsgs))

			if len(result.Errors) > 0 {
				var errMsgs []string
				for _, err := range result.Errors {
					errMsgs = append(errMsgs, err.Error())
				}
				assert.ElementsMatch(t, tc.expectedErrorMsgs, errMsgs)
			}
		})
	}
}
