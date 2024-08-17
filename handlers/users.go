package handlers

import (
	"fmt"

	"RescueSupport.sv/database"
	"RescueSupport.sv/encrypt"
	"RescueSupport.sv/idgenerator"
	"RescueSupport.sv/model"
)

type Users struct {
	store   database.Mongodb
	encrypt encrypt.Encryptor
	idGen   idgenerator.IdGenerator
}

func NewUsers(store database.Mongodb, idGen idgenerator.IdGenerator) Users {
	return Users{
		store: store,
		idGen: idGen,
	}
}

func (u Users) CompanySignUp(data model.SignUp) error {
	//Check if the email is already exist
	_, err := u.store.GetCompanyByEmail(data.Email)
	if err != nil {
		return fmt.Errorf("error company's email already exist %v", data.Email)
	}

	//Check if the password and confirm password is the same
	if data.Password != data.ConfirmPassword {
		return fmt.Errorf("error password and confirm password is not the same %v", data.ConfirmPassword)
	}

	//Encrypt the password and confirm password
	hashPassword, err := u.encrypt.HashPassword(data.Password)
	if err != nil {
		return fmt.Errorf("error encrypting password")
	}

	hashConfirmPassword, err := u.encrypt.HashPassword(data.ConfirmPassword)
	if err != nil {
		return fmt.Errorf("error encrypting password")
	}

	//Pass the Company email to messaging system Kafka

	//Persist the data to the db
	rec := &model.Company{
		ID:              u.idGen.Generate(), //Fix me
		Email:           data.Email,
		Password:        hashPassword,
		ConfirmPassword: hashConfirmPassword,
	}
	if err := u.store.CreateCompany(rec); err != nil {
		return fmt.Errorf("error creating company %v", err)
	}
	return nil
}
