package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cybre/salesforge-assignment/internal/database"
	"github.com/cybre/salesforge-assignment/internal/sequence"
	transporthttp "github.com/cybre/salesforge-assignment/internal/transport/http"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestServer struct {
	Address    string
	Repository sequence.Repository
}

func NewTestServer(t *testing.T) *TestServer {
	ctx := context.Background()

	newNetwork, err := network.New(ctx, network.WithCheckDuplicate())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := newNetwork.Remove(ctx); err != nil {
			t.Fatalf("failed to remove network: %s", err)
		}
	})

	networkName := newNetwork.Name

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres"),
		postgres.WithDatabase("salesforge"),
		postgres.WithUsername("salesforge"),
		postgres.WithPassword("salesforge"),
		network.WithNetwork([]string{"postgres"}, newNetwork),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start postgres: %s", err)
	}

	t.Cleanup(func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate postgres container: %s", err)
		}
	})

	postgresContainerHost, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get postgres container host: %v", err)
	}

	postgresContainerPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get postgres container port: %v", err)
	}

	database, err := database.NewPostgresDB(database.Config{
		Host:     postgresContainerHost,
		Port:     postgresContainerPort.Port(),
		Name:     "salesforge",
		User:     "salesforge",
		Password: "salesforge",
	})
	if err != nil {
		t.Fatalf("failed to create database connection: %v", err)
	}

	repository := sequence.NewPostgresRepository(database)

	seqContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context: "../",
			},
			Networks: []string{
				networkName,
			},
			ExposedPorts: []string{"8080/tcp"},
			Env: map[string]string{
				"PORT":              "8080",
				"DATABASE_HOST":     "postgres",
				"DATABASE_PORT":     "5432",
				"DATABASE_NAME":     "salesforge",
				"DATABASE_USER":     "salesforge",
				"DATABASE_PASSWORD": "salesforge",
			},
			WaitingFor: wait.ForHTTP("/health").WithStartupTimeout(5 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("failed to start sequence container: %v", err)
	}

	t.Cleanup(func() {
		if err := seqContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate service container: %s", err)
		}
	})

	seqContainerHost, err := seqContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get sequence container host: %v", err)
	}
	seqContainerPort, err := seqContainer.MappedPort(ctx, "8080")
	if err != nil {
		t.Fatalf("failed to get sequence container port: %v", err)
	}

	return &TestServer{
		Address:    fmt.Sprintf("http://%s:%s", seqContainerHost, seqContainerPort.Port()),
		Repository: repository,
	}
}

func (ts *TestServer) CreateSequence(t *testing.T, request transporthttp.CreateSequenceRequest) *http.Response {
	validPayload, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, ts.Address+"/sequence", bytes.NewReader(validPayload))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	return res
}

func (ts *TestServer) GetSequence(t *testing.T, id int) *http.Response {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/sequence/%d", ts.Address, id), nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	return res
}

func (ts *TestServer) PatchSequence(t *testing.T, patch transporthttp.PatchSequenceRequest) *http.Response {
	payload, err := json.Marshal(patch)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/sequence/%d", ts.Address, patch.ID), bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	return res
}

func (ts *TestServer) PutStep(t *testing.T, request transporthttp.UpdateStepRequest) *http.Response {
	payload, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/step/%d", ts.Address, request.ID), bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	return res
}

func (ts *TestServer) DeleteStep(t *testing.T, id int) *http.Response {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/step/%d", ts.Address, id), nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	return res
}
