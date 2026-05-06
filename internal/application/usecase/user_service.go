package usecase

import (
	"tsniper/internal/domain/scope"
	"tsniper/internal/domain/snipe"
	"tsniper/internal/domain/user"
)

type UserService struct {
	snipeCache      snipe.SnipeCache
	snipeRepository snipe.SnipeRepository
	userRepository  user.UserRepository
}

func NewUserService(snipeCache snipe.SnipeCache, snipeRepository snipe.SnipeRepository, userRepository user.UserRepository) *UserService {
	return &UserService{
		snipeCache:      snipeCache,
		snipeRepository: snipeRepository,
		userRepository:  userRepository,
	}
}

// --- Get User ---
type GetUserRequest struct {
	UserID string
}

type GetUserResult struct {
	User *user.User
}

func (s *UserService) Get(req GetUserRequest) (*GetUserResult, error) {
	usr, err := s.userRepository.Get(req.UserID)
	if err != nil {
		return nil, err
	}

	return &GetUserResult{User: usr}, nil
}

// --- Get All Users ---
type GetAllUsersRequest struct{}

type GetAllUsersResult struct {
	Users []*user.User
}

func (s *UserService) GetAll(req GetAllUsersRequest) (*GetAllUsersResult, error) {
	users, err := s.userRepository.GetAll()
	if err != nil {
		return nil, err
	}

	return &GetAllUsersResult{Users: users}, nil
}

// --- Set User Campus ---
type SetUserCampusRequest struct {
	UserID string
	Campus string
}

type SetUserCampusResult struct{}

func (s *UserService) SetUserCampus(req SetUserCampusRequest) (*SetUserCampusResult, error) {
	usr, err := s.userRepository.Get(req.UserID)
	if err != nil {
		return nil, err
	}

	cmp, err := scope.ParseCampus(req.Campus)
	if err != nil {
		return nil, err
	}

	usr.SetCampus(cmp)

	if err := s.userRepository.Save(usr); err != nil {
		return nil, err
	}

	return &SetUserCampusResult{}, nil
}

// --- Clear User Campus ---
type ClearUserCampusRequest struct {
	UserID string
}

type ClearUserCampusResult struct{}

func (s *UserService) ClearUserCampus(req ClearUserCampusRequest) (*ClearUserCampusResult, error) {
	usr, err := s.userRepository.Get(req.UserID)
	if err != nil {
		return nil, err
	}

	usr.ClearCampus()

	if err := s.userRepository.Save(usr); err != nil {
		return nil, err
	}

	return &ClearUserCampusResult{}, nil
}

// --- User Join ---
type UserJoinRequest struct {
	UserID string
}

type UserJoinResult struct{}

func (s *UserService) Join(req UserJoinRequest) (*UserJoinResult, error) {
	usr := user.NewUser(req.UserID)

	if err := s.userRepository.Create(usr); err != nil {
		return nil, err
	}

	return &UserJoinResult{}, nil
}

// --- User Leave ---
type UserLeaveRequest struct {
	UserID string
}

type UserLeaveResult struct{}

func (s *UserService) Leave(req UserLeaveRequest) (*UserLeaveResult, error) {
	usr, err := s.userRepository.Get(req.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.userRepository.Delete(usr); err != nil {
		return nil, err
	}

	snipes, err := s.snipeRepository.ListByUser(req.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.snipeRepository.DeleteByUser(req.UserID); err != nil {
		return nil, err
	}
	s.snipeCache.Clear(snipes)

	return &UserLeaveResult{}, nil
}
