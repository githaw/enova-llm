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
        "/api/pilot/v1/docker/deploy": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "task"
                ],
                "summary": "Get docker task info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Task Name",
                        "name": "task_name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/enova_internal_meta.DetectTaskSpecResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "task"
                ],
                "summary": "Create alert policy object",
                "parameters": [
                    {
                        "description": "request content",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/enova_internal_meta.DockerDeployRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "task"
                ],
                "summary": "Delete docker task",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Task Name",
                        "name": "task_name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/pilot/v1/task/detect/history": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "task"
                ],
                "summary": "Get task detect history list",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Task Name",
                        "name": "task_name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/enova_internal_meta.TaskDetectHistoryResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "enova_internal_meta.AnomalyRecommendResult": {
            "type": "object",
            "properties": {
                "configRecommendResult": {
                    "$ref": "#/definitions/enova_pkg_api.ConfigRecommendResult"
                },
                "currentConfig": {
                    "$ref": "#/definitions/enova_pkg_api.ConfigRecommendResult"
                },
                "isAnomaly": {
                    "type": "boolean"
                },
                "timestamp": {
                    "type": "integer"
                }
            }
        },
        "enova_internal_meta.DetectTaskSpecResponse": {
            "type": "object",
            "properties": {
                "task_infos": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/enova_internal_meta.TaskInfo"
                    }
                },
                "task_spec": {
                    "$ref": "#/definitions/enova_internal_meta.TaskSpec"
                }
            }
        },
        "enova_internal_meta.DockerDeployRequest": {
            "type": "object",
            "properties": {
                "backend": {
                    "type": "string"
                },
                "backendConfig": {
                    "type": "object",
                    "additionalProperties": true
                },
                "backend_extra_config": {
                    "type": "object",
                    "additionalProperties": true
                },
                "envs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/enova_internal_meta.Env"
                    }
                },
                "exporter_endpoint": {
                    "type": "string"
                },
                "exporter_service_name": {
                    "type": "string"
                },
                "host": {
                    "type": "string"
                },
                "model": {
                    "type": "string"
                },
                "modelConfig": {
                    "$ref": "#/definitions/enova_internal_meta.ModelConfig"
                },
                "name": {
                    "type": "string"
                },
                "port": {
                    "type": "integer"
                },
                "replica": {
                    "type": "integer"
                },
                "volumes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/enova_internal_meta.Volume"
                    }
                }
            }
        },
        "enova_internal_meta.Env": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "enova_internal_meta.ModelConfig": {
            "type": "object",
            "properties": {
                "gpu": {
                    "$ref": "#/definitions/enova_pkg_api.Gpu"
                },
                "llm": {
                    "$ref": "#/definitions/enova_pkg_api.Llm"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "enova_internal_meta.TaskDetectHistoryResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/enova_internal_meta.AnomalyRecommendResult"
                    }
                }
            }
        },
        "enova_internal_meta.TaskInfo": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "enova_internal_meta.TaskSpec": {
            "type": "object",
            "properties": {
                "backend": {
                    "type": "string"
                },
                "backendConfig": {},
                "envs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/enova_internal_meta.Env"
                    }
                },
                "exporter_endpoint": {
                    "type": "string"
                },
                "exporter_service_name": {
                    "type": "string"
                },
                "gpus": {
                    "type": "string"
                },
                "host": {
                    "type": "string"
                },
                "model": {
                    "type": "string"
                },
                "modelConfig": {
                    "$ref": "#/definitions/enova_internal_meta.ModelConfig"
                },
                "name": {
                    "type": "string"
                },
                "port": {
                    "type": "integer"
                },
                "replica": {
                    "type": "integer"
                },
                "volumes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/enova_internal_meta.Volume"
                    }
                }
            }
        },
        "enova_internal_meta.Volume": {
            "type": "object",
            "properties": {
                "hostPath": {
                    "type": "string"
                },
                "mountPath": {
                    "type": "string"
                }
            }
        },
        "enova_pkg_api.ConfigRecommendResult": {
            "type": "object",
            "properties": {
                "gpu_memory_utilization": {
                    "type": "number"
                },
                "max_num_seqs": {
                    "type": "integer"
                },
                "replicas": {
                    "type": "integer"
                },
                "tensor_parallel_size": {
                    "type": "integer"
                }
            }
        },
        "enova_pkg_api.Gpu": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "num": {
                    "type": "integer"
                },
                "spec": {
                    "type": "integer"
                }
            }
        },
        "enova_pkg_api.Llm": {
            "type": "object",
            "properties": {
                "framework": {
                    "type": "string"
                },
                "param": {
                    "type": "number"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	// LeftDelim:        "{{",
	// RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
