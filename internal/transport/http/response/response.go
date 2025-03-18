package response

import (
	"backoffice/pkg/validator"
	"backoffice/utils"
	"encoding/csv"
	"fmt"
	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Response struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Meta    interface{} `json:"meta"`
	Data    interface{} `json:"data"`
}

func new(status int, meta interface{}, data interface{}) *Response {
	success := false
	if status >= 200 && status <= 299 {
		success = true
	}

	response := &Response{
		Status:  status,
		Success: success,
		Meta:    meta,
		Data:    data,
	}

	if response.Data == nil {
		response.Data = http.StatusText(status)
	}

	if v, ok := data.(error); ok {
		response.Data = v.Error()
	}

	if v, ok := data.([]error); ok {
		response.Data = lo.Map(v, func(item error, index int) string {
			return item.Error()
		})
	}

	return response
}

func OK(ctx *gin.Context, data interface{}, meta interface{}) {
	r := new(http.StatusOK, meta, data)
	ctx.JSON(r.Status, r)
}

func NoContent(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}

func CSV[T any](ctx *gin.Context, itemName string, items []T) {
	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%v.csv"`, itemName))

	writer := csv.NewWriter(ctx.Writer)
	if err := writer.WriteAll(utils.ExtractTable(items, "csv")); err != nil {
		BadRequest(ctx, err, nil)
	}
}

func XLSX(ctx *gin.Context, file []byte, filename string) {
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%v", filename))
	ctx.Data(http.StatusOK, "application/octet-stream", file)
}

func XLSXFile(c *gin.Context, file *excelize.File, fileName string) {
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	if err := file.Write(c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "Error writing file: %s", err.Error())
	}

	c.Status(http.StatusOK)
}

func BadRequest(ctx *gin.Context, data interface{}, meta interface{}) {
	zap.S().Warn(data)
	r := new(http.StatusBadRequest, meta, data)

	ctx.AbortWithStatusJSON(r.Status, r)
}

func Unauthorized(ctx *gin.Context, data interface{}, meta interface{}) {
	zap.S().Warn(data)
	r := new(http.StatusUnauthorized, meta, data)
	ctx.AbortWithStatusJSON(r.Status, r)
}

func Forbidden(ctx *gin.Context, data interface{}, meta interface{}) {
	zap.S().Warn(data)
	r := new(http.StatusForbidden, meta, data)
	ctx.AbortWithStatusJSON(r.Status, r)
}

func NotFound(ctx *gin.Context, data interface{}, meta interface{}) {
	zap.S().Warn(data)
	r := new(http.StatusNotFound, meta, data)
	ctx.AbortWithStatusJSON(r.Status, r)
}

func ServerError(ctx *gin.Context, data interface{}, meta interface{}) {
	zap.S().Error(data)
	r := new(http.StatusInternalServerError, meta, data)
	ctx.AbortWithStatusJSON(r.Status, r)
}

func Code(ctx *gin.Context, code int, data interface{}, meta interface{}) {
	zap.S().Warn(data)
	zap.S().Warn("code: " + strconv.Itoa(code))
	r := new(code, meta, data)
	ctx.AbortWithStatusJSON(r.Status, r)
}

func ValidationFailed(ctx *gin.Context, err error) {
	data := make([]string, 0)

	for _, taggedError := range validator.CheckValidationErrors(err) {
		e := taggedError.Err
		data = append(data, e.Error())
	}

	r := new(http.StatusUnprocessableEntity, nil, data)
	ctx.AbortWithStatusJSON(r.Status, r)
}

func Conflict(ctx *gin.Context, data interface{}, meta interface{}) {
	zap.S().Error(data)
	r := new(http.StatusConflict, meta, data)
	ctx.AbortWithStatusJSON(r.Status, r)
}
