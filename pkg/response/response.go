package response

import (
	"github.com/gin-gonic/gin"
)

type Meta struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Page       *int   `json:"page,omitempty"`
	Limit      *int   `json:"limit,omitempty"`
	TotalRows  *int   `json:"total_rows,omitempty"`
	TotalPages *int   `json:"total_pages,omitempty"`
}

type Response struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

// SendSuccess sends a standard success response
func SendSuccess(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Data: data,
		Meta: Meta{
			Code:    code,
			Message: message,
		},
	})
}

// SendPagination sends a paginated success response
func SendPagination(c *gin.Context, code int, message string, data interface{},
	page, limit, totalRows, totalPages int) {

	c.JSON(code, Response{
		Data: data,
		Meta: Meta{
			Code:       code,
			Message:    message,
			Page:       &page,
			Limit:      &limit,
			TotalRows:  &totalRows,
			TotalPages: &totalPages,
		},
	})
}

// SendError sends an error response
func SendError(c *gin.Context, code int, message string, err interface{}) {
	c.JSON(code, Response{
		Data: err,
		Meta: Meta{
			Code:    code,
			Message: message,
		},
	})
}
