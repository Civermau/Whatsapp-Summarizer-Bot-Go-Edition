-- SQLite
DROP TABLE IF EXISTS messages;

CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chat_id TEXT NOT NULL,
    sender TEXT NOT NULL,
    message TEXT,
    message_type TEXT,
    timestamp DATETIME
);