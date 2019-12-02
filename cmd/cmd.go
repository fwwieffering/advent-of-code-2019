package cmd

import (
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
}

// Execute is the entrypoint for all commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
