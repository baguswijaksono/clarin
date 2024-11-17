package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
  "net/url"
	"strings"
	"time"
 	"os"
  
	"github.com/tucnak/telebot"
  "github.com/baguswijaksono/clarin/vuln"
)

const apiBase = "https://api.baguswinaksono.my.id"
var userStates = make(map[int]string)

func main() {
  token := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	RegisterKakeiboCommands(bot)
	RegisterUrlsCommands(bot)
	RegisterUtilsCommands(bot)
  vuln.RegisterVulnCommands(bot)
  RegisterBkrmkCommands(bot)
	log.Println("Bot is running...")
	bot.Start()
}


func RegisterKakeiboCommands(bot *telebot.Bot) {
  bot.Handle("/kd", func(m *telebot.Message) {
  	parts := strings.Split(m.Text, " ")
  	if len(parts) < 2 {
  		bot.Send(m.Sender, "Usage: /delete <id>")
  		return
  	}
  
  	id := parts[1]
  
  	resp, err := http.Get(fmt.Sprintf("%s/kakeibo/d.php?id=%s", apiBase, id))
  	if err != nil {
  		bot.Send(m.Sender, "Error deleting transaction.")
  		return
  	}
  	defer resp.Body.Close()
  
  	var result map[string]interface{}
  	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
  		bot.Send(m.Sender, "Error parsing response.")
  		return
  	}
  
  	if errMsg, exists := result["error"]; exists {
  		bot.Send(m.Sender, fmt.Sprintf("Error: %s", errMsg))
  		return
  	}
  
  	msg, hasMessage := result["message"]
  	deletedData, hasDeletedData := result["data"]
  
  	if !hasMessage {
  		bot.Send(m.Sender, "Unknown error occurred.")
  		return
  	}
  
  	if hasDeletedData {
  		bot.Send(m.Sender, fmt.Sprintf("Transaction deleted successfully: %s\nDeleted Data: %v", msg, deletedData))
  	} else {
  		bot.Send(m.Sender, fmt.Sprintf("Transaction deleted successfully: %s", msg))
  	}
  })


	bot.Handle("/kg", func(m *telebot.Message) {
		parts := strings.Split(m.Text, " ")
		if len(parts) < 2 {
			bot.Send(m.Sender, "Usage: /get <id>")
			return
		}

		id := parts[1]
		resp, err := http.Get(fmt.Sprintf("%s/kakeibo/r.php?id=%s", apiBase, id))
		if err != nil {
			bot.Send(m.Sender, "Error fetching transaction.")
			return
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		bot.Send(m.Sender, fmt.Sprintf("Transaction Details: %+v", result))
	})

	bot.Handle("/ku", func(m *telebot.Message) {
		parts := strings.Split(m.Text, " ")
		if len(parts) < 5 {
			bot.Send(m.Sender, "Usage: /update <id> <type> <category> <amount> [date]")
			return
		}

		id := parts[1]
		typeVal := parts[2]
		category := parts[3]
		amount := parts[4]

		date := ""
		if len(parts) >= 6 {
			date = parts[5]
		} else {
			date = time.Now().Format("2006-01-02")
		}

		updateData := map[string]interface{}{
			"type":     typeVal,
			"category": category,
			"amount":   amount,
			"date":     date,
		}

		updatePayload, err := json.Marshal(updateData)
		if err != nil {
			bot.Send(m.Sender, "Error preparing update data.")
			return
		}

		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/kakeibo/u.php?id=%s", apiBase, id), bytes.NewBuffer(updatePayload))
		if err != nil {
			bot.Send(m.Sender, "Error creating update request.")
			return
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		updateResp, err := client.Do(req)
		if err != nil {
			bot.Send(m.Sender, "Error updating transaction.")
			return
		}
		defer updateResp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(updateResp.Body).Decode(&result); err != nil {
			bot.Send(m.Sender, "Error parsing update response.")
			return
		}

		if msg, exists := result["message"]; exists {
			bot.Send(m.Sender, fmt.Sprintf("Transaction updated successfully: %s\nUpdated Data: %v", msg, updateData))
		} else if errMsg, exists := result["error"]; exists {
			bot.Send(m.Sender, fmt.Sprintf("Error: %s", errMsg))
		} else {
			bot.Send(m.Sender, "Unknown error occurred during update.")
		}
	})

	bot.Handle("/kc", func(m *telebot.Message) {
		parts := strings.Split(m.Text, " ")
		if len(parts) < 4 {
			bot.Send(m.Sender, "Usage: /create <type> <category> <amount> [date]")
			return
		}

		date := ""
		if len(parts) >= 5 {
			date = parts[4]
		} else {
			date = time.Now().Format("2006-01-02")
		}

		transactionData := map[string]interface{}{
			"type":     parts[1],
			"category": parts[2],
			"amount":   parts[3],
			"date":     date,
		}

		jsonData, err := json.Marshal(transactionData)
		if err != nil {
			bot.Send(m.Sender, "Error preparing data.")
			return
		}

		resp, err := http.Post(fmt.Sprintf("%s/kakeibo/c.php", apiBase), "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			bot.Send(m.Sender, "Error creating transaction.")
			return
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			bot.Send(m.Sender, "Error parsing response.")
			return
		}

		if msg, exists := result["message"]; exists {
			bot.Send(m.Sender, fmt.Sprintf("Transaction created successfully: %s\nData: %v", msg, transactionData))
		} else if errMsg, exists := result["error"]; exists {
			bot.Send(m.Sender, fmt.Sprintf("Error: %s", errMsg))
		} else {
			bot.Send(m.Sender, "Unknown error occurred.")
		}
	})
}



func RegisterUrlsCommands(bot *telebot.Bot) {
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


func RegisterBkrmkCommands(bot *telebot.Bot) {
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


func RegisterUtilsCommands(bot *telebot.Bot){
	bot.Handle("/menu", func(m *telebot.Message) {
		menuText := "Here are the available commands:\n" +
			"/kd <id> - Delete a transaction by ID\n" +
			"/kg <id> - Get details of a transaction by ID\n" +
			"/ku <id> <type> <category> <amount> [date] - Update a transaction\n" +
			"/kc <type> <category> <amount> [date] - Create a new transaction\n" +
			"/uc <shorturl> <longurl> - Create a shortened URL\n" +
			"/uu <short_key> <new_url> - Update a shortened URL\n" +
			"/ud <short_key> - Delete a shortened URL\n" +
			"/menu - Show this menu\n" +
			"/man - Get detailed usage instructions for each command"

		bot.Send(m.Sender, menuText)
	})

	bot.Handle("/man", func(m *telebot.Message) {
		manualText := "Here is a detailed manual for each command:\n" +
			"/kd <id> - Delete a transaction by its ID. Example: /kd 1234\n" +
			"Usage: /kd <id>\n" +
			"This command will remove the transaction associated with the given ID.\n\n" +
			"/kg <id> - Get details of a transaction by ID. Example: /kg 1234\n" +
			"Usage: /kg <id>\n" +
			"Retrieves the details of a transaction by ID.\n\n" +
			"/ku <id> <type> <category> <amount> [date] - Update a transaction\n" +
			"Usage: /ku <id> <type> <category> <amount> [date]\n" +
			"Updates a transaction's details such as type, category, amount, and optionally the date.\n\n" +
			"/kc <type> <category> <amount> [date] - Create a new transaction\n" +
			"Usage: /kc <type> <category> <amount> [date]\n" +
			"Creates a new transaction with specified type, category, amount, and optionally a date.\n\n" +
			"/uc <shorturl> <longurl> - Create a shortened URL\n" +
			"Usage: /uc <shorturl> <longurl>\n" +
			"Creates a shortened URL with the specified short key and long URL.\n\n" +
			"/uu <short_key> <new_url> - Update a shortened URL\n" +
			"Usage: /uu <short_key> <new_url>\n" +
			"Updates an existing shortened URL with a new long URL.\n\n" +
			"/ud <short_key> - Delete a shortened URL\n" +
			"Usage: /ud <short_key>\n" +
			"Deletes the shortened URL associated with the provided short key."

		bot.Send(m.Sender, manualText)
	})
}