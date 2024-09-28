package serverutil

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"RescueSupport.sv/config"
	"RescueSupport.sv/database"
	"RescueSupport.sv/encrypt"
	"RescueSupport.sv/handlers"
	"RescueSupport.sv/idgenerator"
	"RescueSupport.sv/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUpDatabase(url, name string) (database.DataStore, *mongo.Client) {
	repo, client, err := database.NewMongo(name, url)
	if err != nil {
		log.Fatalf("Error failed to open MongoDB: %v", err)
	}
	return repo, client
}

/*func SetUpKakifyHandler(brokers []string) *producer.KafkaProducer {
	p, err := producer.NewKafkaProducer(brokers)
	if err != nil {
		log.Fatal("Error failed to start kakify")
	}
	return p
}
*/

func SetUpHandler(store database.DataStore) handlers.Users {
	return handlers.NewUsers(store, idgenerator.New(), encrypt.NewPasswordEncryptor())
}

func SetUpServer(userHandler handlers.Users) server.User {
	return server.NewUser(userHandler)
}

func SetupRouter(server *server.User) *gin.Engine {
	router := gin.Default()

	// Add Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	//List your API endpoints here..
	router.POST("/user", server.SignUp())
	router.POST("/user/login", server.Login())
	router.PUT("/user/change_password", server.ChangePassword())
	router.PUT("/user/update_password", server.UpdatePassword())

	//Oauth Api endpoints
	router.GET("/", server.OauthPage())
	router.GET("/user/oauth_login", server.LoginWithOauth())
	router.GET("/user/oauth_redirect", server.GoogleRedirect())
	return router
}

func StartServer(router *gin.Engine, client *mongo.Client) {
	//var c config.Config
	var c = config.ImportConfig(config.OSSource{})
	interruptHandler := make(chan os.Signal, 1)
	signal.Notify(interruptHandler, syscall.SIGTERM, syscall.SIGINT)

	addr := fmt.Sprintf(":%s", c.ServicePort)
	go func(addr string) {
		log.Println(fmt.Sprintf("RescueSupport.sv API service running on %v. Environment=%s", addr, c.AppEnv))
		if err := http.ListenAndServe(addr, router); err != nil {
			log.Printf("Error starting server: %v", err)
		}
	}(addr)

	<-interruptHandler
	log.Println("Closing application...")
	if err := client.Disconnect(context.Background()); err != nil {
		log.Fatal("Failed to disconnect from database")
	}
}
