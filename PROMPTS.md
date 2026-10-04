# Prompts & AI Engineering Interaction Log

As required in the Sezzle technical evaluation instructions:
> *"Share any prompts that you used in your work"*

This document records the exact iterative prompts and instructions provided by the engineer (**Valeria Tocarruncho Mosquera**) to guide the AI pair programmer through the design, implementation, debugging, UI/UX redesign, and documentation of the full-stack calculator application.

---

## Stage 1: Project Kickoff & Step-by-Step Architecture Planning

### Engineer Prompt:
> *"We will build the following project based on the Sezzle evaluation requirements:*
> - *Objective: Full-stack calculator application with a React frontend and a Go backend microservice.*
> - *Operations: Basic (Addition, Subtraction, Multiplication, Division) and Advanced (Exponentiation, Square Root, Percentage).*
> - *Frontend: React with TypeScript, intuitive UI, input validation, error handling, responsive mobile support.*
> - *Backend: Go REST API microservice, validate inputs, handle edge cases (division by zero, invalid data), return JSON.*
> - *Non-functional: Clean idiomatic code, unit tests covering key functionality with coverage reports, documentation (setup, API examples, design decisions), Dockerfile/Docker Compose.*
> *Let's build this step-by-step."*

### Resulting Work:
- Initialized Git repository and domain-driven Go project structure (`cmd/api`, `internal/calculator`, `internal/handler`, `internal/middleware`).
- Initialized React 19 + TypeScript frontend with Vite.
- Implemented core arithmetic logic with IEEE-754 precision sanitization and 100% statement test coverage in Go.

---

## Stage 2: Docker Environment & Version Compatibility Debugging

### Engineer Prompt:
> *"I encountered an error when trying to run the project with Docker Compose:*
> `go: go.mod requires go >= 1.27.0 (running go 1.24.13; GOTOOLCHAIN=local)`
> `process '/bin/sh -c go mod download' did not complete successfully: exit code: 1`
> *Diagnose the cause of this mismatch and fix the configuration so the full stack builds cleanly in Docker."*

### Resulting Work:
- Adjusted `go.mod` directive to `go 1.22`, which resolved the toolchain mismatch with the `golang:1.24-alpine` builder image.
- Cleared background port conflicts on port 8080 and validated successful container compilation for both backend and frontend.

---

## Stage 3: Functional Quality Assurance & Display State Correction

### Engineer Prompt:
> *"During functional testing, I noticed the following bug in the calculator state:*
> *When I want to calculate an addition such as 55 + 55, when I click on '+', the first 55 does not disappear from the main display, and entering the second 55 causes it to append into 5555.*
> *Fix the state transition so that pressing an operator moves the first operand to the equation header, clears/resets the main display to 0, and cleanly accepts the second number."*

### Resulting Work:
- Refactored `useCalculator.ts` to immediately reset the main numeric display to `'0'` upon selecting an operator while storing the operand in the equation header.
- Added visual highlighting (`active-op`) to the selected operator button.
- Added an automated integration test in `src/App.test.tsx` verifying the exact `55 + 55 = 110` flow without any numeric duplication.

---

## Stage 4: UI/UX Redesign Guided by Visual Reference

### Engineer Prompt:
> *"Now we are going to redesign the user interface. Please guide yourself by the uploaded image (Android Google Calculator / Material You):*
> - *Use circular buttons.*
> - *Use soft, tonal pastel colors (mint green for AC, soft pastel blue for operators, warm light gray for numbers, soft dusty rose/pink for equals).*
> - *Implement the prominent top display layout with a large expression line and secondary result line.*
> - *Include the secondary mini function bar at the top with shortcuts (√, ^, %, ±)."*

### Resulting Work:
- Transformed the keypad from rectangular dark-mode tiles into a clean, circular Material You button grid matching the visual reference.
- Applied exact tonal pastel color variables in `src/index.css`.
- Updated the display hierarchy so expressions appear prominent and bold on top, with results and live values beneath them.

---

## Stage 5: Title & Branding Refinement

### Engineer Prompt:
> *"Update the header branding: I don't want the title to be 'FullStack Calc', rename it simply to 'Calculator'."*

### Resulting Work:
- Updated the header component and HTML title to "Calculator", providing a clean and minimalist look.
- Synchronized all unit test assertions with the new title.

---

## Stage 6: REST API Verification & Troubleshooting

### Engineer Prompt:
> *"Provide comprehensive examples of API calls for the REST endpoints. Also, when I run `curl -X GET http://localhost:8080/api/v1/health` in Windows PowerShell, it throws a parameter error."*

### Resulting Work:
- Diagnosed the Windows PowerShell alias collision (`curl -> Invoke-WebRequest`) and documented alternative executions (`curl.exe`, native `Invoke-RestMethod`, and direct browser access).
- Documented full request and response payloads with status codes for all arithmetic operations.

---

## Stage 7: Standardizing API Documentation with Swagger UI

### Engineer Prompt:
> *"Can we implement the API documentation using Swagger? I want interactive OpenAPI documentation rather than custom static text."*

### Resulting Work:
- Embedded the official **OpenAPI 3.0 specification** (`/swagger.json`) directly into the Go microservice.
- Integrated **Swagger UI** (`/swagger`), providing an interactive browser explorer with "Try it out" capabilities to execute live API calls against the backend.
- Linked the frontend "API Online" badge directly to Swagger UI.

---

## Stage 8: Documentation Consolidation & Delivery Verification

### Engineer Prompt:
> *"Ensure that Swagger is the single unified API documentation entrypoint so the project remains clean without redundant documentation files. Review the overall project against all Sezzle requirements and confirm when it is ready for final submission."*

### Resulting Work:
- Consolidated all API entrypoints (`/`, `/api`, `/api/v1`) to redirect directly to Swagger UI.
- Verified test suites: 18 passing tests in frontend and ~96% statement coverage in Go backend.
- Confirmed all functional, non-functional, and deliverable constraints are satisfied.
