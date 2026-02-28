package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"remy/internal/domainErrors"
	"remy/internal/logging"
	"remy/internal/response"
)

var defaultValidationMessage = "Invalid value."

var validatorTagToMessage = map[string]string{
	"required": "This field is required.",
	"min":      "Value is too short.",
	"max":      "Value is too long.",
}

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			var apiErrors []*response.APIError

			for _, err := range c.Errors {
				var domainErr *domainErrors.DomainError
				var validationErrs validator.ValidationErrors

				if errors.As(err.Err, &domainErr) {
					apiErrors = append(apiErrors, mapDomainErrorToAPIError(domainErr))
				} else if errors.As(err.Err, &validationErrs) {
					apiErrors = append(apiErrors, mapValidationErrorsToAPIErrors(validationErrs)...)
				} else {
					logging.Logger.WithFields(logrus.Fields{
						"error": err.Err.Error(),
					}).Error("unexpected error occurred")

					c.JSON(http.StatusInternalServerError, defaultErrorResponse())
					return
				}
			}

			c.JSON(apiErrors[0].Status, response.APIResponse{
				Errors: apiErrors,
			})
		}
	}
}

func defaultErrorResponse() response.APIResponse {
	return response.APIResponse{
		Errors: []*response.APIError{
			defaultError(),
		},
	}
}

func defaultError() *response.APIError {
	return &response.APIError{
		Status: http.StatusInternalServerError,
		Code:   "internal_server_error",
		Title:  "Internal Server Error",
		Detail: "An unexpected error occurred. Please try again later.",
	}
}

func mapDomainErrorToAPIError(domainErr *domainErrors.DomainError) *response.APIError {
	return defaultError()
}

func mapValidationErrorsToAPIErrors(validationErrs validator.ValidationErrors) []*response.APIError {
	var apiErrors []*response.APIError

	for _, fieldErr := range validationErrs {
		apiErrors = append(apiErrors, &response.APIError{
			Status: http.StatusBadRequest,
			Code:   "validation_error",
			Title:  "Validation Error",
			Detail: mapValidationTagToMessage(fieldErr.Tag()),
			Source: &response.Source{
				Pointer: buildPointerFromNamespace(fieldErr.Namespace()),
			},
		})
	}

	return apiErrors
}

func mapValidationTagToMessage(tag string) string {
	if msg, exists := validatorTagToMessage[tag]; exists {
		return msg
	}

	return defaultValidationMessage
}

func buildPointerFromNamespace(namespace string) string {
	parts := strings.Split(namespace, ".")

	if len(parts) > 1 {
		parts = parts[1:]
	}

	var pointer strings.Builder

	for _, part := range parts {
		pointer.WriteString("/")
		pointer.WriteString(strings.ToLower(part))
	}

	return pointer.String()
}
