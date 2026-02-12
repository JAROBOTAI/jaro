package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/JAROBOTAI/jaro/internal/core/domain"
	"github.com/JAROBOTAI/jaro/internal/core/ports"
)

// TaskRepository is an in-memory implementation of the ports.TaskRepository interface.
// It stores tasks in a thread-safe map for local development and testing without a database.
// All data is lost when the application stops (non-persistent).
type TaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

// NewTaskRepository creates a new in-memory task repository.
// Purpose: Factory function for creating the in-memory task storage adapter.
// Inputs: None
// Outputs:
//   - ports.TaskRepository: Initialized repository ready for use
func NewTaskRepository() ports.TaskRepository {
	return &TaskRepository{
		tasks: make(map[string]*domain.Task),
	}
}

// SaveTask persists a task to the in-memory map (insert or update).
// Purpose: Stores or updates task state in memory with thread-safe access.
//          This implementation uses a write lock to ensure concurrent safety.
// Inputs:
//   - ctx: Context for cancellation and timeout control (unused in this implementation)
//   - task: The task to save (must have a valid ID)
// Outputs:
//   - error: Returns error if task is nil or has an empty ID
func (r *TaskRepository) SaveTask(ctx context.Context, task *domain.Task) error {
	if task == nil {
		return fmt.Errorf("task cannot be nil")
	}
	if task.ID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Deep copy to avoid external mutations
	taskCopy := *task
	r.tasks[task.ID] = &taskCopy

	return nil
}

// GetTask retrieves a task by its unique identifier from memory.
// Purpose: Loads task state from in-memory storage with thread-safe read access.
// Inputs:
//   - ctx: Context for cancellation and timeout control (unused in this implementation)
//   - id: Unique identifier of the task to retrieve
// Outputs:
//   - *domain.Task: The retrieved task with all fields populated
//   - error: Returns error if task is not found or id is empty
func (r *TaskRepository) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	if id == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", id)
	}

	// Deep copy to avoid external mutations
	taskCopy := *task
	return &taskCopy, nil
}
