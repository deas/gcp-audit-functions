package cmd

import (
	"context"
	"errors"

	function "github.com/deas/gcp-audit-label/fn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1"
)

var (
	instanceActions map[string]func(context.Context, *assetpb.ResourceSearchResult, map[string]string) error = map[string]func(context.Context, *assetpb.ResourceSearchResult, map[string]string) error{
		"start":      function.Start,
		"stop":       function.Stop,
		"start-stop": function.StartStop,
	}
	/* commands map[string]interface{} = map[string]interface{}{}*/
)

var searchCmd = &cobra.Command{
	Use: "search",

	Short: "search and proces compute instances",

	PreRun: func(cmd *cobra.Command, args []string) {
		// cmd.Help()
		viper.BindPFlags(cmd.Flags())
		// function.Logger.Info(context.Background(), "Starting search")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetString("scope") == "" {
			return errors.New("`scope` not set")
		}
		req := &assetpb.SearchAllResourcesRequest{
			Scope: viper.GetString("scope"),
			Query: viper.GetString("query"),
			AssetTypes: []string{
				"compute.googleapis.com/Instance",
				// "cloudresourcemanager.googleapis.com/Project",
			},
		}
		return function.Search(context.Background(), req, instanceActions[viper.GetString("command")], nil)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().String("query", "", "asset query")
	searchCmd.Flags().String("scope", "", "search scope")
	searchCmd.Flags().String("command", "", "start/stop/start-stop")

}
