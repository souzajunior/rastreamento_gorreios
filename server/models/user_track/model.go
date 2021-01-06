package user_track

import (
	"context"
	"errors"
	"fmt"
	"log"
	"rastreamento_gorreios/botapi"
	"rastreamento_gorreios/httpService"
	"rastreamento_gorreios/models/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserTrack struct {
	ID         int     `json:"id"`
	Code       *string `json:"code"`
	Automatic  bool    `json:"automatic"`
	Delivered  bool    `json:"delivered"`
	ChatID     *int64  `json:"chatid"`
	LastUpdate *string `json:"lastupdate"`
}

const (
	databaseName   = "gorreios"
	collectionName = "user_track"
)

// New is responsible to initialize a new user track data
func New(user *tgbotapi.User, code *string, chatID *int64) *UserTrack {
	return &UserTrack{
		ID:         user.ID,
		Code:       code,
		Automatic:  false,
		LastUpdate: new(string),
		ChatID:     chatID,
	}
}

// GetTracks is responsible to get all order tracks filtering by an user
func (u *UserTrack) GetTracks(mgClient *mongo.Client) (userTracks []UserTrack, err error) {
	cursor, err := mgClient.
		Database(databaseName).
		Collection(collectionName).
		Find(context.Background(), bson.D{
			{"id", u.ID},
			{"delivered", false},
		})
	if err != nil {
		return
	}

	// Defering to close our cursor
	defer func() { _ = cursor.Close(context.Background()) }()

	for cursor.Next(context.Background()) {
		var userData UserTrack
		if err = cursor.Decode(&userData); err != nil {
			return
		}

		userTracks = append(userTracks, userData)
	}

	return
}

// CheckTrack is responsible to check if a track already is being tracked
func (u *UserTrack) CheckTrack(mgClient *mongo.Client) (found bool, err error) {
	cursor, err := mgClient.
		Database(databaseName).
		Collection(collectionName).
		Find(context.Background(), bson.D{
			{"id", u.ID},
			{"code", u.Code},
		})

	// Defering to close our cursor
	defer cursor.Close(context.Background())

	found = cursor.Next(context.Background())

	return
}

// AddTrackAuto is responsible to add a automatically tracking of an order
func (u *UserTrack) AddTrackAuto(mgClient *mongo.Client) (err error) {
	u.Automatic = true
	if _, err = mgClient.
		Database(databaseName).
		Collection(collectionName).
		InsertOne(
			context.Background(),
			u,
		); err != nil {
		return
	}

	return
}

// DeleteTrack is responsible to delete a track from database of an order
func (u *UserTrack) DeleteTrack(mgClient *mongo.Client) error {
	var result, err = mgClient.
		Database(databaseName).
		Collection(collectionName).
		DeleteOne(
			context.Background(),
			bson.M{
				"id":   u.ID,
				"code": u.Code,
			})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("nenhum registro foi encontrado para solicitação")
	}

	return nil
}

// ResumeTracksActives is responsible to resume the tracks that were actives before the bot was closed
// this function should be called when bot is started
func ResumeTracksActives(bot *botapi.TrackBot, mgClient *mongo.Client, errChan chan error) (err error) {
	cursor, err := mgClient.
		Database(databaseName).
		Collection(collectionName).
		Find(context.Background(), bson.D{
			{"automatic", true},
			{"delivered", false},
		})
	if err != nil {
		return
	}

	// Defering to close our cursor
	defer func() { _ = cursor.Close(context.Background()) }()

	var (
		tracksActives int
	)
	for cursor.Next(context.Background()) {
		var userData UserTrack
		if err = cursor.Decode(&userData); err != nil {
			return
		}

		go userData.TrackAutomatic(bot, mgClient, errChan)
		tracksActives++
	}

	log.Println("Actually exists", tracksActives, "tracks actives")

	return
}

// TrackAutomatic is responsible for automatically tracking an order of a user_track
func (u *UserTrack) TrackAutomatic(tBot *botapi.TrackBot, mgClient *mongo.Client, errChan chan error) {
	ticker := time.Tick(time.Duration(tBot.TrackInterval) * time.Minute)
	for t := range ticker {
		data, err := httpService.MakeRequest(&http.TrackRequest{Code: *u.Code})
		if err != nil {
			errChan <- err
			continue
		}

		if len(data.Eventos) == 0 || len(data.Eventos[0].SubStatus) < 1 {
			errChan <- errors.New("falha ao obter dados da API do rastreio: " + *u.Code)
			continue
		}

		if *u.LastUpdate != data.Ultimo.Format("02/01/2006 - 15:04") {
			evt := data.Eventos[0]
			botMsg := fmt.Sprintf("Nova atualização (%s):\n\nData global: %s\nData status: %s\nHora: %s\nStatus: %s\nDetalhe: %s\nBot sended: %s\n",
				data.Codigo, data.Ultimo.Format("02/01/2006"), evt.Data, evt.Hora, evt.Status, evt.SubStatus[0], t.Format("02/01/2006 - 15:04:05"))
			tBot.SendChat(*u.ChatID, botMsg)

			var (
				result *mongo.UpdateResult
				err    error
			)

			if result, err = mgClient.
				Database(databaseName).
				Collection(collectionName).
				UpdateOne(
					context.Background(),
					bson.M{
						"id":   u.ID,
						"code": u.Code,
					},
					bson.D{
						{"$set",
							bson.D{
								{"lastupdate", data.Ultimo.Format("02/01/2006 - 15:04")},
							},
						},
					},
				); err != nil {
				errChan <- err
				continue
			}

			if result.ModifiedCount == 0 {
				errChan <- errors.New("falha na atualização de um registro durante a verificação automática da encomenda: " + *u.Code)
				continue
			}

			lastUpdt := data.Ultimo.Format("02/01/2006 - 15:04")
			u.LastUpdate = &lastUpdt

			if evt.Status == "Objeto entregue ao destinatário" {
				if _, err = mgClient.
					Database(databaseName).
					Collection(collectionName).
					UpdateOne(
						context.Background(),
						bson.M{
							"id":   u.ID,
							"code": u.Code,
						},
						bson.D{
							{"$set",
								bson.D{
									{"delivered", true},
									{"automatic", false},
								},
							},
						},
					); err != nil {
					errChan <- err
				}

				return
			}
		}
	}
	return
}
