package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

const version = "v0.1.6"

func Version(cmd *cobra.Command, args []string) {
	fmt.Println("My CLI tool " + version)
}
