package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the image command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

. <(carpenter completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.bash_profile
. <(carpenter completion)
`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
