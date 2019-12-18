package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/labbcb/wf/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:     "validate",
	Aliases: []string{"check"},
	Short:   "Validate a workflow source file",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workflow := args[0]

		host := viper.GetString("host")
		c := &client.Client{Host: host}
		if len(inputs) > 0 {
			for _, i := range inputs {
				validate(c, workflow, i)
			}
		} else {
			validate(c, workflow, "")
		}
	},
}

func validate(c *client.Client, workflow, inputs string) {
	res, err := c.Describe(workflow, inputs)
	fatalOnErr(err)

	format := viper.GetString("format")
	if format == "json" {
		fatalOnErr(json.NewEncoder(os.Stdout).Encode(&res))
	} else {
		if res.Valid {
			fmt.Println("Valid")
		} else {
			fmt.Println("Invalid")
			for _, e := range res.Errors {
				fmt.Fprintln(os.Stderr, e)
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringArrayVarP(&inputs, "inputs", "i", nil, "JSON or YAML file containing the inputs as an object")
}
