package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {
	// DEBUG LEVELS: DEBUG, INFO, WARN, ERROR, FATAL
	// This determines the level of logging in console it seems
	dbLog := waLog.Stdout("Database", "WARN", true)
	ctx := context.Background()
	container, err := sqlstore.New(ctx, "sqlite3", "file:work.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	// Initialize database connection for message storage
	var err2 error
	db, err2 = initDatabase()
	if err2 != nil {
		panic(err2)
	}
	defer db.Close()

	// Initialize WhatsApp client
	client, err := initWhatsAppClient(ctx, container)
	if err != nil {
		panic(err)
	}

	// Connect to WhatsApp
	if err := connectToWhatsApp(client); err != nil {
		panic(err)
	}

	// Record the time when the bot successfully connects
	// This will be used to filter out historical messages
	botStartTime = time.Now()
	fmt.Printf("Bot started at: %s\n", botStartTime.Format("2006-01-02 15:04:05"))
	fmt.Println("Bot is now listening for NEW messages only...")

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
