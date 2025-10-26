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
		fmt.Printf("[chatio-cat][%s] Wrong time format (need to be a string in format 0h0m0s)", channel)
		return err
	}

	fmt.Printf("[chatio-cat][%s] Going to remove messages older than %s.\n", channel, durationString)

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Printf("[chatio-cat][%s] cant initialize client: %s\n", channel, err)
		return err
	}

	var allmessages []*discordgo.Message
	var before string
	for {
		msgs, err := dg.ChannelMessages(channel, messageCount, before, "", "")
		if err != nil {
			fmt.Printf("[chatio-cat][%s] cant collect messages: %s\n", channel, err)
			return err
		}
		allmessages = append(allmessages, msgs...)

		if len(msgs) == messageCount {
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

	if len(toDelete) > 0 {
		fmt.Printf("[chatio-cat][%s] Removing %d messages.\n", channel, len(toDelete))

		// bulkdelete api only takes 100 messages
		chunked := slices.Chunk(toDelete, 100)
		for messagesToDelete := range chunked {
			err = dg.ChannelMessagesBulkDelete(channel, messagesToDelete)
			if err != nil {
				fmt.Printf("[chatio-cat][%s] cant delete messages: %s", channel, err)
				return err
			}
		}
	} else {
		fmt.Printf("[chatio-cat][%s] Nothing to delete.\n", channel)
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
