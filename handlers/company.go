package handlers

import (
	"fmt"

	"RescueSupport.sv/database"
	"RescueSupport.sv/encrypt"
	"RescueSupport.sv/idgenerator"
	"RescueSupport.sv/model"
)

type Company struct {
	store   database.DataStore
	encrypt encrypt.Encryptor
	idGen   idgenerator.IdGenerator
}

func NewCompany(s database.DataStore, encrypt encrypt.Encryptor, idGen idgenerator.IdGenerator) Company {
	return Company{
		store:   s,
		encrypt: encrypt,
		idGen:   idGen,
	}
}

func (c Company) SignUp(data model.SignUp) error {
	//Ensure that the newly registered company is not already exist
	_, err := c.store.GetCompanyByEmail(data.Email)
	if err == nil {
		return fmt.Errorf("Company detail already exist %v", err)
	}

	rec := &model.Company{
		ID:              c.idGen.Generate(),
		Email:           data.Email,
		Password:        data.Password,
		ConfirmPassword: data.ConfirmPassword,
	}

	//Persist the Company record to the DB
	if err := c.store.CreateCompany(rec); err != nil {
		return fmt.Errorf("cannot persist the Company record %v", err)
	}
	return nil
}
