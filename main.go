package main

import (
	"log"
	"time"
 	"os"
  
	"github.com/tucnak/telebot"
  "github.com/baguswijaksono/clarin/src/vuln"
  "github.com/baguswijaksono/clarin/src/kakeibo"
  "github.com/baguswijaksono/clarin/src/urls"
  "github.com/baguswijaksono/clarin/src/bookmark"
)

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

  vuln.RegisterCommands(bot)
  kakeibo.RegisterCommands(bot)
  urls.RegisterCommands(bot)
  bookmark.RegisterCommands(bot)
  
  RegisterUtilsCommands(bot)
  
	log.Println("Bot is running...")
	bot.Start()
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