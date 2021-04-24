package peef

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "stocks",
			Description: "peef discord bot",
			Options: []*discordgo.ApplicationCommandOption{

				// TODO: group passes as a command type
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "ticker",
					Description: "Ask for ticker",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "VT",
							Value: "vt",
						},
						{
							Name:  "VTSAX",
							Value: "vtsax",
						},
						{
							Name:  "VTI",
							Value: "vti",
						},
						{
							Name:  "VOO",
							Value: "voo",
						},
					},
					Required: true,
				},
			},
		},
		{
			Name:        "responses",
			Description: "Interaction responses",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "resp-type",
					Description: "Response type",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Acknowledge",
							Value: 2,
						},
					},
					Required: true,
				},
			},
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"stocks": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			margs := []interface{}{
				i.Data.Options[0].StringValue(),
			}

			msgformat := `%s: :rocket:`

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, we'll discuss them in "responses" part
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
		"responses": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseType(i.Data.Options[0].IntValue()),
			})
			if err != nil {
				s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something gone wrong",
				})
			}
		},
	}
)

func getTickerCurrentPrice(ticker string, key string) {
	url := fmt.Sprintf("https://financialmodelingprep.com/api/v3/quote-short/%s?apikey=%s", ticker, key)
	response, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseString := string(responseData)

	fmt.Println(responseString)
}
