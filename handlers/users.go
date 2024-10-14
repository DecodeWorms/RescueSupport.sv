package handlers

import (
	"fmt"

	"RescueSupport.sv/database"
	"RescueSupport.sv/encrypt"
	"RescueSupport.sv/idgenerator"
	"RescueSupport.sv/model"
)

type Users struct {
	store   database.DataStore
	encrypt encrypt.Encryptor
	idGen   idgenerator.IdGenerator
	//pub     *producer.KafkaProducer
}

func NewUsers(store database.DataStore, idGen idgenerator.IdGenerator, encrypt encrypt.Encryptor) Users {
	return Users{
		store:   store,
		idGen:   idGen,
		encrypt: encrypt,
		//pub:     p,
	}

}

func (u Users) UserSignUp(data model.UserSignUp) error {
	//Check if the email is already exist
	us, _ := u.store.GetUserByEmail(data.Email)
	if us != nil {
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
	/*if err := u.pub.SendMessage("verify-email", []byte("user-email"), []byte(rec.Email)); err != nil {
		return fmt.Errorf("error publishing user email")
	}
	*/

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

func (u Users) ChangePassword(ID string, data model.UserChangePassword) error {
	//Check if the user's record exist
	us, err := u.store.GetUserByID(ID)
	if err != nil {
		return fmt.Errorf("error getting user's record %v", err)
	}

	//Compare new password and confirm new password
	if data.NewPassword != data.ConfirmNewPassword {
		return fmt.Errorf("error newPassword is not equal to confirmNewPassword %v and %v", data.NewPassword, data.ConfirmNewPassword)
	}

	//Hash both newPassword
	pwd, err := u.encrypt.HashPassword(data.NewPassword)
	if err != nil {
		return fmt.Errorf("error encrypting password %v", err)
	}

	//Update user's password and confirmPassword
	rec := &model.Users{
		Password:        pwd,
		ConfirmPassword: pwd,
	}

	if err := u.store.UpdateUser(us.ID, rec); err != nil {
		return fmt.Errorf("error updating user's password %v", err)
	}

	//Pass the user email to Kafka to mail user that password changed successfully.

	return nil
}

func (u Users) UpdatePassword(ID string, data model.UserUpdatePassword) error {
	//Check if it is a registered user
	us, err := u.store.GetUserByID(ID)
	if err != nil {
		return fmt.Errorf("error user's record does not exist %v", err)
	}

	//Compare oldPassword and hashPassword
	b, err := u.encrypt.CompareHashAndPassword(us.Password, data.OldPassword)
	if !b {
		return fmt.Errorf("error oldPassword is incorrect %v", err)
	}

	//Compare newPassword and confirmPassword
	if data.NewPassword != data.ConfirmNewPassword {
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

func (u Users) StoreRecordWithOauth(data model.SignUpWithOauth) error {
	//Check if the email already exist
	_, err := u.store.GetUserByEmail(data.Email)
	if err == nil {
		return fmt.Errorf("user's record already exist %v", err)
	}

	//Persist the record into the DB
	d := &model.Users{
		ID:        u.idGen.Generate(),
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
	}
	if err := u.store.CreateUser(d); err != nil {
		return fmt.Errorf("error creating data using Oauth %v", err)
	}
	return nil
}

func (u Users) CompleteKyc(ID string, data model.UserKyc) error {
	//Check if the user's record exist
	_, err := u.store.GetUserByID(ID)
	if err != nil {
		return fmt.Errorf("user's record not exist %v", err)
	}

	//Update the user's record
	add := model.UserAddress{
		Name:    data.Address.Name,
		City:    data.Address.City,
		Country: data.Address.Country,
		Code:    data.Address.Code,
	}
	rec := &model.Users{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Age:       data.Age,
		Gender:    data.Gender,
		Address:   add,
	}
	if err := u.store.UpdateUser(ID, rec); err != nil {
		return fmt.Errorf("error updating user's record %v", err)
	}
	return nil
}

func (u Users) UpdateKyc(ID string, data model.UserKyc) error {
	//Check if the user's record exist
	_, err := u.store.GetUserByID(ID)
	if err != nil {
		return fmt.Errorf("user's record not exist %v", err)
	}

	//Update the user's record
	add := model.UserAddress{
		Name:    data.Address.Name,
		City:    data.Address.City,
		Country: data.Address.Country,
		Code:    data.Address.Code,
	}
	rec := &model.Users{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Age:       data.Age,
		Gender:    data.Gender,
		Address:   add,
	}
	if err := u.store.UpdateUser(ID, rec); err != nil {
		return fmt.Errorf("error updating user's record %v", err)
	}
	return nil
}

func (u Users) ViewProfile(ID string) (*model.Users, error) {
	//Check if a user's record already exist using ID
	res, err := u.store.GetUserByID(ID)
	if err != nil {
		return nil, fmt.Errorf("user's record not exist %v", err)
	}

	add := model.UserAddress{
		Name:    res.Address.Name,
		City:    res.Address.City,
		Country: res.Address.Country,
		Code:    res.Address.Code,
	}

	result := &model.Users{
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Email:     res.Email,
		Age:       res.Age,
		Gender:    res.Gender,
		Address:   add,
	}
	return result, nil
}
