package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iter"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/server/adkrest/controllers"
	"google.golang.org/adk/server/adkrest/internal/fakes"
	"google.golang.org/adk/server/adkrest/internal/models"
	"google.golang.org/adk/session"
)

var defaultTriggerConfig = launcher.TriggerConfig{
	MaxConcurrentRuns: 10,
	MaxRetries:        3,
	BaseDelay:         1 * time.Millisecond,
	MaxDelay:          5 * time.Millisecond,
}

func TestPubSubTriggerHandler(t *testing.T) {
	tests := []struct {
		name             string
		mockAgentResults []error
		expectedCode     int
		expectedRunCount int
	}{
		{
			name:             "Success_Immediate",
			mockAgentResults: nil,
			expectedCode:     http.StatusOK,
			expectedRunCount: 1,
		},
		{
			name:             "ResourceExhaustedRetry",
			mockAgentResults: []error{fmt.Errorf("429 ResourceExhausted"), fmt.Errorf("429 ResourceExhausted")},
			expectedCode:     http.StatusOK,
			expectedRunCount: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAgentRunCount := 0
			testAgent := createMockAgent(t, tc.mockAgentResults, &mockAgentRunCount)

			apiController := setupTest(t, testAgent)

			reqObj := models.PubSubTriggerRequest{
				Message: models.PubSubMessage{
					Data: []byte("Hello agent"),
				},
				Subscription: "test-sub",
			}
			reqBytes, err := json.Marshal(reqObj)
			if err != nil {
				t.Fatalf("marshal request: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, "/triggers/pubsub", bytes.NewBuffer(reqBytes))
			if err != nil {
				t.Fatalf("new request: %v", err)
			}
			rr := httptest.NewRecorder()

			apiController.PubSubTriggerHandler(rr, req)

			if rr.Code != tc.expectedCode {
				t.Errorf("expected status %d, got %d. Body: %s", tc.expectedCode, rr.Code, rr.Body.String())
			}

			if mockAgentRunCount != tc.expectedRunCount {
				t.Errorf("expected %d run attempts, got %d", tc.expectedRunCount, mockAgentRunCount)
			}
		})
	}
}

func setupTest(t *testing.T, a agent.Agent) *controllers.TriggersAPIController {
	t.Helper()
	sessionService := &fakes.FakeSessionService{Sessions: make(map[fakes.SessionKey]fakes.TestSession)}
	agentLoader := agent.NewSingleLoader(a)
	return controllers.NewTriggersAPIController(sessionService, agentLoader, nil, nil, runner.PluginConfig{}, defaultTriggerConfig)
}

func createMockAgent(t *testing.T, results []error, runCount *int) agent.Agent {
	t.Helper()
	testAgent, err := agent.New(agent.Config{
		Name: "test-agent",
		Run: func(ctx agent.InvocationContext) iter.Seq2[*session.Event, error] {
			return func(yield func(*session.Event, error) bool) {
				*runCount++
				if *runCount <= len(results) {
					err := results[*runCount-1]
					if err != nil {
						yield(nil, err)
						return
					}
				}
				yield(&session.Event{ID: "success-event"}, nil)
			}
		},
	})
	if err != nil {
		t.Fatalf("agent.New failed: %v", err)
	}
	return testAgent
}
