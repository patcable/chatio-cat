package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
)

type Config struct {
	ChannelID       string `json:"channel_id"`
	MessageDuration string `json:"message_duration"`
}

func handleRequest(ctx context.Context, event json.RawMessage) error {
	var config Config
	if err := json.Unmarshal(event, &config); err != nil {
		fmt.Printf("Failed to unmarshal event: %v", err)
		return err
	}

	token := os.Getenv("CHATIO_CAT_TOKEN")
	channel := config.ChannelID
	messageCount := 100
	durationString := config.MessageDuration

	msgDuration, err := time.ParseDuration(durationString)
	if err != nil {
		fmt.Printf("wrong time format (need to be a string in format 0h0m0s)")
		return err
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Printf("sorry: %s\n", err)
		return err
	}

	var allmessages []*discordgo.Message
	var before string
	for {
		msgs, err := dg.ChannelMessages(channel, messageCount, before, "", "")
		if err != nil {
			fmt.Printf("sorry: %s", err)
			return err
		}
		allmessages = append(allmessages, msgs...)

		if len(msgs) == (messageCount - 1) {
			// there may be more, so i guess get the next batch?
			before = msgs[messageCount-1].ID
		} else {
			// there are less, lets assume we caught em all?
			break
		}
	}

	var toDelete []string
	for _, v := range allmessages {
		if time.Since(v.Timestamp) > msgDuration {
			toDelete = append(toDelete, v.ID)
			// useful for troubleshooting but you dont actually need it
			//fmt.Printf("[%s] @ %s: %s\n", v.Author, v.Timestamp)
		}
	}

	// bulkdelete api only takes 100 messages
	chunked := slices.Chunk(toDelete, 100)
	for messagesToDelete := range chunked {
		err = dg.ChannelMessagesBulkDelete(channel, messagesToDelete)
		if err != nil {
			fmt.Printf("sorry: %s", err)
			return err
		}
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
