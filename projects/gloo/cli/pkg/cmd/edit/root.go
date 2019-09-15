package edit

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/edit"

	editOptions "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/edit/options"
	glooOptions "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/go-utils/cliutils"
	"github.com/spf13/cobra"
)

func RootCmd(opts *glooOptions.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	editFlags := &editOptions.EditOptions{Options: opts}

	cmd := edit.RootCmdWithEditOpts(editFlags, optionsFunc...)

	//cmdutils.MustAddChildCommand(cmd, gloovs.RootCmd(editFlags), virtualservice.RateLimitConfig(editFlags, optionsFunc...))
	//cmd.AddCommand(settings.RootCmd(editFlags, optionsFunc...))
	return cmd
}
