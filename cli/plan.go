package cli

import (
	"github.com/bhoriuchi/optimus/plan"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func initPlanCmd() *cobra.Command {
	in := &plan.Input{}
	noColor := false

	cmd := &cobra.Command{
		Use:   "plan",
		Short: "plan changes implemented in the change file",
		Run: func(cmd *cobra.Command, args []string) {
			color.NoColor = noColor
			err := plan.Plan(in)
			handleError(err, "")
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&in.Plan, "plan", "", "Current plan JSON file")
	flags.StringVar(&in.State, "state", "", "Current state file")
	flags.StringVar(&in.Change, "change", "", "Change definition file")
	flags.StringVar(&in.Out, "out", "", "Output state file")
	flags.BoolVar(&in.HideUpdates, "hide-updates", false, "Hides update output")
	flags.BoolVar(&noColor, "no-color", false, "Disable color coded output")

	return cmd
}
