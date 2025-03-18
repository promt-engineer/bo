package services

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"backoffice/internal/entities"
	"backoffice/internal/errors"
	"backoffice/internal/transport/http/requests"
	"backoffice/pkg/overlord"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type LobbyConfig struct {
	LobbyBaseURL string
}

type LobbyService struct {
	cfg            *LobbyConfig
	gameService    *GameService
	overlordClient overlord.Client
}

func NewLobbyService(
	cfg *LobbyConfig,
	overlordClient overlord.Client,
	gameService *GameService) *LobbyService {
	return &LobbyService{
		cfg:            cfg,
		overlordClient: overlordClient,
		gameService:    gameService,
	}
}

func (s *LobbyService) validate(game *entities.Game, req requests.LobbyRequest) error {
	if *req.ShortLink {
		wagerSetFound := false
		for _, wagerSet := range game.AvailableWagerSets {
			if wagerSet.ID == req.WagerSetID {
				wagerSetFound = true
				break
			}
		}
		if !wagerSetFound {
			return errors.ErrValidationFailed("wager_set_id")
		}

		if req.Volatility != nil {
			volatilityFound := false
			for _, v := range game.AvailableVolatility {
				if *req.Volatility == v {
					volatilityFound = true
					break
				}
			}
			if !volatilityFound {
				return errors.ErrValidationFailed("volatility")
			}
		}

		if req.RTP != nil {
			rtpFound := false
			for _, r := range game.AvailableRTP {
				if *req.RTP == r {
					rtpFound = true
					break
				}
			}
			if !rtpFound {
				return errors.ErrValidationFailed("rtp")
			}
		}
	}

	currencyFound := false
	for _, currency := range game.Currencies {
		if strings.EqualFold(currency, req.Currency) {
			currencyFound = true
			break
		}
	}
	if !currencyFound {
		return errors.ErrValidationFailed("currency")
	}

	localeFound := false
	for _, locale := range game.UserLocales {
		if locale == req.UserLocale {
			localeFound = true
			break
		}
	}
	if !localeFound {
		return errors.ErrValidationFailed("user_locale")
	}

	return nil

}

func (s *LobbyService) getWagerSetByID(wagerSetID uuid.UUID, wagerSets []entities.WagerSet) *entities.WagerSet {
	for _, wagerSet := range wagerSets {
		if wagerSet.ID == wagerSetID {
			return &wagerSet
		}
	}
	return nil
}

func (s *LobbyService) Lobby(ctx context.Context, req requests.LobbyRequest) (string, error) {
	game, err := s.gameService.GetGameByName(ctx, req.Game)
	if err != nil {
		zap.S().Error(err)
		return "", err
	}

	host := getGameHost(game, s.cfg)

	if err = s.validate(game, req); err != nil {
		return "", err
	}

	convertedLocale := strings.Replace(req.UserLocale, "_", "-", 1)

	if *req.ShortLink {
		wagerSet := s.getWagerSetByID(req.WagerSetID, game.AvailableWagerSets)
		msg := &overlord.SaveParamsIn{
			Integrator:   req.Integrator,
			Game:         req.Game,
			Rtp:          req.RTP,
			Volatility:   req.Volatility,
			Wagers:       wagerSet.WagerLevels,
			SessionId:    req.SessionID.String(),
			IsDemo:       true,
			Currency:     strings.ToLower(req.Currency),
			UserId:       req.UserID.String(),
			UserLocale:   convertedLocale,
			DefaultWager: &wagerSet.DefaultWager,
			LobbyUrl:     req.LobbyURL,
			Jurisdiction: req.Jurisdiction,
		}

		if req.ShowCheats != nil {
			msg.ShowCheats = *req.ShowCheats
		}
		if req.LowBalance != nil {
			msg.LowBalance = *req.LowBalance
		}
		if req.ShortLink != nil {
			msg.ShortLink = *req.ShortLink
		}

		_, err = s.overlordClient.SaveParams(ctx, msg)
		if err != nil {
			zap.S().Info(fmt.Sprintf("save params err: %+v", err))
			return "", err
		}

		return fmt.Sprintf("%v/%v/?integrator=%v&session_id=%v",
			host, req.Game, req.Integrator, req.SessionID), nil
	}

	values := url.Values{}
	values.Add("user_id", req.UserID.String())
	values.Add("jurisdiction", req.Jurisdiction)
	values.Add("currency", strings.ToUpper(req.Currency))
	values.Add("user_locale", convertedLocale)
	values.Add("integrator", req.Integrator)

	if req.ShowCheats != nil && *req.ShowCheats {
		values.Add("showcheats", fmt.Sprintf("%v", *req.ShowCheats))
	}
	if req.LowBalance != nil && *req.LowBalance {
		values.Add("low_balance", fmt.Sprintf("%v", *req.LowBalance))
	}
	if req.LobbyURL != "" {
		values.Add("lobbyurl", req.LobbyURL)
	}

	fullURL := fmt.Sprintf("%v/%v/?%v", host, req.Game, values.Encode())

	return fullURL, nil
}

func getGameHost(game *entities.Game, cfg *LobbyConfig) string {
	host := cfg.LobbyBaseURL

	if game != nil {
		s := ""
		if game.ClientUrl != "" {
			s = game.ClientUrl
		} else if game.ApiUrl != "" {
			s = game.ApiUrl
		}

		u, err := url.Parse(s)
		if err != nil {
			return host
		}

		if u.Host == "dev.heronbyte.com" {
			return host
		}

		host = u.Host
	}

	return host
}
