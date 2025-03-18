package handlers

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"backoffice/internal/services"
	"backoffice/internal/transport/http/requests"
	"backoffice/internal/transport/http/response"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

type currencyHandler struct {
	currencyService *services.CurrencyService
	cfgSender       *services.ConfigSenderService
	fileService     *services.FileDownloadingService
}

func (h *currencyHandler) Register(router *gin.RouterGroup) {
	currency := router.Group("currency")

	currency.POST("multiplier", h.create)
	currency.PUT("multiplier", h.update)
	currency.DELETE("multiplier", h.delete)
	currency.POST("multiplier/get", h.get)
	currency.POST("multiplier/search", h.search)
	currency.POST("multiplier/download", h.download)
	currency.POST("multiplier/upload", h.upload)
	currency.GET("", h.currencyGetAll)
	currency.GET(":alias", h.currencyGet)
	currency.POST("", h.currencyCreate)
	currency.DELETE(":alias", h.currencyDelete)
	currency.GET(":alias/exchange", h.currencyExchangeGet)
	currency.POST(":alias/exchange", h.currencyExchangeAdd)
	currency.DELETE(":alias/exchange", h.currencyExchangeDelete)
}

func NewCurrencyHandler(currencyService *services.CurrencyService, cfgSender *services.ConfigSenderService,
	fileService *services.FileDownloadingService) *currencyHandler {
	return &currencyHandler{currencyService: currencyService, cfgSender: cfgSender, fileService: fileService}
}

// @Summary Create new currency for specific provider and integrator.
// @Tags currency
// @Consume application/json
// @Description Create new currency and bound multiplier by provider_id and integrator_id.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.CurrencyMultiplier true "CreateCurrencyMultiplier"
// @Success 200 {object} response.Response{data=entities.CurrencyMultiplier}
// @Router /api/currency/multiplier [post].
func (h *currencyHandler) create(ctx *gin.Context) {
	req := requests.CurrencyMultiplier{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	cm, err := h.currencyService.CreateCurrencyMultiplier(ctx, req.OrganizationPairID, req.Title, req.Multiplier, req.Synonym)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.OK(ctx, cm, nil)
}

// @Summary Get currency multiplier.
// @Tags currency
// @Consume application/json
// @Description Get currency multiplier information.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.CurrencyMultiplierIdentify true "CurrencyMultiplierIdentify"
// @Success 200 {object} response.Response{data=entities.CurrencyMultiplier}
// @Router /api/currency/multiplier/get [post].
func (h *currencyHandler) get(ctx *gin.Context) {
	req := &requests.CurrencyMultiplierIdentify{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	cm, err := h.currencyService.Get(ctx, req.OrganizationPairID, req.Title)
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)

			return
		}

		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, cm, nil)
}

// @Summary Get all currency by filter.
// @Tags currency
// @Consume application/json
// @Description Filtering can be done by organization_pair_id or you can send an empty request body
// @Description and get all currency.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param request body requests.CurrencySearchRequest true "CurrencySearch"
// @Success 200 {object} response.Response{data=[]entities.CurrencyMultiplier}
// @Router /api/currency/multiplier/search [post].
func (h *currencyHandler) search(ctx *gin.Context) {
	req := &requests.CurrencySearchRequest{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	filter := make(map[string]interface{})

	if req.OrganizationPairID != uuid.Nil {
		filter["organization_pair_id"] = req.OrganizationPairID
	}

	cm, err := h.currencyService.Search(ctx, filter)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, cm, nil)
}

// @Summary Updates currency for specific provider and integrator.
// @Tags currency
// @Consume application/json
// @Description Updates currency and bound multiplier by provider_id and integrator_id.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.CurrencyMultiplier true "CreateCurrencyMultiplier"
// @Success 200 {object} response.Response{data=entities.CurrencyMultiplier}
// @Router /api/currency/multiplier [put].
func (h *currencyHandler) update(ctx *gin.Context) {
	req := requests.CurrencyMultiplier{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	cm, err := h.currencyService.UpdateCurrencyMultiplier(ctx, req.OrganizationPairID, req.Title, req.Multiplier, req.Synonym)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.OK(ctx, cm, nil)
}

// @Summary Delete currency for specific provider and integrator.
// @Tags currency
// @Consume application/json
// @Description Delete currency and bound multiplier by provider_id and integrator_id.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.CurrencyMultiplierIdentify true "CurrencyMultiplierIdentify"
// @Success 204
// @Router /api/currency/multiplier [delete].
func (h *currencyHandler) delete(ctx *gin.Context) {
	req := &requests.CurrencyMultiplierIdentify{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	err := h.currencyService.DeleteCurrencyMultiplier(ctx, req.OrganizationPairID, req.Title)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.NoContent(ctx)
}

// @Summary Get all currencies.
// @Tags currency
// @Consume application/json
// @Description Get all currencies` information.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Success 200 {object} response.Response{data=[]entities.Currency}
// @Router /api/currency [get].
func (h *currencyHandler) currencyGetAll(ctx *gin.Context) {
	var filters map[string]interface{}
	currencies, err := h.currencyService.CurrencyGetAll(ctx, filters)
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)

			return
		}

		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, currencies, nil)
}

// @Summary Get currency.
// @Tags currency
// @Consume application/json
// @Description Get currency information.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param title path string true "alias"
// @Success 200 {object} response.Response{data=[]entities.Currency}
// @Router /api/currency/{alias} [get].
func (h *currencyHandler) currencyGet(ctx *gin.Context) {
	alias := ctx.Param("alias")
	if alias == "" {
		response.BadRequest(ctx, fmt.Errorf("alias is empty"), nil)

		return
	}

	currency, err := h.currencyService.CurrencyGet(ctx, alias)
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			response.NotFound(ctx, err, nil)

			return
		}

		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, currency, nil)
}

// @Summary Create new currency.
// @Tags currency
// @Consume application/json
// @Description Create new currency.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.Currency true "CreateCurrency"
// @Success 200 {object} response.Response{data=entities.Currency}
// @Router /api/currency [post].
func (h *currencyHandler) currencyCreate(ctx *gin.Context) {
	req := requests.Currency{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	if req.Type == services.FakeCurrency && req.BaseCurrency == "" {
		response.ValidationFailed(ctx, fmt.Errorf("base currency required for type 'fake'"))

		return
	}

	currency, err := h.currencyService.CreateCurrency(ctx, req.Title, req.Alias, req.Type, req.BaseCurrency, req.Rate)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.OK(ctx, currency, nil)
}

// @Summary Delete currency.
// @Tags currency
// @Consume application/json
// @Description Delete currency.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param title path string true "alias"
// @Success 204
// @Router /api/currency [delete].
func (h *currencyHandler) currencyDelete(ctx *gin.Context) {
	alias := ctx.Param("alias")
	if alias == "" {
		response.BadRequest(ctx, fmt.Errorf("alias is empty"), nil)

		return
	}

	err := h.currencyService.DeleteCurrency(ctx, alias)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	h.cfgSender.Notify(ctx)

	response.NoContent(ctx)
}

// @Summary Download currency data as an Excel file
// @Description This endpoint allows you to search for currency data based on filters and download it as an Excel file.
// @Tags currency
// @Accept  json
// @Produce  application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param request body requests.CurrencySearchRequest true "Search Criteria"
// @Success 200  {file} file "Excel file with currency data"
// @Router /api/currency/multiplier/download [post]
func (h *currencyHandler) download(ctx *gin.Context) {
	req := &requests.CurrencySearchRequest{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	filter := make(map[string]interface{})

	if req.OrganizationPairID != uuid.Nil {
		filter["organization_pair_id"] = req.OrganizationPairID
	}

	cm, err := h.currencyService.Search(ctx, filter)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	currencyInfo, err := h.currencyService.FilterAndFormatCurrency(cm)
	if err != nil {
		response.BadRequest(ctx, "Error filtering and formatting data: %s", err.Error())

		return
	}

	file, nameFile, err := h.fileService.ExportCurrencyXLSX(currencyInfo)
	if err != nil {
		response.ServerError(ctx, "Error exporting data: %s", err.Error())

		return
	}

	response.XLSXFile(ctx, file, nameFile)
}

// @Summary Upload currency data from an Excel file
// @Description This endpoint allows you to upload currency data from an Excel file and save it to the database.
// @Tags currency
// @Accept  multipart/form-data
// @Produce  json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param organization_pair_id formData string true "Organization Pair ID (UUID format)"
// @Param file formData file true "Excel file containing currency data"
// @Success 200 {object} []entities.CurrencyAttributes "Uploaded currency data"
// @Router /api/currency/multiplier/upload [post]
func (h *currencyHandler) upload(ctx *gin.Context) {
	organizationPairID := ctx.PostForm("organization_pair_id")
	if organizationPairID == "" {
		response.ValidationFailed(ctx, errors.New("organization_pair_id is required"))

		return
	}

	organizationID, err := uuid.Parse(organizationPairID)
	if err != nil {
		response.BadRequest(ctx, "Invalid Organization ID format", nil)

		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		response.BadRequest(ctx, "Failed to upload file: %s", err.Error())

		return
	}

	tempFilePath, err := h.currencyService.SaveUploadedFile(ctx, file)
	if err != nil {
		response.ServerError(ctx, "Failed to save file: %s", err.Error())

		return
	}

	currency, err := h.fileService.ImportCurrencyDataFileXLSX(tempFilePath)
	if err != nil {
		response.ServerError(ctx, "Failed to import data: %s", err.Error())

		return
	}

	err = h.currencyService.CreateCurrencies(ctx, organizationID, currency)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, currency, nil)
}

// @Summary Get currency exchange.
// @Tags currency
// @Consume application/json
// @Description Get currency exchange information.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param title path string true "alias"
// @Param   order query string false "ordering format: column decs/acs"
// @Param   limit query string true "spins per page"
// @Param   page query string true "asking page"
// @Success 200 {object} response.Response{data=[]entities.CurrencyExchange}
// @Router /api/currency/{alias}/exchange [get].
func (h *currencyHandler) currencyExchangeGet(ctx *gin.Context) {
	alias := ctx.Param("alias")
	if alias == "" {
		response.BadRequest(ctx, fmt.Errorf("alias is empty"), nil)

		return
	}

	req := entities.CurrencyExchangePagination{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationFailed(ctx, err)

		return
	}

	pagination, err := h.currencyService.PaginateCurrencyExchange(ctx, alias, req.Order, req.Limit, req.Page)
	if err != nil {
		response.BadRequest(ctx, err, nil)
		return
	}

	response.OK(ctx, pagination, nil)
}

// @Summary Add a new currency exchange rate.
// @Tags currency
// @Description Add a new exchange rate between the specified "from" currency and the "to" currency (alias).
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param alias path string true "Alias of the 'to' currency."
// @Param data body requests.CurrencyExchangeRequest true "CurrencyExchangeRequest"
// @Success 200 {object} response.Response{data=entities.CurrencyExchange}
// @Failure 400 {object} response.Response
// @Router /api/currency/{alias}/exchange [post]
func (h *currencyHandler) currencyExchangeAdd(ctx *gin.Context) {
	alias := ctx.Param("alias")
	if alias == "" {
		response.BadRequest(ctx, "alias is empty", nil)
		return
	}

	req := requests.CurrencyExchangeRequest{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)
		return
	}

	currencyExchange, err := h.currencyService.AddCurrencyRate(ctx, req.From, alias, req.Rate)
	if err != nil {
		response.BadRequest(ctx, err, nil)
		return
	}

	response.OK(ctx, currencyExchange, nil)
}

// @Summary Delete a currency exchange rate.
// @Tags currency
// @Description Delete an existing exchange rate between the specified "from" currency and the "to" currency (alias).
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param alias path string true "Alias of the 'to' currency."
// @Param data body requests.DeleteCurrencyExchangeRequest true "DeleteCurrencyExchangeRequest"
// @Success 204
// @Failure 400 {object} response.Response
// @Router /api/currency/{alias}/exchange [delete]
func (h *currencyHandler) currencyExchangeDelete(ctx *gin.Context) {
	alias := ctx.Param("alias")
	if alias == "" {
		response.BadRequest(ctx, "alias is empty", nil)
		return
	}

	req := requests.DeleteCurrencyExchangeRequest{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ValidationFailed(ctx, err)
		return
	}

	err := h.currencyService.DeleteCurrencyRate(ctx, req.From, alias, req.Rate, req.CreatedAt)
	if err != nil {
		response.BadRequest(ctx, err, nil)
		return
	}

	response.NoContent(ctx)
}
