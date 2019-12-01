package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/labbcb/wf/client"
	"github.com/labbcb/wf/models"
	"github.com/spf13/cobra"
	"os"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:     "query [id...]",
	Aliases: []string{"ls", "list", "workflows"},
	Short:   "Get workflows matching some criteria",

	Run: func(cmd *cobra.Command, args []string) {
		params := []*models.WorkflowQueryParameter{{
			Submission: submission,
			Start:      start,
			End:        end,
		}}
		for _, s := range status {
			params = append(params, &models.WorkflowQueryParameter{Status: s})
		}
		for _, n := range name {
			params = append(params, &models.WorkflowQueryParameter{Name: n})
		}
		for _, i := range args {
			params = append(params, &models.WorkflowQueryParameter{ID: i})
		}

		c := &client.Client{Host: host}
		res, err := c.Query(params)
		fatalOnErr(err)

		if format == "json" {
			fatalOnErr(json.NewEncoder(os.Stdout).Encode(&res))
		} else {
			for _, r := range res.Results {
				fmt.Println(r)
			}
		}
	},
}

var submission, start, end string
var status, name []string

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().StringVar(&submission, "submission", "", "Returns only workflows with an equal or later submission datetime")
	queryCmd.Flags().StringVar(&start, "start", "", "Returns only workflows with an equal or later start datetime")
	queryCmd.Flags().StringVar(&end, "end", "", "Returns only workflows with an equal or earlier end datetime")
	queryCmd.Flags().StringArrayVar(&status, "status", nil, "Returns only workflows with the specified status")
	queryCmd.Flags().StringArrayVar(&name, "name", nil, "Returns only workflows with the specified name")
}
