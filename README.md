# JARO - Agentic Framework

**JARO (JAROBOTAI)** is a sophisticated agentic framework built with Go, following strict Hexagonal Architecture principles.

## 🏗️ Architecture

This project implements **Hexagonal Architecture (Ports & Adapters)**:

```
├── cmd/jaro/              # Main application entry point
├── internal/
│   ├── core/
│   │   ├── domain/        # Pure business logic & models
│   │   ├── ports/         # Interface definitions (contracts)
│   │   └── services/      # Business logic implementations
│   └── adapters/
│       ├── memory/        # In-memory implementations (testing)
│       └── primary/http/  # HTTP REST API adapter
```

## 🚀 Features

- ✅ **Task Management** - Create and track agentic tasks
- ✅ **REST API** - HTTP endpoints for task orchestration
- ✅ **Audit Logging** - Complete event tracking
- ✅ **In-Memory Storage** - No external dependencies for development
- ✅ **Extensible Design** - Easy to swap adapters

## 📡 API Endpoints

### Health Check
```bash
GET /health
```

### Create Task
```bash
POST /tasks
Content-Type: application/json

{
  "input": "Find me a two-room apartment in Vracar under 800 EUR",
  "user_id": "user-12345"
}
```

**Response (201 Created):**
```json
{
  "task_id": "0931282d-6164-4be5-be44-457e5ffd1312",
  "status": "NEW",
  "created_at": "2026-02-12T01:28:35+01:00",
  "user_id": "user-12345",
  "input": "Find me a two-room apartment in Vracar under 800 EUR"
}
```

### Get Task Status
```bash
GET /tasks/:id
```

**Response (200 OK):**
```json
{
  "id": "0931282d-6164-4be5-be44-457e5ffd1312",
  "status": "NEW",
  "created_at": "2026-02-12T01:28:35+01:00",
  "updated_at": "2026-02-12T01:28:35+01:00",
  "input": "Find me a two-room apartment in Vracar under 800 EUR",
  "user_id": "user-12345",
  "target_agent": "CORE",
  "artifacts": {},
  "metadata": {}
}
```

## 🛠️ Development

### Prerequisites
- Go 1.25+

### Build
```bash
go build -o jaro.exe ./cmd/jaro
```

### Run
```bash
./jaro.exe
```

Server starts on `http://localhost:8080`

### Test API (PowerShell)
```powershell
# Health check
Invoke-RestMethod -Uri http://localhost:8080/health

# Create task
$body = @{input='Find apartments'; user_id='user-123'} | ConvertTo-Json
Invoke-RestMethod -Uri http://localhost:8080/tasks -Method POST -Body $body -ContentType 'application/json'

# Get task status
Invoke-RestMethod -Uri http://localhost:8080/tasks/TASK_ID
```

## 📦 Components

### Domain Layer
- `Task` - Core task entity with status tracking
- `Plan` - Execution plan with steps
- `Step` - Individual action in a plan
- `AuditEvent` - Event logging for compliance

### Ports Layer
- `Orchestrator` - Primary port for task management
- `TaskRepository` - Task persistence interface
- `AuditRepository` - Audit log interface
- `Planner` - Plan generation interface
- `Executor` - Step execution interface

### Services Layer
- `OrchestratorService` - Core orchestration logic

### Adapters Layer
- **Memory** - In-memory implementations for testing
- **HTTP** - REST API adapter (Gin framework)

## 🔒 Security & Open Core

This project follows "Open Core" strategy:
- ✅ Core domain models are public-safe
- ✅ No secrets hardcoded
- ✅ Configuration-driven design
- ✅ Adapters are swappable

## 📝 Status

**Phase 1: Complete** ✅
- Core domain models
- Ports & interfaces
- Service layer
- In-memory adapters
- HTTP REST API

**Phase 2: Planned**
- Real LLM integration (OpenAI/Anthropic)
- PostgreSQL persistence
- Plan execution loop
- Approval workflow
- Metrics & FinOps tracking

## 🧪 Testing

All components tested:
- ✅ Task creation via API
- ✅ Task status retrieval
- ✅ 404 error handling
- ✅ Audit logging
- ✅ Thread-safe in-memory storage

## 📄 License

[To be determined]

## 👥 Contributing

[Contribution guidelines to be added]
