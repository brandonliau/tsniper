package usecase

import (
	"errors"
	"fmt"

	"tsniper/internal/application/view"
	"tsniper/internal/domain/course"
	"tsniper/internal/domain/scope"
	"tsniper/internal/domain/snipe"
	"tsniper/internal/domain/user"
)

type CourseService struct {
	activeScope      scope.ActiveScope
	snipeRepository  snipe.SnipeRepository
	userRepository   user.UserRepository
	courseRepository course.CourseRepository
}

func NewCourseService(activeScope scope.ActiveScope, snipeRepository snipe.SnipeRepository, userRepository user.UserRepository, courseRepository course.CourseRepository) *CourseService {
	return &CourseService{
		activeScope:      activeScope,
		snipeRepository:  snipeRepository,
		userRepository:   userRepository,
		courseRepository: courseRepository,
	}
}

// --- Search Course ---
type SearchCourseRequest struct {
	UserID string
	Index  string
	Campus *string
	Season *string
}

type SearchCourseResult struct {
	Course *view.CourseView
	Count  int
}

var (
	ErrSearchCourseInvalid = errors.New("invalid search course request")
)

func (s *CourseService) Search(req SearchCourseRequest) (*SearchCourseResult, error) {
	usr, err := s.userRepository.Get(req.UserID)
	if err != nil {
		return nil, err
	}

	var cmp *scope.Campus
	if req.Campus != nil {
		parsed, err := scope.ParseCampus(*req.Campus)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrSearchCourseInvalid, err)
		}
		cmp = &parsed
	} else {
		cmp = usr.DefaultCampus()
	}

	var szn *scope.Season
	if req.Season != nil {
		parsed, err := scope.ParseSeason(*req.Season)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrSearchCourseInvalid, err)
		}
		szn = &parsed
	}

	scp, err := s.activeScope.Resolve(cmp, szn)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrSearchCourseInvalid, err)
	}

	crs, err := s.courseRepository.Get(req.Index, scp)
	if err != nil {
		if errors.Is(err, course.ErrCourseNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrSearchCourseInvalid, err)
		}
		return nil, err
	}

	snipes, err := s.snipeRepository.ListByIndex(req.Index, scp)
	if err != nil {
		return nil, err
	}

	return &SearchCourseResult{Course: view.FromCourse(crs), Count: len(snipes)}, nil
}
