package cmd

import (
	two "advent-of-code-2019/cmd/2"
	five "advent-of-code-2019/cmd/5"
	three "advent-of-code-2019/cmd/3"
	four "advent-of-code-2019/cmd/4"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	one "advent-of-code-2019/cmd/1"
)

var rootCmd = &cobra.Command{
	Use:   "advent",
	Short: "advent runs the advent of code 2019 challenges",
	Long: `advent runs the advent of code 2019 challenges
	https://adventofcode.com/2019/`,
}

func init() {
	rootCmd.AddCommand(one.One)
	rootCmd.AddCommand(two.Two)
	rootCmd.AddCommand(three.Three)
	rootCmd.AddCommand(four.Four)
	rootCmd.AddCommand(five.Five)
}

// Execute is the entrypoint for all commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
