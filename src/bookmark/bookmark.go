package bookmark

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"github.com/tucnak/telebot"
)

const apiBase = "https://api.baguswinaksono.my.id"

func RegisterCommands(bot *telebot.Bot) {
    bot.Handle("/bc", func(m *telebot.Message) {
        parts := strings.SplitN(m.Text, " ", 3)
        if len(parts) < 3 {
            bot.Send(m.Sender, "Usage: /bc <short_key> <url>")
            return
        }

        shortKey := parts[1]
        url := parts[2]
        resp, err := http.Get(fmt.Sprintf("%s/bookmark/c.php?title=%s&url=%s", apiBase, shortKey, url))
        if err != nil {
            log.Printf("Error fetching API: %v\n", err)
            bot.Send(m.Sender, "Error creating bookmark.")
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
            bot.Send(m.Sender, fmt.Sprintf("Bookmark created successfully: %s", msg))
        } else if errMsg, exists := result["error"]; exists {
            bot.Send(m.Sender, fmt.Sprintf("Error: %s", errMsg))
        } else {
            bot.Send(m.Sender, "Unknown error occurred.")
        }
    })

    bot.Handle("/bu", func(m *telebot.Message) {
        parts := strings.Split(m.Text, " ")
        if len(parts) < 3 {
            bot.Send(m.Sender, "Usage: /bu <short_key> <new_url>")
            return
        }

        shortKey := parts[1]
        newURL := parts[2]
        resp, err := http.Get(fmt.Sprintf("%s/bookmark/u.php?short_key=%s&url=%s", apiBase, shortKey, newURL))
        if err != nil {
            bot.Send(m.Sender, "Error updating bookmark.")
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
            bot.Send(m.Sender, fmt.Sprintf("Bookmark updated successfully: %s", msg))
        } else if errMsg, exists := result["error"]; exists {
            bot.Send(m.Sender, fmt.Sprintf("Error: %s", errMsg))
        } else {
            bot.Send(m.Sender, "Unknown error occurred.")
        }
    })

bot.Handle("/bd", func(m *telebot.Message) {
    parts := strings.Split(m.Text, " ")
    if len(parts) < 2 {
        bot.Send(m.Sender, "Usage: /bd <short_key>")
        return
    }

    shortKey := parts[1]

    resp, err := http.Get(fmt.Sprintf("%s/bookmark/d.php?id=%s", apiBase, shortKey))
    if err != nil {
        bot.Send(m.Sender, "Error deleting bookmark.")
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
        bot.Send(m.Sender, fmt.Sprintf("Bookmark deleted successfully: %s", msg))
    } else if errMsg, exists := result["error"]; exists {
        bot.Send(m.Sender, fmt.Sprintf("Error: %s", errMsg))
    } else {
        bot.Send(m.Sender, "Unknown error occurred.")
    }
})

}
