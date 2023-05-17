package metrics

import (
	"fmt"
	"github.com/dejano-with-tie/fantastigo/internal/common/util/collection"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type PrometheusSourcer struct {
	Storer Storer
}

func NewPrometheusSourcer(storer Storer) *PrometheusSourcer {
	measurements := make(chan []Measurement)
	sourcer := PrometheusSourcer{Storer: storer}

	go sourcer.run(measurements, time.Second*20)

	go func() {
		sourcer.store(measurements)
	}()

	return &sourcer
}

func (p *PrometheusSourcer) run(measurements chan []Measurement, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// poll once before ticking
	measurements <- p.poll()

	for range ticker.C {
		measurements <- p.poll()
	}
}

func (p *PrometheusSourcer) store(measurements chan []Measurement) {
	for cv := range measurements {
		for _, value := range cv {
			if err := p.Storer.Put(value); err != nil {
				fmt.Printf("failed storing; error: %v\n", err)
			}
		}
	}
}

var rgx = regexp.MustCompile("(?P<name>^.*){(?P<tags>.*)} (?P<value>.*) (?P<timestamp>.*)")

func (p *PrometheusSourcer) poll() []Measurement {
	lines := p.getMetrics()
	result := make([]Measurement, 0, len(lines))
	for _, line := range lines {
		match := rgx.FindStringSubmatch(line)
		if match == nil {
			fmt.Printf("failed to parse: %v\n", line)
			continue
		}

		measurement, err := p.parse(match)
		if err != nil {
			fmt.Printf("failed to parse: %v\n; error: %v", line, err)
		}
		result = append(result, measurement)
	}

	return result
}

func (p *PrometheusSourcer) parse(match []string) (Measurement, error) {
	named := named(match)
	timestamp, err := strconv.ParseInt(named["timestamp"], 10, 64)
	if err != nil {
		timestamp = time.Now().UnixMilli()
	}
	//host="telegraf",interface="all",swBuild="1.1",swType="sw-type-docker",thingName="vin123"
	tags := make(map[string]string)
	for _, tag := range strings.Split(named["tags"], ",") {
		k, v := func(str string) (string, string) {
			v := strings.Split(tag, "=")
			return v[0], v[1]
		}(tag)

		tags[k] = strings.ReplaceAll(v, "\"", "")
	}

	return Measurement{
		Name:      named["name"],
		Tags:      tags,
		Value:     named["value"],
		Timestamp: time.UnixMilli(timestamp),
	}, nil
}

func (p *PrometheusSourcer) getMetrics() []string {
	resp, err := http.Get("http://localhost:9273/metrics") // prometheus url
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return []string{}
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return []string{}
	}

	lines := strings.Split(string(raw), "\n")
	return collection.Filter(lines, func(s string) bool { return !strings.HasPrefix(s, "#") })
}

func named(match []string) map[string]string {
	result := make(map[string]string)
	for i, name := range rgx.SubexpNames() {
		if i != 0 && name != "" {
			//fmt.Printf("[%s] = %s\n", name, match[i])
			result[name] = match[i]
		}
	}
	return result
}
