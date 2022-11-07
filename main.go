package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/shomali11/slacker"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println(("Command Events"))
		fmt.Println(event.Timestamp)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {
	os.Setenv("SLACK_BOT_TOKEN", "add here")
	os.Setenv("SLACK_APP_TOKEN", "add here")
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	go printCommandEvents(bot.CommandEvents())

	bot.Command("my yob is <year>", &slacker.CommandDefinition{
		Description: "yob calculator",
		Examples:    []string{"yob is 2020"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				response.Reply("Invalid Input, try again")
				return
			}
			currentYear := time.Now().Year()
			age := currentYear - yob
			var r string
			if age >= 0 {
				r = fmt.Sprintf("age is %d", age)
			} else {
				r = "you have not been born yet!"
			}
			response.Reply(r)
		},
	})
	bot.Command("add <n1> + <n2>", &slacker.CommandDefinition{
		Description: "adds 2 numbers",
		Examples:    []string{"5+7"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			num1, err1 := strconv.Atoi(request.Param("n1"))
			num2, err2 := strconv.Atoi(request.Param("n2"))
			if err1 != nil || err2 != nil {
				response.Reply("Invalid Input, try again")
				fmt.Printf("n1 = %v\n", request.Param("n1"))
				fmt.Printf("n2 = %v\n", request.Param("n2"))
				return
			}
			r := strconv.Itoa(num1 + num2)
			response.Reply(r)
		},
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
