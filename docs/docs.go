// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/execlogs": {
            "get": {
                "description": "This endpoint tries to get the logs of an executable that is set in the orchestrator.",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "orchestrator"
                ],
                "summary": "Get the logs of an executable",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "UUID of the executable to get logs",
                        "name": "id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Type of logs to get",
                        "name": "type",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "Offset of the logs to get",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    }
                }
            }
        },
        "/run": {
            "get": {
                "description": "This endpoint tries to run an executable that is set in the orchestrator.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orchestrator"
                ],
                "summary": "Run an executable",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "UUID of the executable to run",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    }
                }
            }
        },
        "/runall": {
            "get": {
                "description": "This endpoint tries to run all the executables that are set in the orchestrator.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orchestrator"
                ],
                "summary": "Run all the executables",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    }
                }
            }
        },
        "/rungroup": {
            "get": {
                "description": "This endpoint tries to run a group of executables that are set in the orchestrator.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orchestrator"
                ],
                "summary": "Run a group of executables",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Group ID to run",
                        "name": "group",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    }
                }
            }
        },
        "/set": {
            "get": {
                "description": "This endpoint sets the executables in the orchestrator. In order to set the executables again, all processes must be stopped and unset.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orchestrator"
                ],
                "summary": "Set the executables",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    }
                }
            }
        },
        "/status": {
            "get": {
                "description": "This endpoint returns the status of the executables that are set in the orchestrator.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orchestrator"
                ],
                "summary": "Return the status of the executables",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/orchestrator.Status"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    }
                }
            }
        },
        "/stop": {
            "get": {
                "description": "This endpoint tries to stop an executable that is set in the orchestrator.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orchestrator"
                ],
                "summary": "Stops an executable",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "UUID of the group to stop",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    }
                }
            }
        },
        "/stopall": {
            "get": {
                "description": "This endpoint tries to stop all the executables that are set in the orchestrator.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orchestrator"
                ],
                "summary": "Stops all the executables",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    }
                }
            }
        },
        "/stopgroup": {
            "get": {
                "description": "This endpoint tries to stop a group of executables that are set in the orchestrator.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orchestrator"
                ],
                "summary": "Stops a group of executables",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Group ID to stop",
                        "name": "group",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    }
                }
            }
        },
        "/unset": {
            "get": {
                "description": "This endpoint unsets the executables in the orchestrator. In order to unset the executables, all processes must be stopped.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orchestrator"
                ],
                "summary": "Unset the executables",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dtos.GenericResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dtos.GenericResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "orchestrator.Status": {
            "type": "object",
            "properties": {
                "auto_restart": {
                    "type": "boolean"
                },
                "group": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "pid": {
                    "type": "integer"
                },
                "running": {
                    "type": "boolean"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.1",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "orchestrator-api",
	Description:      "This is an API that controls running processes.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
