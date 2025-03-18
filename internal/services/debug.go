package services

import (
	"backoffice/internal/repositories"
	"math"
	"sync"
)

var debugService *DebugService
var debugOnce sync.Once

type DebugReport struct {
	DBSizeMB float64
}

type DebugListener func(report DebugReport)

type DebugService struct {
	debugRepo   repositories.DebugRepository
	subscribers []DebugListener
}

func NewDebugService(debugRepo repositories.DebugRepository) *DebugService {
	debugOnce.Do(func() {
		debugService = &DebugService{debugRepo: debugRepo}
	})

	return debugService
}

func (s *DebugService) Subscribe(subs ...DebugListener) {
	s.subscribers = append(s.subscribers, subs...)
}

func (s *DebugService) NotifyAll() error {
	if len(s.subscribers) > 0 {
		size, err := s.debugRepo.SizeMB()
		if err != nil {
			return err
		}

		size = math.Round(size*100) / 100

		rep := DebugReport{DBSizeMB: size}

		for _, sub := range s.subscribers {
			sub(rep)
		}
	}

	return nil
}
