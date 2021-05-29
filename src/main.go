package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type RawCurrency struct {
	Success   bool   `json:"success"`
	Timestamp int64  `json:"timestamp"`
	Base      string `json:"base"`
	Date      string `json:"date"`
	Rates     struct {
		USD float64 `json:"USD"`
		CAD float64 `json:"CAD"`
		IDR float64 `json:"IDR"`
		GBP float64 `json:"GBP"`
		CHF float64 `json:"CHF"`
		SGD float64 `json:"SGD"`
		INR float64 `json:"INR"`
		MYR float64 `json:"MYR"`
		JPY float64 `json:"JPY"`
		KRW float64 `json:"KRW"`
	} `json:"rates"`
}

type ResponseMultiCurrency struct {
	Success   bool   `json:"success"`
	Timestamp int64  `json:"timestamp"`
	Base      string `json:"base"`
	Rates     []Rate `json:"rates"`
}

type Rate struct {
	Name string  `json:"name"`
	Rate float64 `json:"rate"`
}

type ResponseSingleCurrency struct {
	Success   bool    `json:"success"`
	Timestamp int64   `json:"timestamp"`
	Base      string  `json:"base"`
	Target    string  `json:"target"`
	Rate      float64 `json:"rate"`
}

var rawcurrency RawCurrency

func Lists(c *gin.Context) {
	respcurrency := &ResponseMultiCurrency{
		Success:   rawcurrency.Success,
		Timestamp: time.Now().Unix(),
		Base:      rawcurrency.Base,
		Rates:     []Rate{},
	}

	v := reflect.ValueOf(rawcurrency.Rates)
	typeOfS := v.Type()

	// Iterate struct Rates from RawCurrency and append to ResponseMultiCurrency Struct
	for i := 0; i < v.NumField(); i++ {
		name := typeOfS.Field(i).Name
		rate := v.Field(i).Interface().(float64)
		if rate != 0 {
			respcurrency.Rates = append(respcurrency.Rates, Rate{
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
			respcurrency := ResponseSingleCurrency{
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

func main() {
	FetchCurrencyAPI()
	go SchedulerFetchCurrencyAPI()
	r := gin.Default()

	r.GET("/api/currency", Lists)
	r.GET("/api/currency/:symbols", ListItem)
	r.Run()
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
	//We make HTTP request using the Get function
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
