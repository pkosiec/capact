package cmd

import (
	"capact.io/capact/internal/cli/validate"
	"os"

	"capact.io/capact/internal/cli"
	"capact.io/capact/internal/cli/heredoc"
	"github.com/spf13/cobra"
)

// NewValidate returns a cobra.Command for validating Hub Manifests.
func NewValidate() *cobra.Command {
	var opts validate.Options



	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate OCF manifests",
		Example: heredoc.WithCLIName(`
			# Validate interface-group.yaml file with OCF specification in default location
			<cli> validate ocf-spec/0.0.1/examples/interface-group.yaml
			
			# Validate multiple files inside test_manifests directory
			<cli> validate pkg/cli/test_manifests/*.yaml
			
			# Validate interface-group.yaml file with custom OCF specification location 
			<cli> validate -s my/ocf/spec/directory ocf-spec/0.0.1/examples/interface-group.yaml
			
			# Validate all Hub manifests
			<cli> validate ./manifests/**/*.yaml`, cli.Name),
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			validation, err := validate.New(os.Stdout, opts)
			if err != nil {
				return err
			}

			return validation.Run(cmd.Context(), args)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.SchemaLocation, "schemas", "s", "", "Path to the local directory with OCF JSONSchemas. If not provided, built-in JSONSchemas are used.")
	flags.BoolVar(&opts.ServerSide, "server-side", false, "If enabled, the manifests validation is proceeded against Capact Hub.")

	return cmd
}
