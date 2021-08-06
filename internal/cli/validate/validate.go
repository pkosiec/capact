package validate

import (
	"capact.io/capact/internal/cli/client"
	"capact.io/capact/internal/cli/config"
	"capact.io/capact/internal/cli/schema"
	"capact.io/capact/pkg/sdk/manifest"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"strings"
)

type Options struct {
	SchemaLocation string
	ServerSide     bool
}

type ValidationError struct {
	Path   string
	Errors []error
}

func (e *ValidationError) Error() string {
	var errMsgs []string
	for _, err := range e.Errors {
		errMsgs = append(errMsgs, err.Error())
	}

	return fmt.Sprintf("%q:\n\t%s\n", e.Path, strings.Join(errMsgs, "\n\t"))
}

type Validation struct {
	hubCli         client.Hub
	validator 	   manifest.FileSystemValidator
	writer         io.Writer
}

func New(writer io.Writer, opts Options) (*Validation, error) {
	server := config.GetDefaultContext()
	fs, ocfSchemaRootPath := schema.NewProvider(opts.SchemaLocation).FileSystem()

	var (
		hubCli client.Hub
		err    error
		validatorOpts []manifest.ValidatorOption
	)

	if opts.ServerSide {
		hubCli, err = client.NewHub(server)
		if err != nil {
			return nil, errors.Wrap(err, "while creating Hub client")
		}

		validatorOpts = append(validatorOpts, manifest.WithRemoteChecks(hubCli))
	}

	validator := manifest.NewDefaultFilesystemValidator(fs, ocfSchemaRootPath, validatorOpts...)

	return &Validation{
		validator: validator,
		hubCli:         hubCli,
		writer:         writer,
	}, nil
}

func (v *Validation) Run(ctx context.Context, filePaths []string) error {
	fileNoun := properNounFor("file", len(filePaths))

	fmt.Fprintf(v.writer, "Validating %s...\n", fileNoun)

	// TODO: Validate files concurrently

	var errs []error
	for _, filepath := range filePaths {
		result, err := v.validator.Do(ctx, filepath)

		resultErrs := result.Errors
		if err != nil {
			resultErrs = append(resultErrs, err)
		}

		if len(resultErrs) > 0 {
			validationErr := &ValidationError{
				Path:   filepath,
				Errors: resultErrs,
			}
			errs = append(errs, validationErr)
			fmt.Fprintf(v.writer, "- %s\n", validationErr.Error())
			continue
		}
	}

	fmt.Fprintf(v.writer, "Validated %d %s in total.\n", len(filePaths), fileNoun)

	if len(errs) > 0 {
		errNoun := properNounFor("error", len(errs))
		return fmt.Errorf("%d validation %s detected.", len(errs), errNoun)
	}

	fmt.Fprintf(v.writer, "🚀 No errors detected.\n")
	return nil
}

func properNounFor(str string, numberOfItems int) string {
	if numberOfItems == 1 {
		return str
	}

	return str + "s"
}
