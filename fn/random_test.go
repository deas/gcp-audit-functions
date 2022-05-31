package function

import (
	"testing"

	"github.com/deas/gcp-audit-label/fn/test"
)

// Mocking
// https://stackoverflow.com/questions/47643192/how-to-mock-functions-in-golang

// //go:embed <path to file in parent directory> doesn't work #4605
// https://github.com/golang/go/issues/46056
func TestRandmon(t *testing.T) {
	print(test.ServiceAccountCreateJson)
}
