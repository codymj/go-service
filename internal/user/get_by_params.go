package user

import (
	"context"
	"github.com/rs/zerolog/log"
	"go-service.codymj.io/internal/user/dao"
)

// GetByParams returns a filtered list of users using query parameters.
func (s *service) GetByParams(ctx context.Context, params map[string]string) ([]dao.User, error) {
	// Log info.
	log.Info().
		Interface("parameters", params).
		Msg(InfoBeginGetUserByParams)

	// Get all users via repository.
	users, err := s.userdao.GetByParams(ctx, params)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	// Log info.
	log.Info().
		Interface("parameters", params).
		Msg(InfoDoneGetUserByParams)

	return users, nil
}
