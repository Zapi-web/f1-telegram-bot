package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"f1-bot/racing"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

func main() {
	_ = godotenv.Load()

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("Ошибка: Токен не задан")
	}

	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	var (
		menu = &tele.ReplyMarkup{ResizeKeyboard: true}

		btnHello = menu.Text("Привет!👋")
		btnNext  = menu.Text("Следующая гонка🏎️")

		btns = []tele.Btn{btnHello, btnNext}
	)

	menu.Reply(
		menu.Split(2, btns)...,
	)

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Привет! Я F1newsBot🤖\nЯ могу тебе помочь узнать дату следующей гонки формулы 1🏎️\nИ точное время (в днях) до неё🕰️", menu)
	})

	b.Handle(&btnHello, func(c tele.Context) error {
		userName := c.Chat().FirstName

		helloMessage := fmt.Sprintf("Здарова, %s! Я живой! 🤖", userName)

		return c.Send(helloMessage)
	})

	b.Handle(&btnNext, func(c tele.Context) error {
		race, daysBeforeRace, err := racing.GetNewRace()

		if err != nil {
			return c.Send("⚠️ Ошибка: " + err.Error())
		}

		countryEmoji, ok := racing.Countries[race.CountryCode]

		if !ok {
			countryEmoji = "🏳️"
		}

		circuitImage, ok := racing.Circuits[race.CircuitShortName]

		if !ok {
			circuitImage = "https://cdn-6.motorsport.com/images/mgl/68ey3q40/s1100/f1-abu-dhabi-gp-2017-f1-logo.webp"
		}

		answer := fmt.Sprintf("Следующая гонка пройдет в %s%s \nдата: %s \nДней до гонки: %d", race.CountryName, countryEmoji, race.DateStart, daysBeforeRace)

		if daysBeforeRace == 0 {
			answer = fmt.Sprintf("Следующая гонка пройдет в %s%s \nдата: %s \nДней до гонки: RACE DAY!!!", race.CountryName, countryEmoji, race.DateStart)
		}

		fullAnswer := &tele.Photo{
			File:    tele.FromURL(circuitImage),
			Caption: answer,
		}

		return c.Send(fullAnswer)
	})

	log.Println("Бот запущен...")
	b.Start()
}
