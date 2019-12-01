package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/labbcb/wf/client"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit workflow",
	Aliases: []string{"run"},
	Short: "Submits a workflow to Cromwell",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workflow := args[0]

		c := &client.Client{Host: host}
		res, err := c.Submit(workflow, inputs, imports, options)
		fatalOnErr(err)

		if res.Status != "Submitted" {
			log.Fatal(res.Status)
		}

		if format == "json" {
			fatalOnErr(json.NewEncoder(os.Stdout).Encode(&res))
		} else {
			fmt.Println(res.ID)
		}
	},
}

var inputs, imports []string
var options string

func init() {
	rootCmd.AddCommand(submitCmd)

	submitCmd.Flags().StringArrayVarP(&inputs, "inputs", "i", nil, "JSON or YAML file containing the inputs as an object")
	submitCmd.Flags().StringVarP(&options, "options", "o", "", "JSON file containing configuration options for the execution of this workflow")
	submitCmd.Flags().StringArrayVarP(&imports, "imports", "p", nil, "Workflow source files that are used to resolve local imports")
}
