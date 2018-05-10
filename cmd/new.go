package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/humaidq/4DG/models"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new 4D movie script (interactive)",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize scanner and length variable
		scanner := bufio.NewScanner(os.Stdin)
		var length int

		// Get the movie name and length
		fmt.Print("Movie name is usually followed with the year between paranthesis.\nEnter movie name and year: ")
		scanner.Scan()
		name := scanner.Text()

		fmt.Print("Movie length is showtime represented in minutes rounded down.\nEnter movie length: ")
		fmt.Scan(&length)

		// This will contain the generated movie script
		newMovie := models.FDMovie{MovieName: name, MovieLength: length}

		fmt.Print("Now you will have to define the effects, assuming you already have effect names and pins defined in your configuration file.\n" +
			"The position is the movie timestamp in seconds point milliseconds. Such as 1044.8 is a valid position, 43.82 is not.\n" +
			"The effect name will match the effect name given to the pin in your configuration.\nThe length would be in milliseconds.\nEnter 'q' followed by enter when done.\n\n")

		// Get all the positions and effects
		newMovie.Effects = make(map[string]*models.TimestampEffect)
		for i := 1; ; i++ {
			var pos, eff string
			var effLen int
			fmt.Println("Effect #", i)

			fmt.Print("Position: ")
			fmt.Scanln(&pos)
			if pos == "q" {
				break
			}

			fmt.Print("Effect name: ")
			fmt.Scanln(&eff)
			if eff == "q" {
				break
			}

			fmt.Print("Length of effect (in ms): ")
			fmt.Scanln(&effLen)

			ts := models.TimestampEffect{EffectName: eff, EffectLength: effLen}
			newMovie.Effects[pos] = &ts
		}

		// Convert the whole movie script to TOML format
		buf := new(bytes.Buffer)
		if err := toml.NewEncoder(buf).Encode(newMovie); err != nil {
			log.Fatal(err)
		}

		var fileName string
		fmt.Print("Enter movie file name: ")
		fmt.Scan(&fileName)

		// Add .toml extension if the file does not have an extension
		if !strings.Contains(fileName, ".") {
			fileName = fileName + ".toml"
		}

		// Save file
		ioutil.WriteFile("scripts/"+fileName, buf.Bytes(), 0644)
		fmt.Println("File saved at 'scripts/" + fileName + "'")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
