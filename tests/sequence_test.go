package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/cybre/salesforge-assignment/internal/sequence"
	transporthttp "github.com/cybre/salesforge-assignment/internal/transport/http"
)

func TestCreateSequence(t *testing.T) {
	ts := NewTestServer(t)

	// Create a sequence
	request := createSequence(ts, t)

	// Check if the sequence was created in the database
	seq, found, err := ts.Repository.GetSequence(context.Background(), 1)
	if err != nil {
		t.Fatalf("failed to fetch sequence from the database: %v", err)
	}

	if !found {
		t.Fatalf("expected sequence to be found in the database")
	}

	if err = compareSequenceWithRequest(request, seq); err != nil {
		t.Error(err)
	}
}

func TestGetSequence(t *testing.T) {
	ts := NewTestServer(t)

	// Create a sequence
	request := createSequence(ts, t)

	// Get the sequence
	res := ts.GetSequence(t, 1)
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, res.StatusCode)
	}

	var sequence sequence.Sequence
	if err := json.NewDecoder(res.Body).Decode(&sequence); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if err := compareSequenceWithRequest(request, sequence); err != nil {
		t.Error(err)
	}
}

func TestPatchSequence(t *testing.T) {
	ts := NewTestServer(t)

	// Create a sequence
	createSequence(ts, t)

	// Patch the sequence
	patch := transporthttp.PatchSequenceRequest{
		ID:            1,
		OpenTracking:  boolPtr(false),
		ClickTracking: boolPtr(true),
	}
	res := ts.PatchSequence(t, patch)
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, res.StatusCode)
	}

	// Check if the sequence was updated in the database
	seq, found, err := ts.Repository.GetSequence(context.Background(), 1)
	if err != nil {
		t.Fatalf("failed to fetch sequence from the database: %v", err)
	}

	if !found {
		t.Fatalf("expected sequence to be found in the database")
	}

	if seq.OpenTracking != *patch.OpenTracking {
		t.Errorf("expected open tracking to be %t, but got %t", *patch.OpenTracking, seq.OpenTracking)
	}

	if seq.ClickTracking != *patch.ClickTracking {
		t.Errorf("expected click tracking to be %t, but got %t", *patch.ClickTracking, seq.ClickTracking)
	}
}

func TestUpdateStep(t *testing.T) {
	ts := NewTestServer(t)

	// Create a sequence
	createSequence(ts, t)

	// Update a step
	updateStepRequest := transporthttp.UpdateStepRequest{
		ID:      1,
		Subject: "Updated Subject",
		Content: "Updated Content",
	}
	res := ts.PutStep(t, updateStepRequest)
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, res.StatusCode)
	}

	// Check if the step was updated in the database
	seq, found, err := ts.Repository.GetSequence(context.Background(), 1)
	if err != nil {
		t.Fatalf("failed to fetch sequence from the database: %v", err)
	}

	if !found {
		t.Fatalf("expected sequence to be found in the database")
	}

	if seq.Steps[0].Subject != updateStepRequest.Subject {
		t.Errorf("expected step subject to be %s, but got %s", updateStepRequest.Subject, seq.Steps[0].Subject)
	}

	if seq.Steps[0].Content != updateStepRequest.Content {
		t.Errorf("expected step content to be %s, but got %s", updateStepRequest.Content, seq.Steps[0].Content)
	}
}

func TestDeleteStep(t *testing.T) {
	ts := NewTestServer(t)

	// Create a sequence
	createSequence(ts, t)

	// Delete a step
	res := ts.DeleteStep(t, 1)
	if res.StatusCode != http.StatusNoContent {
		t.Errorf("expected status code %d, but got %d", http.StatusNoContent, res.StatusCode)
	}

	// Check if the step was deleted from the database
	seq, found, err := ts.Repository.GetSequence(context.Background(), 1)
	if err != nil {
		t.Fatalf("failed to fetch sequence from the database: %v", err)
	}

	if !found {
		t.Fatalf("expected sequence to be found in the database")
	}

	if len(seq.Steps) != 1 {
		t.Errorf("expected 1 step, but got %d", len(seq.Steps))
	}

	if seq.Steps[0].ID != 2 {
		t.Errorf("expected step ID to be 2, but got %d", seq.Steps[0].ID)
	}
}

func createSequence(ts *TestServer, t *testing.T) transporthttp.CreateSequenceRequest {
	request := transporthttp.CreateSequenceRequest{
		Name:          "Test Sequence",
		OpenTracking:  true,
		ClickTracking: false,
		Steps: []transporthttp.CreateSequenceRequestStep{
			{
				Subject: "Test Subject 1",
				Content: "Test Content 1",
			},
			{
				Subject: "Test Subject 2",
				Content: "Test Content 2",
			},
		},
	}

	err := ts.Repository.CreateSequence(context.Background(), request.BuildSequenceModel())
	if err != nil {
		t.Fatalf("failed to create sequence: %v", err)
	}

	return request
}

func compareSequenceWithRequest(request transporthttp.CreateSequenceRequest, sequence sequence.Sequence) error {
	if sequence.Name != request.Name {
		return fmt.Errorf("expected sequence name to be %s, but got %s", request.Name, sequence.Name)
	}

	if sequence.OpenTracking != request.OpenTracking {
		return fmt.Errorf("expected open tracking to be %t, but got %t", request.OpenTracking, sequence.OpenTracking)
	}

	if sequence.ClickTracking != request.ClickTracking {
		return fmt.Errorf("expected click tracking to be %t, but got %t", request.ClickTracking, sequence.ClickTracking)
	}

	if len(sequence.Steps) != len(request.Steps) {
		return fmt.Errorf("expected %d steps, but got %d", len(request.Steps), len(sequence.Steps))
	}

	for i, step := range sequence.Steps {
		if step.Subject != request.Steps[i].Subject {
			return fmt.Errorf("expected step %d subject to be %s, but got %s", i, request.Steps[i].Subject, step.Subject)
		}

		if step.Content != request.Steps[i].Content {
			return fmt.Errorf("expected step %d content to be %s, but got %s", i, request.Steps[i].Content, step.Content)
		}
	}

	return nil
}

func boolPtr(b bool) *bool {
	return &b
}
