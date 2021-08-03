package validate

import (
	"capact.io/capact/internal/cli/client"
	"capact.io/capact/internal/cli/config"
	"capact.io/capact/internal/cli/schema"
	"capact.io/capact/pkg/sdk/manifest"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"io"
	"strings"
)

type Options struct {
	SchemaLocation string
	ServerSide bool
}

type Validation struct {
	hubCli client.Hub
	schemaProvider *schema.Provider
	writer io.Writer
}

func New(writer io.Writer, opts Options) (*Validation, error) {
	server := config.GetDefaultContext()

	var (
		hubCli client.Hub
		err error
	)
	if opts.ServerSide {
		hubCli, err = client.NewHub(server)
		if err != nil {
			return nil, errors.Wrap(err, "while creating Hub client")
		}
	}

	schemaProvider := schema.NewProvider(opts.SchemaLocation)

	return &Validation{
		schemaProvider: schemaProvider,
		hubCli: hubCli,
		writer: writer,
	}, nil
}

func (v *Validation) Run(ctx context.Context, filePaths []string) error {
	validator := manifest.NewFilesystemValidator(v.schemaProvider.FileSystem())

	fmt.Println("Validating files...")

	var errFilePaths []string
	for _, filepath := range filePaths {
		result, err := validator.Do(filepath)

		resultErrs := result.Errors
		if err != nil {
			resultErrs = append(resultErrs, err)
		}

		// TODO: Improve UX (response)

		if len(resultErrs) > 0 {
			color.Red("- %s: FAILED\n", filepath)

			for _, err := range resultErrs {
				color.Red("\t%v", err)
			}

			errFilePaths = append(errFilePaths, filepath)
			continue
		}

		color.Green("- %s: PASSED\n", filepath)
	}

	if len(errFilePaths) > 0 {
		return fmt.Errorf("the following files failed validation:\n%s", strings.Join(errFilePaths, "\n\t"))
	}

	return nil
}
