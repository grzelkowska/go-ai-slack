package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
	witai "github.com/wit-ai/wit-go/v2"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Evnets")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}
func main() {
	godotenv.Load(".env")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	client := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))

	go printCommandEvents(bot.CommandEvents())

	bot.Command("query for bot - <message>", &slacker.CommandDefinition{
		Description: "send andy question to wolfram",
		Examples: []string{"who is the president of usa"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			query := request.Param("message")
			fmt.Println(query)
			msg, _ := client.Parse(&witai.MessageRequest{
				Query : query,
			})
			data, _ := json.MarshalIndent(msg, "", "    ")
			rough := string(data[:])
			fmt.Println(rough)
			response.Reply("received")

		}, 
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
