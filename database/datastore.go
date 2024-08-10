package database

import "RescueSupport.sv/model"

//go:generate mockgen -source=datastore.go -destination=../mocks/datastore_mock.go -package=mocks
type DataStore interface {
	User
}

type User interface {
	Create(data *model.UserSignUp) error
	Update(data *model.UserKyc) error
	Login(data *model.UserLogin) error
	ChangePassword(data *model.ChangePassword) error
	UpdatePassword(data *model.UpdatePassword) error
}
