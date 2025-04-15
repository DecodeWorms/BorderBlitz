package services

import (
	"sync"
	"time"

	"github.com/DecodeWorms/BorderBlitz/repository"
)

// FXService handles currency exchange rates
type FXService struct {
	stablecoinRepo *repository.StablecoinRepository
	ratesMutex     sync.RWMutex
	rates          map[string]float64
	lastUpdated    time.Time
}

// NewFXService creates a new FX service
func NewFXService(stablecoinRepo *repository.StablecoinRepository) *FXService {
	service := &FXService{
		stablecoinRepo: stablecoinRepo,
		rates:          make(map[string]float64),
	}

	// Initialize rates from the database
	service.updateRates()

	// Start a goroutine to update rates periodically
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			service.updateRates()
		}
	}()

	return service
}

// updateRates updates the exchange rates from the database
func (s *FXService) updateRates() {
	coins, err := s.stablecoinRepo.ListAll()
	if err != nil {
		return
	}

	s.ratesMutex.Lock()
	defer s.ratesMutex.Unlock()

	for _, coin := range coins {
		s.rates[coin.Symbol] = coin.USDRate
	}

	s.lastUpdated = time.Now()
}

// GetUSDRate gets the USD exchange rate for a currency
func (s *FXService) GetUSDRate(symbol string) (float64, bool) {
	s.ratesMutex.RLock()
	defer s.ratesMutex.RUnlock()

	rate, exists := s.rates[symbol]
	return rate, exists
}

// GetExchangeRate gets the exchange rate between two currencies
func (s *FXService) GetExchangeRate(fromSymbol, toSymbol string) (float64, bool) {
	s.ratesMutex.RLock()
	defer s.ratesMutex.RUnlock()

	fromRate, fromExists := s.rates[fromSymbol]
	toRate, toExists := s.rates[toSymbol]

	if !fromExists || !toExists {
		return 0, false
	}

	return fromRate / toRate, true
}

// GetLastUpdated gets the time of the last rates update
func (s *FXService) GetLastUpdated() time.Time {
	s.ratesMutex.RLock()
	defer s.ratesMutex.RUnlock()

	return s.lastUpdated
}
