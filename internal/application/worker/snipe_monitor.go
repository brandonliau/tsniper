package worker

import (
	"time"

	"tsniper/internal/application/event"
	"tsniper/internal/application/ports"
	"tsniper/internal/domain/course"
	"tsniper/internal/domain/scope"
	"tsniper/internal/domain/snipe"
	"tsniper/internal/domain/user"

	"tsniper/pkg/eventbus"
	"tsniper/pkg/logger"
	"tsniper/pkg/multiticker"
	"tsniper/pkg/utils"
)

const (
	refreshInterval = 500 * time.Millisecond
)

type snipeMonitor struct {
	activeScope      scope.ActiveScope
	ticker           *multiticker.MultiTicker
	eventPublisher   eventbus.Publisher[event.CourseOpen]
	sectionsFeed     ports.SectionsFeed
	snipeCache       snipe.SnipeCache
	snipeRepository  snipe.SnipeRepository
	courseRepository course.CourseRepository
	userRepository   user.UserRepository
	logger           logger.Logger
}

func NewMonitorService(
	activeScope scope.ActiveScope,
	eventPublisher eventbus.Publisher[event.CourseOpen],
	sectionsFeed ports.SectionsFeed,
	snipeCache snipe.SnipeCache,
	snipeRepository snipe.SnipeRepository,
	courseRepository course.CourseRepository,
	userRepository user.UserRepository,
	logger logger.Logger,
) *snipeMonitor {
	return &snipeMonitor{
		activeScope:      activeScope,
		ticker:           multiticker.NewMultiTicker(refreshInterval, multiticker.NoOffset),
		eventPublisher:   eventPublisher,
		sectionsFeed:     sectionsFeed,
		snipeCache:       snipeCache,
		snipeRepository:  snipeRepository,
		courseRepository: courseRepository,
		userRepository:   userRepository,
		logger:           logger,
	}
}

func (s *snipeMonitor) Start() error {
	if err := s.pruneSnipes(); err != nil {
		return err
	}

	if err := s.warmCache(); err != nil {
		return err
	}

	for _, scp := range s.activeScope.Scopes() {
		go s.monitorSnipes(s.ticker.Subscribe(), scp)
	}

	return nil
}

func (s *snipeMonitor) Stop() error {
	s.ticker.Stop()
	return nil
}

func (s *snipeMonitor) monitorSnipes(ch <-chan time.Time, scp scope.AcademicScope) {
	previous, _ := s.sectionsFeed.FetchOpenSections(scp)

	for range ch {
		feed, err := s.sectionsFeed.FetchOpenSections(scp)
		if err != nil {
			// todo: remove this error message
			s.logger.Error("Failed to fetch open sections: %v", err)
			continue
		}

		for _, index := range utils.Intersection(feed, s.snipeCache.Tracked(scp)) {
			crs, err := s.courseRepository.Get(index, scp)
			if err != nil {
				s.logger.Error("Failed to get course: %v", err)
				continue
			}

			snipes, err := s.snipeRepository.GetByIndex(index, scp)
			if err != nil {
				s.logger.Error("Failed to get snipes: %v", err)
				continue
			}

			s.handleSnipes(snipes, crs)
		}

		for _, index := range utils.Difference(feed, previous) {
			crs, err := s.courseRepository.Get(index, scp)
			if err != nil {
				continue
			}
			if crs.LastOpen != 0 {
				crs.Open()
				if err := s.courseRepository.Save(crs); err != nil {
					s.logger.Error("Failed to save course: %v", err)
				}
			}
		}

		for _, index := range utils.Difference(previous, feed) {
			crs, err := s.courseRepository.Get(index, scp)
			if err != nil {
				continue
			}
			crs.Close()
			if err := s.courseRepository.Save(crs); err != nil {
				s.logger.Error("Failed to save course: %v", err)
			}
		}

		previous = feed
	}
}

func (s *snipeMonitor) pruneSnipes() error {
	start := time.Now()

	snipes, err := s.snipeRepository.GetAll()
	if err != nil {
		return err
	}

	var cleared int
	for _, snp := range snipes {
		if s.activeScope.Validate(snp.Scope) != nil {
			if err := s.snipeRepository.Delete(snp); err != nil {
				return err
			}
			cleared++
			continue
		}
	}

	s.logger.Info("Pruned %d stale snipes in %v", cleared, time.Since(start))
	return nil
}

func (s *snipeMonitor) warmCache() error {
	start := time.Now()

	snipes, err := s.snipeRepository.GetAll()
	if err != nil {
		return err
	}

	for _, snp := range snipes {
		s.snipeCache.Add(snp)
	}

	s.logger.Info("Synced %d snipes in %v", len(snipes), time.Since(start))
	return nil
}

func (s *snipeMonitor) handleSnipes(snipes []*snipe.Snipe, crs *course.Course) {
	var userIDs []string
	for _, snp := range snipes {
		userIDs = append(userIDs, snp.UserID)
	}

	evt := event.CourseOpen{
		Type:    event.CourseOpenNotification,
		Course:  crs,
		UserIDs: userIDs,
	}
	s.eventPublisher.PublishBlocking(evt)

	for _, snp := range snipes {
		if err := s.snipeRepository.Delete(snp); err != nil {
			s.logger.Error("Failed to delete snipe: %v", err)
		}
	}
	s.snipeCache.Clear(snipes)
}
