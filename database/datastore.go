package database

import "RescueSupport.sv/model"

//go:generate mockgen -source=datastore.go -destination=../mocks/datastore_mock.go -package=mocks
type DataStore interface {
	User
}

type User interface {
	CreateCompany(data *model.Company) error
	CreateUser(data *model.Users) error
	UpdateUser(ID string, data *model.Users) error
	GetUserByID(ID string) (*model.Users, error)
	GetUserByEmail(email string) (*model.Users, error)
	UpdateCompany(ID string, data *model.Company) error
	GetCompanyByName(name string) (*model.Company, error)
	GetCompanyByEmail(email string) (*model.Company, error)
	GetCompanyByID(ID string) (*model.Company, error)
}
