package worker

import (
	"time"

	"tsniper/internal/application/ports"
	"tsniper/internal/domain/course"
	"tsniper/internal/domain/scope"

	"tsniper/pkg/logger"

	"github.com/robfig/cron/v3"
)

type courseIndexer struct {
	cron             *cron.Cron
	activeScope      scope.ActiveScope
	courseFeed       ports.CourseFeed
	courseRepository course.CourseRepository
	logger           logger.Logger
}

func NewCourseIndexer(activeScope scope.ActiveScope, courseFeed ports.CourseFeed, courseRepository course.CourseRepository, logger logger.Logger) *courseIndexer {
	return &courseIndexer{
		cron:             cron.New(),
		activeScope:      activeScope,
		courseFeed:       courseFeed,
		courseRepository: courseRepository,
		logger:           logger,
	}
}

func (s *courseIndexer) Start() error {
	s.indexCourses()

	_, err := s.cron.AddFunc("0 7 * * *", s.indexCourses)
	if err != nil {
		return err
	}
	_, err = s.cron.AddFunc("0 19 * * *", s.indexCourses)
	if err != nil {
		return err
	}
	s.cron.Start()

	return nil
}

func (s *courseIndexer) Stop() error {
	s.cron.Stop()
	return nil
}

func (s *courseIndexer) indexCourses() {
	for _, scp := range s.activeScope.Scopes() {
		start := time.Now()

		courses, err := s.courseFeed.FetchCourses(scp)
		if err != nil {
			s.logger.Error("Failed to get courses: %v", err)
			return
		}

		if err := s.courseRepository.BatchCreate(courses); err != nil {
			s.logger.Error("Failed to batch create courses: %v", err)
			return
		}

		s.logger.Info("Indexed %d courses in %v", len(courses), time.Since(start))
	}
}
