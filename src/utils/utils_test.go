package utils

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestShouldProperlyEncodeLabel(t *testing.T) {
	givenLabels := []string{"feature/test", "feature/test$&", "feature/very-long-name-0123456789-0123456789-0123456789-0123456789"}
	expectedLabels := []string{"initium-feature-test", "initium-feature-test--", "initium-feature-very-long-name-0123456789-0123456789-0123456-z"}

	for idx, label := range givenLabels {
		encodedLabel := EncodeRFC1123(label)
		assert.Check(t, encodedLabel == expectedLabels[idx], "Expected %s, got %s", expectedLabels[idx], encodedLabel)
	}
}
