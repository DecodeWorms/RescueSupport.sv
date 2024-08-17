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

func (u Users) CompanySignUp(data model.UserSignUp) error {
	//Check if the email is already exist
	_, err := u.store.GetUserByEmail(data.Email)
	if err != nil {
		return fmt.Errorf("error user's email already exist %v", data.Email)
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

	//Pass the User email to messaging system Kafka

	//Persist the data to the db
	rec := &model.Users{
		ID:              u.idGen.Generate(), //Fix me
		Email:           data.Email,
		Password:        hashPassword,
		ConfirmPassword: hashConfirmPassword,
	}
	if err := u.store.CreateUser(rec); err != nil {
		return fmt.Errorf("error creating user %v", err)
	}
	return nil
}

func (u Users) Login(data model.UserLogin) error {
	//Check if it is a registered user
	us, err := u.store.GetUserByEmail(data.Email)
	if err != nil {
		return fmt.Errorf("error user's record not found")
	}

	//Compare the hash and non hash password to login a user
	b, err := u.encrypt.CompareHashAndPassword(us.Password, data.Password)
	if !b {
		return fmt.Errorf("error invalid password %v", err)
	}

	//User email and password is correct
	return nil
}
