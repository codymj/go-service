package user

import (
	"context"
	"github.com/rs/zerolog/log"
	"go-service.codymj.io/internal/user/dao"
	"strconv"
)

// GetById returns a single user by ID.
func (s *service) GetById(ctx context.Context, id int64) (dao.User, error) {
	// Log info.
	log.Info().
		Str("id", strconv.Itoa(int(id))).
		Msg(InfoBeginGetUserById)

	// Get users via repository.
	user, err := s.userdao.GetById(ctx, id)
	if err != nil {
		log.Err(err)
		return dao.User{}, err
	}

	// Log info.
	log.Info().
		Str("id", strconv.Itoa(int(id))).
		Msg(InfoDoneGetUserById)

	return user, nil
}
