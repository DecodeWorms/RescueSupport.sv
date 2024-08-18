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

func (u Users) UserSignUp(data model.UserSignUp) error {
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

func (u Users) ChangePassword(data model.UserChangePassword) error {
	//Check if the user's record exist
	us, err := u.store.GetUserByEmail(data.Email)
	if err != nil {
		return fmt.Errorf("error getting user's record %v", err)
	}

	//Compare new password and confirm new password
	if data.NewPassword == data.ConfirmNewPassword {
		return fmt.Errorf("error newPassword is not equal to confirmNewPassword")
	}

	//Hash both newPassword
	pwd, err := u.encrypt.HashPassword(data.NewPassword)
	if err != nil {
		return fmt.Errorf("error encrypting password")
	}

	//Update user's password and confirmPassword
	rec := &model.Users{
		Password:        pwd,
		ConfirmPassword: pwd,
	}

	if err := u.store.UpdateUser(us.ID, rec); err != nil {
		return fmt.Errorf("error updating user's password")
	}

	//Pass the user email to Kafka to mail user that password changed successfully.

	return nil
}

func (u Users) UpdatePassword(data model.UserUpdatePassword) error {
	//Check if it is a registered user
	us, err := u.store.GetUserByEmail(data.Email)
	if err != nil {
		return fmt.Errorf("error user's record does not exist")
	}

	//Compare oldPassword and hashPassword
	b, err := u.encrypt.CompareHashAndPassword(us.Password, data.OldPassword)
	if !b {
		return fmt.Errorf("error oldPassword is incorrect %v", err)
	}

	//Compare newPassword and confirmPassword
	if data.NewPassword == data.ConfirmNewPassword {
		return fmt.Errorf("error comparing newPassword and confirmNewPassword %v", data.NewPassword)
	}

	//Hash the newPassword
	pwd, err := u.encrypt.HashPassword(data.NewPassword)
	if err != nil {
		return fmt.Errorf("error encrypting password %v", err)
	}

	//Update the user's password and confirmPassword records
	rec := &model.Users{
		Password:        pwd,
		ConfirmPassword: pwd,
	}

	if err := u.store.UpdateUser(us.ID, rec); err != nil {
		return fmt.Errorf("error updating user's record %v", err)
	}

	//Send user's email to Kafka to send mail to the user

	return nil

}
