package http_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/cybre/salesforge-assignment/internal/sequence"
	transporthttp "github.com/cybre/salesforge-assignment/internal/transport/http"
	"github.com/cybre/salesforge-assignment/internal/transport/http/testdata"
)

func TestCreateSequence_Success(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request with a JSON payload
	reqBody := `{"name": "Test Sequence", "openTrackingEnabled": true, "clickTrackingEnabled": false, "steps": [{"subject": "Step 1", "content": "Content 1"}]}`
	req := httptest.NewRequest(http.MethodPost, "/sequences", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		CreateSequenceFn: func(ctx context.Context, seq sequence.Sequence) error {
			return nil
		},
	}

	// Create a new server instance with the mock sequence service
	server := transporthttp.NewServer(mockSequenceService)

	// Call the CreateSequence method
	err := server.CreateSequence(c)

	// Check if there was an error
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if the response status code is 201 Created
	if rec.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestCreateSequence_InvalidBody_BadRequest(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request with an invalid JSON payload
	reqBody := `{"}`
	req := httptest.NewRequest(http.MethodPost, "/sequences", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		CreateSequenceFn: func(ctx context.Context, seq sequence.Sequence) error {
			return nil
		},
	}

	// Create a new server instance with the mock sequence
	server := transporthttp.NewServer(mockSequenceService)

	// Call the CreateSequence method
	err := server.CreateSequence(c)

	// Check if there was an error
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if the response status code is 400 Bad Request
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateSequence_ValidationErr_BadRequest(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request with a JSON payload
	reqBody := `{"name": "", "openTrackingEnabled": true, "clickTrackingEnabled": false, "steps": [{"subject": "Step 1", "content": "Content 1"}]}`
	req := httptest.NewRequest(http.MethodPost, "/sequences", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testErr := sequence.ErrSequenceValidation

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		CreateSequenceFn: func(ctx context.Context, seq sequence.Sequence) error {
			return testErr
		},
	}

	// Create a new server instance with the mock sequence
	server := transporthttp.NewServer(mockSequenceService)

	// Call the CreateSequence method
	server.CreateSequence(c)

	// Check if the response status code is 400 Bad Request
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateSequence_UnknownError_InternalServerError(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request with a JSON payload
	reqBody := `{"name": "Test Sequence", "openTrackingEnabled": true, "clickTrackingEnabled": false, "steps": [{"subject": "Step 1", "content": "Content 1"}]}`
	req := httptest.NewRequest(http.MethodPost, "/sequences", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testErr := errors.New("test error")

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		CreateSequenceFn: func(ctx context.Context, seq sequence.Sequence) error {
			return testErr
		},
	}

	// Create a new server instance with the mock sequence
	server := transporthttp.NewServer(mockSequenceService)

	// Call the CreateSequence method
	server.CreateSequence(c)

	// Check if the response status code is 500 Internal Server Error
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestGetSequence_Success(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodGet, "/sequences/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		GetSequenceFn: func(ctx context.Context, id int) (sequence.Sequence, error) {
			return sequence.Sequence{
				ID:            1,
				Name:          "Test Sequence",
				OpenTracking:  false,
				ClickTracking: false,
				Steps: []sequence.Step{
					{ID: 1, Subject: "Step 1", Content: "Content 1"},
				},
			}, nil
		},
	}

	// Create a new server instance with the mock sequence service
	server := transporthttp.NewServer(mockSequenceService)

	// Call the GetSequence method
	err := server.GetSequence(c)

	// Check if there was no error
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if the response status code is 200 OK
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// Check if the response body contains the sequence
	expectedBody := "{\"id\":1,\"name\":\"Test Sequence\",\"openTrackingEnabled\":false,\"clickTrackingEnabled\":false,\"steps\":[{\"id\":1,\"subject\":\"Step 1\",\"content\":\"Content 1\"}]}\n"
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}
}

func TestGetSequence_InvalidParam_BadRequest(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request with an invalid ID
	req := httptest.NewRequest(http.MethodGet, "/sequences/abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{}

	// Create a new server instance with the mock sequence service
	server := transporthttp.NewServer(mockSequenceService)

	// Call the GetSequence method
	err := server.GetSequence(c)

	// Check if there was an error
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if the response status code is 400 Bad Request
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}

	// Check if the response body contains the error message
	expectedBody := "id must be an integer"
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}
}

func TestGetSequence_NotFound(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodGet, "/sequences/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		GetSequenceFn: func(ctx context.Context, id int) (sequence.Sequence, error) {
			return sequence.Sequence{}, sequence.ErrSequenceNotFound
		},
	}

	// Create a new server instance with the mock sequence service
	server := transporthttp.NewServer(mockSequenceService)

	// Call the GetSequence method
	err := server.GetSequence(c)

	// Check if there was no error
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if the response status code is 404 Not Found
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestGetSequence_UnknownError_InternalServerError(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodGet, "/sequences/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	testErr := errors.New("test error")

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		GetSequenceFn: func(ctx context.Context, id int) (sequence.Sequence, error) {
			return sequence.Sequence{}, testErr
		},
	}

	// Create a new server instance with the mock sequence service
	server := transporthttp.NewServer(mockSequenceService)

	// Call the GetSequence method
	err := server.GetSequence(c)

	// Check if there was no error
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if the response status code is 500 Internal Server Error
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}
func TestPatchSequence_Success(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request with a JSON payload
	reqBody := `{"name": "Test Name", "openTracking": true, "clickTracking": false}`
	req := httptest.NewRequest(http.MethodPatch, "/sequence/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		PatchSequenceFn: func(ctx context.Context, seq sequence.SequencePatch) error {
			return nil
		},
	}

	// Create a new server instance with the mock sequence service
	server := transporthttp.NewServer(mockSequenceService)

	// Call the PatchSequence method
	err := server.PatchSequence(c)

	// Check if there was no error
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if the response status code is 200 OK
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestPatchSequence_InvalidBody_BadRequest(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request with an invalid JSON payload
	reqBody := `{"}`
	req := httptest.NewRequest(http.MethodPatch, "/sequence/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		PatchSequenceFn: func(ctx context.Context, seq sequence.SequencePatch) error {
			return nil
		},
	}

	// Create a new server instance with the mock sequence service
	server := transporthttp.NewServer(mockSequenceService)

	// Call the PatchSequence method
	err := server.PatchSequence(c)

	// Check if there was an error
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if the response status code is 400 Bad Request
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestPatchSequence_ValidationErr_BadRequest(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request with a JSON payload
	reqBody := `{"name": "", "openTracking": true, "clickTracking": false}`
	req := httptest.NewRequest(http.MethodPatch, "/sequence/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	testErr := sequence.ErrSequenceValidation

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		PatchSequenceFn: func(ctx context.Context, seq sequence.SequencePatch) error {
			return testErr
		},
	}

	// Create a new server instance with the mock sequence service
	server := transporthttp.NewServer(mockSequenceService)

	// Call the PatchSequence method
	server.PatchSequence(c)

	// Check if the response status code is 400 Bad Request
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestPatchSequence_NotFoundErr_BadRequest(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request with a JSON payload
	reqBody := `{"id": 1, "name": "Test Name", "openTracking": true, "clickTracking": false}`
	req := httptest.NewRequest(http.MethodPatch, "/sequence/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	testErr := sequence.ErrSequenceNotFound

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		PatchSequenceFn: func(ctx context.Context, seq sequence.SequencePatch) error {
			return testErr
		},
	}

	// Create a new server instance with the mock sequence service
	server := transporthttp.NewServer(mockSequenceService)

	// Call the PatchSequence method
	server.PatchSequence(c)

	// Check if the response status code is 400 Bad Request
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestPatchSequence_UnknownError_InternalServerError(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request with a JSON payload
	reqBody := `{"id": 1, "name": "Test Name", "openTracking": true, "clickTracking": false}`
	req := httptest.NewRequest(http.MethodPatch, "/sequence/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	testErr := errors.New("test error")

	// Create a mock sequence service
	mockSequenceService := &testdata.MockSequenceService{
		PatchSequenceFn: func(ctx context.Context, seq sequence.SequencePatch) error {
			return testErr
		},
	}

	// Create a new server instance with the mock sequence service
	server := transporthttp.NewServer(mockSequenceService)

	// Call the PatchSequence method
	server.PatchSequence(c)

	// Check if the response status code is 500 Internal Server Error
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestBuildSequenceModel(t *testing.T) {
	// Create a new CreateSequenceRequest instance
	req := transporthttp.CreateSequenceRequest{
		Name:          "Test Sequence ",
		OpenTracking:  true,
		ClickTracking: false,
		Steps: []transporthttp.CreateSequenceRequestStep{
			{Subject: "Step 1", Content: "Content 1"},
			{Subject: "Step 2", Content: "Content 2"},
			{Subject: "Step 3", Content: "Content 3"},
		},
	}

	// Call the BuildSequenceModel method
	sequenceModel := req.BuildSequenceModel()

	// Check if the sequence model is built correctly
	expectedName := "Test Sequence"
	if sequenceModel.Name != expectedName {
		t.Errorf("expected name %q, got %q", expectedName, sequenceModel.Name)
	}

	expectedOpenTracking := true
	if sequenceModel.OpenTracking != expectedOpenTracking {
		t.Errorf("expected open tracking %v, got %v", expectedOpenTracking, sequenceModel.OpenTracking)
	}

	expectedClickTracking := false
	if sequenceModel.ClickTracking != expectedClickTracking {
		t.Errorf("expected click tracking %v, got %v", expectedClickTracking, sequenceModel.ClickTracking)
	}

	expectedSteps := []sequence.Step{
		{Subject: "Step 1", Content: "Content 1"},
		{Subject: "Step 2", Content: "Content 2"},
		{Subject: "Step 3", Content: "Content 3"},
	}
	if !reflect.DeepEqual(sequenceModel.Steps, expectedSteps) {
		t.Errorf("expected steps %v, got %v", expectedSteps, sequenceModel.Steps)
	}
}

func TestPatchSequenceRequest_BuildSequencePatch(t *testing.T) {
	// Create a new PatchSequenceRequest instance
	req := transporthttp.PatchSequenceRequest{
		ID:            1,
		Name:          stringPtr("Test Name "),
		OpenTracking:  boolPtr(true),
		ClickTracking: boolPtr(false),
	}

	// Call the BuildSequencePatch method
	sequencePatch := req.BuildSequencePatch()

	// Check if the sequence patch is built correctly
	expectedID := 1
	if sequencePatch.ID != expectedID {
		t.Errorf("expected ID %q, got %q", expectedID, sequencePatch.ID)
	}

	expectedName := "Test Name"
	if *sequencePatch.Name != expectedName {
		t.Errorf("expected name %q, got %q", expectedName, *sequencePatch.Name)
	}

	expectedOpenTracking := true
	if *sequencePatch.OpenTracking != expectedOpenTracking {
		t.Errorf("expected open tracking %v, got %v", expectedOpenTracking, *sequencePatch.OpenTracking)
	}

	expectedClickTracking := false
	if *sequencePatch.ClickTracking != expectedClickTracking {
		t.Errorf("expected click tracking %v, got %v", expectedClickTracking, *sequencePatch.ClickTracking)
	}
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
