package manifest

import (
	hubpublicgraphql "capact.io/capact/pkg/hub/api/graphql/public"
	"context"
	"fmt"
	"github.com/pkg/errors"
)

type Hub interface {
	CheckManifestRevisionsExist(ctx context.Context, manifestRefs []hubpublicgraphql.ManifestReference) (map[hubpublicgraphql.ManifestReference]bool, error)
}

func checkManifestRevisionsExist(ctx context.Context, hub Hub, manifestRefsToCheck []hubpublicgraphql.ManifestReference) (ValidationResult, error) {
	if len(manifestRefsToCheck) == 0 {
		return ValidationResult{}, nil
	}

	res, err := hub.CheckManifestRevisionsExist(ctx, manifestRefsToCheck)
	if err != nil {
		return ValidationResult{}, errors.Wrap(err, "while checking if Type revisions exist")
	}

	var validationErrs []error
	for typeRef, exists := range res {
		if exists {
			continue
		}

		validationErrs = append(validationErrs, fmt.Errorf("the '%s:%s' Type revision doesn't exist in Hub", typeRef.Path, typeRef.Revision))
	}

	return ValidationResult{Errors: validationErrs}, nil
}
