package term

import (
	"fmt"

	"github.com/peterh/liner"
	"github.com/spf13/cobra"
)

var exitCmd = &cobra.Command{
	Use:           "exit",
	Short:         "exit",
	Aliases:       []string{"logout"},
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(command *cobra.Command, args []string) error {
		fmt.Println("Bye")
		return liner.ErrPromptAborted
	},
}
