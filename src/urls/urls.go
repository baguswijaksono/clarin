package urls

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"github.com/tucnak/telebot"
)

const apiBase = "https://api.baguswinaksono.my.id"

func RegisterCommands(bot *telebot.Bot) {
    bot.Handle("/uc", func(m *telebot.Message) {
        parts := strings.SplitN(m.Text, " ", 3)
        if len(parts) < 3 {
            bot.Send(m.Sender, "Usage: /uc <shorturl> <longurl>")
            return
        }

        shortURL := url.QueryEscape(parts[1])
        originalURL := url.QueryEscape(parts[2])
        resp, err := http.Get(fmt.Sprintf("%s/urls/c.php?short_key=%s&url=%s", apiBase, shortURL, originalURL))
        if err != nil {
            log.Printf("Error fetching API: %v\n", err)
            bot.Send(m.Sender, "Error creating shortened URL.")
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            bot.Send(m.Sender, "Error: Received non-OK response from the API.")
            return
        }

        var result map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
            log.Printf("Error decoding API response: %v\n", err)
            bot.Send(m.Sender, "Error parsing response.")
            return
        }

        if msg, exists := result["message"]; exists {
            bot.Send(m.Sender, fmt.Sprintf("URL shortened successfully: %s", msg))
        } else if errMsg, exists := result["error"]; exists {
            bot.Send(m.Sender, fmt.Sprintf("Error: %s", errMsg))
        } else {
            bot.Send(m.Sender, "Unknown error occurred.")
        }
    })

    bot.Handle("/uu", func(m *telebot.Message) {
        parts := strings.Split(m.Text, " ")
        if len(parts) < 3 {
            bot.Send(m.Sender, "Usage: /uu <short_key> <new_url>")
            return
        }

        shortKey := parts[1]
        newURL := parts[2]
        resp, err := http.Get(fmt.Sprintf("%s/urls/u.php?short_key=%s&url=%s", apiBase, shortKey, newURL))
        if err != nil {
            bot.Send(m.Sender, "Error updating URL.")
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            bot.Send(m.Sender, "Error: Received non-OK response from the API.")
            return
        }

        var result map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
            log.Printf("Error decoding API response: %v\n", err)
            bot.Send(m.Sender, "Error parsing response.")
            return
        }

        if msg, exists := result["message"]; exists {
            bot.Send(m.Sender, fmt.Sprintf("URL updated successfully: %s", msg))
        } else if errMsg, exists := result["error"]; exists {
            bot.Send(m.Sender, fmt.Sprintf("Error: %s", errMsg))
        } else {
            bot.Send(m.Sender, "Unknown error occurred.")
        }
    })

    bot.Handle("/ud", func(m *telebot.Message) {
        parts := strings.Split(m.Text, " ")
        if len(parts) < 2 {
            bot.Send(m.Sender, "Usage: /ud <short_key>")
            return
        }

        shortKey := parts[1]
        resp, err := http.Get(fmt.Sprintf("%s/urls/d.php?short_key=%s", apiBase, shortKey))
        if err != nil {
            bot.Send(m.Sender, "Error deleting URL.")
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            bot.Send(m.Sender, "Error: Received non-OK response from the API.")
            return
        }

        var result map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
            log.Printf("Error decoding API response: %v\n", err)
            bot.Send(m.Sender, "Error parsing response.")
            return
        }

        if msg, exists := result["message"]; exists {
            bot.Send(m.Sender, fmt.Sprintf("URL deleted successfully: %s", msg))
        } else if errMsg, exists := result["error"]; exists {
            bot.Send(m.Sender, fmt.Sprintf("Error: %s", errMsg))
        } else {
            bot.Send(m.Sender, "Unknown error occurred.")
        }
    })
}