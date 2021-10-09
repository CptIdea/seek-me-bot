package vk

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"log"
	"math/rand"
	"seek-me-bot/handlers"
	"seek-me-bot/service"
	"seek-me-bot/service/pkg"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
)

func GetNewVkHandler(vk *api.VK) handlers.Handler {
	rand.Seed(time.Now().UnixNano())
	return &vkHandler{vk: vk, status: make(map[int]string), controller: service.NewGameController()}
}

type vkHandler struct {
	controller      service.GameController
	vk              *api.VK
	status          map[int]string
	currentPetition pkg.Petition
}

func (v *vkHandler) AsyncHandleMessage(peerId int, message string) {
	switch strings.ToLower(message) {
	case "меню":
		_, err := v.vk.MessagesSend(params.NewMessagesSendBuilder().RandomID(rand.Int()).Message("меню").PeerID(peerId).Keyboard(object.NewMessagesKeyboard(false).AddRow().AddTextButton("про кого", " ", "primary").AddTextButton("кто", " ", "primary").AddRow().AddTextButton("дальше", "", "positive")).Params)
		if err != nil {
			v.logError(peerId, err)
		}
	case "start", "/start", "начать", "сначала", "старт", "сброс":
		v.status[peerId] = "start"
		_, err := v.vk.MessagesSend(params.NewMessagesSendBuilder().RandomID(rand.Int()).Message("Когда будешь готов,отправь мне имя, про кого хочешь написать").PeerID(peerId).Params)
		if err != nil {
			v.logError(peerId, err)
		}
	case "сброс игры":
		err := v.controller.ResetGame()
		if err != nil {
			v.logError(peerId, err)
		}
	case "про кого":
		_, err := v.vk.MessagesSend(params.NewMessagesSendBuilder().RandomID(rand.Int()).Message(v.currentPetition.Answer).PeerID(peerId).Params)
		if err != nil {
			v.logError(peerId, err)
		}
	case "кто":
		_, err := v.vk.MessagesSend(params.NewMessagesSendBuilder().RandomID(rand.Int()).Message(v.currentPetition.Author).PeerID(peerId).Params)
		if err != nil {
			v.logError(peerId, err)
		}
	case "дальше":
		petition, err := v.controller.GetPetition()
		if err != nil {
			if err == service.EmptyGameError {
				_, err := v.vk.MessagesSend(params.NewMessagesSendBuilder().RandomID(rand.Int()).Message("Больше нет описаний").PeerID(peerId).Params)
				if err != nil {
					v.logError(peerId, err)
				}
			} else {
				v.logError(peerId, err)
			}
			break
		}

		v.currentPetition = petition
		_, err = v.vk.MessagesSend(params.NewMessagesSendBuilder().RandomID(rand.Int()).Message(strings.Join(v.currentPetition.Words, ", ")).PeerID(peerId).Params)
		if err != nil {
			v.logError(peerId, err)
		}
	default:
		switch v.status[peerId] {
		case "start":
			v.status[peerId] = message
			_, err := v.vk.MessagesSend(params.NewMessagesSendBuilder().RandomID(rand.Int()).Message("Отлично. Теперь напиши про него ТРИ слова через пробел.").Attachment("photo-207762396_457239018").PeerID(peerId).Params)
			if err != nil {
				v.logError(peerId, err)
			}
		default:
			if v.status[peerId] == "" {
				v.AsyncHandleMessage(peerId, "start")
			}
			users, err := v.vk.UsersGet(params.NewUsersGetBuilder().UserIDs([]string{fmt.Sprint(peerId)}).Params)
			if err != nil || len(users) == 0 {
				v.logError(peerId, err)
				_, err := v.vk.MessagesSend(params.NewMessagesSendBuilder().RandomID(rand.Int()).Message("Что-то пошло не так, попробуй ещё раз").PeerID(peerId).Params)
				if err != nil {
					v.logError(peerId, err)
				}
			}

			err = v.controller.AddPetition(pkg.Petition{
				Author: users[0].FirstName,
				Words:  strings.Split(message, " "),
				Answer: v.status[peerId],
			})
			if err != nil {
				v.logError(peerId, err)
				_, err := v.vk.MessagesSend(params.NewMessagesSendBuilder().RandomID(rand.Int()).Message("Что-то пошло не так, попробуй ещё раз").PeerID(peerId).Params)
				if err != nil {
					v.logError(peerId, err)
				}
			}

			v.AsyncHandleMessage(peerId, "start")
		}
	}
}

func (v vkHandler) logError(peer int, err error) {
	log.Printf("[peer#%d] %s", peer, err)
}
