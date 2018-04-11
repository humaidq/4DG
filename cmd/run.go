package cmd

import (
	"github.com/humaidq/4DG/server"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the 4DG web server",
	Run: func(cmd *cobra.Command, args []string) {
		server.RunServer()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

}
