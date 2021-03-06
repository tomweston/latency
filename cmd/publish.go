/*
Copyright © 2021 Tom Weston weston.tom@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ably/ably-go/ably"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	names "github.com/tomweston/latency/pkg/namesgenerator"
	utils "github.com/tomweston/latency/pkg/utils"
)

const (
	num      = 3
	pubDelay = 5
)

var pubAblyChannel string
var pubAblyEvent string

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a message to a channel",
	Long:  `Publish a message to a channel on the Ably Realtime API.`,
	Run: func(cmd *cobra.Command, args []string) {

		seed := time.Now().UTC().UnixNano()
		ng := names.NewNameGenerator(seed)
		id := ng.Generate()

		pubAblyChannel := cmd.Flag("channel").Value.String()
		pubAblyEvent := cmd.Flag("event").Value.String()
		// numMessages := cmd.Flag("number").Value.String()
		// pubDelay := cmd.Flag("delay").Value.String()
		AblyKey := utils.GetEnv("ABLY_KEY", "")

		// Create a new Ably client
		client, err := ably.NewRealtime(
			ably.WithKey(AblyKey),
			ably.WithClientID(id))
		if err != nil {
			panic(err)
		}

		// Setup a channel to publish to
		channel := client.Channels.Get(pubAblyChannel)

		// Convert the number of messages to an integer
		// n, err := strconv.Atoi(numMessages)
		// if err != nil {
		// 	fmt.Println(err)
		// }

		messages := [3]string{}
		for m := range messages {

			// Assign a timestamp to the message
			t := time.Now()
			tUnixMicro := int64(time.Nanosecond) * t.UnixNano() / int64(time.Microsecond)

			// Splits required data into parts separated by a colon
			CompositMessage := fmt.Sprint(m) + ":" + fmt.Sprint(tUnixMicro)

			// Publish the message
			Publish(channel, pubAblyEvent, CompositMessage)

			// Sleep for prescribed delay

			// delayTime, err := time.ParseDuration(pubDelay)
			// if err != nil {
			// 	log.WithFields(log.Fields{
			// 		"error": err,
			// 	}).Error("Error parsing delay")

			// }

			time.Sleep(pubDelay * time.Second)
		}

	},
}

func init() {
	rootCmd.AddCommand(publishCmd)

	// log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	publishCmd.Flags().StringVarP(&pubAblyChannel, "channel", "c", "", "The channel to subscribe to")
	publishCmd.Flags().StringVarP(&pubAblyEvent, "event", "e", "", "The channel event to subscribe to")
	// publishCmd.Flags().StringVarP(&numMessages, "number", "n", "", "The channel event to subscribe to")
	// publishCmd.Flags().StringVarP(&pubDelay, "delay", "d", "", "The delay between messages")

}

// Publish a message to the channel
func Publish(channel *ably.RealtimeChannel, AblyEvent, message string) {
	// Publish the message to Ably Channel
	err := channel.Publish(context.Background(), AblyEvent, message)
	if err != nil {
		err := fmt.Errorf("error publishing to channel: %w", err)
		fmt.Println(err)
	}
	log.WithFields(log.Fields{
		"channel": channel.Name,
		"event":   AblyEvent,
		"data":    message,
	}).Info("Successfully published message")
}
