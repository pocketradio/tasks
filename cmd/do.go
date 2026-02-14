package cmd

import (
	"fmt"
	"strconv"
	"task/db"

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

		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println("something went wrong :(", err)
			return
		}

		for _, id := range ids {
			if id <= 0 || id > len(tasks) {
				fmt.Println("\nenter a valid task number.\n")
				continue // cuz other task numbers might be valid ( assuming user wants to delete multiple tasks at the same time)
			}

			task := tasks[id-1]
			err := db.DeleteTask(task.Key)
			if err != nil {
				fmt.Printf("failed to mark task '%d' as completed. Error : %s\n", id, err)
			} else {
				fmt.Printf("\nmarked task \"%d\" as completed\n", id)
			}
		}
		// fmt.Println(ids) // prints out all IDS after checking if they're numbers.
	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}
