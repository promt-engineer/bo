package handlers

import (
	"backoffice/internal/entities"
	"backoffice/internal/services"
	"backoffice/internal/transport/http/response"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type reportHandler struct {
	spinService *services.SpinService
	fileService *services.FileDownloadingService
}

func NewReportHandler(spinService *services.SpinService,
	fileService *services.FileDownloadingService,
) *reportHandler {
	return &reportHandler{
		spinService: spinService,
		fileService: fileService,
	}
}

func (h *reportHandler) Register(router *gin.RouterGroup) {
	reports := router.Group("reports")

	reports.GET("financial", h.financial)
	reports.GET("financial/csv", h.financialCSV)
	reports.GET("financial/xlsx", h.financialXLSX)

	spins := reports.Group("spins")
	spins.GET("", h.spins)
	spins.GET("csv", h.spinsCSV)
	spins.GET(":id", h.spin)
	spins.GET("xlsx", h.spinsXLSX)

	sessions := reports.Group("sessions")
	sessions.GET("", h.sessions)
	sessions.GET(":id", h.session)
	sessions.GET("csv", h.sessionCSV)
	sessions.GET("xlsx", h.sessionXLSX)

	users := reports.Group("users")

	users.GET(":id", h.user)

	aggregated := reports.Group("aggregated")

	aggregated.GET("by_game", h.aggregatedByGame)
	aggregated.GET("by_game/:country", h.aggregatedByGamePerCountry)
	aggregated.GET("by_game/csv", h.aggregatedByGameCSV)
	aggregated.GET("by_game/xlsx", h.aggregatedByGameXLSX)
	aggregated.GET("by_country", h.aggregatedByCountry)
	aggregated.GET("by_country/:game", h.aggregatedByCountryPerGame)
	aggregated.GET("by_country/csv", h.aggregatedByCountryCSV)
	aggregated.GET("by_country/xlsx", h.aggregatedByCountryXLSX)

	reports.GET("currencies", h.currencies)
}

// @Summary Get financial report.
// @Tags reports
// @Consume application/json
// @Description Only for admin.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Param   currency query string true "currency"
// @Param   integrator query string false "integrator name"
// @Param   operator query string false "operator name"
// @Param   session_id query string false "session id"
// @Param   round_id query string false "round id"
// @Param   game query string false "game name"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05"
// @Success 200  {object} response.Response{data=entities.FinancialReport}
// @Router /api/reports/financial [get].
func (h *reportHandler) financial(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	req := &entities.FinancialBase{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	rep, err := h.spinService.FinancialReport(ctx, &session.OrganizationID, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, rep.Prettify(), nil)
}

// @Summary Get spin pagination.
// @Tags reports
// @Consume application/json
// @Description Only for admin.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param   currency query string true "currency"
// @Param   integrator query string false "integrator name"
// @Param   operator query string false "operator name"
// @Param   session_id query string false "session id"
// @Param   round_id query string false "round id"
// @Param   game query string false "game name"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05"
// @Param   order query string false "ordering format: column decs/acs"
// @Param   limit query string true "spins per page"
// @Param   page query string true "asking page"
// @Success 200  {object} response.Response{data=[]entities.Spin}
// @Router /api/reports/spins [get].
func (h *reportHandler) spins(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	req := entities.SpinPagination{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	var groupBy []string
	if req.GroupBy != nil && *req.GroupBy != "" {
		groupBy = strings.Split(strings.Trim(*req.GroupBy, " "), ",")
	}

	if len(groupBy) > 0 {
		pagination, err := h.spinService.PaginateGrouped(ctx, &session.OrganizationID, &req.FinancialBase, req.Order, req.Limit, req.Page, groupBy)
		if err != nil {
			response.BadRequest(ctx, err, nil)
			return
		}

		lo.ForEach(pagination.Items, func(item *entities.GroupedSpin, index int) {
			item.Prettify()
		})

		response.OK(ctx, pagination, nil)
	} else {
		pagination, err := h.spinService.Paginate(ctx, &session.OrganizationID, &req.FinancialBase, req.Order, req.Limit, req.Page)
		if err != nil {
			response.BadRequest(ctx, err, nil)
			return
		}

		lo.ForEach(pagination.Items, func(item *entities.Spin, index int) {
			item.Prettify()
		})

		response.OK(ctx, pagination, nil)
	}
}

// @Summary Get spin.
// @Tags spins
// @Consume application/json
// @Description Get spin information.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Param   id path   string true  "spin_id"
// @Param   currency query string true "currency"
// @Success 200  {object} response.Response{data=entities.Spin}
// @Router /api/reports/spins/{id} [get].
func (h *reportHandler) spin(ctx *gin.Context) {
	spin, err := h.spinService.GetSpin(ctx, ctx.Param("id"), ctx.Query("currency"))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, spin.Prettify(), nil)
}

// @Summary Get spin pagination.
// @Tags reports
// @Consume application/json
// @Description Only for admin.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param   currency query string true "currency"
// @Param   integrator query string false "integrator name"
// @Param   operator query string false "operator name"
// @Param   session_id query string false "session id"
// @Param   round_id query string false "round id"
// @Param   game query string false "game name"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05"
// @Param   rtp_from query string false "rtp from"
// @Param   order query string false "ordering format: column decs/acs"
// @Param   limit query string true "spins per page"
// @Param   page query string true "asking page"
// @Success 200  {object} response.Response{data=[]entities.GamingSession}
// @Router /api/reports/sessions [get].
func (h *reportHandler) sessions(ctx *gin.Context) {
	req := entities.SpinPagination{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	session := ctx.Value("session").(*entities.Session)
	pagination, err := h.spinService.PaginateGamingSession(ctx, &session.OrganizationID, &req.FinancialBase, req.Order, req.Limit, req.Page)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	lo.ForEach(pagination.Items, func(item *entities.GamingSession, index int) {
		item.Prettify()
	})

	response.OK(ctx, pagination, nil)
}

// @Summary Get spin information.
// @Tags reports
// @Consume application/json
// @Description Only for admin.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param   id path   string true  "session_id"
// @Param   currency query string true "currency"
// @Success 200  {object} response.Response{data=entities.GamingSession}
// @Router /api/reports/sessions/{id} [get].
func (h *reportHandler) session(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	gamingSession, err := h.spinService.Session(ctx, &session.OrganizationID, uuid.MustParse(ctx.Param("id")), ctx.Query("currency"))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, gamingSession.Prettify(), nil)
}

// @Summary Get financial report csv.
// @Tags reports
// @Consume text/csv
// @Description Only for admin.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Param   currency query string true "currency"
// @Param   integrator query string false "integrator name"
// @Param   operator query string false "operator name"
// @Param   session_id query string false "session id"
// @Param   round_id query string false "round id"
// @Param   game query string false "game name"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05"
// @Success 200  {object} response.Response{data=entities.FileReportResponse}
// @Router /api/reports/financial/csv [get].
func (h *reportHandler) financialCSV(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	req := &entities.FinancialBase{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	id, err := h.fileService.FinancialCSV(ctx, session, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, entities.FileReportResponse{
		ID: id,
	}, nil)
}

// @Summary Get financial report csv.
// @Tags reports
// @Consume text/csv
// @Description Only for admin.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Param   currency query string true "currency"
// @Param   integrator query string false "integrator name"
// @Param   operator query string false "operator name"
// @Param   session_id query string false "session id"
// @Param   round_id query string false "round id"
// @Param   game query string false "game name"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05"
// @Success 200  {object} response.Response{data=entities.FileReportResponse}
// @Router /api/reports/financial/xlsx [get].
func (h *reportHandler) financialXLSX(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	req := &entities.FinancialBase{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	id, err := h.fileService.FinancialXLSX(ctx, session, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, entities.FileReportResponse{
		ID: id,
	}, nil)
}

// @Summary Get spins report csv.
// @Tags reports
// @Consume text/csv
// @Description Only for admin.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Param   currency query string true "currency"
// @Param   integrator query string false "integrator name"
// @Param   operator query string false "operator name"
// @Param   session_id query string false "session id"
// @Param   round_id query string false "round id"
// @Param   game query string false "game name"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05"
// @Success 200  {object} response.Response{data=entities.FileReportResponse}
// @Router /api/reports/spins/csv [get].
func (h *reportHandler) spinsCSV(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	req := &entities.FinancialBase{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	id, err := h.fileService.SpinsCSV(ctx, session, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, entities.FileReportResponse{
		ID: id,
	}, nil)
}

// @Summary Get spins report xlsx.
// @Tags reports
// @Consume text/xlsx
// @Description Only for admin.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Param   currency query string true "currency"
// @Param   integrator query string false "integrator name"
// @Param   operator query string false "operator name"
// @Param   session_id query string false "session id"
// @Param   round_id query string false "round id"
// @Param   game query string false "game name"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05"
// @Success 200  {object} response.Response{data=entities.FileReportResponse}
// @Router /api/reports/spins/xlsx [get].
func (h *reportHandler) spinsXLSX(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	req := &entities.FinancialBase{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	id, err := h.fileService.SpinsXLSX(ctx, session, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, entities.FileReportResponse{
		ID: id,
	}, nil)
}

// @Summary Get gaming session report csv.
// @Tags reports
// @Consume text/csv
// @Description Only for admin.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Param   currency query string true "currency"
// @Param   integrator query string false "integrator name"
// @Param   operator query string false "operator name"
// @Param   session_id query string false "session id"
// @Param   round_id query string false "round id"
// @Param   game query string false "game name"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05"
// @Success 200  {object} response.Response{data=entities.FileReportResponse}
// @Router /api/reports/sessions/csv [get].
func (h *reportHandler) sessionCSV(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	req := &entities.FinancialBase{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	id, err := h.fileService.SessionCSV(ctx, session, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, entities.FileReportResponse{
		ID: id,
	}, nil)
}

// @Summary Get gaming session report xlsx.
// @Tags reports
// @Consume text/xlsx
// @Description Only for admin.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Param   currency query string true "currency"
// @Param   integrator query string false "integrator name"
// @Param   operator query string false "operator name"
// @Param   session_id query string false "session id"
// @Param   round_id query string false "round id"
// @Param   game query string false "game name"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05"
// @Success 200  {object} response.Response{data=entities.FileReportResponse}
// @Router /api/reports/sessions/xlsx [get].
func (h *reportHandler) sessionXLSX(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	req := &entities.FinancialBase{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	id, err := h.fileService.SessionXLSX(ctx, session, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, entities.FileReportResponse{
		ID: id,
	}, nil)
}

// @Summary Get available currencies.
// @Tags reports
// @Consume application/json
// @Description Only for admin.
// @Accept  json
// @Produce  application/json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Param   integrator query string false "integrator name"
// @Param   operator query string false "operator name"
// @Param   session_id query string false "session id"
// @Param   round_id query string false "round id"
// @Param   game query string false "game name"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05"
// @Success 200  {object} response.Response{data=[]string}
// @Router /api/reports/currencies [get].
func (h *reportHandler) currencies(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	req := &entities.FinancialBase{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	currencies, err := h.spinService.Currencies(ctx, &session.OrganizationID, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, currencies, nil)
}

// @Summary Get aggregated report.
// @Tags reports
// @Consume application/json
// @Description Only for admin.
// @Accept  json
// @Produce  application/json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param   currency query string true "currency name"
// @Param   integrator query string false "integrator name"
// @Param   operator query string false "operator name"
// @Param   session_id query string false "session id"
// @Param   round_id query string false "round id"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05-00:00"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05-00:00"
// @Success 200  {object} response.Response{data=[]entities.AggregatedReportByGame}
// @Router /api/reports/aggregated/by_game [get].
func (h *reportHandler) aggregatedByGame(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	req := &entities.AggregateFilters{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	aggregatedReps, err := h.spinService.AggregatedReportByGame(ctx, &session.OrganizationID, *req.Currency, nil, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	lo.ForEach(aggregatedReps, func(item *entities.AggregatedReportByGame, index int) {
		item.Prettify()
	})

	response.OK(ctx, aggregatedReps, nil)
}

// @Summary Get aggregated report.
// @Tags reports
// @Consume application/json
// @Description Only for admin.
// @Accept  json
// @Produce  application/json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param   currency query string true "currency name"
// @Param   integrator query string false "integrator name"
// @Param   operator query string false "operator name"
// @Param   session_id query string false "session id"
// @Param   round_id query string false "round id"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05-00:00"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05-00:00"
// @Param   country path   string true  "country"
// @Success 200  {object} response.Response{data=[]entities.AggregatedReportByGame}
// @Router /api/reports/aggregated/by_game/{country} [get].
func (h *reportHandler) aggregatedByGamePerCountry(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	req := &entities.AggregateFilters{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	country := ctx.Param("country")

	aggregatedReps, err := h.spinService.AggregatedReportByGame(ctx, &session.OrganizationID, *req.Currency, &country, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	lo.ForEach(aggregatedReps, func(item *entities.AggregatedReportByGame, index int) {
		item.Prettify()
	})

	response.OK(ctx, aggregatedReps, nil)
}

// @Summary Get aggregated xlsx.
// @Tags reports
// @Consume text/xlsx
// @Description Only for admin.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Param   currency query string true "currency"
// @Param   integrator query string false "integrator name"
// @Param   game query string false "game name"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05"
// @Success 200  {object} response.Response{data=entities.FileReportResponse}
// @Router /api/reports/aggregated/by_game/xlsx [get].
func (h *reportHandler) aggregatedByGameXLSX(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	req := &entities.AggregateFilters{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	id, err := h.fileService.AggregatedByGameXLSX(ctx, session, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, entities.FileReportResponse{
		ID: id,
	}, nil)
}

// @Summary Get aggregated report csv.
// @Tags reports
// @Consume text/csv
// @Description Only for admin.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Param   currency query string true "currency"
// @Param   integrator query string false "integrator name"
// @Param   game query string false "game name"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05"
// @Success 200  {object} response.Response{data=entities.FileReportResponse}
// @Router /api/reports/aggregated/by_game/csv [get].
func (h *reportHandler) aggregatedByGameCSV(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	req := &entities.AggregateFilters{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	id, err := h.fileService.AggregatedByGameCSV(ctx, session, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, entities.FileReportResponse{
		ID: id,
	}, nil)
}

// @Summary Get aggregated report.
// @Tags reports
// @Consume application/json
// @Description Only for admin.
// @Accept  json
// @Produce  application/json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param   currency query string true "currency name"
// @Param   integrator query string false "integrator name"
// @Param   operator query string false "operator name"
// @Param   session_id query string false "session id"
// @Param   round_id query string false "round id"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05-00:00"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05-00:00"
// @Success 200  {object} response.Response{data=[]entities.AggregatedReportByGame}
// @Router /api/reports/aggregated/by_country [get].
func (h *reportHandler) aggregatedByCountry(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	req := &entities.AggregateFilters{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	aggregatedReps, err := h.spinService.AggregatedReportByCountry(ctx, &session.OrganizationID, *req.Currency, nil, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	lo.ForEach(aggregatedReps, func(item *entities.AggregatedReportByCountry, index int) {
		item.Prettify()
	})

	response.OK(ctx, aggregatedReps, nil)
}

// @Summary Get aggregated report.
// @Tags reports
// @Consume application/json
// @Description Only for admin.
// @Accept  json
// @Produce  application/json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param   currency query string true "currency name"
// @Param   integrator query string false "integrator name"
// @Param   operator query string false "operator name"
// @Param   session_id query string false "session id"
// @Param   round_id query string false "round id"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05-00:00"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05-00:00"
// @Param   country path   string true  "country"
// @Success 200  {object} response.Response{data=[]entities.AggregatedReportByGame}
// @Router /api/reports/aggregated/by_country/{game} [get].
func (h *reportHandler) aggregatedByCountryPerGame(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	req := &entities.AggregateFilters{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	game := ctx.Param("game")

	aggregatedReps, err := h.spinService.AggregatedReportByCountry(ctx, &session.OrganizationID, *req.Currency, &game, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	lo.ForEach(aggregatedReps, func(item *entities.AggregatedReportByCountry, index int) {
		item.Prettify()
	})

	response.OK(ctx, aggregatedReps, nil)
}

// @Summary Get aggregated xlsx.
// @Tags reports
// @Consume text/xlsx
// @Description Only for admin.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Param   currency query string true "currency"
// @Param   integrator query string false "integrator name"
// @Param   game query string false "game name"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05"
// @Success 200  {object} response.Response{data=entities.FileReportResponse}
// @Router /api/reports/aggregated/by_country/xlsx [get].
func (h *reportHandler) aggregatedByCountryXLSX(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	req := &entities.AggregateFilters{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	id, err := h.fileService.AggregatedByCountryXLSX(ctx, session, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, entities.FileReportResponse{
		ID: id,
	}, nil)
}

// @Summary Get aggregated report csv.
// @Tags reports
// @Consume text/csv
// @Description Only for admin.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Param   currency query string true "currency"
// @Param   integrator query string false "integrator name"
// @Param   game query string false "game name"
// @Param   starting_from query string false "time format: 2006-01-02 15:04:05"
// @Param   ending_at query string false "time format: 2006-01-02 15:04:05"
// @Success 200  {object} response.Response{data=entities.FileReportResponse}
// @Router /api/reports/aggregated/by_country/csv [get].
func (h *reportHandler) aggregatedByCountryCSV(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	req := &entities.AggregateFilters{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	id, err := h.fileService.AggregatedByCountryCSV(ctx, session, req)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, entities.FileReportResponse{
		ID: id,
	}, nil)
}

// @Summary Get user.
// @Tags users
// @Consume application/json
// @Description Get user information.
// @Accept  json
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your refresh token"
// @Param   id path   string true  "user_id"
// @Param  currency query string true "currency"
// @Success 200  {object} response.Response{data=entities.UserReport}
// @Router /api/reports/users/{id} [get].
func (h *reportHandler) user(ctx *gin.Context) {
	id := ctx.Param("id")
	currency, ok := ctx.GetQuery("currency")
	if !ok {
		response.BadRequest(ctx, errors.New("currency not found"), nil)
	}

	session := ctx.Value("session").(*entities.Session)

	report, err := h.spinService.UserReport(ctx, &session.OrganizationID, id, currency)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, report.Prettify(), nil)
}
