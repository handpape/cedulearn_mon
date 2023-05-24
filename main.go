package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ResponseFeed struct {
	nick string
	code int
}

func main() {
	bot, err := tgbotapi.NewBotAPI("1863759384:AAHZxKTZqmPM9BRbEvSSi6m2IDVfAF45q1E")
	if err != nil {
		log.Panic(err)
	}
	ch := make(chan ResponseFeed)
	//chv := make(chan bool, 1)
	//chv <- false
	go monloop(ch)
	go shoot(ch, bot)
	telegram_loop(bot)

}

func monloop(ch chan ResponseFeed) {
	for {
		time.Sleep(time.Duration(time.Minute))
		file, err := os.Open("url.txt")
		if err != nil {
			log.Fatalf("Error when opening file: %s", err)
			os.Exit(-1)
		}
		fileScanner := bufio.NewScanner(file)
		fileScanner.Split(bufio.ScanLines)

		var fileLines []string
		for fileScanner.Scan() {
			fileLines = append(fileLines, fileScanner.Text())
		}

		if err := fileScanner.Err(); err != nil {
			log.Fatalf("Error while reading file: %s", err)
		}

		for _, line := range fileLines {
			slice := strings.Split(line, ",")
			if len(slice) != 2 {
				log.Fatalf("pasing error")
				break
			}
			ret := urlcall(slice[0], slice[1])
			rf := ResponseFeed{
				nick: slice[0],
				code: ret,
			}
			ch <- rf

		}

		defer file.Close()
	}
}

func urlcall(nick string, url string) int {

	fmt.Println(url)
	request, _ := http.NewRequest("GET", url, nil)
	//request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		log.Fatalln(error)
		fmt.Println("error")
	}

	defer response.Body.Close()
	return response.StatusCode
}

func shoot(ch chan ResponseFeed, bot *tgbotapi.BotAPI) {
	for rf := range ch {

		if rf.code == 500 || rf.code == 502 {
			text := fmt.Sprintf("%s:%d (확인필요)", rf.nick, rf.code)
			msg := tgbotapi.NewMessage(273439537, text)
			bot.Send(msg)
		}
	}
}

func telegram_loop(bot *tgbotapi.BotAPI) {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			str := update.Message.Text
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, str)
			bot.Send(msg)
		}
	}
}
