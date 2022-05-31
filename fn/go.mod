module github.com/deas/gcp-audit-label/fn

go 1.16

require github.com/cloudevents/sdk-go/v2 v2.6.1

require (
	cloud.google.com/go/compute v1.6.1
	github.com/GoogleCloudPlatform/functions-framework-go v1.5.3
	google.golang.org/api v0.75.0
	google.golang.org/genproto v0.0.0-20220518221133-4f43b3371335
	google.golang.org/protobuf v1.28.0
)

require (
	cdr.dev/slog v1.4.1
	cloud.google.com/go/asset v1.2.0
	github.com/google/uuid v1.3.0 // indirect
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.11.0
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5
)
