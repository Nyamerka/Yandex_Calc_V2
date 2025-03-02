package app

import (
	"Yandex_Calc_V2.0/internal/eval"
	"Yandex_Calc_V2.0/internal/queue"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	swaggerfiles "github.com/swaggo/files"

	"github.com/gin-gonic/gin"
)

// TODO: Вынести всё это добро либо в .env, либо в файл .yaml для Docker
const (
	DEFAULT_PORT            = "8080"
	TIME_ADDITION_MS        = 200
	TIME_SUBTRACTION_MS     = 152
	TIME_MULTIPLICATIONS_MS = 228
	TIME_DIVISIONS_MS       = 300
)

type Expression struct {
	ID     string   `json:"id"`
	Expr   string   `json:"expression"`
	Status string   `json:"status"`
	Result *float64 `json:"result,omitempty"`
	AST    *ASTNode `json:"-"`
}

// ExpressionRequest swagger model
// @Description Математическое выражение для расчёта
type ExpressionRequest struct {
	Expression string `json:"expression" binding:"required" example:"2+3*4-5/2"`
}

// ExpressionResponse swagger model
// @Description Ответ с идентификатором задачи
type ExpressionResponse struct {
	ID string `json:"id" example:"1"`
}

// Error swagger model
// @Description Описание ошибки
type Error struct {
	Error string `json:"error" example:"Invalid expression"`
}

type Task struct {
	ID            string   `json:"id"`
	ExprID        string   `json:"-"`
	Arg1          float64  `json:"arg1"`
	Arg2          float64  `json:"arg2"`
	Operation     string   `json:"operation"`
	OperationTime int      `json:"operation_time"`
	Node          *ASTNode `json:"-"`
}

type OrchestratorConfig struct {
	WorkingPort           string
	TimeForAddition       int
	TimeForSubtraction    int
	TimeForMultiplication int
	TimeForDivision       int
}

func SetDefaultOrchestratorConfig() *OrchestratorConfig {
	return &OrchestratorConfig{
		WorkingPort:           DEFAULT_PORT,
		TimeForAddition:       TIME_ADDITION_MS,
		TimeForSubtraction:    TIME_SUBTRACTION_MS,
		TimeForMultiplication: TIME_MULTIPLICATIONS_MS,
		TimeForDivision:       TIME_DIVISIONS_MS,
	}
}

type Orchestrator struct {
	Config            *OrchestratorConfig
	expressionStore   map[string]*Expression
	taskStorage       map[string]*Task
	taskQueue         queue.Queue
	mutex             sync.Mutex
	expressionCounter int64
	taskCounter       int64
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		Config:          SetDefaultOrchestratorConfig(),
		expressionStore: make(map[string]*Expression),
		taskStorage:     make(map[string]*Task),
		taskQueue:       *queue.New(),
	}
}

// @Summary Schedule mathematical expression calculation
// @Description Parse expression and create a new calculation task
// @Tags calculations
// @Accept json
// @Produce json
// @Param expression body ExpressionRequest true "Mathematical expression to calculate"
// @Success 201 {object} ExpressionResponse "Calculation ID"
// @Failure 400 {object} Error "Invalid request body"
// @Failure 500 {object} Error "Internal server error"
// @Router /calculate [post]
func (o *Orchestrator) handleCalculateRequest(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Wrong Method"})
		return
	}
	var req ExpressionRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Expression == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid Body"})
		return
	}
	ast, err := ParseAST(req.Expression)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid expression"})
		return
	}
	o.mutex.Lock()
	o.expressionCounter++
	exprID := strconv.FormatInt(o.expressionCounter, 10)
	expr := &Expression{
		ID:     exprID,
		Expr:   req.Expression,
		Status: "pending",
		AST:    ast,
	}
	o.expressionStore[exprID] = expr
	o.scheduleTasksForExpression(expr)
	o.mutex.Unlock()
	c.JSON(http.StatusCreated, gin.H{"id": exprID})
}

// @Summary Get all calculated expressions
// @Description Retrieve list of all expressions with their current status
// @Tags calculations
// @Produce json
// @Success 200 {array} ExpressionResponse
// @Router /expressions [get]
func (o *Orchestrator) handleExpressionsRequest(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Wrong Method"})
		return
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	exprs := make([]*Expression, 0, len(o.expressionStore))
	for _, expr := range o.expressionStore {
		if expr.AST != nil && expr.AST.IsLeaf {
			expr.Status = "completed"
			expr.Result = &expr.AST.Value
		}
		exprs = append(exprs, expr)
	}
	c.JSON(http.StatusOK, gin.H{"expressions": exprs})
}

// @Summary Get expression by ID
// @Description Retrieve specific expression details by unique identifier
// @Tags calculations
// @Produce json
// @Param id path string true "Expression ID"
// @Success 200 {object} ExpressionResponse
// @Failure 404 {object} Error "Expression not found"
// @Router /expressions/{id} [get]
func (o *Orchestrator) handleExpressionByIdRequest(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Wrong Method"})
		return
	}
	id := c.Param("id")
	o.mutex.Lock()
	expr, ok := o.expressionStore[id]
	o.mutex.Unlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expression not found"})
		return
	}
	if expr.AST != nil && expr.AST.IsLeaf {
		expr.Status = "completed"
		expr.Result = &expr.AST.Value
	}
	c.JSON(http.StatusOK, gin.H{"expression": expr})
}

// @Summary Fetch next available task
// @Description Get the next task from the calculation queue (internal use)
// @Tags internal
// @Produce json
// @Success 200 {object} TaskResponse
// @Failure 404 {object} Error "No tasks available"
// @Router /internal/task [get]
func (o *Orchestrator) handleGetTaskRequest(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Wrong Method"})
		return
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.taskQueue.Len() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No task available"})
		return
	}
	taskInterface := o.taskQueue.PopFront()
	if taskInterface == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: task is nil"})
		return
	}
	task, ok := taskInterface.(*Task)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: invalid task type"})
		return
	}
	if expr, exists := o.expressionStore[task.ExprID]; exists {
		expr.Status = "in_progress"
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}

// @Summary Submit task result
// @Description Report calculation result for a specific task (internal use)
// @Tags internal
// @Accept json
// @Produce json
// @Param taskResult body TaskResult true "Task result data"
// @Success 200 {object} SuccessResponse "Result accepted"
// @Failure 400 {object} Error "Invalid request body"
// @Failure 404 {object} Error "Task not found"
// @Router /internal/task [post]
func (o *Orchestrator) handlePostTaskRequest(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Wrong Method"})
		return
	}
	var req TaskResult
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid Body"})
		return
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	task, ok := o.taskStorage[req.ID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	task.Node.Value = req.Result
	task.Node.IsLeaf = true

	delete(o.taskStorage, req.ID)
	if expr, exists := o.expressionStore[task.ExprID]; exists {
		o.scheduleTasksForExpression(expr)
		if expr.AST.IsLeaf {
			expr.Status = "completed"
			tmp, err := eval.Eval(expr.Expr)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Result computation incorrect"})
			}
			res := eval.BigratToFloat(tmp)
			if res != expr.AST.Value {
				expr.Result = &expr.AST.Value
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Result computation incorrect"})
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "result accepted"})
}

func (o *Orchestrator) scheduleTasksForExpression(expr *Expression) {
	var traverse func(node *ASTNode)
	traverse = func(node *ASTNode) {
		if node == nil || node.IsLeaf {
			return
		}
		traverse(node.Left)
		traverse(node.Right)
		if node.Left != nil && node.Right != nil && node.Left.IsLeaf && node.Right.IsLeaf {
			if !node.TaskScheduled {
				o.taskCounter++
				taskID := strconv.FormatInt(o.taskCounter, 10)
				var opTime int
				switch node.Operator {
				case "+":
					opTime = o.Config.TimeForAddition
				case "-":
					opTime = o.Config.TimeForSubtraction
				case "*":
					opTime = o.Config.TimeForMultiplication
				case "/":
					opTime = o.Config.TimeForDivision
				default:
					opTime = 222
				}
				task := &Task{
					ID:            taskID,
					ExprID:        expr.ID,
					Arg1:          node.Left.Value,
					Arg2:          node.Right.Value,
					Operation:     node.Operator,
					OperationTime: opTime,
					Node:          node,
				}
				node.TaskScheduled = true
				o.taskStorage[taskID] = task
				o.taskQueue.PushBack(task)
			}
		}
	}
	traverse(expr.AST)
}

func (o *Orchestrator) StartServer() error {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.POST("/api/v1/calculate", o.handleCalculateRequest)
	r.GET("/api/v1/expressions", o.handleExpressionsRequest)
	r.GET("/api/v1/expressions/:id", o.handleExpressionByIdRequest)
	r.GET("/internal/task", o.handleGetTaskRequest)
	r.POST("/internal/task", o.handlePostTaskRequest)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})

	go func() {
		for {
			time.Sleep(4 * time.Second)
			o.mutex.Lock()
			if o.taskQueue.Len() > 0 {
				log.Printf("Collecting tasks in queue: %d", o.taskQueue.Len())
			}
			o.mutex.Unlock()
		}
	}()

	return r.Run(":" + o.Config.WorkingPort)
}
