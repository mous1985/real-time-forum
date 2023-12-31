package service

import "real-time-forum/internal/repository"

type ModeratorsService struct {
	repo repository.Moderators
}

func NewModeratorsService(repo repository.Moderators) *ModeratorsService {
	return &ModeratorsService{
		repo: repo,
	}
}
