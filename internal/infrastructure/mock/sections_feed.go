package mock

import (
	"math/rand"
	"sync"

	"tsniper/internal/application/ports"
	"tsniper/internal/domain/course"
	"tsniper/internal/domain/scope"

	"tsniper/pkg/logger"
)

var _ ports.SectionsFeed = (*MockSectionsFeed)(nil)

type MockSectionsFeed struct {
	mu          sync.Mutex
	courseRepo  course.CourseRepository
	open        map[string]map[string]bool // scope key -> index -> open
	openChance  float64
	closeChance float64
	logger      logger.Logger
}

func NewSectionsFeed(courseRepo course.CourseRepository, openChance float64, closeChance float64, logger logger.Logger) *MockSectionsFeed {
	return &MockSectionsFeed{
		courseRepo:  courseRepo,
		open:        make(map[string]map[string]bool),
		openChance:  openChance,
		closeChance: closeChance,
		logger:      logger,
	}
}

func scopeKey(scp scope.AcademicScope) string {
	return string(scp.Campus) + ":" + string(scp.Term) + ":" + scp.Year
}

func (m *MockSectionsFeed) FetchOpenSections(scp scope.AcademicScope) ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := scopeKey(scp)
	if _, ok := m.open[key]; !ok {
		if err := m.init(scp); err != nil {
			return nil, err
		}
	}

	sections := m.open[key]
	for index, isOpen := range sections {
		if !isOpen && rand.Float64() < m.openChance {
			sections[index] = true
			m.logger.Debug("[mock] %s opened", index)
		} else if isOpen && rand.Float64() < m.closeChance {
			sections[index] = false
			m.logger.Debug("[mock] %s closed", index)
		}
	}

	var result []string
	for index, isOpen := range sections {
		if isOpen {
			result = append(result, index)
		}
	}
	return result, nil
}

func (m *MockSectionsFeed) init(scp scope.AcademicScope) error {
	courses, err := m.courseRepo.GetAllByScope(scp)
	if err != nil {
		return err
	}

	key := scopeKey(scp)
	m.open[key] = make(map[string]bool, len(courses))
	for _, crs := range courses {
		m.open[key][crs.Index] = false
	}
	m.logger.Debug("[mock] Loaded %d courses for scope %s", len(courses), key)
	return nil
}
