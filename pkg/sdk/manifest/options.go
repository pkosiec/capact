package manifest

import "capact.io/capact/pkg/sdk/apis/0.0.1/types"

type ValidatorOption func(validator *FSValidator)

func WithCommonValidators(validators ...JSONValidator) ValidatorOption {
	return func(r *FSValidator) {
		r.commonValidators = append(r.commonValidators, validators...)
	}
}

func WithKindValidators(kind types.ManifestKind, validators ...JSONValidator) ValidatorOption {
	return func(r *FSValidator) {
		r.kindValidators[kind] = append(r.kindValidators[kind], validators...)
	}
}


