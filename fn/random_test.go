package function

import (
	"testing"

	"github.com/deas/gcp-audit-label/fn/test"
)

// https://github.com/golang/go/issues/46056
func TestRandmon(t *testing.T) {
	print(test.ServiceAccountCreateJson)
}
