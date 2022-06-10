package instances_actions_v1

import (
	"encoding/base64"
	"fmt"
	"log"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/stretchr/testify/assert"
)

func TestInstanceActions(t *testing.T) {
	bpt := tft.NewTFBlueprintTest(t)
	bpt.DefineVerify(func(assert *assert.Assertions) {
		bpt.DefaultVerify(assert)

		// gather custom attributes for tests
		project := bpt.GetStringOutput("project_id")
		region := bpt.GetStringOutput("region")
		functionName := bpt.GetStringOutput("function_name")
		data := fmt.Sprintf(`
		{"data":"%s"}`,
			base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`
			{"action": "start",
			 "search": {"scope": "projects/%s",
						"query": "labels.start_daily_does_not_exist:true AND state:TERMINATED",
						"assetTypes": ["compute.googleapis.com/Instance"]}}`,
				project))))
		// call the function directly
		// log.Printf("Calling function %s with data %s", functionName, data)
		op := gcloud.Run(t,
			fmt.Sprintf("functions call %s", functionName),
			gcloud.WithCommonArgs([]string{"--data", data, "--format", "json", "--project", project, "--region", region}),
		)
		assert.NotNil(op)
		log.Printf("Got result %s, error %s", op.Get("result").String(), op.Get("error").String())
		// assert file random string and secret random string is contained in function response
		assert.NotEmpty(op.Get("executionId").String(), "executionId should be present")
		assert.Empty(op.Get("error").String(), "error should be empty")
		// assert.Contains(op.Get("result").String(), randomSecretString, "contains secret random string")
	})

	bpt.Test()
}
