package handler

import (
	"fmt"
	"testing"
)

func TestFormatDurationHMMSS1(t *testing.T) {
	var seconds int64 = 3600
	want := "1:00:00"
	output, err := formatDurationHMMSS(seconds)
	if err != nil {
		t.Errorf("error computing formatted duration: %v", err)
	} else if output != want {
		t.Fail()
		fmt.Printf("want 1:00:00, got %s\n", output)
	}
}

func TestFormatDurationHMMSS2(t *testing.T) {
	var seconds int64 = 120
	want := "02:00"
	output, err := formatDurationHMMSS(seconds)
	if err != nil {
		t.Errorf("error computing formatted duration: %v", err)
	} else if output != want {
		t.Fail()
		fmt.Printf("want 1:00:00, got %s\n", output)
	}
}
