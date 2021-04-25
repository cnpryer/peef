package peef

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
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
			key := os.Getenv("API_KEY")

			if key == "" {
				log.Fatalf("API_KEY not found")
			}

			url := buildPriceUrl(symbol, key)
			data := getSymbolCurrentPriceData(url)
			msg := fmt.Sprintf(`%s: $%f`, symbol, data.Price)

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

func getSymbolCurrentPriceData(url string) Price {
	var prices []Price

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

	return prices[0]
}

func buildPriceUrl(symbol string, key string) string {
	// TODO: see if base can be set using NewRequest

	baseUrl := fmt.Sprintf("https://financialmodelingprep.com/api/v3/quote-short/%s", symbol)
	request, err := http.NewRequest("GET", baseUrl, nil)

	if err != nil {
		log.Fatal(err)
	}

	query := request.URL.Query()
	query.Add("apikey", key)
	request.URL.RawQuery = query.Encode()

	url := request.URL.String()
	log.Info(url)

	return url
}
