package cmd

import (
	"github.com/spf13/cobra"
	"go-service.codymj.io/internal"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run application",
	Run: func(cmd *cobra.Command, args []string) {
		internal.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
