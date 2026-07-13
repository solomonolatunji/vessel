package utils

import (
	"fmt"
)

type DeploymentError struct {
	Message string
	Err     error
}

func (e *DeploymentError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *DeploymentError) Unwrap() error {
	return e.Err
}

func NewDeploymentError(message string, err error) *DeploymentError {
	return &DeploymentError{
		Message: message,
		Err:     err,
	}
}

type RateLimitError struct {
	Message    string
	RetryAfter int // in seconds
}

func (e *RateLimitError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("%s (retry after %ds)", e.Message, e.RetryAfter)
	}
	return e.Message
}

type ProcessError struct {
	Command  string
	ExitCode int
	Stderr   string
	Err      error
}

func (e *ProcessError) Error() string {
	return fmt.Sprintf("process '%s' failed with exit code %d: %s", e.Command, e.ExitCode, e.Stderr)
}

func (e *ProcessError) Unwrap() error {
	return e.Err
}

type NonReportableError struct {
	Message string
	Err     error
}

func (e *NonReportableError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *NonReportableError) Unwrap() error {
	return e.Err
}

type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s not found: %s", e.Resource, e.ID)
	}
	return fmt.Sprintf("%s not found", e.Resource)
}

func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{Resource: resource, ID: id}
}

type EngineError struct {
	Operation string
	Err       error
}

func (e *EngineError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("docker engine operation '%s' failed: %v", e.Operation, e.Err)
	}
	return fmt.Sprintf("docker engine operation '%s' failed", e.Operation)
}

func (e *EngineError) Unwrap() error {
	return e.Err
}

func NewEngineError(operation string, err error) *EngineError {
	return &EngineError{
		Operation: operation,
		Err:       err,
	}
}
