package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/cybre/salesforge-assignment/internal/sequence"
	"github.com/labstack/echo/v4"
)

// CreateSequence is an echo handler for creating a sequence.
func (s Server) CreateSequence(e echo.Context) error {
	request := CreateSequenceRequest{}
	if err := e.Bind(&request); err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	model := request.BuildSequenceModel()
	if err := s.sequenceService.CreateSequence(e.Request().Context(), model); err != nil {
		if errors.Is(err, sequence.ErrSequenceValidation) {
			return e.String(http.StatusBadRequest, err.Error())
		}

		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.NoContent(http.StatusCreated)
}

// GetSequence is an echo handler for getting a sequence.
func (s Server) GetSequence(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "id must be an integer")
	}

	seq, err := s.sequenceService.GetSequence(e.Request().Context(), id)
	if err != nil {
		if errors.Is(err, sequence.ErrSequenceNotFound) {
			return e.String(http.StatusNotFound, err.Error())
		}

		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, seq)
}

// PatchSequence is an echo handler for patching a sequence.
func (s Server) PatchSequence(e echo.Context) error {
	request := PatchSequenceRequest{}
	if err := e.Bind(&request); err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	if err := s.sequenceService.PatchSequence(e.Request().Context(), request.BuildSequencePatch()); err != nil {
		if errors.Is(err, sequence.ErrSequenceValidation) {
			return e.String(http.StatusBadRequest, err.Error())
		}

		if errors.Is(err, sequence.ErrSequenceNotFound) {
			return e.String(http.StatusBadRequest, err.Error())
		}

		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.NoContent(http.StatusOK)
}

// PatchSequenceRequest represents the request body for patching a sequence.
type PatchSequenceRequest struct {
	ID            int     `param:"id"`
	Name          *string `json:"name"`
	OpenTracking  *bool   `json:"openTrackingEnabled"`
	ClickTracking *bool   `json:"clickTrackingEnabled"`
}

// BuildSequencePatch builds a sequence patch domain object from the request.
func (r PatchSequenceRequest) BuildSequencePatch() sequence.SequencePatch {
	if r.Name != nil {
		trimmedName := strings.TrimSpace(*r.Name)
		r.Name = &trimmedName
	}

	return sequence.SequencePatch{
		ID:            r.ID,
		Name:          r.Name,
		OpenTracking:  r.OpenTracking,
		ClickTracking: r.ClickTracking,
	}
}

// CreateSequenceRequest represents the request body for creating a sequence.
type CreateSequenceRequest struct {
	Name          string                      `json:"name"`
	OpenTracking  bool                        `json:"openTrackingEnabled"`
	ClickTracking bool                        `json:"clickTrackingEnabled"`
	Steps         []CreateSequenceRequestStep `json:"steps"`
}

type CreateSequenceRequestStep struct {
	Subject string `json:"subject"`
	Content string `json:"content"`
}

// BuildSequenceModel builds a sequence domain model from the request.
func (r CreateSequenceRequest) BuildSequenceModel() sequence.Sequence {
	steps := make([]sequence.Step, len(r.Steps))
	for i, step := range r.Steps {
		steps[i] = sequence.Step{
			Subject: strings.TrimSpace(step.Subject),
			Content: strings.TrimSpace(step.Content),
		}
	}

	return sequence.Sequence{
		Name:          strings.TrimSpace(r.Name),
		OpenTracking:  r.OpenTracking,
		ClickTracking: r.ClickTracking,
		Steps:         steps,
	}
}
