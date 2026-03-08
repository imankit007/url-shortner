package repository

import "github.com/imankit007/url-shortner-auth-service/model"

type InMemoryUserRepository struct {
	usersByEmail map[string]model.User
}

func NewInMemoryUserRepository() UserRepository {
	users := []model.User{
		{
			UserID:       "user-1",
			TenantID:     "tenant-1",
			Email:        "owner@tenant1.local",
			PasswordHash: "tenant1-password",
		},
		{
			UserID:       "user-2",
			TenantID:     "tenant-2",
			Email:        "owner@tenant2.local",
			PasswordHash: "tenant2-password",
		},
	}

	usersByEmail := make(map[string]model.User, len(users))
	for _, user := range users {
		usersByEmail[user.Email] = user
	}

	return &InMemoryUserRepository{
		usersByEmail: usersByEmail,
	}
}

func (r *InMemoryUserRepository) FindByEmail(email string) (model.User, error) {
	user, found := r.usersByEmail[email]
	if !found {
		return model.User{}, ErrUserNotFound
	}

	return user, nil
}
