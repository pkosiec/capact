package manifest

import "capact.io/capact/pkg/sdk/apis/0.0.1/types"

type ValidatorOption func(validator *FSValidator)

func WithRemoteChecks(hubCli Hub) ValidatorOption {
	return func(r *FSValidator) {
		r.kindValidators[types.TypeManifestKind] = append(r.kindValidators[types.TypeManifestKind], NewRemoteTypeValidator(hubCli))
		r.kindValidators[types.InterfaceManifestKind] = append(r.kindValidators[types.InterfaceManifestKind], NewRemoteInterfaceValidator(hubCli))
		r.kindValidators[types.ImplementationManifestKind] = append(r.kindValidators[types.ImplementationManifestKind], NewRemoteImplementationValidator(hubCli))
	}
}

func WithKindValidators(kindValidators map[types.ManifestKind][]JSONValidator) ValidatorOption {
	return func(r *FSValidator) {
		r.kindValidators = kindValidators
	}
}
