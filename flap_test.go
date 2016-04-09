package main

import (
	"testing"
	"time"
)

func TestAddStateChange(t *testing.T) {
	flapper := newFlapper(5, 30)
	flapper.noteStateChange("test")

	total := flapper.services["test"].Total()
	if total != 1 {
		t.Errorf("Failed to update counter by noting state change. Expected %d, found %d",
			1, total)
	}
}

func TestCheckStateForNonExistantService(t *testing.T) {
	flapper := newFlapper(5, 30)

	if flapper.isFlapping("service-name", true) {
		t.Errorf("Should not report flapping for non-reported service")
	}
}

func TestSimpleFlapDetection(t *testing.T) {
	flapper := newFlapper(5, 30)

	for i := 0; i < flapper.max; i++ {
		flapper.noteStateChange("test")
	}

	if !flapper.isFlapping("test", false) {
		t.Errorf("Sould be flapping when 'max' state changes reported")
	}
}

func TestSpacedFlapDetection(t *testing.T) {
	flapper := newFlapper(10, 10)

	for i := 0; i < flapper.max; i++ {
		flapper.noteStateChange("test")
		time.Sleep(500 * time.Millisecond)
	}

	if !flapper.isFlapping("test", true) {
		t.Errorf("Should be flapping when 'max' state changes reported in less than 'duration'")
	}
}

func TestFlapDetectionResetsOnSlidingOneSecondWindows(t *testing.T) {
	flapper := newFlapper(10, 2)

	for i := 0; i < flapper.max; i++ {
		if i == (flapper.max / 2) {
			time.Sleep(1 * time.Second)
		}
		flapper.noteStateChange("test")
	}

	if !flapper.isFlapping("test", true) {
		t.Errorf("Should be flapping when 'max' state changes reported in less than 'duration'")
	}

	time.Sleep(1200 * time.Millisecond)

	if flapper.isFlapping("test", true) {
		t.Errorf("Should not be flapping when window has moved past counts")
	}
}
