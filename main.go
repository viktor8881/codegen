package main

import (
	"codegen/command"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "My CLI tool",
	Long:  `A CLI tool with version and codegen commands.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the current version number of this CLI tool.`,
	Run:   command.Version,
}

var codegenCmd = &cobra.Command{
	Use:   "codegen",
	Short: "Generate some code",
	Long:  `Generate some code and save it to a file.`,
	Run:   command.CodeGen,
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(codegenCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
