package manifest

import (
	"capact.io/capact/pkg/sdk/apis/0.0.1/types"
	"fmt"
)

// FileValidator is a interface, with the Do method.
// Do validates the manifest in filepath and return a ValidationResult.
// If other, not manifest related errors occur, it will return an error.
type FileValidator interface {
	Do(filepath string) (ValidationResult, error)
}

// ValidationResult hold the result of the manifest validation.
type ValidationResult struct {
	Errors []error
}

type ValidationError struct {
	Path string
	UnderlyingError error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("while validating %q: %s", e.Path, e.UnderlyingError.Error())
}

// Valid returns true, if the manifest contains no errors.
func (r *ValidationResult) Valid() bool {
	return len(r.Errors) == 0
}

func newValidationResult(errs ...error) ValidationResult {
	return ValidationResult{
		Errors: errs,
	}
}

type PartialValidator interface {
	Do(metadata types.ManifestMetadata, yamlBytes []byte) (ValidationResult, error)
	Name() string
}
