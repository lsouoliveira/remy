package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"remy/internal/appErrors"
	"remy/internal/logging"
	"remy/internal/response"
)

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			var apiErrors []*response.APIError

			for _, err := range c.Errors {
				var appErr *appErrors.AppError

				if errors.As(err.Err, &appErr) {
					apiErrors = append(apiErrors, mapAppErrorToAPIError(appErr))
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
			{
				Status: http.StatusInternalServerError,
				Code:   "internal_server_error",
				Title:  "Internal Server Error",
				Detail: "An unexpected error occurred. Please try again later.",
			},
		},
	}
}

func mapAppErrorToAPIError(appErr *appErrors.AppError) *response.APIError {
	return &response.APIError{
		Status: appErr.Status,
		Code:   appErr.Code,
		Title:  http.StatusText(appErr.Status),
		Detail: appErr.Message,
	}
}
