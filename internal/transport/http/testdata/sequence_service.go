package testdata

import (
	"context"

	"github.com/cybre/salesforge-assignment/internal/sequence"
)

type MockSequenceService struct {
	CreateSequenceFn func(ctx context.Context, seq sequence.Sequence) error
	PatchSequenceFn  func(ctx context.Context, patch sequence.SequencePatch) error
	GetSequenceFn    func(ctx context.Context, id int) (sequence.Sequence, error)
	UpdateStepFn     func(ctx context.Context, step sequence.Step) error
	DeleteStepFn     func(ctx context.Context, id int) error
}

func (m MockSequenceService) CreateSequence(ctx context.Context, seq sequence.Sequence) error {
	return m.CreateSequenceFn(ctx, seq)
}

func (m MockSequenceService) PatchSequence(ctx context.Context, patch sequence.SequencePatch) error {
	return m.PatchSequenceFn(ctx, patch)
}

func (m MockSequenceService) GetSequence(ctx context.Context, id int) (sequence.Sequence, error) {
	return m.GetSequenceFn(ctx, id)
}

func (m MockSequenceService) UpdateStep(ctx context.Context, step sequence.Step) error {
	return m.UpdateStepFn(ctx, step)
}

func (m MockSequenceService) DeleteStep(ctx context.Context, id int) error {
	return m.DeleteStepFn(ctx, id)
}
