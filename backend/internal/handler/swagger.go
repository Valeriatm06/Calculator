package handler

import (
	"net/http"
)

// SwaggerSpecJSON is the raw OpenAPI 3.0 specification for the Calculator REST API.
const SwaggerSpecJSON = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Calculator REST API",
    "description": "Production-grade arithmetic microservice built in Go for Sezzle technical evaluation.",
    "version": "1.0.0",
    "contact": {
      "name": "Valeria Tocarruncho Mosquera"
    }
  },
  "servers": [
    {
      "url": "http://localhost:8080",
      "description": "Direct Go backend service"
    },
    {
      "url": "http://localhost:3000",
      "description": "Via Docker / Nginx Proxy"
    }
  ],
  "paths": {
    "/api/v1/health": {
      "get": {
        "summary": "Health Check",
        "description": "Returns service health status and version.",
        "responses": {
          "200": {
            "description": "Service is healthy",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/HealthResponse"
                },
                "example": {
                  "service": "calculator-api",
                  "status": "healthy",
                  "version": "1.0.0"
                }
              }
            }
          }
        }
      }
    },
    "/api/v1/calculate": {
      "post": {
        "summary": "Calculate arithmetic operation",
        "description": "Executes basic (add, subtract, multiply, divide) and advanced (power, sqrt, percentage) arithmetic operations.",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/CalculateRequest"
              },
              "examples": {
                "addition": {
                  "summary": "Addition (55 + 55)",
                  "value": {
                    "operation": "add",
                    "a": 55,
                    "b": 55
                  }
                },
                "division": {
                  "summary": "Division (100 / 4)",
                  "value": {
                    "operation": "divide",
                    "a": 100,
                    "b": 4
                  }
                },
                "power": {
                  "summary": "Exponentiation (2 ^ 8)",
                  "value": {
                    "operation": "power",
                    "a": 2,
                    "b": 8
                  }
                },
                "square_root": {
                  "summary": "Square Root (√49)",
                  "value": {
                    "operation": "sqrt",
                    "a": 49
                  }
                },
                "percentage": {
                  "summary": "Percentage (15%)",
                  "value": {
                    "operation": "percentage",
                    "a": 15
                  }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Successful calculation result",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/CalculationResponse"
                },
                "example": {
                  "operation": "add",
                  "a": 55,
                  "b": 55,
                  "result": 110,
                  "formatted": "55 + 55 = 110"
                }
              }
            }
          },
          "400": {
            "description": "Bad Request (e.g. division by zero, negative square root, missing operands)",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                },
                "examples": {
                  "division_by_zero": {
                    "summary": "Division by zero",
                    "value": {
                      "success": false,
                      "error": "division by zero is undefined",
                      "code": "DIVISION_BY_ZERO"
                    }
                  },
                  "negative_sqrt": {
                    "summary": "Negative square root",
                    "value": {
                      "success": false,
                      "error": "square root of negative number is undefined for real numbers",
                      "code": "NEGATIVE_SQUARE_ROOT"
                    }
                  }
                }
              }
            }
          }
        }
      }
    },
    "/api/v1/operations/{op}": {
      "post": {
        "summary": "Dedicated operation sub-endpoint",
        "description": "Calculates using the operation specified in the URL path (add, subtract, multiply, divide, power, sqrt, percentage).",
        "parameters": [
          {
            "name": "op",
            "in": "path",
            "required": true,
            "description": "The arithmetic operation to perform",
            "schema": {
              "type": "string",
              "enum": ["add", "subtract", "multiply", "divide", "power", "sqrt", "percentage"]
            }
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/OperationSubRequest"
              },
              "example": {
                "a": 25,
                "b": 75
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Operation result",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/CalculationResponse"
                }
              }
            }
          },
          "400": {
            "description": "Invalid calculation request",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "CalculateRequest": {
        "type": "object",
        "required": ["operation", "a"],
        "properties": {
          "operation": {
            "type": "string",
            "description": "Operation name or symbol",
            "enum": ["add", "subtract", "multiply", "divide", "power", "sqrt", "percentage", "+", "-", "*", "/", "^", "%"],
            "example": "add"
          },
          "a": {
            "type": "number",
            "format": "double",
            "description": "First operand (or target for unary operations like sqrt)",
            "example": 55
          },
          "b": {
            "type": "number",
            "format": "double",
            "description": "Second operand (required for binary operations)",
            "example": 55
          }
        }
      },
      "OperationSubRequest": {
        "type": "object",
        "required": ["a"],
        "properties": {
          "a": {
            "type": "number",
            "format": "double",
            "example": 25
          },
          "b": {
            "type": "number",
            "format": "double",
            "example": 75
          }
        }
      },
      "CalculationResponse": {
        "type": "object",
        "properties": {
          "operation": {
            "type": "string",
            "example": "add"
          },
          "a": {
            "type": "number",
            "format": "double",
            "example": 55
          },
          "b": {
            "type": "number",
            "format": "double",
            "example": 55
          },
          "result": {
            "type": "number",
            "format": "double",
            "example": 110
          },
          "formatted": {
            "type": "string",
            "example": "55 + 55 = 110"
          }
        }
      },
      "HealthResponse": {
        "type": "object",
        "properties": {
          "service": {
            "type": "string",
            "example": "calculator-api"
          },
          "status": {
            "type": "string",
            "example": "healthy"
          },
          "version": {
            "type": "string",
            "example": "1.0.0"
          }
        }
      },
      "ErrorResponse": {
        "type": "object",
        "properties": {
          "success": {
            "type": "boolean",
            "example": false
          },
          "error": {
            "type": "string",
            "example": "division by zero is undefined"
          },
          "code": {
            "type": "string",
            "example": "DIVISION_BY_ZERO"
          }
        }
      }
    }
  }
}`

const SwaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <meta name="description" content="SwaggerUI for Calculator REST API" />
  <title>Calculator API - Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.18.2/swagger-ui.css" />
  <style>
    body { margin: 0; background: #fafafa; }
    .swagger-ui .topbar { background-color: #1e293b; }
    .swagger-ui .topbar .download-url-wrapper .download-url-button { background: #3b82f6; }
  </style>
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.18.2/swagger-ui-bundle.js" charset="UTF-8"></script>
<script src="https://unpkg.com/swagger-ui-dist@5.18.2/swagger-ui-standalone-preset.js" charset="UTF-8"></script>
<script>
window.onload = function() {
  window.ui = SwaggerUIBundle({
    url: "/swagger.json",
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });
};
</script>
</body>
</html>`

// SwaggerUI serves the official interactive Swagger UI web interface.
func (h *CalculatorHandler) SwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(SwaggerUIHTML))
}

// SwaggerJSON serves the raw OpenAPI 3.0 specification in JSON format.
func (h *CalculatorHandler) SwaggerJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(SwaggerSpecJSON))
}
