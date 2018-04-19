package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Creates a new configuration file.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Warning: This setup wizard will override your current config.toml file!")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
