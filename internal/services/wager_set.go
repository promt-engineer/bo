package services

import (
	"backoffice/internal/entities"
	"backoffice/internal/repositories"
	"context"
	"github.com/google/uuid"
	"time"
)

type WagerSetService struct {
	wagerRepo repositories.BaseRepository[entities.WagerSet]
}

func NewWagerSetService(wagerRepo repositories.BaseRepository[entities.WagerSet]) *WagerSetService {
	return &WagerSetService{wagerRepo: wagerRepo}
}

func (s *WagerSetService) Paginate(ctx context.Context, organizationID uuid.UUID, filters map[string]interface{}, limit int, page int) (
	pagination entities.Pagination[entities.WagerSet], err error) {
	filters["organization_id"] = organizationID
	return s.wagerRepo.Paginate(ctx, filters, "created_at desc", limit, page)
}

func (s *WagerSetService) Get(ctx context.Context, id uuid.UUID) (*entities.WagerSet, error) {
	return s.wagerRepo.FindBy(ctx, map[string]interface{}{"id": id})
}

func (s *WagerSetService) Create(ctx context.Context, organizationID uuid.UUID, name string, wagerLevels []int64, defaultWager int64) (*entities.WagerSet, error) {
	ws := &entities.WagerSet{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		ID:             uuid.New(),
		OrganizationID: organizationID,
		Name:           name,

		IsActive: true,
	}

	if err := ws.SetNewWagerParams(wagerLevels, defaultWager); err != nil {
		return nil, err
	}

	return s.wagerRepo.Save(ctx, ws)
}

func (s *WagerSetService) Update(ctx context.Context, id uuid.UUID, name string, wagerLevels []int64, defaultWager int64, isActive bool) (*entities.WagerSet, error) {
	ws, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	ws.Name = name
	ws.IsActive = isActive

	if err := ws.SetNewWagerParams(wagerLevels, defaultWager); err != nil {
		return nil, err
	}

	return s.wagerRepo.Save(ctx, ws)
}

func (s *WagerSetService) Delete(ctx context.Context, id uuid.UUID) error {
	ws, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	return s.wagerRepo.Delete(ctx, ws)
}
