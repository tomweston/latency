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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ably/ably-go/ably"
	"github.com/jedib0t/go-pretty/table"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	names "github.com/tomweston/latency/pkg/namesgenerator"
	"github.com/tomweston/latency/pkg/utils"
)

const (
	timeout = 30
)

var (
	AblyChannel string
	AblyKey     string
	AblyEvent   string
)

// subscribeCmd represents the subscribe command
var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to messages on a channel",
	Long:  `Subscribe to messages on a channel on the Ably Realtime API and generate a latency report`,
	Run: func(cmd *cobra.Command, args []string) {

		seed := time.Now().UTC().UnixNano()
		ng := names.NewNameGenerator(seed)
		id := ng.Generate()

		AblyChannel := cmd.Flag("channel").Value.String()
		AblyEvent := cmd.Flag("event").Value.String()

		AblyKey := utils.GetEnv("ABLY_KEY", "")

		// Create a new Ably client
		client, err := ably.NewRealtime(
			ably.WithKey(AblyKey),
			ably.WithClientID(id))
		if err != nil {
			panic(err)
		}

		log.WithFields(log.Fields{
			"channel": AblyChannel,
			"client":  id,
			"event":   AblyEvent,
			"timeout": fmt.Sprint(timeout) + " Seconds",
		}).Info("Listening for Messages...")

		DurationOfTime := time.Duration(1) * time.Second
		f := func() {
			checkSubscribeToEvent(client, AblyChannel, AblyEvent)
		}
		Timer1 := time.AfterFunc(DurationOfTime, f)
		defer Timer1.Stop()
		time.Sleep(timeout * time.Second)

		// TODO: Handle unsubsribe
		// TODO: Handle no messages received
		GenerateLatencyReport(AblyChannel, id)

	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)

	// log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	subscribeCmd.Flags().StringVarP(&AblyChannel, "channel", "c", "", "The channel to subscribe to")
	subscribeCmd.Flags().StringVarP(&AblyEvent, "event", "e", "", "The channel event to subscribe to")
}

// CheckSubscribeToEvent checks if the channel is subscribed to the event
func checkSubscribeToEvent(client *ably.Realtime, AblyChannel, AblyEvent string) func() {
	c := client.Channels.Get(AblyChannel)
	unsubscribe := SubscribeToEvent(c, AblyEvent)
	return unsubscribe
}

// SubscribeToEvent subscribes to the event
func SubscribeToEvent(channel *ably.RealtimeChannel, AblyEvent string) func() {

	// Subscribe to messages sent on the channel with given eventName
	unsubscribe, err := channel.Subscribe(context.Background(), AblyEvent, func(msg *ably.Message) {

		// Save message to parts separate by ":"
		msgParts := strings.Split(fmt.Sprint(msg.Data), ":")
		messageID := msgParts[0]
		sentTimestamp := msgParts[1]
		sentTimestampInt, parseErr := strconv.ParseInt(sentTimestamp, 10, 64)
		if parseErr != nil {
			log.Error("Error parsing sent timestamp")
		}

		// Get timestamp in Microseconds
		t := time.Now()
		tUnixMicro := int64(time.Nanosecond) * t.UnixNano() / int64(time.Microsecond)

		// Calculate Delay
		delay := tUnixMicro - sentTimestampInt

		log.WithFields(log.Fields{
			"channel":  channel.Name,
			"client":   msg.ClientID,
			"sent":     sentTimestamp,
			"recieved": tUnixMicro,
			"id":       messageID,
			"delay":    delay,
		}).Info("Recieved message")

		// Save to file
		reportData := map[string]string{
			"channel":  channel.Name,
			"client":   msg.ClientID,
			"sent":     sentTimestamp,
			"recieved": fmt.Sprint(tUnixMicro),
			"id":       messageID,
			"delay":    fmt.Sprint(delay),
		}

		report, _ := json.Marshal(reportData)
		err := ioutil.WriteFile(messageID+".json", report, 0644)
		if err != nil {
			err := fmt.Errorf("Error writing data to file: %w", err)
			fmt.Println(err)
		}

	})

	if err != nil {
		err := fmt.Errorf("Error subscribing to channel: %w", err)
		fmt.Println(err)
	}
	return unsubscribe
}

// Reports is a struct to hold reports
type Reports struct {
	Report []Reports
}

// Report is a struct to hold the data for the report
type Report struct {
	ClientID string `json:"client"`
	Channel  string `json:"channel"`
	ID       int64  `json:"id,string"`
	Recieved int64  `json:"recieved,string"`
	Sent     int64  `json:"sent,string"`
	Delay    int64  `json:"delay,string"`
}

// GenerateLatencyReport generates a report of the messages recieved
func GenerateLatencyReport(AblyChannel, id string) {

	// Let's first read the `0.json` file
	r0, err := ioutil.ReadFile("0.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Unmarshall the data into `report0`
	var report0 Report
	err = json.Unmarshal(r0, &report0)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	// Let's first read the `0.json` file
	r1, err := ioutil.ReadFile("1.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Unmarshall the data into `report1`
	var report1 Report
	err = json.Unmarshal(r1, &report1)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	// Let's first read the `0.json` file
	r2, err := ioutil.ReadFile("2.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Unmarshall the data into `report2`
	var report2 Report
	err = json.Unmarshal(r2, &report2)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	log.WithFields(log.Fields{
		"channel": AblyChannel,
		"client":  id,
	}).Info("Generating Latency Report...")

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "SENT", "RECEIVED", "LATENCY"})
	t.AppendRows([]table.Row{
		{0, report0.Sent, report0.Recieved, report0.Delay},
		{1, report1.Sent, report1.Recieved, report1.Delay},
		{2, report2.Sent, report2.Recieved, report2.Delay},
	})

	// Get Average Latency
	average := (report0.Delay + report1.Delay + report2.Delay) / 3

	t.AppendFooter(table.Row{"", "", "AVERAGE", average})
	t.SetStyle(table.StyleLight)

	// Render Latency Table
	t.Render()
}
