package peef

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
			Name:        "tests",
			Description: "peef command for testing",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "arg",
					Description: "test arg",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "TEST",
							Value: ":white_check_mark:",
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
				Log.Fatalf("API_KEY not found")
			}

			url := buildPriceUrl(symbol, key)
			data := getSymbolCurrentPriceData(url)
			msg := fmt.Sprintf(`%s: $%f`, symbol, data.Price)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: msg,
				},
			})
		},
		"tests": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			msg := fmt.Sprintf(`test: %s`, i.Data.Options[0].StringValue())

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: msg,
				},
			})
		},
	}
)

func getSymbolCurrentPriceData(url string) Price {
	var prices []Price

	response, err := http.Get(url)

	if err != nil {
		Log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		Log.Fatal(err)
	}

	err = json.Unmarshal(responseData, &prices)

	if err != nil {
		Log.Fatal(err)
	}

	return prices[0]
}

func buildPriceUrl(symbol string, key string) string {
	baseUrl := "https://financialmodelingprep.com/api/v3/quote-short/" + symbol
	request, err := http.NewRequest("GET", baseUrl, nil)

	if err != nil {
		Log.Fatal(err)
	}

	query := request.URL.Query()
	query.Add("apikey", key)
	request.URL.RawQuery = query.Encode()

	url := request.URL.String()
	Log.Info(url)

	return url
}
