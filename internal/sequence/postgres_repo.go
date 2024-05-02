package sequence

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// PostgresRepository is a repository containing sequences using Postgres.
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a new Postgres repository.
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

const createSequenceQuery = `
INSERT INTO sequence (name, open_tracking_enabled, click_tracking_enabled) VALUES ($1, $2, $3) RETURNING id;
`
const createStepQuery = `
INSERT INTO step (sequence_id, subject, content) VALUES ($1, $2, $3);
`

// CreateSequence creates a new sequence.
func (r PostgresRepository) CreateSequence(ctx context.Context, seq Sequence) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	if err := func() error {
		var seqID int64
		if err := tx.QueryRowxContext(ctx, createSequenceQuery, seq.Name, seq.OpenTracking, seq.ClickTracking).Scan(&seqID); err != nil {
			return err
		}

		for _, step := range seq.Steps {
			_, err = tx.ExecContext(ctx, createStepQuery, seqID, step.Subject, step.Content)
			if err != nil {
				return err
			}
		}

		return nil
	}(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

const updateSequenceQuery = `
UPDATE sequence SET name = $1, open_tracking_enabled = $2, click_tracking_enabled = $3 WHERE id = $4;
`

// UpdateSequence updates a sequence.
func (r PostgresRepository) UpdateSequence(ctx context.Context, seq Sequence) (bool, error) {
	res, err := r.db.ExecContext(ctx, updateSequenceQuery, seq.Name, seq.OpenTracking, seq.ClickTracking, seq.ID)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

const getSequenceQuery = `
SELECT sequence.id, sequence.name, sequence.open_tracking_enabled, sequence.click_tracking_enabled, step.id as step_id, step.subject, step.content
FROM sequence
LEFT JOIN step ON sequence.id = step.sequence_id 
WHERE sequence.id = $1;
`

// GetSequence gets a sequence by ID.
func (r PostgresRepository) GetSequence(ctx context.Context, id int) (Sequence, bool, error) {
	seq := GetSequenceRows{}
	err := r.db.SelectContext(ctx, &seq, getSequenceQuery, id)
	if err != nil {
		return Sequence{}, false, err
	}

	if len(seq) == 0 {
		return Sequence{}, false, nil
	}

	return seq.ToSequence(), true, nil
}

const updateStepQuery = `
UPDATE step SET subject = $1, content = $2 WHERE id = $3;
`

// UpdateStep updates a sequence step.
func (r PostgresRepository) UpdateStep(ctx context.Context, step Step) (bool, error) {
	res, err := r.db.ExecContext(ctx, updateStepQuery, step.Subject, step.Content, step.ID)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

const deleteStepQuery = `
DELETE FROM step WHERE id = $1;
`

// DeleteStep deletes a sequence step.
func (r PostgresRepository) DeleteStep(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, deleteStepQuery, id)
	return err
}

// GetSequenceRow represents a row returned from the get sequence query.
type GetSequenceRow struct {
	ID                   int            `db:"id"`
	Name                 string         `db:"name"`
	OpenTrackingEnabled  bool           `db:"open_tracking_enabled"`
	ClickTrackingEnabled bool           `db:"click_tracking_enabled"`
	StepID               sql.NullInt64  `db:"step_id"`
	Subject              sql.NullString `db:"subject"`
	Content              sql.NullString `db:"content"`
}

// GetSequenceRows represents multiple rows returned from the get sequence query.
type GetSequenceRows []GetSequenceRow

// ToSequence converts the rows to a sequence domain model.
func (r GetSequenceRows) ToSequence() Sequence {
	seq := Sequence{
		ID:            r[0].ID,
		Name:          r[0].Name,
		OpenTracking:  r[0].OpenTrackingEnabled,
		ClickTracking: r[0].ClickTrackingEnabled,
	}

	// If the first row has a step ID that is null, there are no steps.
	if !r[0].StepID.Valid {
		seq.Steps = []Step{}
		return seq
	}

	seq.Steps = make([]Step, len(r))
	for i, row := range r {
		seq.Steps[i] = Step{
			ID:      int(row.StepID.Int64),
			Subject: row.Subject.String,
			Content: row.Content.String,
		}
	}

	return seq
}
