package peef

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "peef",
			Description: "peef discord bot",
			Options: []*discordgo.ApplicationCommandOption{

				// TODO: group passes as a command type
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "ticker",
					Description: "Ask for ticker checkup",
					Choices:     getPeefTickers(),
					Required:    false,
				},

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "chat",
					Description: "Chat with peef",
					Choices:     getPeefChatOptions(),
					Required:    false,
				},

				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "message",
					Description: "Test message to send to peef",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "test",
							Value: 1,
						},
					},
					Required: false,
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
						{
							Name:  "Channel message",
							Value: 3,
						},
						{
							Name:  "Channel message with source",
							Value: 4,
						},
						{
							Name:  "Acknowledge with source",
							Value: 5,
						},
					},
					Required: true,
				},
			},
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"chat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			margs := []interface{}{
				i.Data.Options[0].StringValue(),
				i.Data.Options[1].StringValue(),
				i.Data.Options[2].IntValue(),
			}

			msgformat := ` > test: %s `

			if len(i.Data.Options) >= 2 {
				msgformat += "\n> content: %d"
			}

			if len(i.Data.Options) >= 3 {
				msgformat += "\n> content: %d"
			}

			msgformat += "\n" + msgformat

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

func getPeefTickers() []*discordgo.ApplicationCommandOptionChoice {

	// TODO TODO TODO TODO TODO TODO TODO !?!!?@!?!?@?!@?!@ pls no
	return []*discordgo.ApplicationCommandOptionChoice{
		{
			Name:  "VTSAX",
			Value: "moon",
		},
	}

}

func getPeefChatOptions() []*discordgo.ApplicationCommandOptionChoice {

	// TODO TODO TODO TODO TODO TODO TODO !?!!?@!?!?@?!@?!@ pls no
	return []*discordgo.ApplicationCommandOptionChoice{
		{
			Name:  "invest",
			Value: "VT compensated risk",
		},
	}

}
