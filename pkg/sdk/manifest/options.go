package manifest

import "net/http"

type ValidatorOption func(validator *FilesystemManifestValidator)

func WithOCFSchemaValidator(fs http.FileSystem, schemaRootPath string) ValidatorOption {
	return func(r *FilesystemManifestValidator) {
		r.commonValidators = append(r.commonValidators, NewOCFSchemaValidator(fs, schemaRootPath))
	}
}

