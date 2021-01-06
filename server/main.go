package main

import (
	"context"
	"log"
	"rastreamento_gorreios/botapi"
	"rastreamento_gorreios/database"
	"rastreamento_gorreios/handlers"
	"rastreamento_gorreios/models/user_track"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
	Rastreamento_gorreios é um projeto pequeno e feito por hobby para acompanhar encomendas dos correios (integrado com uma API pública)
	Feito em: Mar - 2020
*/

func main() {
	var (
		tBot        = new(botapi.TrackBot)
		updates     tgbotapi.UpdatesChannel
		mongoClient *mongo.Client
		err         error
		errChan     = make(chan error)
	)

	go func() {
		for {
			select {
			case err := <-errChan:
				log.Println("An error occurred during a track:", err.Error())
			}
		}
	}()

	ctx, cancel := context.WithCancel(context.TODO())

	defer cancel()

	if mongoClient, err = database.OpenMongo(ctx); err != nil {
		log.Fatal("Erro ao abrir conexão com o banco (mongo): ", err.Error())
	}

	if tBot, err = botapi.InitBot(); err != nil {
		log.Fatal("Erro inicializando o bot: ", err.Error())
	}

	if err = user_track.ResumeTracksActives(tBot, mongoClient, errChan); err != nil {
		log.Fatal("Erro ao resumir as encomendas ativas: ", err.Error())
	}

	if updates, err = tBot.Bot.GetUpdatesChan(tBot.Updates); err != nil {
		log.Fatal("Erro recebendo updates: ", err.Error())
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// check if the client exists
		// case not, register and return an instance of user

		if err = handlers.MessageHandler(tBot, update, mongoClient, errChan); err != nil {
			log.Println("Error:", err.Error())
		}
	}

	_ = mongoClient.Disconnect(ctx)
}
