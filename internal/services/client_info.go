package services

type ClientInfoConfig struct {
	LogoUrl         string
	EnvironmentName string
}

type ClientInfoService struct {
	cfg *ClientInfoConfig
}

func NewClientInfoService(cfg *ClientInfoConfig) *ClientInfoService {
	return &ClientInfoService{
		cfg: cfg,
	}
}

func (s *ClientInfoService) GetLogoURL() string {
	return s.cfg.LogoUrl
}

func (s *ClientInfoService) GetEnvironmentName() string {
	return s.cfg.EnvironmentName
}
