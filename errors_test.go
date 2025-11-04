package zhinao

import (
	"errors"
	"testing"
)

func TestAPIError(t *testing.T) {
	t.Run("error message", func(t *testing.T) {
		err := &APIError{
			StatusCode: 400,
			Type:       "invalid_request_error",
			Code:       "invalid_api_key",
			Message:    "Invalid API key provided",
		}

		expected := "[400] invalid_request_error: Invalid API key provided (code: invalid_api_key)"
		if err.Error() != expected {
			t.Errorf("Error() = %v, want %v", err.Error(), expected)
		}
	})

	t.Run("is retryable - 5xx errors", func(t *testing.T) {
		err := &APIError{StatusCode: 500}
		if !err.IsRetryable() {
			t.Error("500 errors should be retryable")
		}

		err = &APIError{StatusCode: 503}
		if !err.IsRetryable() {
			t.Error("503 errors should be retryable")
		}
	})

	t.Run("is retryable - 429 rate limit", func(t *testing.T) {
		err := &APIError{StatusCode: 429}
		if !err.IsRetryable() {
			t.Error("429 errors should be retryable")
		}
	})

	t.Run("is retryable - 408 timeout", func(t *testing.T) {
		err := &APIError{StatusCode: 408}
		if !err.IsRetryable() {
			t.Error("408 errors should be retryable")
		}
	})

	t.Run("is not retryable - 4xx errors", func(t *testing.T) {
		err := &APIError{StatusCode: 400}
		if err.IsRetryable() {
			t.Error("400 errors should not be retryable")
		}

		err = &APIError{StatusCode: 401}
		if err.IsRetryable() {
			t.Error("401 errors should not be retryable")
		}

		err = &APIError{StatusCode: 404}
		if err.IsRetryable() {
			t.Error("404 errors should not be retryable")
		}
	})
}

func TestRateLimitError(t *testing.T) {
	t.Run("error message with retry after", func(t *testing.T) {
		err := &RateLimitError{
			APIError: &APIError{
				StatusCode: 429,
				Type:       "rate_limit_error",
				Code:       "rate_limit",
				Message:    "Rate limit exceeded",
			},
			RetryAfter: 60,
		}

		expected := "[429] rate_limit_error: Rate limit exceeded (code: rate_limit) (retry after 60 seconds)"
		if err.Error() != expected {
			t.Errorf("Error() = %v, want %v", err.Error(), expected)
		}
	})

	t.Run("error message without retry after", func(t *testing.T) {
		err := &RateLimitError{
			APIError: &APIError{
				StatusCode: 429,
				Type:       "rate_limit_error",
				Code:       "rate_limit",
				Message:    "Rate limit exceeded",
			},
			RetryAfter: 0,
		}

		expected := "[429] rate_limit_error: Rate limit exceeded (code: rate_limit)"
		if err.Error() != expected {
			t.Errorf("Error() = %v, want %v", err.Error(), expected)
		}
	})
}

func TestValidationError(t *testing.T) {
	t.Run("error message", func(t *testing.T) {
		err := &ValidationError{
			Field:   "model",
			Message: "model is required",
		}

		expected := "validation error: model - model is required"
		if err.Error() != expected {
			t.Errorf("Error() = %v, want %v", err.Error(), expected)
		}
	})
}

func TestPredefinedErrors(t *testing.T) {
	t.Run("ErrMissingAPIKey", func(t *testing.T) {
		if ErrMissingAPIKey == nil {
			t.Error("ErrMissingAPIKey should not be nil")
		}
		if ErrMissingAPIKey.Error() == "" {
			t.Error("ErrMissingAPIKey should have error message")
		}
	})

	t.Run("ErrInvalidModel", func(t *testing.T) {
		if ErrInvalidModel == nil {
			t.Error("ErrInvalidModel should not be nil")
		}
	})

	t.Run("ErrEmptyMessages", func(t *testing.T) {
		if ErrEmptyMessages == nil {
			t.Error("ErrEmptyMessages should not be nil")
		}
	})

	t.Run("ErrStreamClosed", func(t *testing.T) {
		if ErrStreamClosed == nil {
			t.Error("ErrStreamClosed should not be nil")
		}
	})
}

func TestErrorUnwrap(t *testing.T) {
	t.Run("unwrap API error", func(t *testing.T) {
		apiErr := &APIError{
			StatusCode: 400,
			Message:    "Bad request",
		}

		var targetErr *APIError
		if !errors.As(apiErr, &targetErr) {
			t.Error("Should be able to unwrap to APIError")
		}
	})

	t.Run("unwrap rate limit error", func(t *testing.T) {
		rateLimitErr := &RateLimitError{
			APIError: &APIError{
				StatusCode: 429,
				Message:    "Rate limit",
			},
			RetryAfter: 60,
		}

		var targetErr *RateLimitError
		if !errors.As(rateLimitErr, &targetErr) {
			t.Error("Should be able to unwrap to RateLimitError")
		}
	})
}
