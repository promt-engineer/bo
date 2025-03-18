package services

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"backoffice/internal/repositories"
	"backoffice/pkg/exchange"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"mime/multipart"
	"strings"
	"time"
)

const (
	FakeCurrency = "fake"
	Currency     = "Currency"
	Multiplier   = "Multiplier"
	Synonym      = "Synonym"
)

var (
	ErrCanNotFindCurrencyConfig = errors.New("can not find currency config")
	ErrNegativeWager            = errors.New("negative wager found")
	ErrDefaultWagerOutOfList    = errors.New("default wager is out of default wager list")
)

type CurrencyService struct {
	repo                repositories.CurrencyRepository
	organizationService *OrganizationService
	exchangeClient      exchange.Client
}

func NewCurrencyService(repo repositories.CurrencyRepository, organizationService *OrganizationService, exchangeClient exchange.Client) *CurrencyService {
	return &CurrencyService{repo: repo, organizationService: organizationService, exchangeClient: exchangeClient}
}

func (s *CurrencyService) All(ctx context.Context) ([]*entities.CurrencyMultiplier, error) {
	return s.repo.All(ctx)
}

func (s *CurrencyService) UniqueCurrencyNames(ctx context.Context) ([]string, error) {
	return s.repo.UniqueCurrencyNames(ctx)
}

func (s *CurrencyService) Search(ctx context.Context, filter map[string]interface{}) ([]*entities.CurrencyMultiplier, error) {
	return s.repo.Search(ctx, filter)
}

func (s *CurrencyService) CreateCurrencyMultiplier(ctx context.Context, organizationPairID uuid.UUID, title string, multiplier int64, synonym string) (*entities.CurrencyMultiplier, error) {
	res, err := s.repo.CreateCurrencyMultiplier(ctx, &entities.CurrencyMultiplier{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		OrganizationPairID: organizationPairID,

		Title:      strings.ToLower(title),
		Multiplier: multiplier,
		Synonym:    synonym,
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *CurrencyService) UpdateCurrencyMultiplier(ctx context.Context, organizationPairID uuid.UUID, title string, multiplier int64, synonym string) (*entities.CurrencyMultiplier, error) {
	res, err := s.repo.UpdateCurrencyMultiplier(ctx, &entities.CurrencyMultiplier{
		UpdatedAt: time.Now(),

		OrganizationPairID: organizationPairID,
		Title:              strings.ToLower(title),
		Multiplier:         multiplier,
		Synonym:            synonym,
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *CurrencyService) Get(ctx context.Context, organizationPairID uuid.UUID, title string) (*entities.CurrencyMultiplier, error) {
	return s.repo.Get(ctx, map[string]interface{}{"organization_pair_id": organizationPairID, "title": title})
}

func (s *CurrencyService) DeleteCurrencyMultiplier(ctx context.Context, organizationPairID uuid.UUID, title string) error {
	cm, err := s.Get(ctx, organizationPairID, title)
	if err != nil {
		return err
	}

	err = s.repo.DeleteCurrencyMultiplier(ctx, cm)
	if err != nil {
		return err
	}

	return nil
}

func (s *CurrencyService) MergeCurrenciesByProvider(ctx context.Context, providerID uuid.UUID) ([]string, error) {
	ccs, err := s.repo.All(ctx)
	if err != nil {
		return nil, err
	}

	ccs = lo.Filter(ccs, func(item *entities.CurrencyMultiplier, index int) bool {
		return item.ProviderIntegratorPair.ProviderID == providerID
	})

	currencies := []string{}

	for _, cc := range ccs {
		currencies = append(currencies, cc.Title)
	}

	currencies = lo.Uniq(currencies)

	return currencies, nil
}

func (s *CurrencyService) CurrencyGetAll(ctx context.Context, filters map[string]interface{}) ([]*entities.Currency, error) {
	return s.repo.CurrencyGetAll(ctx, filters)
}

func (s *CurrencyService) CurrencyGet(ctx context.Context, alias string) (*entities.Currency, error) {
	return s.repo.CurrencyGet(ctx, alias)
}

func (s *CurrencyService) CreateCurrency(ctx context.Context, title, alias, curType, baseCurrency string, rate float64) (*entities.Currency, error) {
	curr, err := s.repo.CurrencyGet(ctx, alias)
	if err != nil && !errors.Is(err, e.ErrEntityNotFound) {
		return nil, err
	}

	if curr != nil {
		return nil, fmt.Errorf("сurrency with alias %s already exists", alias)
	}

	if curType == FakeCurrency {
		baseCurr, err := s.repo.CurrencyGet(ctx, baseCurrency)
		if err != nil {
			return nil, err
		}

		if baseCurr.Type == FakeCurrency {
			return nil, fmt.Errorf("сurrency with alias %s cannot be with the base currency of type 'fake'", alias)
		}
	}

	res, err := s.repo.CreateCurrency(ctx, &entities.Currency{
		Title:        strings.ToLower(title),
		Alias:        strings.ToLower(alias),
		Type:         curType,
		Rate:         rate,
		BaseCurrency: strings.ToLower(baseCurrency),
	})

	if err != nil {
		return nil, err
	}

	err = s.UpdateCurrencyExchange(ctx, strings.ToLower(alias), strings.ToLower(baseCurrency))
	if err != nil {
		zap.S().Info(err)
	}

	return res, nil
}

func (s *CurrencyService) DeleteCurrency(ctx context.Context, alias string) error {
	curr, err := s.repo.CurrencyGet(ctx, alias)
	if err != nil && !errors.Is(err, e.ErrEntityNotFound) {
		return err
	}

	filter := map[string]interface{}{"title": alias}

	//cm, err := s.repo.GetCurrencyMultiplier(ctx, filter)
	cm, err := s.repo.Search(ctx, filter)
	if err != nil && !errors.Is(err, e.ErrEntityNotFound) {
		return err
	}

	if len(cm) != 0 {
		return fmt.Errorf("there are multipliers associated with the %s currency", alias)
	}

	err = s.repo.DeleteCurrency(ctx, curr)
	if err != nil {
		return err
	}

	return nil
}

func (s *CurrencyService) UpdateCurrencyExchange(ctx context.Context, currency, baseCurrency string) error {
	resp, err := s.exchangeClient.UpdateCurrencies(ctx, currency, baseCurrency)
	if err != nil {
		return err
	}

	if resp.Status != "Ok" {
		return fmt.Errorf("the exchange rate for currency %s has not been updated", currency)
	}
	return nil
}

func (s *CurrencyService) GetCurrencyMultipliersByOrganizationPairs(ctx context.Context, organizationsPairIDs []uuid.UUID, currency string) (cm []*entities.CurrencyMultiplier, err error) {
	filter := map[string]interface{}{"organization_pair_id": organizationsPairIDs, "title": currency}

	return s.repo.Search(ctx, filter)
}

func (s *CurrencyService) FilterAndFormatCurrency(currencyData []*entities.CurrencyMultiplier) (*entities.CurrencyInfo, error) {
	if len(currencyData) == 0 {
		return nil, nil
	}

	table := [][]string{{Currency, Multiplier, Synonym}}

	for _, currency := range currencyData {
		row := []string{
			currency.Title,
			fmt.Sprintf("%d", currency.Multiplier),
			currency.Synonym,
		}
		table = append(table, row)
	}

	integrator := currencyData[0].ProviderIntegratorPair.Integrator.Name
	provider := currencyData[0].ProviderIntegratorPair.Provider.Name

	currencyInfo := &entities.CurrencyInfo{
		Table:      table,
		Integrator: integrator,
		Provider:   provider,
	}

	return currencyInfo, nil
}

func (s *CurrencyService) SaveUploadedFile(c *gin.Context, file *multipart.FileHeader) (string, error) {
	tempFilePath := fmt.Sprintf("/tmp/%s", file.Filename)
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		return "", fmt.Errorf("failed to save file: %s", err.Error())
	}
	return tempFilePath, nil
}

func (s *CurrencyService) CreateCurrencies(ctx *gin.Context, organizationPairID uuid.UUID, currencies []entities.CurrencyAttributes) error {
	for _, currency := range currencies {
		_, err := s.CreateCurrencyMultiplier(ctx, organizationPairID, currency.Title, currency.Multiplier, currency.Synonym)
		if err != nil {
			return fmt.Errorf("failed to save currency %s: %s", currency.Title, err.Error())
		}
	}
	return nil
}

func (s *CurrencyService) PaginateCurrencyExchange(ctx context.Context, currency string, order string, limit int, page int) (
	pagination entities.Pagination[entities.CurrencyExchange], err error) {

	out, err := s.exchangeClient.GetCurrencyRates(ctx, &exchange.AllCurrencyRatesIn{
		Page:     uint64(page),
		Limit:    uint64(limit),
		Order:    order,
		Currency: currency,
	})
	if err != nil {
		return
	}

	lo.ForEach(out.Items, func(item *exchange.CurrencyRates, index int) {
		pagination.Items = append(pagination.Items, entities.CurrencyFromExchange(item))
	})

	pagination.Total = int(out.Total)
	pagination.Limit = int(out.Limit)
	pagination.CurrentPage = int(out.CurrentPage)

	return pagination, nil
}

func (s *CurrencyService) AddCurrencyRate(ctx context.Context, from, to string, rate float64) (*entities.CurrencyExchange, error) {
	resp, err := s.exchangeClient.AddCurrencyRate(ctx, strings.ToLower(from), strings.ToLower(to), rate)
	if err != nil {
		return nil, err
	}

	currencyRate := &entities.CurrencyExchange{
		CreatedAt: resp.CreatedAt.AsTime(),
		From:      resp.From,
		To:        resp.To,
		Rate:      resp.Rate,
	}

	return currencyRate, nil
}

func (s *CurrencyService) DeleteCurrencyRate(ctx context.Context, from, to string, rate float64, createdAt time.Time) error {
	resp, err := s.exchangeClient.DeleteCurrencyRate(ctx, strings.ToLower(from), strings.ToLower(to), rate, createdAt)
	if err != nil {
		return err
	}

	if resp.Status != "Ok" {
		return fmt.Errorf("failed to delete currency rate from %s to %s", from, to)
	}

	return nil
}
