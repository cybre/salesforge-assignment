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

func TestCreateSequence(t *testing.T) {
	tests := []struct {
		name         string
		reqBody      string
		expected     int
		serviceError error
	}{
		{
			name:     "Success",
			reqBody:  `{"name": "Test Sequence", "openTrackingEnabled": true, "clickTrackingEnabled": false, "steps": [{"subject": "Step 1", "content": "Content 1"}]}`,
			expected: http.StatusCreated,
		},
		{
			name:     "Invalid Body",
			reqBody:  `{"}`,
			expected: http.StatusBadRequest,
		},
		{
			name:         "Validation Error",
			reqBody:      `{"name": "", "openTrackingEnabled": true, "clickTrackingEnabled": false, "steps": [{"subject": "Step 1", "content": "Content 1"}]}`,
			expected:     http.StatusBadRequest,
			serviceError: sequence.ErrSequenceValidation,
		},
		{
			name:         "Unknown Error",
			reqBody:      `{"name": "Test Sequence", "openTrackingEnabled": true, "clickTrackingEnabled": false, "steps": [{"subject": "Step 1", "content": "Content 1"}]}`,
			expected:     http.StatusInternalServerError,
			serviceError: errors.New("test error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Echo instance
			e := echo.New()

			// Create a new HTTP request with a JSON payload
			req := httptest.NewRequest(http.MethodPost, "/sequences", strings.NewReader(tt.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Create a mock sequence service
			mockSequenceService := &testdata.MockSequenceService{
				CreateSequenceFn: func(ctx context.Context, seq sequence.Sequence) error {
					return tt.serviceError
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

			// Check if the response status code matches the expected value
			if rec.Code != tt.expected {
				t.Errorf("expected status code %d, got %d", tt.expected, rec.Code)
			}
		})
	}
}
func TestGetSequence(t *testing.T) {
	testCases := []struct {
		name           string
		idParamValue   string
		expectedStatus int
		expectedBody   string
		sequence       sequence.Sequence
		serviceError   error
	}{
		{
			name:           "Success",
			idParamValue:   "1",
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"id\":1,\"name\":\"Test Sequence\",\"openTrackingEnabled\":false,\"clickTrackingEnabled\":false,\"steps\":[{\"id\":1,\"subject\":\"Step 1\",\"content\":\"Content 1\"}]}\n",
			sequence: sequence.Sequence{
				ID:            1,
				Name:          "Test Sequence",
				OpenTracking:  false,
				ClickTracking: false,
				Steps: []sequence.Step{
					{ID: 1, Subject: "Step 1", Content: "Content 1"},
				},
			},
			serviceError: nil,
		},
		{
			name:           "Invalid ID param",
			idParamValue:   "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "id must be an integer",
			sequence:       sequence.Sequence{},
			serviceError:   nil,
		},
		{
			name:           "Not Found Error",
			idParamValue:   "1",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "sequence with given ID not found",
			sequence:       sequence.Sequence{},
			serviceError:   sequence.ErrSequenceNotFound,
		},
		{
			name:           "Unknown Error",
			idParamValue:   "1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "test error",
			sequence:       sequence.Sequence{},
			serviceError:   errors.New("test error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new Echo instance
			e := echo.New()

			// Create a new HTTP request
			req := httptest.NewRequest(http.MethodGet, "/sequences/"+tc.idParamValue, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tc.idParamValue)

			// Create a mock sequence service
			mockSequenceService := &testdata.MockSequenceService{
				GetSequenceFn: func(ctx context.Context, id int) (sequence.Sequence, error) {
					return tc.sequence, tc.serviceError
				},
			}

			// Create a new server instance with the mock sequence service
			server := transporthttp.NewServer(mockSequenceService)

			// Call the GetSequence method
			err := server.GetSequence(c)

			// Check if there was an error
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			// Check if the response status code matches the expected status code
			if rec.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, rec.Code)
			}

			// Check if the response body matches the expected body
			if rec.Body.String() != tc.expectedBody {
				t.Errorf("expected body %q, got %q", tc.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestPatchSequence(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		serviceError   error
	}{
		{
			name:           "Success",
			requestBody:    `{"name": "Test Name", "openTracking": true, "clickTracking": false}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Body",
			requestBody:    `{"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Validation Error",
			requestBody:    `{"name": "", "openTracking": true, "clickTracking": false}`,
			expectedStatus: http.StatusBadRequest,
			serviceError:   sequence.ErrSequenceValidation,
		},
		{
			name:           "Not Found Error",
			requestBody:    `{"id": 1, "name": "Test Name", "openTracking": true, "clickTracking": false}`,
			expectedStatus: http.StatusBadRequest,
			serviceError:   sequence.ErrSequenceNotFound,
		},
		{
			name:           "Unknown Error",
			requestBody:    `{"id": 1, "name": "Test Name", "openTracking": true, "clickTracking": false}`,
			expectedStatus: http.StatusInternalServerError,
			serviceError:   errors.New("test error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Echo instance
			e := echo.New()

			// Create a new HTTP request with a JSON payload
			req := httptest.NewRequest(http.MethodPatch, "/sequence/1", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues("1")

			// Create a mock sequence service
			mockSequenceService := &testdata.MockSequenceService{
				PatchSequenceFn: func(ctx context.Context, seq sequence.SequencePatch) error {
					return tt.serviceError
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

			// Check if the response status code is as expected
			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
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
