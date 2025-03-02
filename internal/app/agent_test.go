package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"
)

type mockOrchestrator struct {
	taskResponse *TaskResponse
	taskResult   chan map[string]interface{}
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name     string
		op       string
		x        float64
		y        float64
		expected float64
		err      error
	}{
		{"add", "+", 1, 2, 3, nil},
		{"subtract", "-", 5, 3, 2, nil},
		{"multiply", "*", 4, 3, 12, nil},
		{"divide", "/", 10, 2, 5, nil},
		{"divide by zero", "/", 10, 0, 0, errors.New("division by zero is not allowed")},
		{"invalid operator", "^", 1, 2, 0, fmt.Errorf("invalid operator: ^")},
		{"add negative numbers", "+", -1, -2, -3, nil},
		{"subtract negative numbers", "-", -5, -3, -2, nil},
		{"multiply negative numbers", "*", -4, -3, 12, nil},
		{"divide negative numbers", "/", -10, -2, 5, nil},
		{"divide positive by negative", "/", 10, -2, -5, nil},
		{"divide negative by positive", "/", -10, 2, -5, nil},
		{"multiply by zero", "*", 10, 0, 0, nil},
		{"add zero", "+", 0, 0, 0, nil},
		{"subtract zero", "-", 0, 0, 0, nil},
		{"divide zero by non-zero", "/", 0, 10, 0, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			agent := NewAgent()
			result, err := agent.Calculate(test.op, test.x, test.y)
			if test.err != nil && err == nil {
				t.Errorf("Expected error %v, got nil", test.err)
			} else if test.err == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			} else if test.err != nil && err != nil && test.err.Error() != err.Error() {
				t.Errorf("Expected error %v, got %v", test.err, err)
			} else if result != test.expected {
				t.Errorf("Expected result %f, got %f", test.expected, result)
			}
		})
	}
}

func (m *mockOrchestrator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/internal/task":
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			if m.taskResponse != nil {
				json.NewEncoder(w).Encode(m.taskResponse)
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
			}
		case http.MethodPost:
			var result map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&result)
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			m.taskResult <- result
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func TestWorker_SuccessfulTask(t *testing.T) {
	mock := &mockOrchestrator{
		taskResponse: &TaskResponse{},
		taskResult:   make(chan map[string]interface{}, 1),
	}
	mock.taskResponse.Task.ID = "task1337"
	mock.taskResponse.Task.Arg1 = 228
	mock.taskResponse.Task.Arg2 = 2
	mock.taskResponse.Task.Operation = "+"
	mock.taskResponse.Task.OperationTime = 100

	server := httptest.NewServer(mock)
	defer server.Close()

	agent := NewAgent()
	agent.OrchestratorURL = server.URL

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		agent.worker(0)
	}()

	time.Sleep(2500 * time.Millisecond)

	select {
	case result := <-mock.taskResult:
		expectedResult := map[string]interface{}{
			"id":     "task1337",
			"result": 230.0,
		}
		if !reflect.DeepEqual(result, expectedResult) {
			t.Errorf("Unexpected result: %+v, expected %+v", result, expectedResult)
		}
	case <-time.After(500 * time.Millisecond):
		t.Errorf("Timeout waiting for task result")
	}

}

func TestWorker_Handle404(t *testing.T) {
	mock := &mockOrchestrator{
		taskResponse: nil,
		taskResult:   make(chan map[string]interface{}, 1),
	}
	server := httptest.NewServer(mock)
	defer server.Close()

	agent := NewAgent()
	agent.OrchestratorURL = server.URL

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		agent.worker(0)
	}()

	time.Sleep(2500 * time.Millisecond)

	select {
	case _, ok := <-mock.taskResult:
		if ok {
			t.Errorf("Unexpected result sent after 404")
		}
	default:
	}
}

func TestWorker_HandleErrorStatusCode(t *testing.T) {
	mock := &mockOrchestrator{
		taskResponse: &TaskResponse{},
		taskResult:   make(chan map[string]interface{}, 1),
	}

	mock.taskResponse.Task.ID = "task228"
	mock.taskResponse.Task.Arg1 = 25
	mock.taskResponse.Task.Arg2 = 27
	mock.taskResponse.Task.Operation = "+"
	mock.taskResponse.Task.OperationTime = 100

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	agent := NewAgent()
	agent.OrchestratorURL = server.URL

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		agent.worker(0)
	}()

	select {
	case <-mock.taskResult:
		t.Errorf("Unexpected result sent after error status code")
	case <-time.After(2500 * time.Millisecond):
	}
}
