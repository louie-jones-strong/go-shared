package apptype

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnit_StringToTime(t *testing.T) {

	sampleTime := time.Date(2000, 1, 23, 4, 5, 6, 7, time.UTC)

	tests := []string{
		sampleTime.Format(time.Layout),
		sampleTime.Format(time.ANSIC),
		sampleTime.Format(time.UnixDate),
		sampleTime.Format(time.RubyDate),
		sampleTime.Format(time.RFC822),
		sampleTime.Format(time.RFC822Z),
		sampleTime.Format(time.RFC850),
		sampleTime.Format(time.RFC1123),
		sampleTime.Format(time.RFC1123Z),
		sampleTime.Format(time.RFC3339),
		sampleTime.Format(time.RFC3339Nano),
		sampleTime.Format(time.Kitchen),
		sampleTime.Format(time.Stamp),
		sampleTime.Format(time.StampMilli),
		sampleTime.Format(time.StampMicro),
		sampleTime.Format(time.StampNano),
		sampleTime.Format(time.DateTime),
		sampleTime.Format(time.DateOnly),
		sampleTime.Format(time.TimeOnly),
		"2025-03-12T08:00:00+00:00",
	}

	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {

			_, err := StringToTime(tc)
			assert.NoError(t, err)

		})
	}
}
