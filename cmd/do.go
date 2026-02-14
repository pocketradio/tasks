package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Marks a task as complete",
	Run: func(cmd *cobra.Command, args []string) {

		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg) // akin to parseint
			if err != nil {
				fmt.Println("failed to parse the arg: ", arg)
			} else {
				ids = append(ids, id)
			}
		}

		fmt.Println(ids) // prints out all IDS after checking if they're numbers.
	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}
