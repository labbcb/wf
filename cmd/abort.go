package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/labbcb/wf/client"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// abortCmd represents the abort command
var abortCmd = &cobra.Command{
	Use:   "abort id...",
	Short: "Abort a running workflow",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := &client.Client{Host:host}
		for _, id := range args {
			res, err := c.Abort(id)
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
	rootCmd.AddCommand(abortCmd)
}
