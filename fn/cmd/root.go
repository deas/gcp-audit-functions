package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"cdr.dev/slog"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/cloudevents/sdk-go/v2/event"
	function "github.com/deas/gcp-audit-label/fn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	debug bool
	// TODO: Should probably do reflection?
	eventFunctions map[string]func(context.Context, event.Event) error = map[string]func(context.Context, event.Event) error{
		"LabelEvent": function.LabelEvent,
	}
	pubsubFunctions map[string]func(context.Context, function.PubSubMessage) error = map[string]func(context.Context, function.PubSubMessage) error{
		"LabelPubSub":   function.LabelPubSub,
		"HardenPubSub":  function.HardenPubSub,
		"ActionsPubSub": function.ActionsPubSub,
		// "StartPubSub":  function.StartPubSub,
		// "StopPubSub":   function.StopPubSub,
	}
)

var rootCmd = &cobra.Command{
	Use: "gcp-housekeeper",

	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlags(cmd.Flags())
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		// cmd.Help()
		os.Setenv(fmt.Sprintf("%s_%s", function.EnvPrefix, "LOGGER"), "human")
		logger := function.NewLogger()
		ctx := context.Background()
		log.SetOutput(slog.Stdlib(ctx, logger).Writer())
		// Serving multiple functions locally from a single server instance #109
		// https://github.com/GoogleCloudPlatform/functions-framework-go/issues/109
		fn := viper.GetString("function")
		if fn == "" {
			return errors.New("`function` not set")
		}
		if strings.HasSuffix(fn, "PubSub") {
			funcframework.RegisterEventFunctionContext(ctx, "/", pubsubFunctions[fn])
		} else {
			funcframework.RegisterCloudEventFunctionContext(ctx, "/", eventFunctions[fn])
		}
		logger.Info(ctx, fmt.Sprintf("Starting function framework service for function %s", fn))
		if err := funcframework.Start(viper.GetString("port")); err != nil {
			log.Fatal(ctx, "funcframework.Start") //: %v\n", err)
		}
		return nil
	},
	// PersistentPreRun: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "")
	rootCmd.Flags().String("function", "", "the function to call")
	rootCmd.Flags().String("port", "8080", "the port to serve on")

	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	viper.SetEnvPrefix(function.EnvPrefix)
	viper.AutomaticEnv()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
