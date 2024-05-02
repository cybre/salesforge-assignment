package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/cybre/salesforge-assignment/internal/sequence"
	"github.com/labstack/echo/v4"
)

// UpdateStep is an echo handler for updating a sequence step.
func (s Server) UpdateStep(e echo.Context) error {
	request := UpdateStepRequest{}
	if err := e.Bind(&request); err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	if request.ID == 0 {
		return e.String(http.StatusBadRequest, "invalid id")
	}

	model := request.BuildStepModel()
	if err := s.sequenceService.UpdateStep(e.Request().Context(), model); err != nil {
		if errors.Is(err, sequence.ErrStepValidation) {
			return e.String(http.StatusBadRequest, err.Error())
		}

		if errors.Is(err, sequence.ErrStepNotFound) {
			return e.String(http.StatusBadRequest, err.Error())
		}

		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.NoContent(http.StatusOK)
}

// DeleteStep is an echo handler for deleting a sequence step.
func (s Server) DeleteStep(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "id must be an integer")
	}

	if err := s.sequenceService.DeleteStep(e.Request().Context(), id); err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.NoContent(http.StatusNoContent)
}

type UpdateStepRequest struct {
	ID      int    `param:"id"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func (r UpdateStepRequest) BuildStepModel() sequence.Step {
	return sequence.Step{
		ID:      r.ID,
		Subject: r.Subject,
		Content: r.Content,
	}
}
