 package model

// Mongodb document model

 type Users struct {
	ID string `json:"id" bson:"id"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName string `json:"last_name" bson:"last_name"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	ConfirmPassword string `json:"confirm_password" bson:"confirm_password"`
	Age int  `json:"age" bson:"age"`
	Gender string `json:"gender" bson:"gender"`
	Address UserAddress `json:"address" bson:"address"`
}

type UserAddress struct{
	Name string `json:"name" bson:"name"`
	City string `json:"city" bson:"city"`
	Country string `json:"country" bson:"country"`
	Code string `json:"code" bson:"code"`
}

//Request body models

type UserSignUp struct {
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	ConfirmPassword string `json:"confirm_password" bson:"confirm_password"`
}

type UserKyc struct{
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Age int `json:"age"`
	Gender string `json:"gender"`
	Address UserAddress `json:"address"`
}

type UserLogin struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

 type UserChangePassword struct{
	 Email string `json:"email"`
	 NewPassword string `json:"new_password"`
	 ConfirmNewPassword string `json:"confirm_new_password"`
 }

 type UserUpdatePassword struct {
	 Email string `json:"email"`
	 OldPassword string `json:"old_password"`
	 NewPassword string `json:"new_password"`
	 ConfirmNewPassword string `json:"confirm_new_password"`
 }


