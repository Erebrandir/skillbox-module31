package usecase

import "skillbox/internal/entity"

type (
	Usecase interface {
		CreateUser(*entity.User) (string, error)
		DeleteUser(string) (string, error)
		UpdateUser(string, int) error
		GetFriends(string) ([]string, error)
		MakeFriends(string, string) (string, string, error)
		GetUsers(*entity.User) []*entity.User
	}

	Repository interface {
		CreateUser(*entity.User) (string, error)
		DeleteUser(string) (string, error)
		UpdateAge(string, int) error
		GetFriends(string) ([]string, error)
		MakeFriends(string, string) (string, string, error)
		GetUsers(*entity.User) []*entity.User
	}
)

type usecase struct {
	repository Repository
}

func NewUsecase(repository Repository) *usecase {
	return &usecase{
		repository: repository,
	}
}

func (u *usecase) CreateUser(user *entity.User) (string, error) {
	uid, err := u.repository.CreateUser(user)
	return uid, err
}

func (u *usecase) DeleteUser(id string) (string, error) {
	name, err := u.repository.DeleteUser(id)
	return name, err
}

func (u *usecase) GetUsers(user *entity.User) []*entity.User {
	allUsers := u.repository.GetUsers(user)
	return allUsers
}

func (u *usecase) UpdateUser(id string, newAge int) error {
	err := u.repository.UpdateAge(id, newAge)
	return err
}

func (u *usecase) MakeFriends(target string, source string) (string, string, error) {
	name1, name2, err := u.repository.MakeFriends(target, source)
	return name1, name2, err
}

func (u *usecase) GetFriends(userId string) ([]string, error) {
	allUsers, err := u.repository.GetFriends(userId)
	return allUsers, err
}
