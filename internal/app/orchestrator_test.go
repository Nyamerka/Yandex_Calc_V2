package app

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func TestHandleCalculateRequest(t *testing.T) {
	orchestrator := NewOrchestrator()
	router := gin.Default()
	router.POST("/api/v1/calculate", orchestrator.handleCalculateRequest)

	tests := []struct {
		name           string
		inputBody      string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Expression",
			inputBody:      `{"expression": "1 + 2"}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"1"}`,
		},
		{
			name:           "Invalid Expression",
			inputBody:      `{"expression": "1 + "}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":"Invalid expression"}`,
		},
		{
			name:           "Empty Body",
			inputBody:      ``,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":"Invalid Body"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBufferString(test.inputBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			if recorder.Code != test.expectedStatus {
				t.Errorf("Expected status %d, got %d", test.expectedStatus, recorder.Code)
			}

			if recorder.Body.String() != test.expectedBody {
				t.Errorf("Expected body %s, got %s", test.expectedBody, recorder.Body.String())
			}
		})
	}
}

func TestHandleExpressionsRequest(t *testing.T) {
	orchestrator := NewOrchestrator()
	router := gin.Default()
	router.GET("/api/v1/expressions", orchestrator.handleExpressionsRequest)

	orchestrator.mutex.Lock()
	orchestrator.expressionCounter++
	exprID := strconv.FormatInt(orchestrator.expressionCounter, 10)
	expr := &Expression{
		ID:     exprID,
		Expr:   "1 + 2",
		Status: "pending",
		AST:    &ASTNode{IsLeaf: true, Value: 3},
	}
	orchestrator.expressionStore[exprID] = expr
	orchestrator.mutex.Unlock()

	req, err := http.NewRequest("GET", "/api/v1/expressions", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	expectedBody := `{"expressions":[{"id":"1","expression":"1 + 2","status":"completed","result":3}]}`

	if recorder.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, recorder.Body.String())
	}
}

func TestHandleExpressionByIdRequest(t *testing.T) {
	orchestrator := NewOrchestrator()
	router := gin.Default()
	router.GET("/api/v1/expressions/:id", orchestrator.handleExpressionByIdRequest)

	orchestrator.mutex.Lock()
	orchestrator.expressionCounter++
	exprID := strconv.FormatInt(orchestrator.expressionCounter, 10)
	expr := &Expression{
		ID:     exprID,
		Expr:   "1 + 2",
		Status: "pending",
		AST:    &ASTNode{IsLeaf: true, Value: 3},
	}
	orchestrator.expressionStore[exprID] = expr
	orchestrator.mutex.Unlock()

	tests := []struct {
		name           string
		expressionID   string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Expression ID",
			expressionID:   exprID,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"expression":{"id":"1","expression":"1 + 2","status":"completed","result":3}}`,
		},
		{
			name:           "Invalid Expression ID",
			expressionID:   "999",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Expression not found"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/v1/expressions/"+test.expressionID, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			if recorder.Code != test.expectedStatus {
				t.Errorf("Expected status %d, got %d", test.expectedStatus, recorder.Code)
			}

			if recorder.Body.String() != test.expectedBody {
				t.Errorf("Expected body %s, got %s", test.expectedBody, recorder.Body.String())
			}
		})
	}
}

func TestHandleGetTaskRequest(t *testing.T) {
	orchestrator := NewOrchestrator()
	router := gin.Default()
	router.GET("/internal/task", orchestrator.handleGetTaskRequest)

	orchestrator.mutex.Lock()
	orchestrator.taskCounter++
	taskID := strconv.FormatInt(orchestrator.taskCounter, 10)
	task := &Task{
		ID:            taskID,
		ExprID:        "1",
		Arg1:          1,
		Arg2:          2,
		Operation:     "+",
		OperationTime: 100,
	}
	orchestrator.taskStorage[taskID] = task
	orchestrator.taskQueue.PushBack(task)
	orchestrator.mutex.Unlock()

	req, err := http.NewRequest("GET", "/internal/task", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	expectedBody := `{"task":{"id":"1","arg1":1,"arg2":2,"operation":"+","operation_time":100}}`

	if recorder.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, recorder.Body.String())
	}
}

func TestHandlePostTaskRequest(t *testing.T) {
	orchestrator := NewOrchestrator()

	exprID := "1"
	expr := &Expression{
		ID:     exprID,
		Expr:   "1+2",
		Status: "in_progress",
		AST: &ASTNode{
			IsLeaf:   false,
			Operator: "+",
			Left: &ASTNode{
				IsLeaf: true,
				Value:  1,
			},
			Right: &ASTNode{
				IsLeaf: true,
				Value:  2,
			},
		},
	}
	orchestrator.mutex.Lock()
	orchestrator.expressionStore[exprID] = expr

	taskID := "1"
	task := &Task{
		ID:            taskID,
		ExprID:        exprID,
		Arg1:          1,
		Arg2:          2,
		Operation:     "+",
		OperationTime: 100,
		Node:          expr.AST,
	}
	orchestrator.taskStorage[taskID] = task
	orchestrator.mutex.Unlock()

	router := gin.Default()
	router.POST("/internal/task", orchestrator.handlePostTaskRequest)

	tests := []struct {
		name           string
		inputBody      string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Task Result",
			inputBody:      `{"id":"1","result":3}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"result accepted"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/internal/task", strings.NewReader(test.inputBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.expectedStatus {
				t.Errorf("Expected status %d, got %d", test.expectedStatus, w.Code)
			}

			if strings.TrimSpace(w.Body.String()) != test.expectedBody {
				t.Errorf("Expected body %s, got %s", test.expectedBody, w.Body.String())
			}
		})
	}
}
