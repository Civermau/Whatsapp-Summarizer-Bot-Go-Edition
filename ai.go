package main

import (
	"context"
	"fmt"

	"github.com/go-deepseek/deepseek"
	"github.com/go-deepseek/deepseek/request"
	"go.mau.fi/whatsmeow"
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// Summarize messages using DeepSeek API
func summarizeMessages(client *whatsmeow.Client, chat types.JID, msgID types.MessageID, opts SummarizeOptions) {
	go func() {
		messages, err := getMessages(chat.User, opts.Count)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		lengthPrompt := ""
		switch opts.Style {
		case "short":
			lengthPrompt = "Summarize should be short, it should contain most of the information from the messages. it should contain the most important information from the messages. Make it short without losing any information."
		case "medium":
			lengthPrompt = "Summarize should be medium, it should contain most of the information from the messages. it should contain the most important information from the messages. Make it medium without losing any information. Don't make it too short. Don't make it too long. Make it just the right length."
		case "long":
			lengthPrompt = "Summarize should not be short, it should contain most of the information from the messages. Length does not matter, you can write as much as you want to make the summary as long as it contains most of the information from the messages."
		}

		dsRequest := &request.ChatCompletionsRequest{
			Model: deepseek.DEEPSEEK_CHAT_MODEL,
			Messages: []*request.Message{
				{
					Role: "system",
					Content: "You must answer in the same language as the users messages. You are Todoroki Hajime from hololive DEV_IS, but people call you bancho, you aim to be the #1 badass in the universe, you speak in first person\n" +
						"You will be given a list of messages, you must summarize them. " + lengthPrompt,
				},
				{
					Role:    "user",
					Content: "Summarize the following messages: " + messages,
				},
			},
		}

		response, err := dsClient.CallChatCompletionsChat(context.Background(), dsRequest)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		resp, err := client.SendMessage(context.Background(), chat, client.BuildEdit(chat, msgID, &waE2E.Message{
			Conversation: proto.String(response.Choices[0].Message.Content),
		}))
		if err != nil {
			fmt.Println("Error:", err)
		}

		fmt.Println("Response sent:", resp.ID)
	}()
}

// Handle normal AI requests asynchronously
func normalRequestAsync(prompt string, client *whatsmeow.Client, chat types.JID, msgID types.MessageID) {
	go func() {
		dsRequest := &request.ChatCompletionsRequest{
			Model: deepseek.DEEPSEEK_CHAT_MODEL,
			Messages: []*request.Message{
				{
					Role:    "system",
					Content: "You must answer in the same language as the user's message. You are Todoroki Hajime from hololive DEV_IS, but people call you bancho, you aim to be the #1 badass in the universe, you speak in first person",
				},
				{
					Role:    "user",
					Content: prompt,
				},
			},
		}

		response, err := dsClient.CallChatCompletionsChat(context.Background(), dsRequest)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		resp, err := client.SendMessage(context.Background(), chat, client.BuildEdit(chat, msgID, &waE2E.Message{
			Conversation: proto.String(response.Choices[0].Message.Content),
		}))
		if err != nil {
			fmt.Println("Error:", err)
		}

		fmt.Println("Response sent:", resp.ID)
	}()
}
