package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.AddCommand(initPlanCmd())
	return cmd
}

// handles the error
func handleError(err error, format string, data ...interface{}) {
	if err == nil {
		return
	}

	out := ""

	if format == "" {
		out = fmt.Sprintf("%s", err)
	} else {
		f := fmt.Sprintf("%s\n", format)
		out = fmt.Sprintf(f, data...)
	}

	fmt.Println(out)
	os.Exit(1)
}
