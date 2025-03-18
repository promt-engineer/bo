package handlers

import (
	"backoffice/internal/entities"
	"backoffice/internal/services"
	"backoffice/internal/transport/http/response"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type fileHandler struct {
	fileService *services.FileDownloadingService
	upgrader    *websocket.Upgrader
}

func NewFileHandler(
	fileService *services.FileDownloadingService,
) *fileHandler {
	return &fileHandler{
		fileService: fileService,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024 * 32,
			WriteBufferSize: 1024 * 32,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *fileHandler) Register(router *gin.RouterGroup) {
	files := router.Group("files")

	files.GET("", h.files)
	files.GET(":id", h.download)
	files.GET("ws/:id", h.subscribeHandler)
}

// @Summary Get list of files by organizationId.
// @Tags files
// @Description Only for admin.
// @Accept  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Success 200  {object} response.Response{data=[]entities.FileResponse}
// @Router /api/files [get].
func (h *fileHandler) files(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	files, err := h.fileService.GetFiles(ctx, session.ID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, lo.Map(files, func(item entities.File, index int) entities.FileResponse {
		return item.Response()
	}), nil)
}

// @Summary Download file by id.
// @Tags files
// @Description Only for admin.
// @Accept  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Param   id path   string true  "file_id"
// @Success 200  {} null "file"
// @Router /api/files/{id} [get].
func (h *fileHandler) download(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	id := ctx.Param("id")

	file, err := h.fileService.GetFile(ctx, session.ID, uuid.MustParse(id))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	switch file.Type {
	case entities.FileXLSX:
		response.XLSX(ctx, file.Data, file.Name)
	case entities.FileCSV:
		response.CSV(ctx, file.Name, file.Array)
	default:
		response.BadRequest(ctx, errors.New("unknown file type"), nil)
	}
}

// @Summary Check file generation status.
// @Tags files
// @Description Only for admin. Websocket.
// @Accept  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Param   id path   string true  "file_id"
// @Router /api/files/ws/{id} [get].
func (h *fileHandler) subscribeHandler(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	id := ctx.Param("id")

	conn, err := h.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		response.BadRequest(ctx, err, nil)
		return
	}

	for {
		file, err := h.fileService.GetFile(ctx, session.ID, uuid.MustParse(id))
		if err != nil {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(err.Error())); err != nil {
				zap.S().Error(err)
			}
			conn.Close()
		}

		if err := conn.WriteMessage(websocket.TextMessage, []byte(file.Status)); err != nil {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(err.Error())); err != nil {
				zap.S().Error(err)
			}
			conn.Close()
		}

		if file.Status != entities.FileStatusInProgress {
			conn.Close()
			return
		}

		time.Sleep(time.Second)
	}
}
