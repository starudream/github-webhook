package config

import (
	"testing"

	"github.com/starudream/go-lib/core/v2/utils/testutil"
)

func TestSave(t *testing.T) {
	testutil.Nil(t, Save())
}
