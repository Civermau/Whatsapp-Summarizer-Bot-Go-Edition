package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// Initialize WhatsApp client
func initWhatsAppClient(ctx context.Context, container *sqlstore.Container) (*whatsmeow.Client, error) {
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, err
	}

	clientLog := waLog.Stdout("Client", "WARN", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(func(evt interface{}) {
		eventHandler(evt, client)
	})

	return client, nil
}

// Connect to WhatsApp (new login or existing session)
func connectToWhatsApp(client *whatsmeow.Client) error {
	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, err := client.GetQRChannel(context.Background())
		if err != nil {
			return err
		}

		err = client.Connect()
		if err != nil {
			return err
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code in the terminal
				fmt.Println("Scan this QR code with WhatsApp:")
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err := client.Connect()
		if err != nil {
			return err
		}
		fmt.Println("Connected to WhatsApp")
	}

	return nil
}
