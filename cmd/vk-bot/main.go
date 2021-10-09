package main

import (
	"context"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
	"os"
	"seek-me-bot/handlers/vk"
	"strconv"
)

func main() {
	token := os.Getenv("VK_TOKEN")
	if token == "" {
		log.Fatal("нет VK_TOKEN")
	}

	groupId, err := strconv.Atoi(os.Getenv("VK_GROUP_ID"))
	if err != nil {
		log.Fatal("нет VK_GROUP_ID")
	}

	bot := api.NewVK(token)
	handler := vk.GetNewVkHandler(bot)

	lp, err := longpoll.NewLongPoll(bot, groupId)
	if err != nil {
		panic(err)
	}

	lp.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		handler.AsyncHandleMessage(obj.Message.PeerID, obj.Message.Text)
	})

	log.Println("Запуск лонгпул обработчика")
	log.Fatal(lp.Run())
}
