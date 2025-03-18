package services

import (
	"backoffice/internal/entities"
	"backoffice/internal/repositories"
	"context"
	"github.com/google/uuid"
	"time"
)

type CurrencySetService struct {
	currencyRepo    repositories.BaseRepository[entities.CurrencySet]
	currencyService *CurrencyService
}

func NewCurrencySetService(currencyService *CurrencyService, currencyRepo repositories.BaseRepository[entities.CurrencySet]) *CurrencySetService {
	return &CurrencySetService{
		currencyService: currencyService,
		currencyRepo:    currencyRepo,
	}
}

func (s *CurrencySetService) Paginate(ctx context.Context, organizationID uuid.UUID, filters map[string]interface{}, limit int, page int) (
	pagination entities.Pagination[entities.CurrencySet], err error) {
	filters["organization_id"] = organizationID

	return s.currencyRepo.Paginate(ctx, filters, "created_at desc", limit, page)
}

func (s *CurrencySetService) Get(ctx context.Context, id uuid.UUID) (*entities.CurrencySet, error) {
	return s.currencyRepo.FindBy(ctx, map[string]interface{}{"id": id})
}

func (s *CurrencySetService) Create(ctx context.Context, organizationID uuid.UUID, name string, currencies []string) (*entities.CurrencySet, error) {
	allCurrencies, err := s.currencyService.MergeCurrenciesByProvider(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	cs := &entities.CurrencySet{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		ID:             uuid.New(),
		OrganizationID: organizationID,
		Name:           name,

		IsActive: true,
	}

	cs.SetCurrencies(allCurrencies, currencies)

	return s.currencyRepo.Save(ctx, cs)
}

func (s *CurrencySetService) Update(ctx context.Context, id uuid.UUID, name string, currencies []string, isActive bool) (*entities.CurrencySet, error) {
	cs, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	allCurrencies, err := s.currencyService.MergeCurrenciesByProvider(ctx, cs.OrganizationID)
	if err != nil {
		return nil, err
	}

	cs.Name = name
	cs.IsActive = isActive
	cs.SetCurrencies(allCurrencies, currencies)

	return s.currencyRepo.Save(ctx, cs)
}

func (s *CurrencySetService) Delete(ctx context.Context, id uuid.UUID) error {
	ws, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	return s.currencyRepo.Delete(ctx, ws)
}
