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

type SnipeService struct {
	activeScope      scope.ActiveScope
	snipeCache       snipe.SnipeCache
	snipeRepository  snipe.SnipeRepository
	courseRepository course.CourseRepository
	userRepository   user.UserRepository
}

func NewSnipeService(activeScope scope.ActiveScope, snipeCache snipe.SnipeCache, snipeRepository snipe.SnipeRepository, courseRepository course.CourseRepository, userRepository user.UserRepository) *SnipeService {
	return &SnipeService{
		activeScope:      activeScope,
		snipeCache:       snipeCache,
		snipeRepository:  snipeRepository,
		courseRepository: courseRepository,
		userRepository:   userRepository,
	}
}

// --- Add Snipe ---
type AddSnipeRequest struct {
	UserID string
	Index  string
	Campus *string
	Season *string
}

type AddSnipeResult struct {
	Course *view.CourseView
}

var (
	ErrAddSnipeInvalid   = errors.New("invalid add snipe request")
	ErrAddSnipeDuplicate = errors.New("snipe already exists")
)

func (s *SnipeService) Add(req AddSnipeRequest) (*AddSnipeResult, error) {
	usr, err := s.userRepository.Get(req.UserID)
	if err != nil {
		return nil, err
	}

	var cmp *scope.Campus
	if req.Campus != nil {
		parsed, err := scope.ParseCampus(*req.Campus)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrAddSnipeInvalid, err)
		}
		cmp = &parsed
	} else {
		cmp = usr.DefaultCampus()
	}

	var szn *scope.Season
	if req.Season != nil {
		parsed, err := scope.ParseSeason(*req.Season)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrAddSnipeInvalid, err)
		}
		szn = &parsed
	}

	scp, err := s.activeScope.Resolve(cmp, szn)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrAddSnipeInvalid, err)
	}

	crs, err := s.courseRepository.Get(req.Index, scp)
	if err != nil {
		if errors.Is(err, course.ErrCourseNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrAddSnipeInvalid, err)
		}
		return nil, err
	}

	snp := snipe.NewSnipe(req.UserID, req.Index, scp)
	if err := s.snipeRepository.Create(snp); err != nil {
		if errors.Is(err, snipe.ErrSnipeDuplicate) {
			return nil, fmt.Errorf("%w: %w", ErrAddSnipeDuplicate, err)
		}
		return nil, err
	}
	s.snipeCache.Add(snp)

	return &AddSnipeResult{Course: view.FromCourse(crs)}, nil
}

// --- Re-add Snipe ---
type ReAddSnipeRequest struct {
	UserID string
	Index  string
	Campus string
	Term   string
	Year   string
}

type ReAddSnipeResult struct {
	Course *view.CourseView
}

var (
	ErrReAddSnipeInvalid   = errors.New("invalid re-add snipe request")
	ErrReAddSnipeDuplicate = errors.New("snipe already exists")
)

func (s *SnipeService) ReAdd(req ReAddSnipeRequest) (*ReAddSnipeResult, error) {
	cmp, err := scope.ParseCampus(req.Campus)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReAddSnipeInvalid, err)
	}

	term, err := scope.ParseTerm(req.Term)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReAddSnipeInvalid, err)
	}

	scp := scope.AcademicScope{Campus: cmp, Term: term, Year: req.Year}
	if err := s.activeScope.Validate(scp); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReAddSnipeInvalid, err)
	}

	crs, err := s.courseRepository.Get(req.Index, scp)
	if err != nil {
		if errors.Is(err, course.ErrCourseNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrReAddSnipeInvalid, err)
		}
		return nil, err
	}

	snp := snipe.NewSnipe(req.UserID, req.Index, scp)
	if err := s.snipeRepository.Create(snp); err != nil {
		if errors.Is(err, snipe.ErrSnipeDuplicate) {
			return nil, fmt.Errorf("%w: %w", ErrReAddSnipeDuplicate, err)
		}
		return nil, err
	}
	s.snipeCache.Add(snp)

	return &ReAddSnipeResult{Course: view.FromCourse(crs)}, nil
}

// --- Remove Snipe ---
type RemoveSnipeRequest struct {
	UserID string
	Index  string
	Campus *string
	Season *string
}

type RemoveSnipeResult struct {
	Course *view.CourseView
}

var (
	ErrRemoveSnipeInvalid = errors.New("invalid remove snipe request")
)

func (s *SnipeService) Remove(req RemoveSnipeRequest) (*RemoveSnipeResult, error) {
	usr, err := s.userRepository.Get(req.UserID)
	if err != nil {
		return nil, err
	}

	var cmp *scope.Campus
	if req.Campus != nil {
		parsed, err := scope.ParseCampus(*req.Campus)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrRemoveSnipeInvalid, err)
		}
		cmp = &parsed
	} else {
		cmp = usr.DefaultCampus()
	}

	var szn *scope.Season
	if req.Season != nil {
		parsed, err := scope.ParseSeason(*req.Season)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrRemoveSnipeInvalid, err)
		}
		szn = &parsed
	}

	scp, err := s.activeScope.Resolve(cmp, szn)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRemoveSnipeInvalid, err)
	}

	crs, err := s.courseRepository.Get(req.Index, scp)
	if err != nil {
		if errors.Is(err, course.ErrCourseNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrRemoveSnipeInvalid, err)
		}
		return nil, err
	}

	snp, err := s.snipeRepository.Get(req.UserID, req.Index, scp)
	if err != nil {
		if errors.Is(err, snipe.ErrSnipeNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrRemoveSnipeInvalid, err)
		}
		return nil, err
	}

	if err = s.snipeRepository.Delete(snp); err != nil {
		return nil, err
	}
	s.snipeCache.Remove(snp)

	return &RemoveSnipeResult{Course: view.FromCourse(crs)}, nil
}

// --- Clear Snipes ---
type ClearSnipeRequest struct {
	UserID string
}

type ClearSnipeResult struct {
	Count int
}

func (s *SnipeService) Clear(req ClearSnipeRequest) (*ClearSnipeResult, error) {
	snipes, err := s.snipeRepository.ListByUser(req.UserID)
	if err != nil {
		return nil, err
	}

	if err = s.snipeRepository.DeleteByUser(req.UserID); err != nil {
		return nil, err
	}
	s.snipeCache.Clear(snipes)

	return &ClearSnipeResult{Count: len(snipes)}, nil
}

// --- Check Snipes ---
type CheckSnipeRequest struct {
	UserID string
}

type CheckSnipeResult struct {
	Courses []*view.CourseView
	Counts  []int
}

func (s *SnipeService) Check(req CheckSnipeRequest) (*CheckSnipeResult, error) {
	snipes, err := s.snipeRepository.ListByUser(req.UserID)
	if err != nil {
		return nil, err
	}

	var courses []*course.Course
	var counts []int
	for _, snp := range snipes {
		crs, err := s.courseRepository.Get(snp.Index, snp.Scope)
		if err != nil {
			return nil, err
		}
		courses = append(courses, crs)

		indexed, err := s.snipeRepository.ListByIndex(snp.Index, snp.Scope)
		if err != nil {
			return nil, err
		}
		counts = append(counts, len(indexed))
	}

	return &CheckSnipeResult{Courses: view.FromCourses(courses), Counts: counts}, nil
}
