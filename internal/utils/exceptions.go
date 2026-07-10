package utils

import (
	"fmt"
)

// DeploymentError represents an expected deployment failure caused by user/application errors.
// These are not Vessel bugs and should not necessarily cause panic or critical log levels.
// Examples: detection failures, missing Dockerfiles, invalid configs, etc.
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

// NewDeploymentError creates a new DeploymentError.
func NewDeploymentError(message string, err error) *DeploymentError {
	return &DeploymentError{
		Message: message,
		Err:     err,
	}
}

// RateLimitError represents an error when a rate limit is exceeded.
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

// ProcessError represents an error during external process execution.
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

// NonReportableError represents an error that shouldn't trigger external error tracking (like Sentry).
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
