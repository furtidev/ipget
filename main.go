package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/akamensky/argparse"
)

type Response struct {
	Address      string `json:"query"`
	ISP          string `json:"isp"`
	AS           string `json:"as"`
	Country      string `json:"country"`
	Region       string `json:"regionName"`
	Organization string `json:"org"`
	Status       string `json:"status"`
	Message      string `json:"message"`
}

var client *http.Client

func fetchData(ip string, inter interface{}) error {
	resp, err := client.Get(fmt.Sprintf("http://ip-api.com/json/%s", ip))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(inter)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client = &http.Client{Timeout: 10 * time.Second}
	parser := argparse.NewParser("ipget", "fetch information about ip addresses")
	ip := parser.String("a", "address", &argparse.Options{Required: true, Help: "ip to check"})
	if err := parser.Parse(os.Args); err != nil {
		log.Fatal(parser.Usage(err))
	}

	var apiResp Response

	if err := fetchData(*ip, &apiResp); err != nil {
		log.Fatal(parser.Usage(err))
	}

	if apiResp.Status == "fail" {
		fmt.Printf("ERROR: %s\n", apiResp.Message)
		os.Exit(1)
	}

	fmt.Printf("IP:      %s \nCountry: %s \nRegion:  %s \nISP:     %s\nAS:      %s \nORG:     %s", apiResp.Address, apiResp.Country, apiResp.Region, apiResp.ISP, apiResp.AS, apiResp.Organization)
}
