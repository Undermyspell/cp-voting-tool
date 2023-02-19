// Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/question/answer/{id}": {
            "put": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Answers a question of the current session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Question"
                ],
                "summary": "Answers a question",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id of question to answer",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/question/new": {
            "post": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Adds a new question to the current session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Question"
                ],
                "summary": "Adds a new question",
                "parameters": [
                    {
                        "description": "Question JSON",
                        "name": "question",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.NewQuestionDto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/question/session": {
            "get": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Gets the questions of the current session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Question"
                ],
                "summary": "Gets the questions of the current session",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dtos.QuestionDto"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/question/session/start": {
            "post": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Starts a new questions session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Question"
                ],
                "summary": "Starts a new questions session",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/question/session/stop": {
            "post": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Stops the current questions session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Question"
                ],
                "summary": "Stops the current questions session",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/question/upvote/{id}": {
            "put": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Upvotes a question of the current session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Question"
                ],
                "summary": "Upvotes a question",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id of question to upvote",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dtos.NewQuestionDto": {
            "type": "object",
            "required": [
                "text"
            ],
            "properties": {
                "anonymous": {
                    "type": "boolean"
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "dtos.QuestionDto": {
            "type": "object",
            "properties": {
                "anonymous": {
                    "type": "boolean"
                },
                "answered": {
                    "type": "boolean"
                },
                "creator": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "owned": {
                    "type": "boolean"
                },
                "text": {
                    "type": "string"
                },
                "votes": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "JWT": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Voting tool api",
	Description:      "A voting tool API in Go using Gin framework.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
