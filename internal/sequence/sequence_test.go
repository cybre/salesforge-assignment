package sequence_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/cybre/salesforge-assignment/internal/sequence"
	"github.com/cybre/salesforge-assignment/internal/sequence/testdata"
)

func TestSequence_Validate(t *testing.T) {
	testCases := []struct {
		name     string
		sequence sequence.Sequence
		expected error
	}{
		{
			name: "Empty name",
			sequence: sequence.Sequence{
				Name:  "",
				Steps: []sequence.Step{{}},
			},
			expected: errors.New("name is required"),
		},
		{
			name: "Missing steps",
			sequence: sequence.Sequence{
				Name:  "Sequence 2",
				Steps: []sequence.Step{},
			},
			expected: errors.New("steps are required"),
		},
		{
			name: "Invalid step 1",
			sequence: sequence.Sequence{
				Name: "Sequence 3",
				Steps: []sequence.Step{
					{Subject: "", Content: "Content 2"},
				},
			},
			expected: fmt.Errorf("%w: %s", sequence.ErrStepValidation, "subject is required"),
		},
		{
			name: "Valid sequence",
			sequence: sequence.Sequence{
				Name: "Sequence 4",
				Steps: []sequence.Step{
					{Subject: "Subject 1", Content: "Content 1"},
					{Subject: "Subject 2", Content: "Content 2"},
				},
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.sequence.Validate()
			if err == nil && tc.expected != nil {
				t.Errorf("Expected error: %v, got: nil", tc.expected)
			} else if err != nil && tc.expected == nil {
				t.Errorf("Expected no error, got: %v", err)
			} else if err != nil && tc.expected != nil && err.Error() != tc.expected.Error() {
				t.Errorf("Expected error: %v, got: %v", tc.expected, err)
			}
		})
	}
}

func TestStep_Validate(t *testing.T) {
	testCases := []struct {
		name     string
		step     sequence.Step
		expected error
	}{
		{
			name:     "Valid step",
			step:     sequence.Step{Subject: "Subject 1", Content: "Content 1"},
			expected: nil,
		},
		{
			name:     "Invalid step with empty subject",
			step:     sequence.Step{Subject: "", Content: "Content 2"},
			expected: errors.New("subject is required"),
		},
		{
			name:     "Invalid step with empty content",
			step:     sequence.Step{Subject: "Subject 2", Content: ""},
			expected: errors.New("content is required"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.step.Validate()
			if err == nil && tc.expected != nil {
				t.Errorf("Expected error: %v, got: nil", tc.expected)
			} else if err != nil && tc.expected == nil {
				t.Errorf("Expected no error, got: %v", err)
			} else if err != nil && tc.expected != nil && err.Error() != tc.expected.Error() {
				t.Errorf("Expected error: %v, got: %v", tc.expected, err)
			}
		})
	}
}
func TestSequencePatch_Patch(t *testing.T) {
	testCases := []struct {
		name          string
		sequencePatch sequence.SequencePatch
		sequence      sequence.Sequence
		expected      sequence.Sequence
	}{
		{
			name: "Patch name",
			sequencePatch: sequence.SequencePatch{
				Name: stringPtr("New Name"),
			},
			sequence: sequence.Sequence{
				Name: "Old Name",
			},
			expected: sequence.Sequence{
				Name: "New Name",
			},
		},
		{
			name: "Patch open tracking",
			sequencePatch: sequence.SequencePatch{
				OpenTracking: boolPtr(false),
			},
			sequence: sequence.Sequence{
				Name:          "Sequence",
				OpenTracking:  true,
				ClickTracking: false,
			},
			expected: sequence.Sequence{
				Name:          "Sequence",
				OpenTracking:  false,
				ClickTracking: false,
			},
		},
		{
			name: "Patch click tracking",
			sequencePatch: sequence.SequencePatch{
				ClickTracking: boolPtr(true),
			},
			sequence: sequence.Sequence{
				Name:          "Sequence",
				OpenTracking:  true,
				ClickTracking: false,
			},
			expected: sequence.Sequence{
				Name:          "Sequence",
				OpenTracking:  true,
				ClickTracking: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.sequencePatch.Patch(&tc.sequence)
			if !reflect.DeepEqual(tc.sequence, tc.expected) {
				t.Errorf("Expected: %v, got: %v", tc.expected, tc.sequence)
			}
		})
	}
}

func TestSequencePatch_Validate(t *testing.T) {
	testCases := []struct {
		name     string
		sequence sequence.SequencePatch
		expected error
	}{
		{
			name: "Valid sequence patch",
			sequence: sequence.SequencePatch{
				ID:   1,
				Name: stringPtr("Sequence Patch 1"),
			},
			expected: nil,
		},
		{
			name: "Invalid sequence patch with empty ID",
			sequence: sequence.SequencePatch{
				ID:   0,
				Name: stringPtr("Sequence Patch 2"),
			},
			expected: errors.New("id is required"),
		},
		{
			name: "Invalid sequence patch with empty name",
			sequence: sequence.SequencePatch{
				ID:   2,
				Name: stringPtr(""),
			},
			expected: errors.New("name cannot be empty"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.sequence.Validate()
			if err == nil && tc.expected != nil {
				t.Errorf("Expected error: %v, got: nil", tc.expected)
			} else if err != nil && tc.expected == nil {
				t.Errorf("Expected no error, got: %v", err)
			} else if err != nil && tc.expected != nil && err.Error() != tc.expected.Error() {
				t.Errorf("Expected error: %v, got: %v", tc.expected, err)
			}
		})
	}
}

func TestService_CreateSequence(t *testing.T) {
	ctx := context.Background()

	repo := testdata.MockRepo{
		CreateSequenceFn: func(ctx context.Context, seq sequence.Sequence) error {
			return nil
		},
	}

	repoErr := errors.New("repository error")

	testCases := []struct {
		name        string
		sequence    sequence.Sequence
		expectedErr error
		repository  sequence.Repository
	}{
		{
			name: "Valid sequence",
			sequence: sequence.Sequence{
				Name: "Sequence 1",
				Steps: []sequence.Step{
					{Subject: "Subject 1", Content: "Content 1"},
					{Subject: "Subject 2", Content: "Content 2"},
				},
			},
			expectedErr: nil,
			repository:  repo,
		},
		{
			name: "Invalid sequence",
			sequence: sequence.Sequence{
				Name:  "",
				Steps: []sequence.Step{{}},
			},
			expectedErr: sequence.ErrSequenceValidation,
			repository:  repo,
		},
		{
			name: "Repository error",
			sequence: sequence.Sequence{
				Name: "Sequence 2",
				Steps: []sequence.Step{
					{Subject: "Subject 1", Content: "Content 1"},
					{Subject: "Subject 2", Content: "Content 2"},
				},
			},
			expectedErr: repoErr,
			repository: testdata.MockRepo{
				CreateSequenceFn: func(ctx context.Context, seq sequence.Sequence) error {
					return repoErr
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := sequence.NewService(tc.repository)
			err := svc.CreateSequence(ctx, tc.sequence)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestService_PatchSequence(t *testing.T) {
	ctx := context.Background()

	repo := testdata.MockRepo{
		GetSequenceFn: func(ctx context.Context, id int) (sequence.Sequence, bool, error) {
			return sequence.Sequence{}, false, nil
		},
		UpdateSequenceFn: func(ctx context.Context, seq sequence.Sequence) (bool, error) {
			return false, nil
		},
	}

	repoErr := errors.New("repository error")

	testCases := []struct {
		name        string
		patch       sequence.SequencePatch
		expectedErr error
		repository  sequence.Repository
	}{
		{
			name: "Valid patch",
			patch: sequence.SequencePatch{
				ID:   1,
				Name: stringPtr("New Name"),
			},
			expectedErr: nil,
			repository: testdata.MockRepo{
				GetSequenceFn: func(ctx context.Context, id int) (sequence.Sequence, bool, error) {
					return sequence.Sequence{
						ID:   1,
						Name: "Old Name",
					}, true, nil
				},
				UpdateSequenceFn: func(ctx context.Context, seq sequence.Sequence) (bool, error) {
					return true, nil
				},
			},
		},
		{
			name: "Invalid patch",
			patch: sequence.SequencePatch{
				ID:   0,
				Name: stringPtr(""),
			},
			expectedErr: sequence.ErrSequenceValidation,
			repository:  repo,
		},
		{
			name: "Sequence not found",
			patch: sequence.SequencePatch{
				ID:   2,
				Name: stringPtr("New Name"),
			},
			expectedErr: sequence.ErrSequenceNotFound,
			repository: testdata.MockRepo{
				GetSequenceFn: func(ctx context.Context, id int) (sequence.Sequence, bool, error) {
					return sequence.Sequence{}, false, nil
				},
			},
		},
		{
			name: "Failed to update sequence",
			patch: sequence.SequencePatch{
				ID:   3,
				Name: stringPtr("New Name"),
			},
			expectedErr: repoErr,
			repository: testdata.MockRepo{
				GetSequenceFn: func(ctx context.Context, id int) (sequence.Sequence, bool, error) {
					return sequence.Sequence{
						ID:   3,
						Name: "Old Name",
					}, true, nil
				},
				UpdateSequenceFn: func(ctx context.Context, seq sequence.Sequence) (bool, error) {
					return false, repoErr
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := sequence.NewService(tc.repository)
			err := svc.PatchSequence(ctx, tc.patch)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestService_GetSequence(t *testing.T) {
	ctx := context.Background()

	repo := testdata.MockRepo{
		GetSequenceFn: func(ctx context.Context, id int) (sequence.Sequence, bool, error) {
			return sequence.Sequence{
				ID:   id,
				Name: "Test Sequence",
			}, true, nil
		},
	}

	repoErr := errors.New("repository error")

	testCases := []struct {
		name        string
		id          int
		expected    sequence.Sequence
		expectedErr error
		repository  sequence.Repository
	}{
		{
			name: "Valid sequence",
			id:   1,
			expected: sequence.Sequence{
				ID:   1,
				Name: "Test Sequence",
			},
			expectedErr: nil,
			repository:  repo,
		},
		{
			name:        "Sequence not found",
			id:          2,
			expected:    sequence.Sequence{},
			expectedErr: sequence.ErrSequenceNotFound,
			repository: testdata.MockRepo{
				GetSequenceFn: func(ctx context.Context, id int) (sequence.Sequence, bool, error) {
					return sequence.Sequence{}, false, nil
				},
			},
		},
		{
			name:        "Repository error",
			id:          3,
			expected:    sequence.Sequence{},
			expectedErr: repoErr,
			repository: testdata.MockRepo{
				GetSequenceFn: func(ctx context.Context, id int) (sequence.Sequence, bool, error) {
					return sequence.Sequence{}, false, repoErr
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := sequence.NewService(tc.repository)
			seq, err := svc.GetSequence(ctx, tc.id)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}

			if !reflect.DeepEqual(seq, tc.expected) {
				t.Errorf("Expected: %v, got: %v", tc.expected, seq)
			}
		})
	}
}

func TestService_UpdateStep(t *testing.T) {
	ctx := context.Background()

	repo := testdata.MockRepo{
		UpdateStepFn: func(ctx context.Context, step sequence.Step) (bool, error) {
			return true, nil
		},
	}

	repoErr := errors.New("repository error")

	testCases := []struct {
		name        string
		step        sequence.Step
		expectedErr error
		repository  sequence.Repository
	}{
		{
			name: "Valid step",
			step: sequence.Step{
				Subject: "Subject 1",
				Content: "Content 1",
			},
			expectedErr: nil,
			repository:  repo,
		},
		{
			name: "Invalid step",
			step: sequence.Step{
				Subject: "",
				Content: "Content 2",
			},
			expectedErr: sequence.ErrStepValidation,
			repository:  repo,
		},
		{
			name: "Step not found",
			step: sequence.Step{
				Subject: "Subject 3",
				Content: "Content 3",
			},
			expectedErr: sequence.ErrStepNotFound,
			repository: testdata.MockRepo{
				UpdateStepFn: func(ctx context.Context, step sequence.Step) (bool, error) {
					return false, nil
				},
			},
		},
		{
			name: "Failed to update step",
			step: sequence.Step{
				Subject: "Subject 4",
				Content: "Content 4",
			},
			expectedErr: repoErr,
			repository: testdata.MockRepo{
				UpdateStepFn: func(ctx context.Context, step sequence.Step) (bool, error) {
					return false, repoErr
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := sequence.NewService(tc.repository)
			err := svc.UpdateStep(ctx, tc.step)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestService_DeleteStep(t *testing.T) {
	ctx := context.Background()

	repo := testdata.MockRepo{
		DeleteStepFn: func(ctx context.Context, id int) error {
			return nil
		},
	}

	repoErr := errors.New("repository error")

	testCases := []struct {
		name        string
		id          int
		expectedErr error
		repository  sequence.Repository
	}{
		{
			name:        "Valid step deletion",
			id:          1,
			expectedErr: nil,
			repository:  repo,
		},
		{
			name:        "Failed to delete step",
			id:          2,
			expectedErr: repoErr,
			repository: testdata.MockRepo{
				DeleteStepFn: func(ctx context.Context, id int) error {
					return repoErr
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := sequence.NewService(tc.repository)
			err := svc.DeleteStep(ctx, tc.id)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
