package vuln

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
    bot.Handle("/vc", func(m *telebot.Message) {
        parts := strings.SplitN(m.Text, " ", 5)
        if len(parts) < 5 {
            bot.Send(m.Sender, "Usage: /vc <title> <description> <severity> <status>")
            return
        }

        title := parts[1]
        description := parts[2]
        severity := parts[3]
        status := parts[4]

        resp, err := http.Get(fmt.Sprintf("%s/vuln/c.php?title=%s&description=%s&severity=%s&status=%s",
            apiBase, url.QueryEscape(title), url.QueryEscape(description), url.QueryEscape(severity), url.QueryEscape(status)))
        if err != nil {
            log.Printf("Error fetching API: %v\n", err)
            bot.Send(m.Sender, "Error creating vulnerability report.")
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
            bot.Send(m.Sender, fmt.Sprintf("Vulnerability report created successfully: %s", msg))
        } else if errMsg, exists := result["error"]; exists {
            bot.Send(m.Sender, fmt.Sprintf("Error: %s", errMsg))
        } else {
            bot.Send(m.Sender, "Unknown error occurred.")
        }
    })

bot.Handle("/vg", func(m *telebot.Message) {
    parts := strings.Split(m.Text, " ")
    if len(parts) < 2 {
        bot.Send(m.Sender, "Usage: /vg <id>")
        return
    }

    id := parts[1]
    resp, err := http.Get(fmt.Sprintf("%s/vuln/r.php?id=%s", apiBase, id))
    if err != nil {
        bot.Send(m.Sender, "Error fetching vulnerability report.")
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

    if report, exists := result["report"]; exists {
        bot.Send(m.Sender, fmt.Sprintf("Vulnerability Report: %v", report))
    } else if errMsg, exists := result["error"]; exists {
        bot.Send(m.Sender, fmt.Sprintf("Error: %s", errMsg))
    } else {
        bot.Send(m.Sender, "Unknown error occurred.")
    }
})


bot.Handle("/vu", func(m *telebot.Message) {
    parts := strings.Split(m.Text, " ")
    if len(parts) < 3 {
        bot.Send(m.Sender, "Usage: /vu <id> <new_status>")
        return
    }

    id := parts[1]
    newStatus := parts[2]
    resp, err := http.Get(fmt.Sprintf("%s/vuln/u.php?id=%s&status=%s", apiBase, id, url.QueryEscape(newStatus)))
    if err != nil {
        bot.Send(m.Sender, "Error updating vulnerability status.")
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
        bot.Send(m.Sender, fmt.Sprintf("Vulnerability status updated successfully: %s", msg))
    } else if errMsg, exists := result["error"]; exists {
        bot.Send(m.Sender, fmt.Sprintf("Error: %s", errMsg))
    } else {
        bot.Send(m.Sender, "Unknown error occurred.")
    }
})


    bot.Handle("/vd", func(m *telebot.Message) {
        parts := strings.Split(m.Text, " ")
        if len(parts) < 2 {
            bot.Send(m.Sender, "Usage: /vd <id>")
            return
        }

        id := parts[1]
        resp, err := http.Get(fmt.Sprintf("%s/vuln/d.php?id=%s", apiBase, id))
        if err != nil {
            bot.Send(m.Sender, "Error deleting vulnerability report.")
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
            bot.Send(m.Sender, fmt.Sprintf("Vulnerability report deleted successfully: %s", msg))
        } else if errMsg, exists := result["error"]; exists {
            bot.Send(m.Sender, fmt.Sprintf("Error: %s", errMsg))
        } else {
            bot.Send(m.Sender, "Unknown error occurred.")
        }
    })
}