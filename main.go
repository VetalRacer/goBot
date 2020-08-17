package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type GetMeT struct {
	Ok     bool         `json:"ok"`
	Result GetMeResultT `json:"result"`
}

type GetMeResultT struct {
	Id        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type SendMessageT struct {
	Ok     bool     `json:"ok"`
	Result MessageT `json:"result"`
}

type GetUpdatesT struct {
	Ok     bool                `json:"ok"`
	Result []GetUpdatesResultT `json:"result"`
}

type GetUpdatesResultT struct {
	UpdateID int                `json:"update_id"`
	Message  GetUpdatesMessageT `json:"message,omitempty"`
}

type GetUpdatesMessageT struct {
	MessageID int `json:"message_id"`
	From      struct {
		ID           int    `json:"id"`
		IsBot        bool   `json:"is_bot"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Username     string `json:"username"`
		LanguageCode string `json:"language_code"`
	} `json:"from"`
	Chat struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
		Type      string `json:"type"`
	} `json:"chat"`
	Date     int    `json:"date"`
	Text     string `json:"text"`
	Entities []struct {
		Offset int    `json:"offset"`
		Length int    `json:"length"`
		Type   string `json:"type"`
	} `json:"entities"`
}

type MessageT struct {
	MessageID int                          `json:"message_id"`
	From      GetUpdatesResultMessageFromT `json:"from"`
	Chat      GetUpdatesResultMessageChatT `json:"chat"`
	Date      int                          `json:"date"`
	Text      string                       `json:"text"`
}

type GetUpdatesResultMessageFromT struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type GetUpdatesResultMessageChatT struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type GetRandomJokeT struct {
	Content string `json" "content"`
}

const telegramBaseUrl = "https://api.telegram.org/bot"
const telegramToken = "1379645763:AAHJ9TTeL5DL2rZF9wzV3Rr5Fv0YetGUou8"

const methodGetMe = "getMe"
const methodGetUpdates = "getUpdates"
const methodSendMessage = "sendMessage"

const sleep = 5

var body []byte
var offset int

func main() {
	for true {
		if len(body) == len(getBodyByUrlAndData(getUrlByMethod(methodGetUpdates))) {
			time.Sleep(5)
		} else {
			body = getBodyByUrlAndData(getUrlByMethod(methodGetUpdates))
			getUpdates := GetUpdatesT{}
			err := json.Unmarshal(body, &getUpdates)
			if err != nil {
				fmt.Println(err.Error())
			}

			sendMessageUrl := getUrlByMethod(methodSendMessage)
			for _, item := range getUpdates.Result {
				userMessage := strings.ToLower(item.Message.Text)
				offset = item.UpdateID + 1
				switch {
				case checkUserMessage(userMessage, "start"):
					sendMessageByUser(&item, &sendMessageUrl, "Привет! Я - бот для задания SkillBox!%0A%0AЯ умею:%0A/HeadsOrTails - Орел или Решка%0A/Dise - игральные кости%0A/Joke - рандомный анекдот%0A%0AИли могу немного поболтать, но пока знаю мало слов)")
				case checkUserMessage(userMessage, "уме"):
					sendMessageByUser(&item, &sendMessageUrl, "Привет :) ! Да много разного, нажми /start и сам все узнаешь))")
				case checkUserMessage(userMessage, "test"):
					sendMessageByUser(&item, &sendMessageUrl, "Да-да, все работает, не переживай ;) а теперь нажми-ка /start !")
				case checkUserMessage(userMessage, "привет"):
					sendMessageByUser(&item, &sendMessageUrl, "И тебе привет!")
				case checkUserMessage(userMessage, "как") || checkUserMessage(userMessage, "дела"):
					sendMessageByUser(&item, &sendMessageUrl, "Хорошо, а у тебя?")
				case checkUserMessage(userMessage, "отлично"):
					sendMessageByUser(&item, &sendMessageUrl, "Здорово!")
				case checkUserMessage(userMessage, "heads"):
					sendMessageByUser(&item, &sendMessageUrl, gameHeadsOrTails())
				case checkUserMessage(userMessage, "dise"):
					sendMessageByUser(&item, &sendMessageUrl, gameDise())
				case checkUserMessage(userMessage, "joke"):
					sendMessageByUser(&item, &sendMessageUrl, getRandomJoke())
				default:
					sendMessageByUser(&item, &sendMessageUrl, "Прости, я этого не понял :(%0AПопробуй жамкнуть /start и посмотри чем я могу тебе пригодиться ;)")
				}
			}
		}

	}

}

func checkUserMessage(userMessage, findText string) bool {
	return strings.Contains(userMessage, findText)
}

func sendMessageByUser(item *GetUpdatesResultT, sendMessageUrl *string, textMessage string) {

	chatId := strconv.Itoa(item.Message.Chat.ID)
	targetUrl := *sendMessageUrl + "?chat_id=" + chatId + "&text=" + textMessage
	getBodyByUrlAndData(targetUrl)
}

func getUrlByMethod(methodName string) string {
	if methodName == methodGetUpdates {
		return telegramBaseUrl + telegramToken + "/" + methodName + "?offset=" + strconv.Itoa(offset)
	} else {
		return telegramBaseUrl + telegramToken + "/" + methodName
	}
}

func gameHeadsOrTails() string {
	rand.Seed(time.Now().UTC().UnixNano())
	if rand.Intn(500)%2 == 0 {
		return "Вам выпал: Орел%0A/repeat_HeadsOrTails"
	} else {
		return "Вам выпало: Решка%0A/repeat_HeadsOrTails"
	}
}

func gameDise() string {
	rand.Seed(time.Now().UTC().UnixNano())
	diseOne := rand.Intn(5) + 1
	diseTwo := rand.Intn(5) + 1
	text := "Вам выпало: %0AКость №1 - число " + strconv.Itoa(diseOne) + "%0AКость №2 - число " + strconv.Itoa(diseTwo) + "%0A------%0AСумма - " + strconv.Itoa(diseOne+diseTwo) + "%0A/repeat_Dise"
	return text
}

func getRandomJoke() string {

	joke := getBodyByUrlAndData("http://rzhunemogu.ru/RandJSON.aspx?CType=1")

	getJoke := GetRandomJokeT{}
	err := json.Unmarshal(joke, &getJoke)
	if err != nil {
		fmt.Println("Ошибка:", err.Error())
	}
	a := getJoke.Content
	fmt.Println(a)
	return "К сожалению данный функционал временно не работает! Попробуйте что-то другое :)"
}

func getBodyByUrlAndData(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	return body
}
