package detectlanguage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetect(t *testing.T) {
	detections, err := client.Detect(context.Background(), "labas rytas")

	if assert.NoError(t, err) {
		assert.Equal(t, "lt", detections[0].Language)
		assert.True(t, detections[0].Reliable)
		assert.Greater(t, detections[0].Confidence, float32(0))
	}
}

func TestDetectCode(t *testing.T) {
	code, err := client.DetectCode(context.Background(), "labas rytas")

	if assert.NoError(t, err) {
		assert.Equal(t, "lt", code)
	}
}

func TestDetectCodeFailure(t *testing.T) {
	code, err := client.DetectCode(context.Background(), "")

	assert.EqualError(t, err, "Language not detected")
	assert.Equal(t, code, "")
}

func TestDetectBatch(t *testing.T) {
	query := []string{"labas rytas", "good morning"}
	detections, err := client.DetectBatch(context.Background(), query)

	if assert.NoError(t, err) {
		assert.Equal(t, "lt", detections[0][0].Language)
		assert.Equal(t, "en", detections[1][0].Language)
	}
}
