package helpers

import (
	"fmt"
	"time" // Optional: for timestamping errors

	"github.com/gofiber/fiber/v2"
)

// ServiceError defines a structured error for services
type ServiceError struct {
	ServiceName string    `json:"serviceName"` // Name of the service where the error occurred
	Err         error     `json:"-"`           // The original error, might not be suitable for direct JSON response
	Message     string    `json:"message"`     // A user-friendly or log-friendly message
	Timestamp   time.Time `json:"timestamp"`   // Optional: when the error occurred
	// You can add more fields like:
	// Operation   string    `json:"operation"`   // e.g., "GetAccess", "UpdateOrCreateAccess"
	// StatusCode  int       `json:"statusCode"`  // If you want to suggest an HTTP status
}

// Error makes ServiceError satisfy the error interface
func (se *ServiceError) Error() string {
	return fmt.Sprintf("service: %s, message: %s, original_error: %v", se.ServiceName, se.Message, se.Err)
}

// NewServiceError creates a new ServiceError
func NewServiceError(serviceName string, message string, originalErr error) *ServiceError {
	return &ServiceError{
		ServiceName: serviceName,
		Message:     message,
		Err:         originalErr,
		Timestamp:   time.Now(),
	}
}

func InternalError(c *fiber.Ctx, err error) error {
	// Log the detailed error internally
	// In a real application, replace this with a structured logger (e.g., logrus, zap)
	fmt.Printf("Internal error occurred: %s\n", err.Error())

	// Check if it's a ServiceError to provide more context,
	// otherwise, a generic message.
	if se, ok := err.(*ServiceError); ok {
		// For client-facing errors, you might not want to expose se.Message directly
		// if it contains internal details. Consider a more generic message or
		// a specific client-friendly message field in ServiceError.
		return c.
			Status(fiber.StatusInternalServerError). // Or se.StatusCode if you add it
			JSON(fiber.Map{"error": "Internal Server Error", "service_context": se.ServiceName, "timestamp": se.Timestamp})
	}

	return c.
		Status(fiber.StatusInternalServerError).
		JSON(fiber.Map{"error": "Internal Server Error"})
}
