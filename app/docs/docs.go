// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/auth/login": {
            "post": {
                "description": "Check username \u0026 password by json",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Check the username \u0026 password correctness for login",
                "operationId": "check-username-password-by-json",
                "parameters": [
                    {
                        "description": "Login Form",
                        "name": "loginForm",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.AuthLogin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/route.CommonTokenResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/route.CommonErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/route.CommonErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/route.CommonErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/route.CommonErrorResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/route.CommonErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/auth/renew": {
            "post": {
                "description": "Renew authToken to avoid expirations by old authToken",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Renew authToken to avoid expirations",
                "operationId": "renew-auth-token-by-auth-token",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Insert your access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/route.CommonTokenResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/route.CommonErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/route.CommonErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/route.CommonErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/route.CommonErrorResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/route.CommonErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.AuthLogin": {
            "type": "object",
            "properties": {
                "account_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "user_name": {
                    "type": "string"
                }
            }
        },
        "model.AuthToken": {
            "type": "object",
            "properties": {
                "aud": {
                    "type": "string"
                },
                "exp": {
                    "type": "integer"
                },
                "iat": {
                    "type": "integer"
                },
                "iss": {
                    "type": "string"
                },
                "jti": {
                    "type": "string"
                },
                "nbf": {
                    "type": "integer"
                },
                "sub": {
                    "type": "string"
                },
                "ttl": {
                    "type": "integer"
                },
                "userID": {
                    "type": "string"
                },
                "user_name": {
                    "type": "string"
                }
            }
        },
        "route.CommonErrorResponse": {
            "type": "object",
            "properties": {
                "authToken": {
                    "$ref": "#/definitions/model.AuthToken"
                },
                "code": {
                    "type": "integer"
                },
                "count": {
                    "type": "integer"
                },
                "data": {
                    "type": "object"
                },
                "details": {
                    "type": "object"
                },
                "error": {
                    "type": "string"
                },
                "ip": {
                    "type": "string"
                },
                "request_id": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "status_message": {
                    "type": "string"
                }
            }
        },
        "route.CommonTokenResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "count": {
                    "type": "integer"
                },
                "data": {
                    "type": "object"
                },
                "request_id": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "status_message": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}