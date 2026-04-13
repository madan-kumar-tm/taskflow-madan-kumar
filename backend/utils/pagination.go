package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Limit  int
	Offset int
}

func ParsePagination(c *gin.Context) (Pagination, error) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return Pagination{}, fmt.Errorf("page must be a positive integer")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		return Pagination{}, fmt.Errorf("limit must be between 1 and 100")
	}

	offset := (page - 1) * limit
	return Pagination{Limit: limit, Offset: offset}, nil
}
