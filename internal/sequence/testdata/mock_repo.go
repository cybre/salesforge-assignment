package testdata

import (
	"context"

	"github.com/cybre/salesforge-assignment/internal/sequence"
)

type MockRepo struct {
	GetSequenceFn    func(ctx context.Context, id int) (sequence.Sequence, bool, error)
	CreateSequenceFn func(ctx context.Context, seq sequence.Sequence) error
	UpdateSequenceFn func(ctx context.Context, seq sequence.Sequence) (bool, error)
	UpdateStepFn     func(ctx context.Context, step sequence.Step) (bool, error)
	DeleteStepFn     func(ctx context.Context, id int) error
}

func (m MockRepo) GetSequence(ctx context.Context, id int) (sequence.Sequence, bool, error) {
	return m.GetSequenceFn(ctx, id)
}

func (m MockRepo) CreateSequence(ctx context.Context, seq sequence.Sequence) error {
	return m.CreateSequenceFn(ctx, seq)
}

func (m MockRepo) UpdateSequence(ctx context.Context, seq sequence.Sequence) (bool, error) {
	return m.UpdateSequenceFn(ctx, seq)
}

func (m MockRepo) UpdateStep(ctx context.Context, step sequence.Step) (bool, error) {
	return m.UpdateStepFn(ctx, step)
}

func (m MockRepo) DeleteStep(ctx context.Context, id int) error {
	return m.DeleteStepFn(ctx, id)
}
