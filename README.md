# peef

A discord bot for the [FIRE discord guild](https://discord.gg/SYE2JCzsES).

This project is under development.

# Requirements
Go v1.16

## Quickstart
 - Clone the repository
 - Make sure you have go installed and on your PATH
 - Set your [FMP](https://financialmodelingprep.com/developer/docs/) `API_KEY` environment variable
 - Run `go run main.go -token "<Your-bot-token>"`

## Debugging
 - Run with argument `-debug` to set the log level to debug
 - Run with argument `-guild <your-discord-server-id>` to limit the bot commands to a specific discord guild. 

# Using a envFile

Create a `.env` file at the root of the project and add the following values:
```envFile
GUILD_ID=your guild ID
BOT_TOKEN=your bot's token
API_KEY=your FMP API key
```

## Planned Features

- `/stock symbol: [symbol]` for current prices
- Guild moderation
  - `/ban`
  - `/kick`
  - `/warn`
  - `/assign`
  - automated #mod-log channel
  - managed #change-log channel
- Peef dialogue
  - No, that's uncompensated risk
  - You should really invest in VT
- General FIRE content
  - managed #fire-faq channel
