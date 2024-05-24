package grafana

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Neniel/gotennis/config"
	"github.com/Neniel/gotennis/log"
)

type Metric struct {
	Name     string   `json:"name"`
	Interval uint64   `json:"interval"`
	Value    float64  `json:"value"`
	Tags     []string `json:"tags"`
	Time     int64    `json:"time"`
}

var token string

func init() {
	c, err := config.LoadConfiguration()
	if err != nil {
		log.Logger.Error(err.Error())
		os.Exit(1)
	}

	token = c.Grafana.GraphiteToken

}

func SendMetric(name string, interval uint64, value float64, tags map[string]interface{}) {
	go func() {
		if token != "" {
			if appEnvironment := os.Getenv("APP_ENVIRONMENT"); appEnvironment == "" {
				log.Logger.Warn("environment variable APP_ENVIRONMENT is not set")
			} else {
				tags["environment"] = appEnvironment
			}

			hostname, _ := os.Hostname()
			tags["hostname"] = hostname

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
				log.Logger.Error(fmt.Errorf("error when marshaling: %w", err).Error())
				return
			}

			req, err := http.NewRequest(http.MethodPost, "https://graphite-prod-13-prod-us-east-0.grafana.net/graphite/metrics", strings.NewReader(string(bs)))
			if err != nil {
				log.Logger.Info(fmt.Errorf("error when creating request: %w", err).Error())
				return
			}

			req.Header.Add("Authorization", token)
			req.Header.Add("Content-Type", "application/json")

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Logger.Error(fmt.Errorf("error when sending metric: %w", err).Error())
				return
			}

			if res.StatusCode != http.StatusOK {
				log.Logger.Error(fmt.Errorf("error when sending metric: %v", res.StatusCode).Error())
			}

		}

	}()
}