package usecase

import (
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
	Course *course.Course
}

func (s *SnipeService) Add(req AddSnipeRequest) (*AddSnipeResult, error) {
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
		parsed, err := s.activeScope.ParseSeason(*req.Season)
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

	snp := snipe.NewSnipe(req.UserID, req.Index, scp)
	if err := s.snipeRepository.Create(snp); err != nil {
		return nil, err
	}
	s.snipeCache.Add(snp)

	return &AddSnipeResult{Course: crs}, nil
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
	Course *course.Course
}

func (s *SnipeService) ReAdd(req ReAddSnipeRequest) (*ReAddSnipeResult, error) {
	campus, err := scope.ParseCampus(req.Campus)
	if err != nil {
		return nil, err
	}

	term, err := scope.ParseTerm(req.Term)
	if err != nil {
		return nil, err
	}

	scp := scope.AcademicScope{Campus: campus, Term: term, Year: req.Year}
	if err := s.activeScope.Validate(scp); err != nil {
		return nil, err
	}

	crs, err := s.courseRepository.Get(req.Index, scp)
	if err != nil {
		return nil, err
	}

	snp := snipe.NewSnipe(req.UserID, req.Index, scp)
	if err := s.snipeRepository.Create(snp); err != nil {
		return nil, err
	}
	s.snipeCache.Add(snp)

	return &ReAddSnipeResult{Course: crs}, nil
}

// --- Remove Snipe ---
type RemoveSnipeRequest struct {
	UserID string
	Index  string
	Campus *string
	Season *string
}

type RemoveSnipeResult struct {
	Course *course.Course
}

func (s *SnipeService) Remove(req RemoveSnipeRequest) (*RemoveSnipeResult, error) {
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
		parsed, err := s.activeScope.ParseSeason(*req.Season)
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

	snp, err := s.snipeRepository.Get(req.UserID, req.Index, scp)
	if err != nil {
		return nil, err
	}

	if err = s.snipeRepository.Delete(snp); err != nil {
		return nil, err
	}
	s.snipeCache.Remove(snp)

	return &RemoveSnipeResult{Course: crs}, nil
}

// --- Clear Snipes ---
type ClearSnipeRequest struct {
	UserID string
}

type ClearSnipeResult struct {
	Count int
}

func (s *SnipeService) Clear(req ClearSnipeRequest) (*ClearSnipeResult, error) {
	snipes, err := s.snipeRepository.GetByUser(req.UserID)
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
	Courses []*course.Course
}

func (s *SnipeService) Check(req CheckSnipeRequest) (*CheckSnipeResult, error) {
	snipes, err := s.snipeRepository.GetByUser(req.UserID)
	if err != nil {
		return nil, err
	}

	var courses []*course.Course
	for _, snp := range snipes {
		crs, err := s.courseRepository.Get(snp.Index, snp.Scope)
		if err != nil {
			return nil, err
		}
		courses = append(courses, crs)
	}

	return &CheckSnipeResult{Courses: courses}, nil
}
