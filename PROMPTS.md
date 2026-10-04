# Prompts & AI Interaction Log

As requested in the Sezzle assessment instructions:
> *"Share any prompts that you used in your work"*

This document documents the prompt workflow, architectural guidance, and iterative instructions used to build the full-stack calculator application.

---

## 1. Project Initialization & Architecture Strategy Prompt

```text
Objective: Build a full-stack calculator application with a React frontend (TypeScript) and a Go backend microservice.

Key requirements:
1. Operations: Basic (Add, Subtract, Multiply, Divide) and Advanced (Exponentiation, Square Root, Percentage).
2. Clean, idiomatic code in both Go and React/TypeScript.
3. Unit tests with high statement coverage and edge case handling (division by zero, negative square root, floating-point precision, overflow).
4. Responsive modern UI with input validation, error handling, session calculation history, and physical keyboard support.
5. Dockerfile and docker-compose for one-command deployment.
6. Documentation with setup, API usage examples (cURL), design rationale, and test reporting.

Plan:
- Go backend: Domain-driven structure (cmd/api, internal/calculator, internal/handler, internal/middleware).
- React frontend: Vite + React 19 + TypeScript + Vitest + React Testing Library.
- Multi-stage Docker builds.
```

---

## 2. Go Backend Design Prompt

```text
Write a robust, idiomatic Go REST API microservice:
- Service layer (internal/calculator):
  - Add, Subtract, Multiply, Divide, Power, SquareRoot, Percentage.
  - Sane floating-point precision handling (eliminating IEEE-754 binary floating micro-artifacts like 0.1 + 0.2 = 0.30000000000000004).
  - Explicit domain errors (ErrDivisionByZero, ErrNegativeSquareRoot, ErrInvalidOperation, ErrMissingOperand, ErrInvalidResult).
- Transport layer (internal/handler):
  - Unified endpoint: POST /api/v1/calculate with JSON payload {"operation": "...", "a": ..., "b": ...}
  - Convenience REST endpoints: POST /api/v1/operations/{op}
  - Health check endpoint: GET /api/v1/health
  - Standard JSON response format and structured error responses with HTTP status codes (200, 400, 422).
- Middlewares:
  - CORS with full preflight support.
  - Request logging with latency and status codes.
- Unit tests:
  - Full table-driven tests for service, handlers, and middlewares targeting >80-100% test coverage.
```

---

## 3. Frontend Architecture & React Hook Prompt

```text
Build a React + TypeScript frontend with Vite:
- State management hook (useCalculator):
  - Encapsulate the calculator state machine: input digits, decimals, toggle sign (+/-), delete backspace, AC clear.
  - Chaining calculations (e.g., 5 + 5 + resolves the first addition and keeps the next pending).
  - Unary execution (√, %) with instantaneous backend API evaluation.
  - Calculation history list with localStorage persistence and one-click reloading into the active display.
  - Proper closure ref management to avoid stale state during rapid keystrokes.
- UI Design:
  - Modern fintech aesthetic with glassmorphism display, responsive keypad, and status indicator.
  - Physical keyboard event listeners (0-9, +, -, *, /, ^, %, Enter, Backspace, Esc).
  - Vitest + React Testing Library integration tests covering all critical user flows.
```

---

## 4. Packaging & Docker Deployment Prompt

```text
Create containerization and documentation assets:
- Backend: Multi-stage Dockerfile using golang:1.24-alpine as builder and alpine:3.20 as minimal runner.
- Frontend: Multi-stage Dockerfile building static assets with node:22-alpine and serving with nginx:1.27-alpine.
- docker-compose.yml connecting both services with healthcheck readiness.
- Comprehensive README.md with setup instructions, API curl examples, architectural decisions, and test instructions.
```
