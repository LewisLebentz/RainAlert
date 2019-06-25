package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"reflect"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

type GoogleMaps struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Bounds struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"bounds"`
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

type DarkSky struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Currently struct {
		Time                 int64   `json:"time"`
		Summary              string  `json:"summary"`
		Icon                 string  `json:"icon"`
		NearestStormDistance int     `json:"nearestStormDistance"`
		NearestStormBearing  int     `json:"nearestStormBearing"`
		PrecipIntensity      float64 `json:"precipIntensity"`
		PrecipProbability    float64 `json:"precipProbability"`
		Temperature          float64 `json:"temperature"`
		ApparentTemperature  float64 `json:"apparentTemperature"`
		DewPoint             float64 `json:"dewPoint"`
		Humidity             float64 `json:"humidity"`
		Pressure             float64 `json:"pressure"`
		WindSpeed            float64 `json:"windSpeed"`
		WindGust             float64 `json:"windGust"`
		WindBearing          int     `json:"windBearing"`
		CloudCover           float64 `json:"cloudCover"`
		UvIndex              int     `json:"uvIndex"`
		Visibility           float64 `json:"visibility"`
		Ozone                float64 `json:"ozone"`
	} `json:"currently"`
	Minutely struct {
		Summary string `json:"summary"`
		Icon    string `json:"icon"`
		Data    []struct {
			Time                 int64   `json:"time"`
			PrecipIntensity      float64 `json:"precipIntensity"`
			PrecipProbability    float64 `json:"precipProbability"`
			PrecipIntensityError float64 `json:"precipIntensityError,omitempty"`
			PrecipType           string  `json:"precipType,omitempty"`
		} `json:"data"`
	} `json:"minutely"`
	Hourly struct {
		Summary string `json:"summary"`
		Icon    string `json:"icon"`
		Data    []struct {
			Time                int     `json:"time"`
			Summary             string  `json:"summary"`
			Icon                string  `json:"icon"`
			PrecipIntensity     float64 `json:"precipIntensity"`
			PrecipProbability   float64 `json:"precipProbability"`
			PrecipType          string  `json:"precipType,omitempty"`
			Temperature         float64 `json:"temperature"`
			ApparentTemperature float64 `json:"apparentTemperature"`
			DewPoint            float64 `json:"dewPoint"`
			Humidity            float64 `json:"humidity"`
			Pressure            float64 `json:"pressure"`
			WindSpeed           float64 `json:"windSpeed"`
			WindGust            float64 `json:"windGust"`
			WindBearing         int     `json:"windBearing"`
			CloudCover          float64 `json:"cloudCover"`
			UvIndex             int     `json:"uvIndex"`
			Visibility          float64 `json:"visibility"`
			Ozone               float64 `json:"ozone"`
		} `json:"data"`
	} `json:"hourly"`
	Daily struct {
		Summary string `json:"summary"`
		Icon    string `json:"icon"`
		Data    []struct {
			Time                        int     `json:"time"`
			Summary                     string  `json:"summary"`
			Icon                        string  `json:"icon"`
			SunriseTime                 int     `json:"sunriseTime"`
			SunsetTime                  int     `json:"sunsetTime"`
			MoonPhase                   float64 `json:"moonPhase"`
			PrecipIntensity             float64 `json:"precipIntensity"`
			PrecipIntensityMax          float64 `json:"precipIntensityMax"`
			PrecipIntensityMaxTime      int     `json:"precipIntensityMaxTime"`
			PrecipProbability           float64 `json:"precipProbability"`
			PrecipType                  string  `json:"precipType"`
			TemperatureHigh             float64 `json:"temperatureHigh"`
			TemperatureHighTime         int     `json:"temperatureHighTime"`
			TemperatureLow              float64 `json:"temperatureLow"`
			TemperatureLowTime          int     `json:"temperatureLowTime"`
			ApparentTemperatureHigh     float64 `json:"apparentTemperatureHigh"`
			ApparentTemperatureHighTime int     `json:"apparentTemperatureHighTime"`
			ApparentTemperatureLow      float64 `json:"apparentTemperatureLow"`
			ApparentTemperatureLowTime  int     `json:"apparentTemperatureLowTime"`
			DewPoint                    float64 `json:"dewPoint"`
			Humidity                    float64 `json:"humidity"`
			Pressure                    float64 `json:"pressure"`
			WindSpeed                   float64 `json:"windSpeed"`
			WindGust                    float64 `json:"windGust"`
			WindGustTime                int     `json:"windGustTime"`
			WindBearing                 int     `json:"windBearing"`
			CloudCover                  float64 `json:"cloudCover"`
			UvIndex                     int     `json:"uvIndex"`
			UvIndexTime                 int     `json:"uvIndexTime"`
			Visibility                  float64 `json:"visibility"`
			Ozone                       float64 `json:"ozone"`
			TemperatureMin              float64 `json:"temperatureMin"`
			TemperatureMinTime          int     `json:"temperatureMinTime"`
			TemperatureMax              float64 `json:"temperatureMax"`
			TemperatureMaxTime          int     `json:"temperatureMaxTime"`
			ApparentTemperatureMin      float64 `json:"apparentTemperatureMin"`
			ApparentTemperatureMinTime  int     `json:"apparentTemperatureMinTime"`
			ApparentTemperatureMax      float64 `json:"apparentTemperatureMax"`
			ApparentTemperatureMaxTime  int     `json:"apparentTemperatureMaxTime"`
		} `json:"data"`
	} `json:"daily"`
	Alerts []struct {
		Title       string   `json:"title"`
		Regions     []string `json:"regions"`
		Severity    string   `json:"severity"`
		Time        int64    `json:"time"`
		Expires     int      `json:"expires"`
		Description string   `json:"description"`
		URI         string   `json:"uri"`
	} `json:"alerts"`
	Flags struct {
		Sources           []string `json:"sources"`
		MeteoalarmLicense string   `json:"meteoalarm-license"`
		NearestStation    float64  `json:"nearest-station"`
		Units             string   `json:"units"`
	} `json:"flags"`
	Offset int `json:"offset"`
}

func geocode(location string) (float64, float64) {
	apiKey := os.Getenv("MapsKey")
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s", location, apiKey)
	// fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}

	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	var record GoogleMaps

	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	for _, v := range record.Results[0].AddressComponents {
		if v.Types[0] == "postal_town" {
			fmt.Println(v.LongName)
		}
		if v.Types[0] == "route" {
			fmt.Println(v.LongName)
		}
	}

	lat := record.Results[0].Geometry.Location.Lat
	lng := record.Results[0].Geometry.Location.Lng

	return lat, lng
}

func weather(lat, lng float64) string {
	apiKey := os.Getenv("DarkSkyKey")
	url := fmt.Sprintf("https://api.darksky.net/forecast/%s/%f,%f", apiKey, lat, lng)
	// fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}

	defer resp.Body.Close()

	var record DarkSky

	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	summary := record.Minutely.Summary
	fmt.Println(summary)
	actualTime := time.Unix(record.Minutely.Data[0].Time, 0)

	for _, element := range record.Minutely.Data {
		if element.PrecipProbability > 0.2 && record.Currently.PrecipProbability < 0.2 {
			precipProb := element.PrecipProbability
			precipInt := element.PrecipIntensity
			precipType := element.PrecipType
			precipTime := time.Unix(element.Time, 0)
			if actualTime != precipTime {
				fmt.Println(precipProb, precipInt, precipType, precipTime)
				summary := fmt.Sprint(strings.Title(precipType), " starting in ", subtractTime(precipTime, actualTime), " mins")
				return summary
			}
		}
	}

	return summary
}

func subtractTime(time1, time2 time.Time) float64 {
	diff := time2.Sub(time1).Minutes()
	return -diff
}

func Handler(ctx context.Context) (Response, error) {
	lat, lng := geocode("rg248jz")
	summary := weather(lat, lng)
	fmt.Println(lat, lng)
	fmt.Println(summary)
	fmt.Println(reflect.TypeOf(summary))
	message := fmt.Sprintf(" { \"Message\" : \"%s\" } ", summary)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            message,
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "rainalert-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
