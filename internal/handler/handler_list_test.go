package handler

import (
	"testing"
)

func TestFormatDurationHMMSS(t *testing.T) {
	var seconds int64 = 120
	want := "02:00"
	output, err := formatDurationHMMSS(seconds)
	if err != nil {
		t.Errorf("error computing formatted duration: %v", err)
	} else if output != want {
		t.Fail()
	}
}
