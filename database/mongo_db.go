package database

import (
	"context"
	"log"
	"time"

	"RescueSupport.sv/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	supporterCollection = "supporters"
)

type Mongodb struct {
	Client       *mongo.Client
	databaseName string
}

func NewMongo(address, url string) (DataStore, *mongo.Client, error) {
	log.Println("Connecting to Mongodb store")

	//Config the datastore environment
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//Connect to the mongodb client
	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, nil, err
	}

	//Ping the database to know if it was connected
	if err := cli.Ping(ctx, readpref.Primary()); err != nil {
		return nil, nil, err
	}

	log.Println("Connected to Mongodb successfully")
	return &Mongodb{Client: cli, databaseName: address}, cli, nil
}

// CreateCompany creates or record Company's record to the DB
func (repo Mongodb) CreateCompany(data *model.Company) error {
	_, err := repo.col(supporterCollection).InsertOne(context.Background(), data)
	if err != nil {
		return err
	}
	return nil
}

// CreateUser persists new user records to the DB
func (repo Mongodb) CreateUser(data *model.Users) error {
	_, err := repo.col(supporterCollection).InsertOne(context.Background(), data)
	if err != nil {
		return err
	}
	return nil
}

// UpdateCompany updates the company's record
func (repo Mongodb) UpdateCompany(ID string, data *model.Company) error {
	filter := bson.M{"id": ID}

	//Fetch existing records using filter
	var old *model.Company
	if err := repo.col(supporterCollection).FindOne(context.Background(), filter).Decode(&old); err != nil {
		return err
	}

	//Replace the new record with the existing one
	_, err := repo.col(supporterCollection).ReplaceOne(context.Background(), filter, buildCompany(old, data))
	if err != nil {
		return err
	}
	return nil
}

func (repo Mongodb) UpdateUser(ID string, data *model.Users) error {
	filter := bson.M{"id": ID}
	//Get existing record using filter
	var old *model.Users
	if err := repo.col(supporterCollection).FindOne(context.Background(), filter).Decode(&old); err != nil {
		return err
	}

	_, err := repo.col(supporterCollection).ReplaceOne(context.Background(), filter, buildUser(old, data))
	if err != nil {
		return err
	}
	return nil
}

func (repo Mongodb) GetUserByID(ID string) (*model.Users, error) {
	filter := bson.M{"id": ID}
	var user *model.Users
	if err := repo.col(supporterCollection).FindOne(context.Background(), filter).Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (repo Mongodb) GetUserByEmail(email string) (*model.Users, error) {
	filter := bson.M{"email": email}
	var user *model.Users
	if err := repo.col(supporterCollection).FindOne(context.Background(), filter).Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (repo Mongodb) GetCompanyByName(name string) (*model.Company, error) {
	filter := bson.M{"name": name}
	var cmp *model.Company
	if err := repo.col(supporterCollection).FindOne(context.Background(), filter).Decode(&cmp); err != nil {
		return nil, err
	}
	return cmp, nil
}

func (repo Mongodb) GetCompanyByID(ID string) (*model.Company, error) {
	filter := bson.M{"id": ID}
	var cmp *model.Company
	if err := repo.col(supporterCollection).FindOne(context.Background(), filter).Decode(&cmp); err != nil {
		return nil, err
	}
	return cmp, nil
}

func (repo Mongodb) GetCompanyByEmail(email string) (*model.Company, error) {
	filter := bson.M{"email": email}
	var cmp *model.Company
	if err := repo.col(supporterCollection).FindOne(context.Background(), filter).Decode(&cmp); err != nil {
		return nil, err
	}
	return cmp, nil
}

func (repo Mongodb) col(collectionName string) *mongo.Collection {
	return repo.Client.Database(repo.databaseName).Collection(collectionName)
}

func buildCompany(old, new *model.Company) *model.Company {
	if new == nil {
		return old
	}
	if new.Email != "" {
		old.Email = new.Email
	}
	if new.Address.Name != "" {
		old.Address.Name = new.Address.Name
	}
	if new.Address.Code != "" {
		old.Address.Code = new.Address.Code
	}
	if new.Address.Country != "" {
		old.Address.Code = new.Address.Code
	}
	if new.Address.City != "" {
		old.Address.City = new.Address.City
	}
	if new.Password != "" {
		old.Password = new.Password
	}
	if new.ConfirmPassword != "" {
		old.ConfirmPassword = new.ConfirmPassword
	}
	if new.PhoneNumber != "" {
		old.PhoneNumber = new.PhoneNumber
	}
	if new.NumberOfEmployees != 0 {
		old.NumberOfEmployees = new.NumberOfEmployees
	}
	return old
}

func buildUser(old, new *model.Users) *model.Users {
	if new == nil {
		return old
	}

	if new.FirstName != "" {
		old.FirstName = new.FirstName
	}
	if new.LastName != "" {
		old.LastName = new.LastName
	}
	if new.Email != "" {
		old.Email = new.Email
	}
	if new.Gender != "" {
		old.Gender = new.Gender
	}
	if new.Age != 0 {
		old.Age = new.Age
	}
	if new.Password != "" {
		old.Password = new.Password
	}
	if new.ConfirmPassword != "" {
		old.ConfirmPassword = new.ConfirmPassword
	}
	if new.Address.Code != "" {
		old.Address.Code = new.Address.Code
	}
	if new.Address.Name != "" {
		old.Address.Name = new.Address.Name
	}
	if new.Address.City != "" {
		old.Address.City = new.Address.City
	}
	if new.Address.Country != "" {
		old.Address.Country = new.Address.Country
	}
	return old
}

var _ DataStore = &Mongodb{}
