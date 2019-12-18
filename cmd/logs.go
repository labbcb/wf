package cmd

import (
	"fmt"
	"github.com/labbcb/wf/client"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs id",
	Short: "Get the logs for a workflow",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("host")
		c := client.Client{Host: host}
		json, err := c.Logs(args[0])
		fatalOnErr(err)
		fmt.Println(json)
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
