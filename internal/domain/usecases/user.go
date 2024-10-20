package usecases

import (
	"net/url"
	"seeker/internal/domain/entities"
	errs "seeker/internal/domain/errors"
	"seeker/internal/domain/repositories"
)

type UserUsecase interface {
	FindUser(url.Values) (entities.User, error)
}

type userUsecase struct {
	userRepository repositories.UserRepository
}

func NewUserUsecase(userRepository repositories.UserRepository) UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
	}
}

func (u userUsecase) FindUser(values url.Values) (entities.User, error) {
	userId := values.Get("id")
	var user entities.User

	if userId != "" {
		user, err := u.userRepository.FindByID(userId)
		if err != nil {
			return user, errs.ErrUserDoesNotExist
		}
		return user, nil
	}

	userEmail := values.Get("email")

	if userEmail != "" {
		user, err := u.userRepository.FindByEmail(userEmail)
		if err != nil {
			return user, errs.ErrUserDoesNotExist
		}

		return user, nil
	}

	return user, errs.ErrNoParams
}
