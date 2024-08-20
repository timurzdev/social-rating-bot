CREATE TABLE IF NOT EXISTS levels (
    "level" SERIAL NOT NULL,
    "media" BOOLEAN NOT NULL,
    "messages" BOOLEAN NOT NULL,
    "description" TEXT NOT NULL,
    PRIMARY KEY (level)
);

CREATE TABLE IF NOT EXISTS ratings (
    "user_id" INTEGER NOT NULL,
    "chat_id" INTEGER NOT NULL,
    "rating" INTEGER NOT NULL,
    "level" INTEGER NOT NULL,
    UNIQUE ("user_id", "chat_id"),
    FOREIGN KEY (level) REFERENCES levels (level)
);

CREATE TABLE IF NOT EXISTS events (
    "message_id" INTEGER NOT NULL,
    "user_id" INTEGER NOT NULL,
    "created_at" TEXT NOT NULL,
    "is_deleted" BOOLEAN NOT NULL DEFAULT FALSE
);
