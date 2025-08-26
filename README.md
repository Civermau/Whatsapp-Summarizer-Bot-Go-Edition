
# WhatsApp Summarizer Bot - Go Edition

So, my friends talk a lot in the group, but A LOT, and pretty often I find 200+ unread messages, so I got tired of it and decided to make a bot to summarize the chat. I made one in JS using whatsapp-web.js, but there are not enough words in english nor spanish to express my eternal hate to nodeJS, and after investigating a little more, found whatsmeow!

So yeah, this is a new one, but better, and in Go

## Features

- Send commands in chat like in a linux command line
    - --summarize <n>
    - -- info
    - -- version
    ---
    It has it's shortened version too!
    - -s <n>
    - -i
    - -v
- Mention everyone in chat using `@everyone`
- Stores messages in a .db sqlite (since whatsmeow doesn't have a way to fetch messages unlike whatsapp-web.js) 
- Summarize the conversation using deepseek api
- Notifies bot owner via DM when someone uses summarize command (I need to be aware of the api usage)

## Usage/Examples
In config.go you must place your API key and your phone number so the bot can do the summaries correctly and inform you when someone uses it
```go
var (
	botStartTime time.Time
	db           *sql.DB
	dsClient, _  = deepseek.NewClient("YOUR_API_KEY_HERE")
	ownerJID     = "YOUR_PHONE_NUMBER_HERE"
)

```

on first run and logging in, whatsmeow creates a .db file, edit that file to add a new table for storing the messages
```sql
CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chat_id TEXT NOT NULL,
    sender TEXT NOT NULL,
    message TEXT,
    message_type TEXT,
    timestamp DATETIME
);
```


## Deployment

I build and run the project with

```bash
  # This is for building and running in the same command
  clear ; go build && ./Whatsapp-summarizer-Bot-Go-Edition 

  # You can just do
  go build
  ./Whatsapp-Bot-Go-Edition

  # As my other projects, if you are planning on running This
  # Please give it a proper name haha
```


To deploy on an OrangePi 5B I had to do

```bash
  CC=aarch64-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o ./OrangePiDeploy/Whatsapp-Summarizer-Bot .

  #Yeah, that's a long command
```

IMPORTANT

The project structure should look like 
```
Whatsapp-Summarizer/
│── Whatsapp-Summarizer     # compiled binary
│── Media/                  # Media folder for stickers
│── example.db              # Sqlite db
```
## License

This project and its contained wallpapers are licensed under the Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License - see the [LICENSE](LICENSE) file for details.

Check also [this link](https://creativecommons.org/licenses/by-nc-sa/4.0/)

## Acknowledgements

 - [Whatsmeow library for Golang](https://github.com/tulir/whatsmeow)
 - [Golang, I guess](https://go.dev/)
 - [Todoroki Hajime](https://hololive.hololivepro.com/en/talents/todoroki-hajime/)
    - (the personality of the bot is an inside joke with a couple friends that comes from [this video](https://www.youtube.com/watch?v=DZTXaq23534&list=RDDZTXaq23534&start_radio=1) from Hajime)
- [This sticker pack](https://store.line.me/stickershop/product/29303803/en)
    - This are the stickers the bot use
    - Since they are paid stickers i won't upload them to the repo, I will pay for them when I have the money I swear
- [This readme maker](https://readme.so/es/editor)