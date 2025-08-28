package main

import (
	"context"
	"fmt"
	"os"

	"go.mau.fi/whatsmeow"
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

type SummarizeOptions struct {
	Count int    // number of messages to summarize
	Style string // short, medium, long
	Media bool   // include media in the summary, this is for future update, so if you are reading the code, surprise I guess haha.
}

func sendSticker(client *whatsmeow.Client, chat types.JID, stickerID string, mentions ...string) {
	data, err := os.ReadFile("Media/Bancho-" + stickerID + ".webp")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	upload, err := client.Upload(context.Background(), data, whatsmeow.MediaImage)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	sticker := &waE2E.StickerMessage{
		URL:           proto.String(upload.URL),
		DirectPath:    proto.String(upload.DirectPath),
		MediaKey:      upload.MediaKey,
		Mimetype:      proto.String("image/webp"),
		FileEncSHA256: upload.FileEncSHA256,
		FileSHA256:    upload.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(data))),
	}

	// Only add ContextInfo if mentions are provided
	if len(mentions) > 0 {
		sticker.ContextInfo = &waE2E.ContextInfo{
			MentionedJID: mentions,
		}
	}

	msg := &waE2E.Message{
		StickerMessage: sticker,
	}

	_, err = client.SendMessage(context.Background(), chat, msg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func sendEveryoneSticker(client *whatsmeow.Client, chat types.JID, stickerID string, mentions []string) {
	sendSticker(client, chat, stickerID, mentions...)
}

func sendMessageToOwner(client *whatsmeow.Client, message string) {
	ownerJIDParsed := types.NewJID(ownerJID, types.DefaultUserServer)
	client.SendMessage(context.Background(), ownerJIDParsed, &waE2E.Message{
		Conversation: proto.String(message),
	})
}

func getMessageConversation(message *waE2E.Message) string {
	if message.GetConversation() != "" {
		return message.GetConversation()
	}

	if message.GetExtendedTextMessage() != nil {
		return message.GetExtendedTextMessage().GetContextInfo().GetQuotedMessage().GetConversation()
	}

	if message.GetImageMessage() != nil {
		if message.GetImageMessage().GetCaption() != "" {
			return message.GetImageMessage().GetCaption()
		}
	}

	if message.GetVideoMessage() != nil {
		if message.GetVideoMessage().GetCaption() != "" {
			return message.GetVideoMessage().GetCaption()
		}
	}

	if message.GetDocumentMessage() != nil {
		if message.GetDocumentMessage().GetCaption() != "" {
			return message.GetDocumentMessage().GetCaption()
		}
	}

	if message.GetReactionMessage() != nil {
		return message.GetReactionMessage().GetText()
	}
	return "[Does not have text]"
}
