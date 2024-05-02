package http_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/cybre/salesforge-assignment/internal/sequence"
	transporthttp "github.com/cybre/salesforge-assignment/internal/transport/http"
	"github.com/cybre/salesforge-assignment/internal/transport/http/testdata"
	"github.com/labstack/echo/v4"
)

func TestUpdateStep(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		serviceError   error
		idParamValue   string
	}{
		{
			name:           "Success",
			requestBody:    `{ "subject": "Test Subject", "content": "Test Content" }`,
			expectedStatus: http.StatusOK,
			idParamValue:   "1",
		},
		{
			name:           "Invalid Request Body",
			requestBody:    `{"}`,
			expectedStatus: http.StatusBadRequest,
			idParamValue:   "1",
		},
		{
			name:           "Validation Error",
			requestBody:    `{"subject": "Test Subject"}`,
			expectedStatus: http.StatusBadRequest,
			serviceError:   sequence.ErrStepValidation,
			idParamValue:   "1",
		},
		{
			name:           "Invalid ID param",
			requestBody:    `{"subject": "Test Subject", "content": "Test Content"}`,
			expectedStatus: http.StatusBadRequest,
			idParamValue:   "abc",
		},
		{
			name:           "Not Found Error",
			requestBody:    `{"subject": "Test Subject", "content": "Test Content"}`,
			expectedStatus: http.StatusBadRequest,
			idParamValue:   "1",
			serviceError:   sequence.ErrStepNotFound,
		},
		{
			name:           "Unknown Error",
			requestBody:    `{"subject": "Test Subject", "content": "Test Content"}`,
			expectedStatus: http.StatusInternalServerError,
			serviceError:   errors.New("test error"),
			idParamValue:   "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Echo instance
			e := echo.New()

			// Create a new HTTP request with a JSON payload
			req := httptest.NewRequest(http.MethodPut, "/sequence/1", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.idParamValue)

			// Create a mock sequence service
			mockSequenceService := &testdata.MockSequenceService{
				UpdateStepFn: func(ctx context.Context, seq sequence.Step) error {
					return tt.serviceError
				},
			}

			// Create a new server instance with the mock sequence service
			server := transporthttp.NewServer(mockSequenceService)

			// Call the UpdateStep method
			err := server.UpdateStep(c)

			// Check if there was an error
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			// Check if the response status code matches the expected status code
			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestUpdateStepRequest_BuildStepModel(t *testing.T) {
	r := transporthttp.UpdateStepRequest{
		ID:      1,
		Subject: "Test Subject",
		Content: "Test Content",
	}

	expected := sequence.Step{
		ID:      r.ID,
		Subject: r.Subject,
		Content: r.Content,
	}

	result := r.BuildStepModel()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestDeleteStep(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		expectedStatus int
		serviceError   error
	}{
		{
			name:           "Success",
			id:             "1",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Invalid ID",
			id:             "abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Unknown Error",
			id:             "1",
			expectedStatus: http.StatusInternalServerError,
			serviceError:   errors.New("test error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Echo instance
			e := echo.New()

			// Create a new HTTP request
			req := httptest.NewRequest(http.MethodDelete, "/step/"+tt.id, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)

			// Create a mock sequence service
			mockSequenceService := &testdata.MockSequenceService{
				DeleteStepFn: func(ctx context.Context, id int) error {
					return tt.serviceError
				},
			}

			// Create a new server instance with the mock sequence service
			server := transporthttp.NewServer(mockSequenceService)

			// Call the DeleteStep method
			err := server.DeleteStep(c)

			// Check if there was an error
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			// Check if the response status code matches the expected status code
			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}
