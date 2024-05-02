package sequence

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrSequenceNotFound is returned when a sequence with the given ID is not found.
	ErrSequenceNotFound = errors.New("sequence with given ID not found")

	// ErrSequenceValidation is returned when a sequence model fails validation.
	ErrSequenceValidation = errors.New("sequence model is invalid")

	// ErrStepNotFound is returned when a step with the given ID is not found.
	ErrStepNotFound = errors.New("step with given ID not found")

	// ErrStepValidation is returned when a step model fails validation.
	ErrStepValidation = errors.New("step model is invalid")
)

// Sequence represents a sequence of emails.
type Sequence struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	OpenTracking  bool   `json:"openTrackingEnabled"`
	ClickTracking bool   `json:"clickTrackingEnabled"`
	Steps         []Step `json:"steps"`
}

// Validate validates the sequence model.
func (s Sequence) Validate() error {
	if s.Name == "" {
		return errors.New("name is required")
	}

	if len(s.Steps) == 0 {
		return errors.New("steps are required")
	}

	for _, step := range s.Steps {
		if err := step.Validate(); err != nil {
			return fmt.Errorf("%w: %s", ErrStepValidation, err)
		}
	}

	return nil
}

// Step represents an email in a sequence.
type Step struct {
	ID      int    `json:"id"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

// Validate validates the step model.
func (s Step) Validate() error {
	if s.Subject == "" {
		return errors.New("subject is required")
	}

	if s.Content == "" {
		return errors.New("content is required")
	}

	return nil
}

// SequencePatch represents a patch for a sequence.
type SequencePatch struct {
	ID            int
	Name          *string
	OpenTracking  *bool
	ClickTracking *bool
}

// Patch applies the patch to the given sequence.
func (s SequencePatch) Patch(seq *Sequence) {
	if s.Name != nil {
		seq.Name = *s.Name
	}

	if s.OpenTracking != nil {
		seq.OpenTracking = *s.OpenTracking
	}

	if s.ClickTracking != nil {
		seq.ClickTracking = *s.ClickTracking
	}
}

// Validate validates the sequence patch.
func (s SequencePatch) Validate() error {
	if s.ID == 0 {
		return errors.New("id is required")
	}

	if s.Name != nil && *s.Name == "" {
		return errors.New("name cannot be empty")
	}

	return nil
}

// Repository represents a sequence repository.
type Repository interface {
	CreateSequence(ctx context.Context, seq Sequence) error
	UpdateSequence(ctx context.Context, seq Sequence) (bool, error)
	GetSequence(ctx context.Context, id int) (Sequence, bool, error)
	UpdateStep(ctx context.Context, step Step) (bool, error)
	DeleteStep(ctx context.Context, id int) error
}

// Service contains the business logic for handling sequences.
type Service struct {
	repo Repository
}

// NewService creates a new sequence service.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateSequence creates a new sequence.
func (s Service) CreateSequence(ctx context.Context, seq Sequence) error {
	if err := seq.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrSequenceValidation, err)
	}

	if err := s.repo.CreateSequence(ctx, seq); err != nil {
		return fmt.Errorf("failed to create sequence: %w", err)
	}

	return nil
}

// PatchSequence patches a sequence using the given patch.
func (s Service) PatchSequence(ctx context.Context, patch SequencePatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrSequenceValidation, err)
	}

	seq, exists, err := s.repo.GetSequence(ctx, patch.ID)
	if err != nil {
		return err
	}

	if !exists {
		return ErrSequenceNotFound
	}

	patch.Patch(&seq)

	updated, err := s.repo.UpdateSequence(ctx, seq)
	if err != nil {
		return fmt.Errorf("failed to patch sequence: %w", err)
	}

	if !updated {
		return ErrSequenceNotFound
	}

	return nil
}

// GetSequence gets a sequence by ID.
func (s Service) GetSequence(ctx context.Context, id int) (Sequence, error) {
	seq, exists, err := s.repo.GetSequence(ctx, id)
	if err != nil {
		return Sequence{}, err
	}

	if !exists {
		return Sequence{}, ErrSequenceNotFound
	}

	return seq, nil
}

// UpdateStep updates a sequence step.
func (s Service) UpdateStep(ctx context.Context, step Step) error {
	if step.ID == 0 {
		return fmt.Errorf("%w: id is required", ErrStepValidation)
	}

	if err := step.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrStepValidation, err)
	}

	updated, err := s.repo.UpdateStep(ctx, step)
	if err != nil {
		return fmt.Errorf("failed to update step: %w", err)
	}

	if !updated {
		return ErrStepNotFound
	}

	return nil
}

// DeleteStep deletes a sequence step.
func (s Service) DeleteStep(ctx context.Context, id int) error {
	if err := s.repo.DeleteStep(ctx, id); err != nil {
		return fmt.Errorf("failed to delete step: %w", err)
	}

	return nil
}
