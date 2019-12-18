package cmd

import (
	"fmt"
	"github.com/labbcb/wf/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// outputsCmd represents the outputs command
var outputsCmd = &cobra.Command{
	Use:   "outputs id",
	Short: "Get the outputs for a workflow",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("host")
		c := client.Client{Host: host}
		json, err := c.Outputs(args[0])
		fatalOnErr(err)
		fmt.Println(json)
	},
}

func init() {
	rootCmd.AddCommand(outputsCmd)
}
