package kakeibo

import (
	"encoding/json"
	"fmt"
	"net/http"
 	"bytes"
 	"time"
	"strings"
	"github.com/tucnak/telebot"
)

const apiBase = "https://api.baguswinaksono.my.id"

func RegisterCommands(bot *telebot.Bot) {
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
