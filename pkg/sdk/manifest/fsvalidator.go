package manifest

import (
	"capact.io/capact/pkg/sdk/apis/0.0.1/types"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"sigs.k8s.io/yaml"
)

// FilesystemManifestValidator validates manifests using a OCF specification, which is read from a filesystem.
type FilesystemManifestValidator struct {
	commonValidators []JSONValidator
	kindValidators map[types.ManifestKind][]JSONValidator
}

// NewFilesystemValidator returns a new FilesystemManifestValidator.
func NewFilesystemValidator(opts ...ValidatorOption) FileValidator {
	fsValidator := &FilesystemManifestValidator{}

	for _, opt := range opts {
		opt(fsValidator)
	}

	return fsValidator
}

// Do validates a manifest.
func (v *FilesystemManifestValidator) Do(path string) (ValidationResult, error) {
	yamlBytes, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return newValidationResult(), err
	}

	metadata, err := loadManifestMetadata(yamlBytes)
	if err != nil {
		return newValidationResult(errors.Wrap(err, "failed to read manifest metadata")), err
	}

	jsonBytes, err := yaml.YAMLToJSON(yamlBytes)
	if err != nil {
		return newValidationResult(errors.Wrap(err, "cannot convert YAML manifest to JSON")), err
	}

	validators := append (v.commonValidators, v.kindValidators[metadata.Kind]...)

	var validationErrs []error
	for _, validator := range validators {
		res, err := validator.Do(metadata, jsonBytes)
		if err != nil {
			validationErrs = append(validationErrs, errors.Wrapf(err, "while running validator %s", validator.Name()))
		}

		validationErrs = append(validationErrs, res.Errors...)
	}

	return newValidationResult(validationErrs...), nil
}

func loadManifestMetadata(yamlBytes []byte) (types.ManifestMetadata, error) {
	mm := types.ManifestMetadata{}
	err := yaml.Unmarshal(yamlBytes, &mm)
	if err != nil {
		return mm, err
	}
	return mm, nil
}
