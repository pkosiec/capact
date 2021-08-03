package manifest

import "capact.io/capact/pkg/sdk/apis/0.0.1/types"

type ValidatorOption func(validator *FilesystemManifestValidator)

func WithServerSideCheck() ValidatorOption {
	return func(r *FilesystemManifestValidator) {
		r.inputTypeInstances = typeInstances
	}
}
