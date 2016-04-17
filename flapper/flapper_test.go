package flapper

import (
	"testing"
	"time"
)

func TestAddStateChange(t *testing.T) {
	flapper := NewFlapper(5, 30)
	flapper.NoteStateChange("test")

	total := flapper.services["test"].Total()
	if total != 1 {
		t.Errorf("Failed to update counter by noting state change. Expected %d, found %d",
			1, total)
	}
}

func TestCheckStateForNonExistantService(t *testing.T) {
	flapper := NewFlapper(5, 30)

	if flapper.IsFlapping("service-name", true) {
		t.Errorf("Should not report flapping for non-reported service")
	}
}

func TestSimpleFlapDetection(t *testing.T) {
	flapper := NewFlapper(5, 30)

	for i := uint(0); i < flapper.max; i++ {
		flapper.NoteStateChange("test")
	}

	if !flapper.IsFlapping("test", false) {
		t.Errorf("Sould be flapping when 'max' state changes reported")
	}
}

func TestSpacedFlapDetection(t *testing.T) {
	flapper := NewFlapper(10, 10)

	for i := uint(0); i < flapper.max; i++ {
		flapper.NoteStateChange("test")
		time.Sleep(500 * time.Millisecond)
	}

	if !flapper.IsFlapping("test", true) {
		t.Errorf("Should be flapping when 'max' state changes reported in less than 'duration'")
	}
}

func TestFlapDetectionResetsOnSlidingOneSecondWindows(t *testing.T) {
	flapper := NewFlapper(10, 2)

	for i := uint(0); i < flapper.max; i++ {
		if i == (flapper.max / uint(2)) {
			time.Sleep(1 * time.Second)
		}
		flapper.NoteStateChange("test")
	}

	if !flapper.IsFlapping("test", true) {
		t.Errorf("Should be flapping when 'max' state changes reported in less than 'duration'")
	}

	time.Sleep(1200 * time.Millisecond)

	if flapper.IsFlapping("test", true) {
		t.Errorf("Should not be flapping when window has moved past counts")
	}
}
