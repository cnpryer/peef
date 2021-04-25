package peef

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Price struct {
	Symbol string  `json:"symbol"`
	Price  float32 `json:"price"`
	Volume int32   `json:"volume"`
}

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "stocks",
			Description: "peef discord bot",
			Options: []*discordgo.ApplicationCommandOption{

				// TODO: open up for non-choice method with unlimited API
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "symbol",
					Description: "Stock ticker symbol",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "VT",
							Value: "VT",
						},
						{
							Name:  "VTSAX",
							Value: "VTSAX",
						},
						{
							Name:  "VTI",
							Value: "VTI",
						},
						{
							Name:  "VOO",
							Value: "VOO",
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
			symbol := i.Data.Options[0].StringValue()
			msg := fmt.Sprintf(`%s: $%f`, symbol, getSymbolCurrentPrice(symbol, os.Getenv("API_KEY")))

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, we'll discuss them in "responses" part
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: msg,
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

func getSymbolCurrentPrice(symbol string, key string) float32 {
	var prices []Price

	url := fmt.Sprintf("https://financialmodelingprep.com/api/v3/quote-short/%s?apikey=%s", symbol, key)
	fmt.Println(url)

	response, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(responseData, &prices)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(prices)

	return prices[0].Price
}
