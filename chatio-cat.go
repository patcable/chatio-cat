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
		fmt.Printf("[chatio-cat][%s][] Wrong time format (need to be a string in format 0h0m0s)", channel)
		return err
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Printf("[chatio-cat][%s][] cant initialize client: %s\n", channel, err)
		return err
	}

	channelInfo, err := dg.Channel(channel)
	if err != nil {
		fmt.Printf("[chatio-cat][%s][] cant resolve channel: %s\n", channel, err)
	}

	fmt.Printf("[chatio-cat][%s][%s] Going to remove messages older than %s.\n", channel, channelInfo.Name, durationString)

	var toDelete []string
	var before string
	for {
		msgs, err := dg.ChannelMessages(channel, messageCount, before, "", "")
		if err != nil {
			fmt.Printf("[chatio-cat][%s][%s] cant collect messages: %s\n", channel, channelInfo.Name, err)
			return err
		}

		for _, v := range msgs {
			if time.Since(v.Timestamp) > msgDuration {
				toDelete = append(toDelete, v.ID)
				// useful for troubleshooting but you dont actually need it
				//fmt.Printf("[%s] @ %s: %s\n", v.Author, v.Timestamp)
			}
		}

		if len(msgs) == messageCount {
			// there may be more, so i guess get the next batch?
			before = msgs[messageCount-1].ID
		} else {
			// there are less, lets assume we caught em all?
			break
		}
	}

	if len(toDelete) > 0 {
		fmt.Printf("[chatio-cat][%s][%s] Removing %d messages.\n", channel, channelInfo.Name, len(toDelete))

		// bulkdelete api only takes 100 messages
		chunked := slices.Chunk(toDelete, 100)
		for messagesToDelete := range chunked {
			err = dg.ChannelMessagesBulkDelete(channel, messagesToDelete)
			if err != nil {
				fmt.Printf("[chatio-cat][%s][%s] cant delete messages: %s", channel, channelInfo.Name, err)
				return err
			}
		}
	} else {
		fmt.Printf("[chatio-cat][%s][%s] Nothing to delete.\n", channel, channelInfo.Name)
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
