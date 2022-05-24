package cmd

import (
	"context"
	"fmt"
	stdlog "log"
	"os"

	"cdr.dev/slog"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	function "github.com/deas/gcp-audit-label/fn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	debug bool
)

// TODO: Cleanup
var rootCmd = &cobra.Command{
	Use: "gcp-audit-label",

	Run: func(cmd *cobra.Command, args []string) {
		// cmd.Help()
		os.Setenv(fmt.Sprintf("%s_%s", function.EnvPrefix, "LOGGER"), "human")
		log := function.NewLogger()
		ctx := context.Background()
		stdlog.SetOutput(slog.Stdlib(ctx, log).Writer())
		log.Info(ctx, "Starting gcp-audit-label")
		if err := funcframework.RegisterEventFunctionContext(ctx, "/label-event", function.LabelEvent); err != nil {
			// 	log.Critical(ctx, "funcframework.RegisterEventFunctionContext")
			log.Fatal(ctx, "funcframework.RegisterEventFunctionContext") //: %v\n", err)
		}
		if err := funcframework.RegisterEventFunctionContext(ctx, "/label-pubsub", function.LabelPubSub); err != nil {
			log.Fatal(ctx, "funcframework.RegisterEventFunctionContext") //: %v\n", err)
		}
		// Use PORT environment variable, or default to 8080.
		port := "8080"
		if envPort := os.Getenv("PORT"); envPort != "" {
			port = envPort
		}
		if err := funcframework.Start(port); err != nil {
			log.Fatal(ctx, "funcframework.Start") //: %v\n", err)
		}
	},
	// PersistentPreRun: func(cmd *cobra.Command, args []string) {},
}

func init() {
	// rootCmd.PersistentFlags().String("foo-provider", "foo", "")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "")

	// viper.BindPFlag("foo-provider", rootCmd.PersistentFlags().Lookup("foo-provider"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	viper.SetEnvPrefix(function.EnvPrefix)
	viper.AutomaticEnv()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		stdlog.Println(err)
		os.Exit(1)
	}
}
