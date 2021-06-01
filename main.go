package main

import (
	"encoding/json"
	"github.com/chrisvinsen/go-code/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

var rawcurrency models.RawCurrency

func main() {
	FetchCurrencyAPI()
	go SchedulerFetchCurrencyAPI()
	server := gin.Default()

	server.GET("/api/currency", Lists)
	server.GET("/api/currency/:symbols", ListItem)
	server.Run()
}

func Lists(c *gin.Context) {
	respcurrency := &models.ResponseMultiCurrency{
		Success:   rawcurrency.Success,
		Timestamp: time.Now().Unix(),
		Base:      rawcurrency.Base,
		Rates:     []models.Rate{},
	}

	v := reflect.ValueOf(rawcurrency.Rates)
	typeOfS := v.Type()

	// Iterate struct Rates from RawCurrency and append to ResponseMultiCurrency Struct
	for i := 0; i < v.NumField(); i++ {
		name := typeOfS.Field(i).Name
		rate := v.Field(i).Interface().(float64)
		if rate != 0 {
			respcurrency.Rates = append(respcurrency.Rates, models.Rate{
				Name: name,
				Rate: rate,
			})
		}
	}

	c.JSON(http.StatusOK, respcurrency)
}

func ListItem(c *gin.Context) {
	requested_symbols := strings.ToUpper(c.Param("symbols"))

	v := reflect.ValueOf(rawcurrency.Rates)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		name := typeOfS.Field(i).Name
		rate := v.Field(i).Interface().(float64)
		if name == requested_symbols {
			respcurrency := models.ResponseSingleCurrency{
				Success:   rawcurrency.Success,
				Timestamp: time.Now().Unix(),
				Base:      rawcurrency.Base,
				Target:    name,
				Rate:      rate,
			}
			c.JSON(http.StatusOK, respcurrency)
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": "Target currency not found"})
}

func SchedulerFetchCurrencyAPI() {
	for _ = range time.Tick(time.Second * 60) {
		FetchCurrencyAPI()
	}
}

func FetchCurrencyAPI() {
	//Build The URL string
	API_KEY := "b5f8c66b6ec3d2c964771b6b00fbae7c"
	SUPPORTED_SYMBOLS := "USD,CAD,IDR,GBP,CHF,SGD,INR,MYR,JPY,KRW"
	URL := "http://api.exchangeratesapi.io/v1/latest?access_key=" + API_KEY + "&symbols=" + SUPPORTED_SYMBOLS
	//We make HTTP request using the Get
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal("ooopsss an error occurred, please try again")
	}
	defer resp.Body.Close()
	//Create a variable of the same type as our model
	// var cResp Currency
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&rawcurrency); err != nil {
		log.Fatal("ooopsss! an error occurred, please try again")
	}
}
