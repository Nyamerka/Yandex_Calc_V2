package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// TODO: Вынести всё это добро либо в .env, либо в файл .yaml для Docker
const (
	COMPUTING_POWER  = 1
	ORCHESTRATOR_URL = "http://localhost:8080"
)

// TaskResponse swagger model
// @Description Информация о задаче
type TaskResponse struct {
	Task struct {
		ID            string  `json:"id" example:"1"`
		Arg1          float64 `json:"arg1" example:"2.0"`
		Arg2          float64 `json:"arg2" example:"3.0"`
		Operation     string  `json:"operation" example:"+"`
		OperationTime int     `json:"operation_time" example:"200"`
	} `json:"task"`
}

// TaskResult swagger model
// @Description Результат задачи
type TaskResult struct {
	ID     string  `json:"id" example:"1"`
	Result float64 `json:"result" example:"5.0"`
}

type Agent struct {
	ComputingPower  int
	OrchestratorURL string
}

func SetDefaultAgent() *Agent {
	return &Agent{
		ComputingPower:  COMPUTING_POWER,
		OrchestratorURL: ORCHESTRATOR_URL,
	}
}

func NewAgent() *Agent {
	return SetDefaultAgent()
}

func (a *Agent) Run() {
	for i := 0; i < a.ComputingPower; i++ {
		log.Printf("Starting worker %d", i)
		go a.worker(i)
	}
	select {}
}

func (a *Agent) Calculate(op string, x, y float64) (float64, error) {
	switch op {
	case "+":
		return x + y, nil
	case "-":
		return x - y, nil
	case "*":
		return x * y, nil
	case "/":
		if y == 0 {
			return 0, errors.New("division by zero is not allowed")
		}
		return x / y, nil
	default:
		return 0, errors.New(fmt.Sprintf("invalid operator: %s", op))
	}
}

func (a *Agent) worker(id int) {
	for {
		resp, err := http.Get(a.OrchestratorURL + "/internal/task")
		if err != nil {
			log.Printf("Worker %d: error getting task: %v", id, err)
			time.Sleep(2 * time.Second)
			continue
		}
		if resp.StatusCode == http.StatusNotFound {
			err = resp.Body.Close()
			if err != nil {
				log.Printf("Worker %d: error closing task body: %v", id, err)
			}
			time.Sleep(2 * time.Second)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Worker %d: error reading response body: %v", id, err)
			} else {
				log.Printf("Worker %d: unexpected status code: %d, response: %s", id, resp.StatusCode, string(body))
			}
			err = resp.Body.Close()
			if err != nil {
				log.Printf("Worker %d: error closing task body: %v", id, err)
			}
			time.Sleep(2 * time.Second)
			continue
		}
		var taskResp TaskResponse
		err = json.NewDecoder(resp.Body).Decode(&taskResp)
		if err != nil {
			log.Printf("Worker %d: error decoding task: %v", id, err)
			time.Sleep(2 * time.Second)
			continue
		}
		err = resp.Body.Close()
		if err != nil {
			log.Printf("Worker %d: error closing task body: %v", id, err)
		}
		task := taskResp.Task
		log.Printf("Worker %d: received task %s: %f %s %f, simulating computation %d ms", id, task.ID, task.Arg1, task.Operation, task.Arg2, task.OperationTime)
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		result, err := a.Calculate(task.Operation, task.Arg1, task.Arg2)
		if err != nil {
			log.Printf("Worker %d: error computing task %s: %v", id, task.ID, err)
			continue
		}
		resultPayload := &TaskResult{
			ID:     task.ID,
			Result: result,
		}
		payloadBytes, err := json.Marshal(resultPayload)
		if err != nil {
			log.Printf("Worker %d: error marshaling result for task %s: %v", id, task.ID, err)
			continue
		}
		respPost, err := http.Post(a.OrchestratorURL+"/internal/task", "application/json", bytes.NewReader(payloadBytes))
		if err != nil {
			log.Printf("Worker %d: error posting result for task %s: %v", id, task.ID, err)
			continue
		}
		if respPost.StatusCode != http.StatusOK {
			body, err := io.ReadAll(respPost.Body)
			if err != nil {
				log.Printf("Worker %d: error reading response body for task %s: %v", id, task.ID, err)
				continue
			}
			log.Printf("Worker %d: error response posting result for task %s: %s", id, task.ID, string(body))
		} else {
			log.Printf("Worker %d: successfully completed task %s with result %f", id, task.ID, result)
		}
		err = respPost.Body.Close()
		if err != nil {
			log.Printf("Worker %d: error closing task %s body: %v", id, task.ID, err)
		}
	}
}
