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

func TestUpdateStep_Success(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request with a JSON payload
	reqBody := `{ "id": 1, "subject": "Test Subject", "content": "Test Content" }`
	req := httptest.NewRequest(http.MethodPut, "/sequences", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		UpdateStepFn: func(ctx context.Context, seq sequence.Step) error {
			return nil
		},
	}

	// Create a new server instance with the mock sequence service
	server := transporthttp.NewServer(mockSequenceService)

	// Call the UpdateStep method
	err := server.UpdateStep(c)

	// Check if there was no error
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if the response status code is 200 OK
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestUpdateStep_InvalidBody_BadRequest(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request with an invalid JSON payload
	reqBody := `{"}`
	req := httptest.NewRequest(http.MethodPut, "/step/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		UpdateStepFn: func(ctx context.Context, seq sequence.Step) error {
			return nil
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

	// Check if the response status code is 400 Bad Request
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateStep_ValidationErr_BadRequest(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request with a JSON payload
	reqBody := `{"subject": "Test Subject", "content": "Test Content"}`
	req := httptest.NewRequest(http.MethodPut, "/sequences", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testErr := sequence.ErrStepValidation

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		UpdateStepFn: func(ctx context.Context, seq sequence.Step) error {
			return testErr
		},
	}

	// Create a new server instance with the mock sequence service
	server := transporthttp.NewServer(mockSequenceService)

	// Call the UpdateStep method
	server.UpdateStep(c)

	// Check if the response status code is 400 Bad Request
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateStep_NotFoundErr_BadRequest(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request with a JSON payload
	reqBody := `{"subject": "Test Subject", "content": "Test Content"}`
	req := httptest.NewRequest(http.MethodPut, "/sequences", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	testErr := sequence.ErrStepNotFound

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		UpdateStepFn: func(ctx context.Context, seq sequence.Step) error {
			return testErr
		},
	}

	// Create a new server instance with the mock sequence service
	server := transporthttp.NewServer(mockSequenceService)

	// Call the UpdateStep method
	server.UpdateStep(c)

	// Check if the response status code is 400 Bad Request
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateStep_UnknownError_InternalServerError(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request with a JSON payload
	reqBody := `{"subject": "Test Subject", "content": "Test Content"}`
	req := httptest.NewRequest(http.MethodPut, "/step/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	testErr := errors.New("test error")

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		UpdateStepFn: func(ctx context.Context, seq sequence.Step) error {
			return testErr
		},
	}

	// Create a new server instance with the mock sequence service
	server := transporthttp.NewServer(mockSequenceService)

	// Call the UpdateStep method
	server.UpdateStep(c)

	// Check if the response status code is 500 Internal Server Error
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rec.Code)
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

func TestDeleteStep_Success(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodDelete, "/step/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		DeleteStepFn: func(ctx context.Context, id int) error {
			return nil
		},
	}

	// Create a new server instance with the mock sequence service
	server := transporthttp.NewServer(mockSequenceService)

	// Call the DeleteStep method
	err := server.DeleteStep(c)

	// Check if there was no error
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if the response status code is 204 No Content
	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status code %d, got %d", http.StatusNoContent, rec.Code)
	}
}

func TestDeleteStep_InvalidID_BadRequest(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request with an invalid ID
	req := httptest.NewRequest(http.MethodDelete, "/step/abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		DeleteStepFn: func(ctx context.Context, id int) error {
			return nil
		},
	}

	// Create a new server instance with the mock sequence service
	server := transporthttp.NewServer(mockSequenceService)

	// Call the DeleteStep method
	server.DeleteStep(c)

	// Check if the response status code is 400 Bad Request
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestDeleteStep_UnknownError_InternalServerError(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodDelete, "/step/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	testErr := errors.New("test error")

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		DeleteStepFn: func(ctx context.Context, id int) error {
			return testErr
		},
	}

	// Create a new server instance with the mock sequence service
	server := transporthttp.NewServer(mockSequenceService)

	// Call the DeleteStep method
	server.DeleteStep(c)

	// Check if the response status code is 500 Internal Server Error
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}
