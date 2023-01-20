package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client   //ref to mongo.Client
	Db     *mongo.Database //ref to mongo.Database
}

var mg MongoInstance

type Employee struct {
	ID     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    float64 `json:"age"`
}

// connect DB func
func connectDB() error {
	//load the .env file
	err := godotenv.Load()

	//check err
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//get mongoURI value
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")

	//create a newclient connection to mongo returns an instance and err
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	//define timeout to avoid code blocking if code takes time
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//make a connection
	err = client.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {
		return err
	}

	mg = MongoInstance{
		Client: client,
		Db:     db,
	}

	return nil
}

func main() {

	//get return
	if err := connectDB(); err != nil {
		//handle error using log library
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, world!!")
	})

	app.Get("/employee", func(c *fiber.Ctx) error {

		//define query
		query := bson.D{{}}

		//send back cursor and error
		cursor, err := mg.Db.Collection("employees").Find(c.Context(), query)

		//handle err
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		//creata a slice with struct of type Employee
		var employees = make([]Employee, 0)

		//handle error
		if err := cursor.All(c.Context(), &employees); err != nil {
			c.Status(500).SendString(err.Error())
		}

		//return employees
		return c.JSON(employees)
	})

	app.Listen(":8000")
}
