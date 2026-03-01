package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetPageParam(c *gin.Context) int {
	pageStr := c.Query("page")
	return max(ParseInt(pageStr, 1), 1)
}

func GetPageSizeParam(c *gin.Context) int {
	pageSizeStr := c.Query("page_size")
	pageSize := ParseInt(pageSizeStr, 10)

	if pageSize < 1 {
		pageSize = 10
	}

	return pageSize
}

func ParseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return value
}

func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
