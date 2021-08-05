package manifest

import (
	hubpublicgraphql "capact.io/capact/pkg/hub/api/graphql/public"
	"fmt"
	"github.com/pkg/errors"
)

type Hub interface {
	CheckTypeRevisionsExist(typeRefs []hubpublicgraphql.TypeReference) (map[hubpublicgraphql.TypeReference]bool, error)
}

func checkTypeRevisionsExist(hub Hub, typeRefsToCheck []hubpublicgraphql.TypeReference) (ValidationResult, error) {
	res, err := hub.CheckTypeRevisionsExist(typeRefsToCheck)
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
