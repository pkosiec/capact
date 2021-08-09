package validate

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"io"
	"strings"
	"sync"

	"capact.io/capact/internal/cli/client"
	"capact.io/capact/internal/cli/config"
	"capact.io/capact/internal/cli/schema"
	"capact.io/capact/pkg/sdk/manifest"
	"github.com/pkg/errors"
)

// Options struct defines validation options for OCF manifest validation.
type Options struct {
	SchemaLocation string
	ServerSide     bool
	Verbose        bool
	MaxConcurrency int
}

// ValidationResult defines a validation error.
type ValidationResult struct {
	Path   string
	Errors []error
}

// IsSuccess returns if there were any validation errors.
func (r *ValidationResult) IsSuccess() bool {
	return len(r.Errors) == 0
}

// Error returns error message based on the ValidationResult data.
func (r *ValidationResult) Error() string {
	if r == nil || len(r.Errors) == 0 {
		return ""
	}

	var errMsgs []string
	for _, err := range r.Errors {
		errMsgs = append(errMsgs, err.Error())
	}

	return fmt.Sprintf("%q:\n    * %s\n", r.Path, strings.Join(errMsgs, "\n    * "))
}

// Validation defines OCF manifest validation operation.
type Validation struct {
	hubCli      client.Hub
	writer      io.Writer
	verbose     bool
	maxWorkers  int
	validatorFn func() manifest.FileSystemValidator
}

// New creates new Validation.
func New(writer io.Writer, opts Options) (*Validation, error) {
	server := config.GetDefaultContext()
	fs, ocfSchemaRootPath := schema.NewProvider(opts.SchemaLocation).FileSystem()

	var (
		hubCli        client.Hub
		err           error
		validatorOpts []manifest.ValidatorOption
	)

	if opts.ServerSide {
		hubCli, err = client.NewHub(server)
		if err != nil {
			return nil, errors.Wrap(err, "while creating Hub client")
		}

		validatorOpts = append(validatorOpts, manifest.WithRemoteChecks(hubCli))
	}

	return &Validation{
		validatorFn: func() manifest.FileSystemValidator {
			return manifest.NewDefaultFilesystemValidator(fs, ocfSchemaRootPath, validatorOpts...)
		},
		hubCli:     hubCli,
		writer:     writer,
		verbose:    opts.Verbose,
		maxWorkers: opts.MaxConcurrency,
	}, nil
}

// Run runs validation across all JSON validators.
func (v *Validation) Run(ctx context.Context, filePaths []string) error {
	fileNoun := properNounFor("file", len(filePaths))

	fmt.Fprintf(v.writer, "Validating %s...\n", fileNoun)

	var workersCount = v.maxWorkers
	if len(filePaths) < workersCount {
		workersCount = len(filePaths)
	}

	jobsCh := make(chan string, len(filePaths))
	resultsCh := make(chan ValidationResult, len(filePaths))

	var wg sync.WaitGroup
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		wrker := newWorker(&wg, v.validatorFn())
		go wrker.Do(ctx, jobsCh, resultsCh)
	}

	for _, filepath := range filePaths {
		jobsCh <- filepath
	}
	close(jobsCh)

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	var errsNumber int
	for res := range resultsCh {
		resultErrs := res.Errors
		if len(resultErrs) > 0 {
			errsNumber += len(resultErrs)
			var prefix string
			if v.verbose {
				prefix = fmt.Sprintf("%s ", color.RedString("✗"))
			}
			fmt.Fprintf(v.writer, "- %s%s\n", prefix, res.Error())
			continue
		}

		if v.verbose {
			fmt.Fprintf(v.writer, "- %s %q\n", color.GreenString("✓"), res.Path)
		}
	}

	fmt.Fprintf(v.writer, "Validated %d %s in total.\n", len(filePaths), fileNoun)

	if errsNumber > 0 {
		errNoun := properNounFor("error", errsNumber)
		return fmt.Errorf("detected %d validation %s", errsNumber, errNoun)
	}

	fmt.Fprintf(v.writer, "🚀 No errors detected.\n")
	return nil
}

type worker struct {
	wg        *sync.WaitGroup
	validator manifest.FileSystemValidator
}

func newWorker(wg *sync.WaitGroup, validator manifest.FileSystemValidator) *worker {
	return &worker{wg: wg, validator: validator}
}

func (w *worker) Do(ctx context.Context, jobCh <-chan string, resultCh chan<- ValidationResult) {
	defer w.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case filePath, ok := <-jobCh:
			if !ok {
				return
			}
			res, err := w.validator.Do(ctx, filePath)
			resultErrs := res.Errors
			if err != nil {
				resultErrs = append(resultErrs, err)
			}

			resultCh <- ValidationResult{
				Path:   filePath,
				Errors: resultErrs,
			}
		}
	}
}

func properNounFor(str string, numberOfItems int) string {
	if numberOfItems == 1 {
		return str
	}

	return str + "s"
}
