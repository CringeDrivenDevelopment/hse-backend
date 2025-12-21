package service

import (
	"context"

	"go.uber.org/zap"
)

type ParticipantService struct {
	logger *zap.Logger
}

func NewParticipantService(logger *zap.Logger) *ParticipantService {
	return &ParticipantService{logger: logger}
}

func (s *ParticipantService) Add(ctx context.Context, chatID int64, userIDs []int64) error {
	s.logger.Info("adding participants", zap.Int64("chatID", chatID), zap.Int64s("userIDs", userIDs))
	return nil
}

func (s *ParticipantService) Remove(ctx context.Context, chatID int64, userID int64) error {
	s.logger.Info("removing participant", zap.Int64("chatID", chatID), zap.Int64("userID", userID))
	return nil
}
