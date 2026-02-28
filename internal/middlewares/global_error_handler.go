package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"remy/internal/domainErrors"
	infraErrors "remy/internal/infrastructure/errors"
	"remy/internal/logging"
	"remy/internal/response"
)

var defaultValidationMessage = "Invalid value."

var validatorTagToMessage = map[string]func(fieldErr validator.FieldError) string{
	"required": func(fieldErr validator.FieldError) string {
		return "This field is required."
	},
	"min": func(fieldErr validator.FieldError) string {
		param := fieldErr.Param()
		return "Value must be at least " + param + " characters long."
	},
	"max": func(fieldErr validator.FieldError) string {
		param := fieldErr.Param()
		return "Value must be at most " + param + " characters long."
	},
	"oneof": func(fieldErr validator.FieldError) string {
		param := fieldErr.Param()
		options := strings.Split(param, " ")
		return "Value must be one of the following: " + strings.Join(options, ", ") + "."
	},
}

type ErrorMeta struct {
	Status int
	Title  string
}

var infraErrorMapping = map[string]ErrorMeta{
	"infra.version_conflict": {
		Status: http.StatusConflict,
		Title:  "Version Conflict",
	},
}

var domainErrorMapping = map[string]ErrorMeta{
	"srs.invalid_repetitions": {
		Status: http.StatusBadRequest,
		Title:  "Invalid Repetitions",
	},
	"srs.invalid_interval": {
		Status: http.StatusBadRequest,
		Title:  "Invalid Interval",
	},
	"srs.invalid_ease_factor": {
		Status: http.StatusBadRequest,
		Title:  "Invalid Ease Factor",
	},
	"srs.invalid_quality": {
		Status: http.StatusBadRequest,
		Title:  "Invalid Quality",
	},
	"general.not_found": {
		Status: http.StatusNotFound,
		Title:  "Resource Not Found",
	},
}

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			var apiErrors []*response.APIError

			for _, err := range c.Errors {
				var domainErr *domainErrors.DomainError
				var validationErrs validator.ValidationErrors
				var queryValidationErr *infraErrors.QueryValidationError
				var infraErr *infraErrors.InfrastructureError

				if errors.As(err.Err, &domainErr) {
					apiErrors = append(apiErrors, mapDomainErrorToAPIError(domainErr))
				} else if errors.As(err.Err, &queryValidationErr) {
					apiErrors = append(apiErrors, mapQueryValidationErrorToAPIError(queryValidationErr)...)
				} else if errors.As(err.Err, &infraErr) {
					apiErrors = append(apiErrors, mapInfraErrorToAPIError(infraErr))
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
	return envelopeDomainError(domainErr)
}

func mapInfraErrorToAPIError(infraErr *infraErrors.InfrastructureError) *response.APIError {
	meta, exists := infraErrorMapping[infraErr.Code]
	if !exists {
		return defaultError()
	}

	return &response.APIError{
		Status: meta.Status,
		Code:   infraErr.Code,
		Title:  meta.Title,
		Detail: strings.ToUpper(infraErr.Message[:1]) + infraErr.Message[1:],
	}
}

func mapValidationErrorsToAPIErrors(validationErrs validator.ValidationErrors) []*response.APIError {
	var apiErrors []*response.APIError

	for _, fieldErr := range validationErrs {
		apiErrors = append(apiErrors, &response.APIError{
			Status: http.StatusBadRequest,
			Code:   "validation_error",
			Title:  fmt.Sprintf("Validation failed for field '%s'", fieldErr.Field()),
			Detail: mapValidationTagToMessage(fieldErr),
			Source: &response.Source{
				Pointer: buildPointerFromNamespace(fieldErr.Namespace()),
			},
		})
	}

	return apiErrors
}

func mapQueryValidationErrorToAPIError(queryValidationErr *infraErrors.QueryValidationError) []*response.APIError {
	var apiErrors []*response.APIError

	for _, fieldErr := range queryValidationErr.OriginalError {
		apiErrors = append(apiErrors, &response.APIError{
			Status: http.StatusBadRequest,
			Code:   "validation_error",
			Title:  fmt.Sprintf("Validation failed for query parameter '%s'", fieldErr.Field()),
			Detail: mapValidationTagToMessage(fieldErr),
			Source: &response.Source{
				Parameter: fieldErr.Field(),
			},
		})
	}

	return apiErrors
}

func mapValidationTagToMessage(fieldErr validator.FieldError) string {
	if msgFunc, exists := validatorTagToMessage[fieldErr.Tag()]; exists {
		return msgFunc(fieldErr)
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

func envelopeDomainError(domainErr *domainErrors.DomainError) *response.APIError {
	meta, exists := domainErrorMapping[domainErr.Code]

	if !exists {
		return defaultError()
	}

	return &response.APIError{
		Status: meta.Status,
		Code:   domainErr.Code,
		Title:  meta.Title,
		Detail: strings.ToUpper(domainErr.Message[:1]) + domainErr.Message[1:],
	}
}
