package main

import (
	"cloud.google.com/go/civil"
	"reflect"
	"testing"
)

var wantTime, _ = civil.ParseDateTime("2019-10-13T15:52:41")

func Test_getDecimalLessDatetime(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Expect to have removed timezone and microseconds",
			struct{ s string }{s: "2019-10-13T15:52:41.045111+00:00"},
			"2019-10-13T15:52:41",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDecimalLessDatetime(tt.args.s); got != tt.want {
				t.Errorf("getDecimalLessDatetime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCivilDateTime(t *testing.T) {
	type args struct {
		lastChanged string
	}
	tests := []struct {
		name string
		args args
		want civil.DateTime
	}{
		{
			"Expect to be able to convert string to Civil datetime",
			struct{ lastChanged string }{lastChanged: "2019-10-13T15:52:41"},
			wantTime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCivilDateTime(tt.args.lastChanged); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCivilDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

var timeToInsert, _ = civil.ParseDateTime("2019-10-13T18:10:05")

func Test_getItemToInsert(t *testing.T) {
	type args struct {
		message HomeAssistantMessage
	}
	tests := []struct {
		name string
		args args
		want TemperatureTable
	}{
		{
			name: "",
			args: args{message:HomeAssistantMessage{
				Entity_id:   "sensor.landing_sensor_temperature",
				State:       "18.98",
				Attributes:  Attributes{
					BatteryLevel: 69,
					FriendlyName: "Landing sensor temperature",
					DeviceClass:  "temperature",
				},
				LastChanged: "\"2019-10-13T18:10:05.044067+00:00\"",
			}},
			want: TemperatureTable{
				Datetime:    timeToInsert,
				Temperature: 18.98,
				Room:        "sensor.landing_sensor_temperature",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getItemToInsert(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getItemToInsert() = %v, want %v", got, tt.want)
			}
		})
	}
}
