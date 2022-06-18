package function

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/deas/gcp-audit-label/fn/test"
	"github.com/stretchr/testify/assert"
)

// Mocking
// https://stackoverflow.com/questions/47643192/how-to-mock-functions-in-golang

// //go:embed <path to file in parent directory> doesn't work #4605
// https://github.com/golang/go/issues/46056
func TestUnmarshal(t *testing.T) {
	actionsSearch := &ActionsSearch{}
	json.Unmarshal([]byte(test.AssetSearchStartStopJson), &actionsSearch)
	log.Printf("%+v,n", actionsSearch)
}

func TestStartStop(t *testing.T) {
	ssLabel := "-_2-4_4-6_6-8_10-12_12-23_-" // Monday is 0! (not Sunday)
	// now, _ := Now("Europe/Berlin")
	// nowDay := int(now.Weekday())
	// nowHour := now.Hour()
	// act, _ := LabelAction(context.Background(), ssLabel, nowDay, nowHour)
	act, _ := LabelAction(context.Background(), ssLabel, 0, 2)
	assert.Equal(t, "stop", act)
	act, _ = LabelAction(context.Background(), ssLabel, 1, 3)
	assert.Equal(t, "start", act)
	// fmt.Printf("day %d, hour = %d, action = %s\n", nowDay, nowHour, act)
}
