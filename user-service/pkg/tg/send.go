package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func loadLocation(locationStr string) (time.Time, error) {
	location, err := time.LoadLocation(locationStr)
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().In(location), nil
}

func SendError(errorMessage any, route string) {
	tokenBot := "5944731024:AAGCikQUMGDzBsN3ZvjD4AYPOStYXaFjA80"
	chatID := 484031907

	location, err := loadLocation("Asia/Yekaterinburg")
	if err != nil {
		log.Fatal(err)
		return
	}

	var infoString string
	switch v := errorMessage.(type) {
	case string:
		infoString = v
	case error:
		infoString = v.Error()
	default:
		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		encoder.SetIndent("", "  ")
		err := encoder.Encode(v)
		if err != nil {
			infoString = fmt.Sprintf("%+v", v)
		} else {
			infoString = buf.String()
		}
	}

	errorStr := fmt.Sprintf("Ошибка\n")
	errorStr += fmt.Sprintf("Сервис ошибки: %s\n", os.Getenv("app.name"))
	errorStr += fmt.Sprintf("Маршрут: %s\n", route)
	errorStr += fmt.Sprintf("Дата: %d:%d:%d\n", location.Day(), location.Month(), location.Year())
	errorStr += fmt.Sprintf("Время: %s\n", location.Format("15:04.05"))
	errorStr += fmt.Sprintf("Подробности:\n%s", infoString)

	encodedMessage := url.QueryEscape(errorStr)

	urlToSend := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d", tokenBot, chatID)

	urlAPI := fmt.Sprintf("%s&text=%s", urlToSend, encodedMessage)

	resp, err := http.Get(urlAPI)
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println(fmt.Sprintf("Ошибка при отправке сообщения: %s", string(body)))
		return
	}
}

func SendInfo(info any, route string) {
	tokenBot := "5944731024:AAGCikQUMGDzBsN3ZvjD4AYPOStYXaFjA80"
	chatID := 484031907

	location, err := loadLocation("Asia/Yekaterinburg")
	if err != nil {
		log.Fatal(err)
		return
	}

	var infoString string
	switch v := info.(type) {
	case string:
		infoString = v
	case error:
		infoString = v.Error()
	default:
		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		encoder.SetIndent("", "  ")
		err := encoder.Encode(v)
		if err != nil {
			infoString = fmt.Sprintf("%+v", v)
		} else {
			infoString = buf.String()
		}
	}

	infoStr := fmt.Sprintf("Информация\n")
	infoStr += fmt.Sprintf("Сервис информации: %s\n", os.Getenv("app.name"))
	infoStr += fmt.Sprintf("Маршрут: %s\n", route)
	infoStr += fmt.Sprintf("Дата: %d:%d:%d\n", location.Day(), location.Month(), location.Year())
	infoStr += fmt.Sprintf("Время: %s\n", location.Format("15:04.05"))
	infoStr += fmt.Sprintf("Подробности: \n%s", infoString)

	encodedMessage := url.QueryEscape(infoStr)

	urlToSend := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d", tokenBot, chatID)

	urlAPI := fmt.Sprintf("%s&text=%s", urlToSend, encodedMessage)

	resp, err := http.Get(urlAPI)
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println(fmt.Sprintf("Ошибка при отправке сообщения: %s", string(body)))
		return
	}
}
