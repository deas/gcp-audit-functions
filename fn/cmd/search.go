package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"cdr.dev/slog"
	function "github.com/deas/gcp-audit-label/fn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1"
)

var (
	instanceActions map[string]func(context.Context, *assetpb.ResourceSearchResult) error = map[string]func(context.Context, *assetpb.ResourceSearchResult) error{
		"start": function.Start,
		"stop":  function.Stop,
	}
	/* commands map[string]interface{} = map[string]interface{}{}*/
)

// TODO: Cleanup
var searchCmd = &cobra.Command{
	Use: "search",

	Short: "search and proces compute instances",

	PreRun: func(cmd *cobra.Command, args []string) {
		// cmd.Help()
		viper.BindPFlags(cmd.Flags())
		os.Setenv(fmt.Sprintf("%s_%s", function.EnvPrefix, "LOGGER"), "human")
		logger := function.NewLogger()
		ctx := context.Background()
		log.SetOutput(slog.Stdlib(ctx, logger).Writer())
		logger.Info(ctx, "Starting search")
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
		return function.Search(context.Background(), req, instanceActions[viper.GetString("command")])
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().String("query", "", "asset query")
	searchCmd.Flags().String("scope", "", "search scope")
	searchCmd.Flags().String("command", "", "start or stop")

}
