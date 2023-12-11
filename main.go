package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)


var (
	accountSid string
	authToken string
	fromPhone string
	toPhone string
	client *twilio.RestClient
	wg sync.WaitGroup
)

var (
	pollingInterval = time.Second * 5
)

const (
	endpoint = "https://api.open-meteo.com/v1/forecast" //?latitude=52.52&longitude=13.41&hourly=temperature_2m
)

type WeatherData struct {
    Elevation float64               `json:"elevation"`
    Hourly    map[string]interface{} `json:"hourly"`
    Rain      map[string]interface{} `json:"rain"`
    Showers   map[string]interface{} `json:"showers"`
    Sunrise   map[string]interface{} `json:"sunrise"`
    Sunset    map[string]interface{} `json:"sunset"`
    UVIndex   map[string]interface{} `json:"uv_index_max"`
    UVIndexClearSky map[string]interface{} `json:"uv_index_clear_sky_max"`
}


type Sender interface {
	Send(*WeatherData) error
}

type SMSSender struct {
	number string
}

func NewSMSSender(number string) *SMSSender {
	return &SMSSender{
		number: number,
	}
}

func (s *SMSSender) Send(data *WeatherData) error {
	fmt.Println("Sending SMS to", s.number)
	return nil
}

type WeatherPoller struct {
	closeCh chan struct{}
	senders []Sender
}


func NewWeatherPoller(senders ...Sender) *WeatherPoller {
	return &WeatherPoller{
		closeCh: make(chan struct{}),
		senders: senders,
	}
}

func main() {
	smssender := NewSMSSender(fromPhone)
	weatherpoller := NewWeatherPoller(smssender)

	wg.Add(1)
	go func() {
		weatherpoller.start(&wg)
	}()

	wg.Wait()

	//time.Sleep(time.Second * 5)
	//weatherpoller.Close()

	// select {}
}

func (wp *WeatherPoller) Close() {
	close(wp.closeCh)
}

func (wp *WeatherPoller) start(wg *sync.WaitGroup) {
	fmt.Println("Starting Weather Poller")
	ticker := time.NewTicker(pollingInterval)
outer:
	for {
		select {
		case <-ticker.C:
			data, err := getWeatherResult(12.5828, 77.0429)
			if err != nil {
				log.Fatal(err)
			}
			if err := wp.handleWeatherData(data); err != nil {
				log.Fatal(err)
			}
			// Send message after handling weather data
			sendMessage(data, wg)
		case <-wp.closeCh:
			break outer
		}
	}
	fmt.Println("Stopping Weather Poller Gracefully")
}

func (wp *WeatherPoller) handleWeatherData(data *WeatherData) error {

	fmt.Println(data)
	for _, sender := range wp.senders {
		if err := sender.Send(data); err != nil {

			fmt.Println(err)
		}
	}
	return nil
}

func getWeatherResult(lat, long float64) (*WeatherData, error) {

	// req, err := http.NewRequest("GET", endpoint, nil)
	// if err != nil {
	//      log.Fatal(err)
	// }
	// req.Header.Set("Accept", "application/json")
	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	//      log.Fatal(err)
	// }
	// defer resp.Body.Close()

	uri := fmt.Sprintf("%s?latitude=%0.2f&longitude=%0.2f&hourly=temperature_2m,rain,showers&daily=sunrise,sunset,uv_index_max,uv_index_clear_sky_max&timezone=auto", endpoint, lat, long)
	fmt.Println("--------------")
	fmt.Println(uri)
	fmt.Println("--------------")
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	var data WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal(err)
	}
	return &data, nil
}

// Implementing twilio
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error Loading .env: %s\n", err.Error())
		os.Exit(1)
	}

	accountSid = os.Getenv("ACCOUNT_SID") // Use global variables, not creating new ones
	authToken = os.Getenv("AUTH_TOKEN")
	fromPhone = os.Getenv("FROM_PHONE")
	toPhone = os.Getenv("TO_PHONE")

	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
}

func sendMessage(data *WeatherData, wg *sync.WaitGroup) {
	params := openapi.CreateMessageParams{}
	params.SetTo(toPhone)
	params.SetFrom(fromPhone)

	// Check if temperature_2m exists and is a slice of interfaces
	temperatures, ok := data.Hourly["temperature_2m"].([]interface{})
	if !ok {
		fmt.Printf("Failed to convert temperature_2m to []interface{}. Type is %T\n", data.Hourly["temperature_2m"])
		return
	}

	// Assuming you want to send the first temperature from the array
	if len(temperatures) > 0 {
		// Assuming each temperature in the array is a float64
		temperature, ok := temperatures[0].(float64)
		if !ok {
			fmt.Printf("Failed to convert temperature to float64. Type is %T\n", temperatures[0])
			return
		}
	
		//

		params.SetBody(fmt.Sprintf("-----Today's Temperature is-----\n: %0.2f", temperature))
		_, err := client.Api.CreateMessage(&params)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Message Sent")
	} else {
		fmt.Println("No temperature data available.")
	}
	wg.Done()
}

