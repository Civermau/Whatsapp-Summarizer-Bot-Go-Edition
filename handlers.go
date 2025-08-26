package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go.mau.fi/whatsmeow"
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

// * Event handler for WhatsApp events
func eventHandler(evt interface{}, client *whatsmeow.Client) {
	switch v := evt.(type) {
	case *events.Message:
		// * Get message info
		chatID := v.Info.Chat.User
		pushName := v.Info.PushName
		sender := v.Info.Sender.User
		message := v.Message.GetConversation()
		messageType := "text"

		// * Determine message type
		if v.Message.GetConversation() != "" {
			messageType = "text"
		} else if v.Message.GetImageMessage() != nil {
			messageType = "image"
			message = "[Image Message]"
		} else if v.Message.GetVideoMessage() != nil {
			messageType = "video"
			message = "[Video Message]"
		} else if v.Message.GetAudioMessage() != nil {
			messageType = "audio"
			message = "[Audio Message]"
		} else if v.Message.GetDocumentMessage() != nil {
			messageType = "document"
			message = "[Document Message]"
		} else if v.Message.GetStickerMessage() != nil {
			messageType = "sticker"
			message = "[Sticker Message]"
		} else if v.Message.GetReactionMessage() != nil {
			messageType = "reaction"
			message = "[Reaction Message]"
		} else if v.Message.GetViewOnceMessage() != nil {
			messageType = "view_once"
			message = "[View Once Message]"
		} else if v.Message.GetLiveLocationMessage() != nil {
			messageType = "live_location"
			message = "[Live Location Message]"
		} else if v.Message.GetLocationMessage() != nil {
			messageType = "location"
			message = "[Location Message]"
		} else {
			messageType = "unknown"
			message = "[Unknown Message Type]"
		}

		timestamp := v.Info.Timestamp

		if pushName != "" {
			sender = pushName
		}

		// * Insert message into database
		if err := insertMessage(chatID, sender, message, messageType, timestamp); err != nil {
			fmt.Printf("Error inserting message into database: %v\n", err)
		} else {
			fmt.Printf("Message saved to database: %s from %s in %s (%s)\n", messageType, sender, chatID, timestamp.Format("2006-01-02 15:04:05"))
		}

		// * Ignore messages from myself
		if v.Info.IsFromMe {
			fmt.Println("Ignoring message from myself")
			return
		}

		// This prevents historical messages from triggering commands during sync
		// Best approach I found as an alternative client.on("message_create") on whatsapp-web.js
		if v.Info.Timestamp.Before(botStartTime) {
			return
		}

		// * Handle owner messages
		if v.Info.Chat.User == ownerJID {
			fmt.Println("Sending request to DeepSeek API for ", v.Info.Chat.User)
			msgID, _ := client.SendMessage(context.Background(), v.Info.Chat, &waE2E.Message{
				Conversation: proto.String("Processing..."),
			})

			normalRequestAsync(v.Message.GetConversation(), client, v.Info.Chat, msgID.ID)
			return
		}

		// TODO: Handle direct messages, next version
		if !v.Info.MessageSource.IsGroup {
			client.SendMessage(context.Background(), v.Info.Chat, &waE2E.Message{
				Conversation: proto.String("Sending direct messages to the bot is not supported yet"),
			})
			client.MarkRead([]types.MessageID{v.Info.ID}, v.Info.Timestamp, v.Info.Chat, v.Info.Sender, types.ReceiptTypeRead)
			fmt.Println("Direct message from: ", v.Info.PushName, "("+v.Info.Chat.User+")")
			return
		}

		// * Handle @everyone command
		if strings.Contains(v.Message.GetConversation(), "@everyone") {
			handleEveryoneCommand(client, v.Info.Chat)
			return
		}

		// * Handle summarize command
		words := strings.Split(v.Message.GetConversation(), " ")

		// * Handle summarize command
		if words[0] == "--summarize" || words[0] == "-s" {
			opts, err := parseSummarizeCommand(words)
			if err != nil {
				sendSticker(client, v.Info.Chat, "11")
				fmt.Println("Error parsing summarize command: ", err)
				client.SendMessage(context.Background(), v.Info.Chat, &waE2E.Message{
					Conversation: proto.String("Usage: --summarize <number of messages>"),
				})
				return
			}

			handleSummarizeCommand(client, v.Info.Chat, opts, v.Info.PushName)
			return
		}

		// * Handle info command
		if words[0] == "--info" || words[0] == "-i" || words[0] == "-h" || words[0] == "--help" {
			handleInfoCommand(client, v.Info.Chat)
			return
		}

		if words[0] == "--version" || words[0] == "-v" {
			handleVersionCommand(client, v.Info.Chat)
			return
		}
	}
}

func handleVersionCommand(client *whatsmeow.Client, chat types.JID) {
	client.SendMessage(context.Background(), chat, &waE2E.Message{
		Conversation: proto.String(
			"*Bot version 3.0.0!*\n" +
				"Bot is now running on Go! Quicker, more stable and more efficient!\n" +
				"Now it responds with a message while processing the request, then it edits it with the result!\n" +
				"\n" +
				"*Future updates:* \n" +
				"- Add support for direct messages (I still don't if this is a good idea)\n" +
				"- Add support for group messages (Make bancho a part of the group!)\n" +
				"- Add support for media messages (Image recognition, Speech to text for audio messages, etc. So it can understand even more!)\n" +
				"\n" +
				"Check out the code: https://github.com/Civermau/Whatsapp-Summarizer-Bot-Go-Edition\n" +
				"Also check out my website: https://civermau.dev",
		),
	})
}

func handleEveryoneCommand(client *whatsmeow.Client, chat types.JID) {
	groupInfo, err := client.GetGroupInfo(chat)
	if err != nil {
		fmt.Println("Error getting group info: ", err)
		return
	}

	members := groupInfo.Participants
	mentions := []string{}

	for _, member := range members {
		jid := member.JID.String()
		mentions = append(mentions, jid)
	}
	sendEveryoneSticker(client, chat, "9", mentions)
}

func handleInfoCommand(client *whatsmeow.Client, chat types.JID) {
	client.SendMessage(context.Background(), chat, &waE2E.Message{
		Conversation: proto.String(
			"Bot created by *Civer_mau*!\n" +
				"\n" +
				"Summarizes messages via DeepSeek API (I have to pay for that, please don't abuse it)\n" +
				"\n" +
				"*Commands:* \n" +
				"- --summarize <number of messages> (Summarizes the last <number of messages> messages)\n" +
				"- --info (Shows info about the bot)\n" +
				"- --version (Shows the version of the bot)\n" +
				"\n" +
				"*Summarize Command Flags:*\n" +
				"- --short (Creates a short summary)\n" +
				"- --medium (Creates a medium-length summary - default)\n" +
				"- --long (Creates a long, detailed summary)\n" +
				"\n" +
				"*Examples:*\n" +
				"- --summarize 50 --short (Summarize last 50 messages in short format)\n" +
				"- -s 100 --long (Summarize last 100 messages in long format)\n" +
				"\n" +
				"Check out the code: https://github.com/Civermau/Whatsapp-Summarizer-Bot-Go-Edition\n" +
				"Also check out my website: https://civermau.dev",
		),
	})
}

// Handle summarize command logic
func handleSummarizeCommand(client *whatsmeow.Client, chat types.JID, opts SummarizeOptions, pushName string) {
	if opts.Count <= 10 && opts.Count > 0 {
		sendSticker(client, chat, "11")
		client.SendMessage(context.Background(), chat, &waE2E.Message{
			Conversation: proto.String("I will NOT summarize less than 10 messages."),
		})
		return
	}
	if opts.Count <= 0 {
		sendSticker(client, chat, "21")
		client.SendMessage(context.Background(), chat, &waE2E.Message{
			Conversation: proto.String("Tu te crees muy chistosito verdad w?"),
		})
		return
	}

	if opts.Count >= 100 && opts.Count < 500 {
		sendSticker(client, chat, "14")
	}
	if opts.Count >= 500 {
		sendSticker(client, chat, "5")
	}

	msgID, _ := client.SendMessage(context.Background(), chat, &waE2E.Message{
		Conversation: proto.String("Reading " + strconv.Itoa(opts.Count) + " messages..."),
	})

	fmt.Println("Summarizing " + strconv.Itoa(opts.Count) + " messages via DeepSeek API for " + chat.User)

	groupName := chat.User
	groupInfo, err := client.GetGroupInfo(chat)
	if err == nil {
		groupName = groupInfo.Name
	}

	sendMessageToOwner(client, pushName+" requested to summarize "+strconv.Itoa(opts.Count)+" messages in "+groupName)
	summarizeMessages(client, chat, msgID.ID, opts)
}

func parseSummarizeCommand(words []string) (SummarizeOptions, error) {
	opts := SummarizeOptions{
		Count: 0,
		Style: "medium",
		Media: false, //Implement later
	}

	i := 0
	for i < len(words) {
		switch words[i] {
		case "-s", "--summarize":
			if i+1 < len(words) {
				if num, err := strconv.Atoi(words[i+1]); err == nil {
					opts.Count = num
					i++ // skip number
				} else {
					// fmt.Println("Invalid number after ", words[i], " : ", err)
					return opts, fmt.Errorf("invalid number after %s", words[i])
				}
			} else {
				// fmt.Println("Missing number after ", words[i])
				return opts, fmt.Errorf("missing number after %s", words[i])
			}

		case "--short":
			opts.Style = "short"
		case "--medium":
			opts.Style = "medium"
		case "--long":
			opts.Style = "long"
		case "--media":
			opts.Media = true
		}
		i++
	}

	// fmt.Println("Summarize options:")
	// fmt.Println("Count: ", opts.Count)
	// fmt.Println("Style: ", opts.Style)
	// fmt.Println("Media: ", opts.Media)
	return opts, nil
}
