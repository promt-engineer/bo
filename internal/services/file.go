package services

import (
	"backoffice/internal/entities"
	"backoffice/internal/repositories"
	"backoffice/pkg/file"
	"backoffice/utils"
	"bytes"
	"context"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

const ()

type FileDownloadingService struct {
	cfg         *file.Config
	cfgClient   *ClientInfoConfig
	spinService *SpinService
	fileRepo    repositories.FileRepository
}

func NewFileDownloadingService(cfg *file.Config, cfgClient *ClientInfoConfig, spinService *SpinService, fileRepo repositories.FileRepository) *FileDownloadingService {
	return &FileDownloadingService{
		cfg:         cfg,
		cfgClient:   cfgClient,
		spinService: spinService,
		fileRepo:    fileRepo,
	}
}

func (s *FileDownloadingService) GetFiles(ctx context.Context, organizationID uuid.UUID) ([]entities.File, error) {
	files, err := s.fileRepo.Find(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].CreatedAt.After(files[j].CreatedAt)
	})

	return files, nil
}

func (s *FileDownloadingService) GetFile(ctx context.Context, organizationID, id uuid.UUID) (*entities.File, error) {
	file, err := s.fileRepo.Get(ctx, organizationID, id)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (s *FileDownloadingService) FinancialXLSX(ctx context.Context, session *entities.Session, req *entities.FinancialBase) (uuid.UUID, error) {
	id := uuid.New()
	name := time.Now().UTC().Format("finacial_report-20060102150405.xlsx")

	file := &entities.File{
		ID:        id,
		CreatedAt: time.Now(),
		Status:    entities.FileStatusInProgress,
		Type:      entities.FileXLSX,
		Name:      name,
	}

	if err := s.fileRepo.Create(ctx, session.ID, file, s.cfg.TTL); err != nil {
		return id, err
	}

	go s.generateFinancialXLSX(ctx, session, req, file)

	return id, nil
}

func (s *FileDownloadingService) SpinsXLSX(ctx context.Context, session *entities.Session, req *entities.FinancialBase) (uuid.UUID, error) {
	id := uuid.New()
	name := time.Now().UTC().Format("spins_report-20060102150405.xlsx")

	file := &entities.File{
		ID:        id,
		CreatedAt: time.Now(),
		Status:    entities.FileStatusInProgress,
		Type:      entities.FileXLSX,
		Name:      name,
	}

	if err := s.fileRepo.Create(ctx, session.ID, file, s.cfg.TTL); err != nil {
		return id, err
	}

	go s.generateSpinsXLSX(ctx, session, req, file)

	return id, nil
}

func (s *FileDownloadingService) SessionXLSX(ctx context.Context, session *entities.Session, req *entities.FinancialBase) (uuid.UUID, error) {
	id := uuid.New()
	name := time.Now().UTC().Format("sessions_report-20060102150405.xlsx")

	file := &entities.File{
		ID:        id,
		CreatedAt: time.Now(),
		Status:    entities.FileStatusInProgress,
		Type:      entities.FileXLSX,
		Name:      name,
	}

	if err := s.fileRepo.Create(ctx, session.ID, file, s.cfg.TTL); err != nil {
		return id, err
	}

	go s.generateSessionXLSX(ctx, session, req, file)

	return id, nil
}

func (s *FileDownloadingService) FinancialCSV(ctx context.Context, session *entities.Session, req *entities.FinancialBase) (uuid.UUID, error) {
	id := uuid.New()

	file := &entities.File{
		ID:        id,
		CreatedAt: time.Now(),
		Status:    entities.FileStatusInProgress,
		Type:      entities.FileCSV,
		Name:      "financial",
	}

	if err := s.fileRepo.Create(ctx, session.ID, file, s.cfg.TTL); err != nil {
		return id, err
	}

	go s.generateFinancialCVS(ctx, session, req, file)

	return id, nil
}

func (s *FileDownloadingService) SpinsCSV(ctx context.Context, session *entities.Session, req *entities.FinancialBase) (uuid.UUID, error) {
	id := uuid.New()

	file := &entities.File{
		ID:        id,
		CreatedAt: time.Now(),
		Status:    entities.FileStatusInProgress,
		Type:      entities.FileCSV,
		Name:      "spins",
	}

	if err := s.fileRepo.Create(ctx, session.ID, file, s.cfg.TTL); err != nil {
		return id, err
	}

	go s.generateSpinsCVS(ctx, session, req, file)

	return id, nil
}

func (s *FileDownloadingService) SessionCSV(ctx context.Context, session *entities.Session, req *entities.FinancialBase) (uuid.UUID, error) {
	id := uuid.New()

	file := &entities.File{
		ID:        id,
		CreatedAt: time.Now(),
		Status:    entities.FileStatusInProgress,
		Type:      entities.FileCSV,
		Name:      "gaming-sessions",
	}

	if err := s.fileRepo.Create(ctx, session.ID, file, s.cfg.TTL); err != nil {
		return id, err
	}

	go s.generateSessionCSV(ctx, session, req, file)

	return id, nil
}

func (s *FileDownloadingService) AggregatedByGameCSV(ctx context.Context, session *entities.Session, req *entities.AggregateFilters) (uuid.UUID, error) {
	id := uuid.New()

	file := &entities.File{
		ID:        id,
		CreatedAt: time.Now(),
		Status:    entities.FileStatusInProgress,
		Type:      entities.FileCSV,
		Name:      "report-by-game",
	}

	if err := s.fileRepo.Create(ctx, session.ID, file, s.cfg.TTL); err != nil {
		return id, err
	}

	go s.generateAggregatedByGameCSV(ctx, session, req, file)

	return id, nil
}

func (s *FileDownloadingService) AggregatedByGameXLSX(ctx context.Context, session *entities.Session, req *entities.AggregateFilters) (uuid.UUID, error) {
	id := uuid.New()
	name := time.Now().UTC().Format("report_by_game-20060102150405.xlsx")
	file := &entities.File{
		ID:        id,
		CreatedAt: time.Now(),
		Status:    entities.FileStatusInProgress,
		Type:      entities.FileXLSX,
		Name:      name,
	}

	if err := s.fileRepo.Create(ctx, session.ID, file, s.cfg.TTL); err != nil {
		return id, err
	}

	go s.generateAggregatedByGameXLSX(ctx, session, req, file)

	return id, nil
}

func (s *FileDownloadingService) AggregatedByCountryCSV(ctx context.Context, session *entities.Session, req *entities.AggregateFilters) (uuid.UUID, error) {
	id := uuid.New()

	file := &entities.File{
		ID:        id,
		CreatedAt: time.Now(),
		Status:    entities.FileStatusInProgress,
		Type:      entities.FileCSV,
		Name:      "report-by-country",
	}

	if err := s.fileRepo.Create(ctx, session.ID, file, s.cfg.TTL); err != nil {
		return id, err
	}

	go s.generateAggregatedByCountryCSV(ctx, session, req, file)

	return id, nil
}

func (s *FileDownloadingService) AggregatedByCountryXLSX(ctx context.Context, session *entities.Session, req *entities.AggregateFilters) (uuid.UUID, error) {
	id := uuid.New()
	name := time.Now().UTC().Format("report_by_country-20060102150405.xlsx")
	file := &entities.File{
		ID:        id,
		CreatedAt: time.Now(),
		Status:    entities.FileStatusInProgress,
		Type:      entities.FileXLSX,
		Name:      name,
	}

	if err := s.fileRepo.Create(ctx, session.ID, file, s.cfg.TTL); err != nil {
		return id, err
	}

	go s.generateAggregatedByCountryXLSX(ctx, session, req, file)

	return id, nil
}
func (s *FileDownloadingService) ExportCurrencyXLSX(currencyInfo *entities.CurrencyInfo) (*excelize.File, string, error) {

	if len(currencyInfo.Table) == 0 {
		return nil, "", fmt.Errorf("no currency data provided")
	}

	file, err := utils.ExportXLSX(currencyInfo.Table)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create Excel file: %s", err.Error())
	}

	fileName := generateFileName(s.cfgClient.EnvironmentName, currencyInfo.Integrator, currencyInfo.Provider)

	return file, fileName, nil
}

func (s *FileDownloadingService) ImportCurrencyDataFileXLSX(filePath string) ([]entities.CurrencyAttributes, error) {
	rows, err := utils.ImportDataXLSX(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to import Excel file: %s", err.Error())
	}

	currencies, err := parseRowsToCurrencies(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Excel file: %s", err.Error())
	}

	return currencies, nil
}

func parseRowsToCurrencies(rows [][]string) ([]entities.CurrencyAttributes, error) {
	var currencies []entities.CurrencyAttributes
	for idx, row := range rows {
		if idx == 0 {
			continue
		}

		if len(row) < 3 {
			continue
		}

		multiplier, err := strconv.ParseInt(row[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid multiplier in row %d: %s", idx+1, err.Error())
		}

		currency := entities.CurrencyAttributes{
			Title:      row[0],
			Multiplier: multiplier,
			Synonym:    row[2],
		}

		currencies = append(currencies, currency)
	}

	return currencies, nil
}

func generateFileName(params ...string) string {
	if len(params) == 0 {
		return fmt.Sprintf("%s.xlsx", time.Now().Format("20060102150405"))
	}
	return fmt.Sprintf("%s.xlsx", strings.Join(params, "_"))
}

func (s *FileDownloadingService) generateFinancialXLSX(ctx context.Context, session *entities.Session, req *entities.FinancialBase, file *entities.File) {
	groupedReport, err := s.spinService.FinancialReport(ctx, &session.OrganizationID, req)
	if err != nil {
		s.saveFileWithError(ctx, session.ID, file, err)

		return
	}

	spins, err := s.spinService.AllSpins(ctx, &session.OrganizationID, req)
	if err != nil {
		s.saveFileWithError(ctx, session.ID, file, err)

		return
	}

	lo.ForEach(spins, func(item *entities.Spin, index int) {
		item.Prettify()
	})

	var pages []utils.Page
	separator := make([][]string, 3)

	exchangeInfo := make([][]string, 2)
	exchangeInfo[0] = []string{"Exchange currency", *req.Currency}
	totalPageData := []*entities.FinancialReport{groupedReport.Prettify()}

	table := exchangeInfo
	table = append(table, utils.ExtractTable(totalPageData, "xlsx")...)
	table = append(table, separator...)
	table = append(table, utils.ExtractTable(spins, "xlsx")...)

	pages = append([]utils.Page{{
		Name:  "total",
		Table: table,
	}}, pages...)

	xlsx, err := utils.ExportMultiPageXLSX(pages)
	if err != nil {
		s.saveFileWithError(ctx, session.ID, file, err)

		return
	}

	s.saveFileXLSX(ctx, session.ID, file, xlsx)
}

func (s *FileDownloadingService) generateSpinsXLSX(ctx context.Context, session *entities.Session, req *entities.FinancialBase, file *entities.File) {
	spins, err := s.spinService.AllSpins(ctx, &session.OrganizationID, req)
	if err != nil {
		s.saveFileWithError(ctx, session.ID, file, err)

		return
	}

	data := lo.Map(spins, func(item *entities.Spin, index int) *entities.Spin {
		return item.Prettify()
	})

	table := utils.ExtractTable(data, "xlsx")
	xlsx, err := utils.ExportXLSX(table)
	if err != nil {
		s.saveFileWithError(ctx, session.ID, file, err)

		return
	}

	s.saveFileXLSX(ctx, session.ID, file, xlsx)
}

func (s *FileDownloadingService) generateSessionXLSX(ctx context.Context, session *entities.Session, req *entities.FinancialBase, file *entities.File) {
	gamingSessions, err := s.spinService.AllGamingSessions(ctx, &session.OrganizationID, req)
	if err != nil {
		s.saveFileWithError(ctx, session.ID, file, err)

		return
	}

	lo.ForEach(gamingSessions, func(item *entities.GamingSession, index int) {
		item.Prettify()
	})

	table := utils.ExtractTable(gamingSessions, "xlsx")
	xlsx, err := utils.ExportXLSX(table)
	if err != nil {
		s.saveFileWithError(ctx, session.ID, file, err)

		return
	}

	s.saveFileXLSX(ctx, session.ID, file, xlsx)
}

func (s *FileDownloadingService) generateAggregatedByGameXLSX(ctx context.Context, session *entities.Session, req *entities.AggregateFilters, file *entities.File) {
	aggregatedReps, err := s.spinService.AggregatedReportByGame(ctx, &session.OrganizationID, *req.Currency, nil, req)
	if err != nil {
		s.saveFileWithError(ctx, session.ID, file, err)

		return
	}

	lo.ForEach(aggregatedReps, func(item *entities.AggregatedReportByGame, index int) {
		item.Prettify()
	})

	table := utils.ExtractTable(aggregatedReps, "xlsx")
	xlsx, err := utils.ExportXLSX(table)
	if err != nil {
		s.saveFileWithError(ctx, session.ID, file, err)

		return
	}

	s.saveFileXLSX(ctx, session.ID, file, xlsx)
}

func (s *FileDownloadingService) generateAggregatedByCountryXLSX(ctx context.Context, session *entities.Session, req *entities.AggregateFilters, file *entities.File) {
	aggregatedReps, err := s.spinService.AggregatedReportByCountry(ctx, &session.OrganizationID, *req.Currency, nil, req)
	if err != nil {
		s.saveFileWithError(ctx, session.ID, file, err)

		return
	}

	lo.ForEach(aggregatedReps, func(item *entities.AggregatedReportByCountry, index int) {
		item.Prettify()
	})

	table := utils.ExtractTable(aggregatedReps, "xlsx")
	xlsx, err := utils.ExportXLSX(table)
	if err != nil {
		s.saveFileWithError(ctx, session.ID, file, err)

		return
	}

	s.saveFileXLSX(ctx, session.ID, file, xlsx)
}

func (s *FileDownloadingService) generateFinancialCVS(ctx context.Context, session *entities.Session, req *entities.FinancialBase, file *entities.File) {
	rep, err := s.spinService.FinancialReport(ctx, &session.OrganizationID, req)
	if err != nil {
		s.saveFileWithError(ctx, session.ID, file, err)

		return
	}

	file.Status = entities.FileStatusReady
	file.Array = append(file.Array, rep.Prettify())
	file.ReflectType = reflect.TypeOf(file.Array[0]).Name()
	err = s.fileRepo.Update(ctx, session.ID, file, s.cfg.TTL)
	if err != nil {
		zap.S().Error(err)
	}
}

func (s *FileDownloadingService) generateAggregatedByGameCSV(ctx context.Context, session *entities.Session, req *entities.AggregateFilters, file *entities.File) {
	aggregatedReps, err := s.spinService.AggregatedReportByGame(ctx, &session.OrganizationID, *req.Currency, nil, req)
	if err != nil {
		s.saveFileWithError(ctx, session.ID, file, err)

		return
	}

	lo.ForEach(aggregatedReps, func(item *entities.AggregatedReportByGame, index int) {
		item.Prettify()
	})

	file.Status = entities.FileStatusReady

	for _, elem := range aggregatedReps {
		file.Array = append(file.Array, elem)
	}

	file.ReflectType = reflect.TypeOf(file.Array[0]).Name()
	err = s.fileRepo.Update(ctx, session.ID, file, s.cfg.TTL)
	if err != nil {
		zap.S().Error(err)
	}
}
func (s *FileDownloadingService) generateAggregatedByCountryCSV(ctx context.Context, session *entities.Session, req *entities.AggregateFilters, file *entities.File) {
	aggregatedReps, err := s.spinService.AggregatedReportByCountry(ctx, &session.OrganizationID, *req.Currency, nil, req)
	if err != nil {
		s.saveFileWithError(ctx, session.ID, file, err)

		return
	}

	lo.ForEach(aggregatedReps, func(item *entities.AggregatedReportByCountry, index int) {
		item.Prettify()
	})

	file.Status = entities.FileStatusReady

	for _, elem := range aggregatedReps {
		file.Array = append(file.Array, elem)
	}

	file.ReflectType = reflect.TypeOf(file.Array[0]).Name()
	err = s.fileRepo.Update(ctx, session.ID, file, s.cfg.TTL)
	if err != nil {
		zap.S().Error(err)
	}
}

func (s *FileDownloadingService) generateSpinsCVS(ctx context.Context, session *entities.Session, req *entities.FinancialBase, file *entities.File) {
	spins, err := s.spinService.AllSpins(ctx, &session.OrganizationID, req)
	if err != nil {
		s.saveFileWithError(ctx, session.ID, file, err)

		return
	}

	file.Status = entities.FileStatusReady

	for _, item := range spins {
		file.Array = append(file.Array, item.Prettify())
	}

	file.ReflectType = reflect.TypeOf(file.Array[0]).Name()

	err = s.fileRepo.Update(ctx, session.ID, file, s.cfg.TTL)
	if err != nil {
		zap.S().Error(err)
	}
}

func (s *FileDownloadingService) generateSessionCSV(ctx context.Context, session *entities.Session, req *entities.FinancialBase, file *entities.File) {
	gamingSessions, err := s.spinService.AllGamingSessions(ctx, &session.OrganizationID, req)
	if err != nil {
		s.saveFileWithError(ctx, session.ID, file, err)

		return
	}

	file.Status = entities.FileStatusReady

	for _, item := range gamingSessions {
		file.Array = append(file.Array, item.Prettify())
	}

	file.ReflectType = reflect.TypeOf(file.Array[0]).Name()

	err = s.fileRepo.Update(ctx, session.ID, file, s.cfg.TTL)
	if err != nil {
		zap.S().Error(err)
	}
}

func (s *FileDownloadingService) saveFileWithError(ctx context.Context, sessionID uuid.UUID, file *entities.File, err error) {
	file.Status = entities.FileStatusError
	file.Data = []byte(err.Error())

	err = s.fileRepo.Update(ctx, sessionID, file, s.cfg.TTL)
	if err != nil {
		zap.S().Error(err)
	}
}

func (s *FileDownloadingService) saveFileXLSX(ctx context.Context, sessionID uuid.UUID, file *entities.File, xlsx *excelize.File) {
	var b bytes.Buffer
	if err := xlsx.Write(&b); err != nil {
		s.saveFileWithError(ctx, sessionID, file, err)
		return
	}

	file.Status = entities.FileStatusReady
	file.Data = b.Bytes()

	err := s.fileRepo.Update(ctx, sessionID, file, s.cfg.TTL)
	if err != nil {
		zap.S().Error(err)
	}
}
