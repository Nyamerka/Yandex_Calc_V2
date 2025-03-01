// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Nyamerka"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/calculate": {
            "post": {
                "description": "Parse expression and create a new calculation task",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "calculations"
                ],
                "summary": "Schedule mathematical expression calculation",
                "parameters": [
                    {
                        "description": "Mathematical expression to calculate",
                        "name": "expression",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/app.ExpressionRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Calculation ID",
                        "schema": {
                            "$ref": "#/definitions/app.ExpressionResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "$ref": "#/definitions/app.Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/app.Error"
                        }
                    }
                }
            }
        },
        "/expressions": {
            "get": {
                "description": "Retrieve list of all expressions with their current status",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "calculations"
                ],
                "summary": "Get all calculated expressions",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/app.ExpressionsResponse"
                            }
                        }
                    }
                }
            }
        },
        "/expressions/{id}": {
            "get": {
                "description": "Retrieve specific expression details by unique identifier",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "calculations"
                ],
                "summary": "Get expression by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Expression ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.ExpressionResponse"
                        }
                    },
                    "404": {
                        "description": "Expression not found",
                        "schema": {
                            "$ref": "#/definitions/app.Error"
                        }
                    }
                }
            }
        },
        "/internal/task": {
            "get": {
                "description": "Get the next task from the calculation queue (internal use)",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "internal"
                ],
                "summary": "Fetch next available task",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.TaskResponse"
                        }
                    },
                    "404": {
                        "description": "No tasks available",
                        "schema": {
                            "$ref": "#/definitions/app.Error"
                        }
                    }
                }
            },
            "post": {
                "description": "Report calculation result for a specific task (internal use)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "internal"
                ],
                "summary": "Submit task result",
                "parameters": [
                    {
                        "description": "Task result data",
                        "name": "taskResult",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/app.TaskResult"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Result accepted",
                        "schema": {
                            "$ref": "#/definitions/app.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "$ref": "#/definitions/app.Error"
                        }
                    },
                    "404": {
                        "description": "Task not found",
                        "schema": {
                            "$ref": "#/definitions/app.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "app.Error": {
            "description": "Описание ошибки",
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Invalid expression"
                }
            }
        },
        "app.ExpressionRequest": {
            "description": "Математическое выражение для расчёта",
            "type": "object",
            "required": [
                "expression"
            ],
            "properties": {
                "expression": {
                    "type": "string",
                    "example": "2+3*4-5/2"
                }
            }
        },
        "app.ExpressionResponse": {
            "description": "Ответ с идентификатором задачи",
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "1"
                }
            }
        },
        "app.ExpressionsResponse": {
            "description": "Ответ с идентификатором задачи",
            "type": "object",
            "properties": {
                "expression": {
                    "type": "string",
                    "example": "2+3*4-5/2"
                },
                "id": {
                    "type": "string",
                    "example": "1"
                },
                "result": {
                    "type": "number"
                },
                "status": {
                    "type": "string",
                    "example": "completed"
                }
            }
        },
        "app.SuccessResponse": {
            "description": "Успешный ответ",
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "result accepted"
                }
            }
        },
        "app.TaskResponse": {
            "description": "Информация о задаче",
            "type": "object",
            "properties": {
                "task": {
                    "type": "object",
                    "properties": {
                        "arg1": {
                            "type": "number",
                            "example": 2
                        },
                        "arg2": {
                            "type": "number",
                            "example": 3
                        },
                        "id": {
                            "type": "string",
                            "example": "1"
                        },
                        "operation": {
                            "type": "string",
                            "example": "+"
                        },
                        "operation_time": {
                            "type": "integer",
                            "example": 200
                        }
                    }
                }
            }
        },
        "app.TaskResult": {
            "description": "Результат задачи",
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "1"
                },
                "result": {
                    "type": "number",
                    "example": 5
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Yandex Calculator API",
	Description:      "Calculator service with distributed computation",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
