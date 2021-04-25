package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
	"github.com/firediscordchat/peef/peef"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	APIKey         = flag.String("api-key", "", "FMP API Key for financial data")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
	Debug          = flag.Bool("debug", false, "Run in debug mode")
	LogFile        = flag.String("logfile", "peef-bot.log", "Log file name")
)

var session *discordgo.Session

func init() {
	flag.Parse()

	config := peef.InitConfig(GuildID, BotToken, APIKey)

	// Setup loggers
	logFile, err := os.OpenFile(*LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		peef.Log.Fatal(err)
	}
	peef.Log.Out = logFile

	if *Debug {
		peef.Log.Level = logrus.DebugLevel
	} else {
		peef.Log.Level = logrus.InfoLevel
	}

	session, err = discordgo.New("Bot " + config.BotToken)
	if err != nil {
		peef.Log.Fatalf("Invalid bot parameters: %v", err)
	}

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := peef.CommandHandlers[i.Data.Name]; ok {
			h(s, i)
		}
	})
}

func main() {

	// TODO: there's probably a better way to set config
	config := peef.InitConfig(GuildID, BotToken, APIKey)

	session.AddHandler(
		func(s *discordgo.Session, r *discordgo.Ready) {
			peef.Log.Info("Bot is up!")
		},
	)

	err := session.Open()
	if err != nil {
		peef.Log.Fatalf("Cannot open the session: %v", err)
	}

	for _, v := range peef.Commands {

		_, err := session.ApplicationCommandCreate(session.State.User.ID, config.GuildID, v)

		if err != nil {
			peef.Log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	defer session.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	peef.Log.Info("Gracefully shutting down")
}
