package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/labbcb/wf/client"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// inputsCmd represents the inputs command
var inputsCmd = &cobra.Command{
	Use:   "inputs",
	Short: "Generate and output a new inputs JSON for this workflow",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := &client.Client{Host: host}
		res, err := c.Describe(args[0], "")
		fatalOnErr(err)

		if !res.Valid {
			log.Fatal(res.Errors)
		}

		inputs := make(map[string]string)
		for _, i := range res.Inputs {
			key := fmt.Sprintf("%s.%s", res.Name, i.Name)

			var extra string
			if i.Optional && i.Default != "" {
				extra = fmt.Sprintf("(optional, default = %s)", i.Default)
			} else if i.Optional {
				extra = "(optional)"
			} else if i.Default != "" {
				extra = fmt.Sprintf("default = %s", i.Default)
			}

			value := fmt.Sprintf("%s%s", i.TypeDisplayName, extra)
			inputs[key] = value
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		fatalOnErr(encoder.Encode(&inputs))
	},
}

func init() {
	rootCmd.AddCommand(inputsCmd)
}
