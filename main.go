package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
	"github.com/firediscordchat/peef/peef"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
	Debug          = flag.Bool("debug", false, "Run in debug mode")
	LogFile        = flag.String("logfile", "peef-bot.log", "Log file name")
)

var s *discordgo.Session
var log *logrus.Logger

func init() {
	log = logrus.New()
	godotenv.Load()
	flag.Parse()

	// Setup loggers
	log_file, err := os.OpenFile(*LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatal(err)
	}
	log.Out = log_file

	if *Debug {
		log.Level = logrus.DebugLevel
	} else {
		log.Level = logrus.InfoLevel
	}

	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := peef.CommandHandlers[i.Data.Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(
		func(s *discordgo.Session, r *discordgo.Ready) {
			log.Info("Bot is up!")
		},
	)

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	for _, v := range peef.Commands {

		_, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)

		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	defer s.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Info("Gracefully shutting down")
}
