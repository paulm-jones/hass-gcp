package main

// Remember to go get "cloud.google.com/go/bigquery"
import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type HomeAssistantMessage struct {
	Entity_id   string     `json:"entity_id"` //e.g. "entity_id": "sensor.kitchen_sensor_temperature"
	State       string     `json:"state"`     //e.g. "state": "18.06",
	Attributes  Attributes `json:"attributes"`
	LastChanged string     `json:"last_changed"`
}

type Attributes struct {
	BatteryLevel int    `json:"battery_level"` //48
	FriendlyName string `json:"friendly_name"` //Kitchen sensor temperature
	DeviceClass  string `json:"device_class"`  //temperature
}

type TemperatureTable struct {
	Datetime    civil.DateTime
	Temperature float64
	Room        string
}

const bigQueryDatasetName = "homeassistant"
const bigQueryTableName = "temperature"

func InsertTemperatureIntoBigQuery(ctx context.Context, m PubSubMessage) error {

	message := getHomeAssistantMessage(m.Data)

	if message.Entity_id == "sensor.kitchen_sensor_temperature" || message.Entity_id == "sensor.landing_sensor_temperature" {

		client, err := bigquery.NewClient(ctx, "proven-solstice-255019")
		if err != nil {
			log.Fatalf("Failed to create BigQuery client because [%v]", err)
		}

		table := client.Dataset(bigQueryDatasetName).Table(bigQueryTableName)
		inserter := table.Inserter()

		item := getItemToInsert(message)
		if err := inserter.Put(ctx, item); err != nil {
			log.Fatalf("Failed to insert [%+v] into BigQuery table because [%v]", item, err)
		}

	}
	return nil
}

func getHomeAssistantMessage(m []byte) HomeAssistantMessage {

	log.Printf("Received raw args [%s]", m)

	var message HomeAssistantMessage
	if err := json.Unmarshal(m, &message); err != nil {
		log.Fatalf("Failed to unmarshal byte array to JSON %v", err)
	}

	log.Printf("Unmarshalled args as [%+v]", message)

	return message
}

func getItemToInsert(message HomeAssistantMessage) TemperatureTable {
	temp, _ := strconv.ParseFloat(message.State, 64)
	datetime := getCivilDateTime(message.LastChanged)
	item := TemperatureTable{
		Datetime:    datetime,
		Temperature: temp,
		Room:        message.Entity_id,
	}
	return item
}

func getCivilDateTime(lastChanged string) civil.DateTime {

	cleanString := strings.Replace(lastChanged, "\"", "", 2)
	dateTimeNoMicroseconds := getDecimalLessDatetime(cleanString)
	datetime, err := civil.ParseDateTime(dateTimeNoMicroseconds)
	if err != nil {
		log.Fatalf("Failed to parse civil Datetime from datetime string [%s] because of [%v]", dateTimeNoMicroseconds, err)
	}
	return datetime
}

func getDecimalLessDatetime(s string) string {
	return s[0:19]
}
