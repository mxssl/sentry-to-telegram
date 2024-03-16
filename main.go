package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"net/http"
	"os"
)

type sentryWebHook struct {
	Project string `json:"project"`
	URL     string `json:"url"`
	Event   struct {
		Title string `json:"title"`
	} `json:"event"`
}

type tgMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("setting Telegram bot token")
	tgBotToken, ok := os.LookupEnv("TELEGRAM_BOT_TOKEN")
	if !ok {
		slog.Error("cannot get the env variable TELEGRAM_BOT_TOKEN")
		os.Exit(1)
	}

	slog.Info("setting Telegram group ID")
	tgGroupID, ok := os.LookupEnv("TELEGRAM_GROUP_ID")
	if !ok {
		slog.Error("cannot get the env variable TELEGRAM_GROUP_ID")
		os.Exit(1)
	}

	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request", slog.String("remote_address", r.RemoteAddr))
		var wh sentryWebHook
		err := json.NewDecoder(r.Body).Decode(&wh)
		if err != nil {
			slog.Error("cannot decode request body", slog.Any("error_message", err))
			http.Error(w, "cannot decode request body", http.StatusInternalServerError)
			return
		}

		msg := fmt.Sprintf("<b>Project:</b> %s\n<b>Title:</b> %s\n<b>URL:</b> %s",
			html.EscapeString(wh.Project),
			html.EscapeString(wh.Event.Title),
			html.EscapeString(wh.URL))

		err = sendToTelegram(tgBotToken, tgGroupID, msg)
		if err != nil {
			slog.Error("cannot send message to telegram", slog.Any("error_message", err))
			http.Error(w, "cannot send message to telegram", http.StatusInternalServerError)
			return
		}
	})

	slog.Info("server is starting on 0.0.0.0:9999")
	if err := http.ListenAndServe("0.0.0.0:9999", nil); err != nil {
		slog.Error("cannot start the server", slog.Any("error_message", err))
		os.Exit(1)
	}
}

func sendToTelegram(token, group, message string) error {
	b, _ := json.Marshal(tgMessage{
		ChatID:    group,
		Text:      message,
		ParseMode: "HTML",
	})
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}
	return nil
}
