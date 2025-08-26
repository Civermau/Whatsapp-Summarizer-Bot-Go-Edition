package main

import (
	"database/sql"
	"time"

	"github.com/go-deepseek/deepseek"
)

// Global variables
var (
	botStartTime time.Time
	db           *sql.DB
	dsClient, _  = deepseek.NewClient("YOUR_API_KEY_HERE")
	ownerJID     = "YOUR_PHONE_NUMBER_HERE"
)
