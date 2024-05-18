package grafana

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Metric struct {
	Name     string   `json:"name"`
	Interval uint64   `json:"interval"`
	Value    float64  `json:"value"`
	Tags     []string `json:"tags"`
	Time     int64    `json:"time"`
}

type Client struct {
	http.Client
}

var once sync.Once
var client *http.Client
var token string

func SendMetric(name string, interval uint64, value float64, tags map[string]string) {
	once.Do(func() {
		client = &http.Client{}
		bsGrafanaGraphiteToken, err := os.ReadFile(os.Getenv("GRAFANA_GRAPHITE_TOKEN"))
		if err != nil {
			log.Fatalln(err.Error())
		}
		token = strings.Replace(string(bsGrafanaGraphiteToken), "\n", "", -1)
	})
	go func() {
		appEnvironment := os.Getenv("APP_ENVIRONMENT")
		if appEnvironment != "" {
			tags["APP_ENVIRONMENT"] = appEnvironment

			t := make([]string, 0)
			for k, v := range tags {
				t = append(t, fmt.Sprintf("%s=%s", k, v))
			}

			m := Metric{
				Name:     name,
				Interval: interval,
				Value:    value,
				Tags:     t,
				Time:     time.Now().UTC().Unix(),
			}

			mm := make([]Metric, 0)
			mm = append(mm, m)
			bs, err := json.Marshal(&mm)
			if err != nil {
				log.Println(fmt.Errorf("error when marshaling: %w", err))
				return
			}

			req, err := http.NewRequest(http.MethodPost, "https://graphite-prod-13-prod-us-east-0.grafana.net/graphite/metrics", strings.NewReader(string(bs)))
			if err != nil {
				log.Println(fmt.Errorf("error when creating request: %w", err))
				return
			}

			req.Header.Add("Authorization", token)
			req.Header.Add("Content-Type", "application/json")

			res, err := client.Do(req)
			if err != nil {
				log.Println(fmt.Errorf("error when sending metric: %w", err))
				return
			}

			log.Println(res.StatusCode)
		}

	}()
}
