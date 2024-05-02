package sequence_test

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/cybre/salesforge-assignment/internal/sequence"
)

func TestGetSequenceRows_ToSequence(t *testing.T) {
	// Create a sample GetSequenceRows instance
	rows := sequence.GetSequenceRows{
		{
			ID:                   1,
			Name:                 "Sequence 1",
			OpenTrackingEnabled:  true,
			ClickTrackingEnabled: false,
			StepID:               sql.NullInt64{Int64: 1, Valid: true},
			Subject:              sql.NullString{String: "Step 1 Subject", Valid: true},
			Content:              sql.NullString{String: "Step 1 Content", Valid: true},
		},
		{
			ID:                   1,
			Name:                 "Sequence 1",
			OpenTrackingEnabled:  true,
			ClickTrackingEnabled: false,
			StepID:               sql.NullInt64{Int64: 2, Valid: true},
			Subject:              sql.NullString{String: "Step 2 Subject", Valid: true},
			Content:              sql.NullString{String: "Step 2 Content", Valid: true},
		},
	}

	expectedResult := sequence.Sequence{
		ID:            1,
		Name:          "Sequence 1",
		OpenTracking:  true,
		ClickTracking: false,
		Steps: []sequence.Step{
			{
				ID:      1,
				Subject: "Step 1 Subject",
				Content: "Step 1 Content",
			},
			{
				ID:      2,
				Subject: "Step 2 Subject",
				Content: "Step 2 Content",
			},
		},
	}

	// Call the ToSequence method
	seq := rows.ToSequence()

	// Compare the result with the expected value
	if !reflect.DeepEqual(seq, expectedResult) {
		t.Errorf("Expected %v, but got %v", expectedResult, seq)
	}
}

func TestGetSequenceRows_ToSequence_NoSteps(t *testing.T) {
	// Create a sample GetSequenceRows instance
	rows := sequence.GetSequenceRows{
		{
			ID:                   1,
			Name:                 "Sequence 1",
			OpenTrackingEnabled:  true,
			ClickTrackingEnabled: false,
			StepID:               sql.NullInt64{Valid: false},
			Subject:              sql.NullString{Valid: false},
			Content:              sql.NullString{Valid: false},
		},
	}

	// Call the ToSequence method
	seq := rows.ToSequence()

	// Assert the sequence properties
	if seq.ID != 1 {
		t.Errorf("Expected sequence ID to be 1, but got %d", seq.ID)
	}
	if seq.Name != "Sequence 1" {
		t.Errorf("Expected sequence name to be 'Sequence 1', but got '%s'", seq.Name)
	}
	if seq.OpenTracking != true {
		t.Errorf("Expected open tracking to be true, but got false")
	}
	if seq.ClickTracking != false {
		t.Errorf("Expected click tracking to be false, but got true")
	}

	// Assert the steps
	if len(seq.Steps) != 0 {
		t.Errorf("Expected 0 steps, but got %d", len(seq.Steps))
	}
}
