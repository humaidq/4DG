package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var verCmd = &cobra.Command{
	Use:   "ver",
	Short: "Displays the program version and license",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("4DG v0.1 released under the MIT license.\nCopyright (c) 2018 Humaid AlQassimi.\nhttps://github.com/humaidq/4DG")
	},
}

func init() {
	rootCmd.AddCommand(verCmd)
}
