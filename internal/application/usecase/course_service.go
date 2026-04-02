package usecase

import (
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
	Course *course.Course
	Count  int
}

func (s *CourseService) Search(req SearchCourseRequest) (*SearchCourseResult, error) {
	usr, err := s.userRepository.Get(req.UserID)
	if err != nil {
		return nil, err
	}

	var cmp *scope.Campus
	if req.Campus != nil {
		parsed, err := scope.ParseCampus(*req.Campus)
		if err != nil {
			return nil, err
		}
		cmp = &parsed
	} else {
		cmp = usr.DefaultCampus()
	}

	var szn *scope.Season
	if req.Season != nil {
		parsed, err := scope.ParseSeason(*req.Season)
		if err != nil {
			return nil, err
		}
		szn = &parsed
	}

	scp, err := s.activeScope.Resolve(cmp, szn)
	if err != nil {
		return nil, err
	}

	crs, err := s.courseRepository.Get(req.Index, scp)
	if err != nil {
		return nil, err
	}

	snipes, err := s.snipeRepository.GetByIndex(req.Index, scp)
	if err != nil {
		return nil, err
	}

	return &SearchCourseResult{Course: crs, Count: len(snipes)}, nil
}
