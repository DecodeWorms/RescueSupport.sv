package database

import (
	"RescueSupport.sv/model"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
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

func (repo Mongodb) Create(data *model.UserSignUp) error {
	//TODO implement me
	panic("implement me")
}

func (repo Mongodb) Update(data *model.UserKyc) error {
	//TODO implement me
	panic("implement me")
}

func (repo Mongodb) Login(data *model.UserLogin) error {
	//TODO implement me
	panic("implement me")
}

func (repo Mongodb) ChangePassword(data *model.ChangePassword) error {
	//TODO implement me
	panic("implement me")
}

func (repo Mongodb) UpdatePassword(data *model.UpdatePassword) error {
	//TODO implement me
	panic("implement me")
}

func (repo Mongodb) col(collectionName string) *mongo.Collection {
	return repo.Client.Database(repo.databaseName).Collection(collectionName)
}

var _ DataStore = &Mongodb{}
