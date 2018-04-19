package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/humaidq/4DG/models"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Creates a new configuration file.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Warning: This setup wizard will override your current config.toml file!\nEnter label for effect followed by pin numbers.")

		var effAmt int
		fmt.Print("Amount of effects: ")
		fmt.Scan(&effAmt)
		effects := make([]models.Effect, effAmt)

		for i := 0; i < len(effects); i++ {
			var name string
			var pinAmt int
			fmt.Println("Label #", (i + 1))
			fmt.Print("   Name: ")
			fmt.Scanln(&name)
			if name == "" {
				break
			}
			fmt.Print("   Amount of pins: ")
			fmt.Scan(&pinAmt)
			pins := make([]int, pinAmt)
			for j := 0; j < len(pins); j++ {
				fmt.Print("   Pin: ")
				fmt.Scan(&pins[j])
			}
			effects[i] = models.Effect{EffectName: name, Pins: pins}
		}

		config := models.Config{Effects: effects}

		// Convert config struct to TOML format
		buf := new(bytes.Buffer)
		if err := toml.NewEncoder(buf).Encode(config); err != nil {
			log.Fatal(err)
		}

		// Save file
		ioutil.WriteFile("config.toml", buf.Bytes(), 0644)
		fmt.Println("Saved new config.toml")

	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
