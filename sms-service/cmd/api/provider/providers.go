package provider

import (
	"encoding/json"
	"log"
	"os"
	"sms-service/cmd/api/utils"
)

var logEvent utils.LogPayload

type Configuration struct {
	Providers struct {
		Rahyab struct {
			Username string
			Password string
			Company  string
			BaseUrl  string
		}
		PG struct {
			Username string
			Password string
			BaseUrl  string
		}
		Kavenegar struct {
			ApiKey  string
			BaseUrl string
		}
		HamyarSMS struct {
			Username string
			Password string
			BaseUrl  string
		}
	}
}

var ProviderConfig Configuration

func init() {
	// read config.json file
	file, _ := os.Open("/app/config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&ProviderConfig)
	if err != nil {
		log.Panic(err)
	}
}
