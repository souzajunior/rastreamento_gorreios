package handlers

import (
	"errors"
	"fmt"
	"rastreamento_gorreios/botapi"
	"rastreamento_gorreios/httpService"
	"rastreamento_gorreios/models/http"
	"rastreamento_gorreios/models/user_track"

	"go.mongodb.org/mongo-driver/mongo"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// MessageHandler is responsible to handle all messages that are received from updates
func MessageHandler(tBot *botapi.TrackBot, u tgbotapi.Update, mgClient *mongo.Client, errChan chan error) error {
	command, msg := u.Message.Command(), u.Message.CommandArguments()
	if command == "atualizar" {
		if msg == "" {
			tBot.Send(u, "Informe o código de rastreio, exemplo: /atualizar OLX01BSCEWQBR", true)
			return nil
		}

		// trackData is responsible to create a instance of track data object
		var trackData = &http.TrackRequest{
			Code: msg,
		}

		// Getting data from API
		data, err := httpService.MakeRequest(trackData)
		if err != nil {
			tBot.Send(u, "Falha ao realizar a consulta! :(", true)
			return err
		}

		if len(data.Eventos) == 0 || len(data.Eventos[0].SubStatus) < 1 {
			tBot.Send(u, "Falha ao obter dados do pacote! :(", true)
			return errors.New("dados retornados são inválidos! (talvez apenas foi postado)")
		}

		evt := data.Eventos[0]
		botMsg := fmt.Sprintf("Última atualização (%s):\n\nData global: %s\nData status: %s\nHora: %s\nStatus: %s\nDetalhe: %s\n",
			data.Codigo, data.Ultimo.Format("02/01/2006"), evt.Data, evt.Hora, evt.Status, evt.SubStatus[0])
		tBot.Send(u, botMsg, true)
	} else if command == "codigos" {
		var user = user_track.New(u.Message.From, nil, nil)
		tracks, err := user.GetTracks(mgClient)
		if err != nil {
			tBot.Send(u, "Falha ao consultar seus códigos no banco de dados. :(", true)
			return err
		}

		if len(tracks) == 0 {
			tBot.Send(u, "Você não possui nenhum código cadastrado para o acompanhamento automático, utilize /acompanhe {código}", true)
			return nil
		}

		var msgReply = "Seus códigos:\n"
		for i := range tracks {
			var acompanhando = "Não"
			if tracks[i].Automatic {
				acompanhando = "Sim"
			}
			msgReply += *tracks[i].Code + " - Acompanhando: " + acompanhando + "\n"
		}

		tBot.Send(u, msgReply, true)
	} else if command == "acompanhe" {
		if msg == "" {
			tBot.Send(u, "Informe o código de rastreio, exemplo: /acompanhe OLX01BSCEWQBR", true)
			return nil
		}

		var (
			user  = user_track.New(u.Message.From, &msg, &u.Message.Chat.ID)
			err   error
			found bool
		)

		if found, err = user.CheckTrack(mgClient); err != nil {
			tBot.Send(u, "Falha ao tentar buscar informações do seu pacote. Tente novamente mais tarde", true)
			return nil
		}

		if found {
			tBot.Send(u, "Já existe um registro desta encomenda! Pode ser que já tenha sido entregue.", true)
			return nil
		}

		if err := user.AddTrackAuto(mgClient); err != nil {
			tBot.Send(u, "Ocorreu um erro na hora de registrar sua solicitação, tente novamente mais tarde.", true)
			return err
		}

		tBot.Send(u, "Estou agora acompanhando a encomenda "+msg, false)

		// Send our go routine to track the code
		go user.TrackAutomatic(tBot, mgClient, errChan)

	} else if command == "remova" {
		if msg == "" {
			tBot.Send(u, "Informe o código de rastreio que deseja que remova dos registros, exemplo: /remova OLX01BSCEWQBR", true)
			return nil
		}

		var (
			user  = user_track.New(u.Message.From, &msg, nil)
			found bool
			err   error
		)

		if found, err = user.CheckTrack(mgClient); err != nil {
			tBot.Send(u, "Falha ao tentar buscar informações do seu pacote. Tente novamente mais tarde", true)
			return nil
		}

		if !found {
			tBot.Send(u, "Não foi possível localizar este pacote, confira se você possui esse código ativo", true)
			return nil
		}

		if err = user.DeleteTrack(mgClient); err != nil {
			tBot.Send(u, "Falha na remoção do seu pacote, tente novamente mais tarde", true)
			return err
		}

		tBot.Send(u, "Código deletado com sucesso!", true)
		return nil
	} else if command == "comandos" {
		msg := `
/atualizar {codigo_rastreio} - Retorna o atual status da encomenda atualizado
/acompanhe {codigo_rastreio} - Solicita para que o bot fique acompanhando o status da encomenda e notifique caso ela atualize
/remova {codigo_rastreio} - Solicita que o bot remova um código do rastreio automático que já/não foi entregue
/comandos - Lista todos os comandos disponíveis
`
		tBot.Send(u, msg, true)
	} else {
		tBot.Send(u, "Comando não reconhecido, utilize /comandos para listagem de todos os comandos", true)
	}

	return nil
}
