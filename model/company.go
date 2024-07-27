package model

//Mongodb document model for Company Supporter type

type Company struct {
	ID string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	ConfirmPassword string `json:"confirm_password" bson:"confirm_password"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
	Address Address `json:"address" bson:"address"`
	NumberOfEmployees int `json:"number_of_employees" bson:"number_of_employees"`
}

type Address struct{
	Name string `json:"name" bson:"name"`
	City string `json:"city" bson:"city"`
	Country string `json:"country" bson:"country"`
	Code string `json:"code" bson:"code"`
}

// Request body models

type SignUp struct{
	Email string `json:"email"`
	Password string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type Kyc struct {
	PhoneNumber string `json:"phone_number"`
	NumberOfEmployees int `json:"number_of_employees"`
	Address Address `json:"address"`
}

type Login struct{
	Email string `json:"email"`
	Password string `json:"password"`
}