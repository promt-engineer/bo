package handlers

import (
	"backoffice/internal/entities"
	"backoffice/internal/services"
	"backoffice/internal/transport/http/requests"
	"backoffice/internal/transport/http/response"
	"github.com/gin-gonic/gin"
	"strings"
)

type dictionaryHandler struct {
	gameService         *services.GameService
	spinService         *services.SpinService
	currencyService     *services.CurrencyService
	organizationService *services.OrganizationService
}

func NewDictionaryHandler(
	gameService *services.GameService, spinService *services.SpinService,
	currencyService *services.CurrencyService, organizationService *services.OrganizationService) *dictionaryHandler {
	return &dictionaryHandler{
		gameService:         gameService,
		spinService:         spinService,
		currencyService:     currencyService,
		organizationService: organizationService,
	}
}

func (h *dictionaryHandler) Register(route *gin.RouterGroup) {
	dictionary := route.Group("dictionaries")
	dictionary.GET("games", h.games)
	dictionary.GET("hosts", h.hosts)
	dictionary.GET("currencies", h.currencies)
	dictionary.GET("jurisdictions", h.jurisdictions)
	dictionary.GET("languages", h.languages)
	dictionary.GET("locales", h.locales)
	dictionary.GET("integrators", h.integrators)
	dictionary.GET("integrator-operators", h.integratorOperators)
	dictionary.GET("main-currencies", h.mainCurrencies)
	dictionary.POST("locales", h.addLocale)
	dictionary.POST("jurisdictions", h.addJurisdiction)
	dictionary.POST("languages", h.addLanguage)
	dictionary.DELETE("locales/:locale", h.delLocale)
	dictionary.DELETE("jurisdictions/:jurisdiction", h.delJurisdiction)
	dictionary.DELETE("languages/:language", h.delLanguage)
}

// @Summary Get available games.
// @Tags dictionaries
// @Consume application/json
// @Description Get list of available games for filters.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Success 200 {object} response.Response{data=[]entities.Game}
// @Router /api/dictionaries/games [get].
func (h *dictionaryHandler) games(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	games, err := h.gameService.AllForStat(ctx, &session.OrganizationID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, games, nil)
}

// @Summary Get available hosts.
// @Tags dictionaries
// @Consume application/json
// @Description Get list of available hosts for filters.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/dictionaries/hosts [get].
func (h *dictionaryHandler) hosts(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	games, err := h.spinService.Hosts(ctx, &session.OrganizationID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, games, nil)
}

// @Summary Get available currencies.
// @Tags dictionaries
// @Consume application/json
// @Description Get list of available currencies.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/dictionaries/currencies [get].
func (h *dictionaryHandler) currencies(ctx *gin.Context) {
	currencies, err := h.currencyService.UniqueCurrencyNames(ctx)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, currencies, nil)
}

// @Summary Get available jurisdictions.
// @Tags dictionaries
// @Consume application/json
// @Description Get list of available jurisdictions.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/dictionaries/jurisdictions [get].
func (h *dictionaryHandler) jurisdictions(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	dict, err := h.gameService.GetDictionaries(ctx, &session.OrganizationID, "jurisdictions")
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, dict, nil)
}

// @Summary Get available languages.
// @Tags dictionaries
// @Consume application/json
// @Description Get list of available languages.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/dictionaries/languages [get].
func (h *dictionaryHandler) languages(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	dict, err := h.gameService.GetDictionaries(ctx, &session.OrganizationID, "languages")
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, dict, nil)
}

// @Summary Get available locales.
// @Tags dictionaries
// @Consume application/json
// @Description Get list of available locales.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/dictionaries/locales [get].
func (h *dictionaryHandler) locales(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	dict, err := h.gameService.GetDictionaries(ctx, &session.OrganizationID, "user_locales")
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, dict, nil)
}

// @Summary Get available integrators.
// @Tags dictionaries
// @Consume application/json
// @Description Get list of available integrators.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/dictionaries/integrators [get].
func (h *dictionaryHandler) integrators(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	integrators, err := h.organizationService.GetIntegratorNames(ctx, session.Account)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, integrators, nil)
}

// @Summary Get available integrator/operator pair.
// @Tags dictionaries
// @Consume application/json
// @Description Get list of available integrator/operator pair.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Success 200 {object} response.Response{data=map[string][]string}
// @Router /api/dictionaries/integrator-operators [get].
func (h *dictionaryHandler) integratorOperators(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)

	integrators, err := h.spinService.IntegratorOperatorsMap(ctx, &session.OrganizationID)
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, integrators, nil)
}

// @Summary Get main currencies.
// @Tags dictionaries
// @Consume application/json
// @Description Get list of main currencies.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/dictionaries/main-currencies [get].
func (h *dictionaryHandler) mainCurrencies(ctx *gin.Context) {
	response.OK(ctx, []string{"usd", "eur"}, nil)
}

// @Summary Add new locale.
// @Tags dictionaries
// @Consume application/json
// @Description Add new locale.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.Dictionary true "AddDictionaryValue"
// @Success 200 {object} response.Response{data=string}
// @Router /api/dictionaries/locales [post].
func (h *dictionaryHandler) addLocale(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	var requestData requests.Dictionary
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.BadRequest(ctx, err, nil)
		return
	}
	dict, err := h.gameService.AddValueToTheDictionary(ctx, &session.OrganizationID, "user_locales", requestData.Value)
	if err != nil {
		response.BadRequest(ctx, err, nil)
		return
	}

	response.OK(ctx, dict, nil)
}

// @Summary Add new jurisdiction.
// @Tags dictionaries
// @Consume application/json
// @Description Add new jurisdiction.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.Dictionary true "AddDictionaryValue"
// @Success 200 {object} response.Response{data=string}
// @Router /api/dictionaries/jurisdictions [post].
func (h *dictionaryHandler) addJurisdiction(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	var requestData requests.Dictionary
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.BadRequest(ctx, err, nil)
		return
	}
	dict, err := h.gameService.AddValueToTheDictionary(ctx, &session.OrganizationID, "jurisdictions", strings.ToUpper(requestData.Value))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, dict, nil)
}

// @Summary Add new language.
// @Tags dictionaries
// @Consume application/json
// @Description Add new language.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param data body requests.Dictionary true "AddDictionaryValue"
// @Success 200 {object} response.Response{data=string}
// @Router /api/dictionaries/languages [post].
func (h *dictionaryHandler) addLanguage(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	var requestData requests.Dictionary
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.BadRequest(ctx, err, nil)
		return
	}
	dict, err := h.gameService.AddValueToTheDictionary(ctx, &session.OrganizationID, "languages", strings.ToUpper(requestData.Value))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, dict, nil)
}

// @Summary Delete locale.
// @Tags dictionaries
// @Consume application/json
// @Description Delete locale.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param locale path string true "Locale"
// @Success 200 {object} response.Response{data=string}
// @Router /api/dictionaries/locales/{locale} [delete].
func (h *dictionaryHandler) delLocale(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	locale := ctx.Param("locale")

	err := h.gameService.DeleteValueFromTheDictionary(ctx, &session.OrganizationID, "user_locales", locale)
	if err != nil {
		response.BadRequest(ctx, err, nil)
		return
	}

	response.OK(ctx, nil, nil)
}

// @Summary Delete jurisdiction.
// @Tags dictionaries
// @Consume application/json
// @Description Delete jurisdiction.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param jurisdiction path string true "Jurisdiction"
// @Success 200 {object} response.Response{data=string}
// @Router /api/dictionaries/jurisdictions/{jurisdiction} [delete].
func (h *dictionaryHandler) delJurisdiction(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	jurisdiction := ctx.Param("jurisdiction")

	err := h.gameService.DeleteValueFromTheDictionary(ctx, &session.OrganizationID, "jurisdictions", strings.ToUpper(jurisdiction))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, nil, nil)
}

// @Summary Delete language.
// @Tags dictionaries
// @Consume application/json
// @Description Delete language.
// @Accept json
// @Produce json
// @Security X-Authenticate
// @Param X-Auth header string true "Insert your access token"
// @Param language path string true "Language"
// @Success 200 {object} response.Response{data=string}
// @Router /api/dictionaries/languages/{language} [delete].
func (h *dictionaryHandler) delLanguage(ctx *gin.Context) {
	session := ctx.Value("session").(*entities.Session)
	language := ctx.Param("language")

	err := h.gameService.DeleteValueFromTheDictionary(ctx, &session.OrganizationID, "languages", strings.ToUpper(language))
	if err != nil {
		response.BadRequest(ctx, err, nil)

		return
	}

	response.OK(ctx, nil, nil)
}
