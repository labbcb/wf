package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/labbcb/wf/client"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status id...",
	Short: "Retrieves the current state for a workflow",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := &client.Client{Host:host}
		for _, id := range args {
			res, err := c.Status(id)
			if err != nil {
				log.Println(err)
			}

			if format == "json" {
				fatalOnErr(json.NewEncoder(os.Stdout).Encode(&res))
			} else {
				fmt.Println(res.Status)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
