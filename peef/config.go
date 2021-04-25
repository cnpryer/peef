package peef

import (
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// .env-like configuration
type Config struct {
	GuildID  string `mapstructure:"GUILD_ID"`
	BotToken string `mapstructure:"BOT_TOKEN"`
	APIKey   string `mapstructure:"API_KEY"`
}

var Log = logrus.New()

// Get Config struct from params passed or using .env file.
// guildID (string)  - Optional
// botToken (string) - Optional but if passed apiKey must also be passed
//                     in order to prefer your configuration over .env.
//				       This is the token assigned to the bot through the
//					   Discord Developer Portal.
// apiKey (string)   - Optional but if passed apiKey must also be passed
//                     in order to prefer your configuration over .env.
//					   This is the FMP API Key used to give peef access
//					   to financial data.
// returns Config
func InitConfig(guildID *string, botToken *string, apiKey *string) Config {
	var result map[string]interface{}
	var config Config

	viper.SetConfigFile(".env")

	// TODO: https://github.com/spf13/viper#working-with-flags
	if *botToken != "" && *apiKey != "" {
		config.GuildID = *guildID
		os.Setenv("GUILD_ID", *guildID)

		config.BotToken = *botToken
		os.Setenv("BOT_TOKEN", *botToken)

		config.APIKey = *apiKey
		os.Setenv("API_KEY", *apiKey)

		return config
	}

	if err := viper.ReadInConfig(); err != nil {
		Log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&result)
	if err != nil {
		Log.Fatalf("Unable to decode into map, %v", err)
	}

	decErr := mapstructure.Decode(result, &config)

	if decErr != nil {
		Log.Fatalf("error decoding")
	}

	return config
}
