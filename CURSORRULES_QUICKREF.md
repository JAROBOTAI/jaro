# JARO .cursorrules - Quick Reference

## 🚨 FORBIDDEN in `internal/core/*`

### Direct Stdlib Calls:
```go
❌ time.Now()           → ✅ clock.Now()
❌ uuid.New()           → ✅ idGen.Generate()
❌ fmt.Printf()         → ✅ logger.Info/Warn/Error()
❌ log.Printf()         → ✅ logger.Info/Warn/Error()
```

### SDK Imports:
```go
❌ github.com/openai/openai-go
❌ cloud.google.com/*
❌ github.com/aws/*
❌ github.com/lib/pq
❌ gorm.io/gorm
❌ go.mongodb.org/mongo-driver
❌ github.com/gin-gonic/gin
```

### Adapter Type Leaks:
```go
❌ func DoWork(c *gin.Context) error
❌ func Process(db *sql.DB) error
❌ func Handle(req *http.Request) error

✅ func DoWork(ctx context.Context, data WorkData) error
```

---

## ✅ REQUIRED Patterns

### Service Constructor:
```go
func NewMyService(
    // Business dependencies
    repo ports.Repository,
    
    // Cross-cutting concerns
    clock ports.Clock,
    idGen ports.IDGenerator,
    logger ports.Logger,
) ports.MyService {
    return &MyServiceImpl{...}
}
```

### Timestamp Generation:
```go
// In Service:
now := s.clock.Now()

// In Adapter:
now := time.Now()  // OK in adapters
```

### ID Generation:
```go
// In Service:
id := s.idGen.Generate()

// In Adapter:
id := uuid.New().String()  // OK in adapters
```

### Logging:
```go
// In Service:
s.logger.Info("task created", map[string]interface{}{
    "task_id": taskID,
    "user_id": userID,
})

s.logger.Error("failed to save", err, map[string]interface{}{
    "task_id": taskID,
})

// In Adapter:
fmt.Println("[INFO] task created")  // OK for simple adapters
```

---

## 📋 Definition of Done Checklist

Before claiming task complete:
- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes (if tests exist)
- [ ] No linter errors (`ReadLints`)
- [ ] All exported types have GoDoc (Purpose, Inputs, Outputs)
- [ ] No adapter-specific types in Core signatures
- [ ] Audit/Log hooks for all I/O operations
- [ ] No magic numbers (timeouts, limits in Config)
- [ ] No hardcoded secrets

---

## 🔒 Security Checklist

- [ ] No API keys, passwords in code
- [ ] Input validation at adapter boundaries
- [ ] Max request body size enforced
- [ ] MIME types whitelisted
- [ ] HTML sanitized before Core

---

## 📦 Domain Model Stability

### Can Add:
```go
✅ NewField *string `json:"new_field,omitempty"`
   // Optional field, backward compatible
```

### Cannot Change Without v2:
```go
❌ Renaming: UserID → UserIdentifier
❌ Removing: Status field
❌ Type change: CreatedAt time.Time → int64
```

### Migration Pattern:
```go
// Old: internal/core/domain/task.go
// New: internal/core/domain/v2/task.go
// Adapter: maps v1 ↔ v2 at boundary
```

---

## ⚠️ AI Warning Triggers

AI will STOP and WARN if you write:

### Critical:
1. `time.Now()` in `internal/core/services/`
2. `uuid.New()` in `internal/core/services/`
3. SDK import in `internal/core/`
4. `*gin.Context` in service signature
5. "Google Sheets" mentioned in domain logic

### Smells:
6. Hardcoded secret: `apiKey := "sk-..."`
7. Magic number: `timeout := 30 * time.Second`
8. Missing audit: Task created without event
9. Breaking change without migration

---

## 📚 Example: Adding New Feature

### ❌ Wrong Way:
```go
// internal/core/services/task_service.go
import "github.com/openai/openai-go"  // VIOLATION

func (s *Service) Process() error {
    client := openai.NewClient("sk-...")  // VIOLATION
    time.Sleep(30 * time.Second)          // VIOLATION
    
    id := uuid.New().String()             // VIOLATION
    fmt.Println("Processing...")          // VIOLATION
    // ...
}
```

### ✅ Right Way:
```go
// internal/core/ports/llm.go
type LLMProvider interface {
    Generate(ctx context.Context, prompt string) (string, error)
}

// internal/core/services/task_service.go
func NewService(
    llm ports.LLMProvider,  // Depend on interface
    clock ports.Clock,
    idGen ports.IDGenerator,
    logger ports.Logger,
) ports.Service {
    return &ServiceImpl{...}
}

func (s *ServiceImpl) Process(ctx context.Context) error {
    id := s.idGen.Generate()
    s.logger.Info("processing", map[string]interface{}{"id": id})
    
    result, err := s.llm.Generate(ctx, "prompt")
    if err != nil {
        return fmt.Errorf("llm failed: %w", err)
    }
    // ...
}

// internal/adapters/llm/openai.go
import "github.com/openai/openai-go"  // OK in adapter

type OpenAIAdapter struct {
    client *openai.Client
}
```

---

## 🎯 Quick Commands

### Build & Test:
```bash
go build ./...              # Check compile
go test ./...               # Run tests
go build -o jaro.exe ./cmd/jaro  # Build binary
```

### Run & Test:
```bash
./jaro.exe                          # Start server
powershell -File test_api.ps1       # Run API tests
```

### Verify Compliance:
```bash
# Check for violations in services
grep -r "time\.Now\|uuid\.New" internal/core/services/

# Should return: nothing (0 results)
```

---

## 📖 Full Rules

See `.cursorrules` file for complete documentation.

---

**Last Updated:** 2026-02-12  
**Version:** 1.0 (Sharpened Edition)
