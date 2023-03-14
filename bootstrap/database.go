package bootstrap

import (
	"context"
	"fmt"
	"go-mirayway/mongodbImplement"
	"log"
	"time"
)

func NewMongoClient(env *Env) mongodbImplement.Client {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	databaseUrl := fmt.Sprintf("mongodb://%v:%v/?replicaSet=rs0&directConnection=true", env.DBHost, env.DBPort)
	if env.DBUser != "" {
		databaseUrl = fmt.Sprintf("mongodb://%v:%v@%v:%v/?replicaSet=rs0&directConnection=true", env.DBUser, env.DBPass, env.DBHost, env.DBPort)
	}

	client, err := mongodbImplement.NewClient(databaseUrl)
	if err != nil {
		log.Fatalln("Can't connect to database", err)
	}

	if err := client.Connect(ctx); err != nil {
		log.Fatalln(err)
	}

	if err := client.Ping(ctx); err != nil {
		log.Fatalln(err)
	}

	return client
}

func CLoseMongoClient(client mongodbImplement.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Disconnect to database")
}
