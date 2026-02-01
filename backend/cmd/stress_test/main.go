package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/seka/reci-pin/backend/cmd/stress_test/scenario"
	"github.com/seka/reci-pin/backend/config"
)

var (
	// Flags
	targetURL  string
	configFile string

	// Global stats
	totalRequests int64
	totalErrors   int64
)

func init() {
	flag.StringVar(&targetURL, "target-url", "http://localhost:8080/api", "Target API URL for requests")
	flag.StringVar(&configFile, "config", "", "Path to configuration file (JSON)")
}

func main() {
	flag.Parse()

	var configs []config.ScenarioConfig
	if configFile != "" {
		// Load from config file
		data, err := os.ReadFile(configFile)
		if err != nil {
			log.Fatalf("Failed to read config file: %v", err)
		}

		if len(configs) == 0 {
			if err := json.Unmarshal(data, &configs); err != nil {
				log.Fatalf("Failed to parse config file (checked struct and array): %v", err)
			}
		}
	} else {
		log.Fatal("Please provide -config")
	}

	log.Printf("Target URL: %s", targetURL)
	var wg sync.WaitGroup

	for _, cfg := range configs {
		scenFunc, ok := scenario.Get(cfg.Name)
		if !ok {
			log.Printf("Unknown scenario: %s (skipping)", cfg.Name)
			continue
		}

		duration, err := time.ParseDuration(cfg.Duration)
		if err != nil {
			log.Printf("Invalid duration %s for %s: %v", cfg.Duration, cfg.Name, err)
			continue
		}

		log.Printf("Starting scenario '%s': Concurrency=%d, Duration=%s, Rate=%d",
			cfg.Name, cfg.Concurrency, duration, cfg.Rate)

		for i := 0; i < cfg.Concurrency; i++ {
			wg.Add(1)
			go func(c config.ScenarioConfig, fn scenario.Scenario, d time.Duration) {
				defer wg.Done()
				runWorker(c, fn, d, targetURL)
			}(cfg, scenFunc, duration)
		}
	}

	wg.Wait()
	log.Printf("Test Completed. Total Requests: %d, Errors: %d", atomic.LoadInt64(&totalRequests), atomic.LoadInt64(&totalErrors))
}

func runWorker(cfg config.ScenarioConfig, scen scenario.Scenario, duration time.Duration, targetURL string) {
	client := &http.Client{Timeout: 10 * time.Second}
	jar, _ := cookiejar.New(nil)
	client.Jar = jar

	ctx, cancel := context.WithTimeout(context.Background(), duration+1*time.Second) // Check deadline in loop
	defer cancel()

	endTime := time.Now().Add(duration)

	// Rate limiting logic
	var ticker *time.Ticker
	if cfg.Rate > 0 {
		// Very simple rate limiter per worker: Rate / Concurrency
		workerRate := float64(cfg.Rate) / float64(cfg.Concurrency)
		if workerRate < 1 {
			workerRate = 1
		} // Min 1
		interval := time.Duration(1e9 / workerRate)
		ticker = time.NewTicker(interval)
		defer ticker.Stop()
	}

	for time.Now().Before(endTime) {
		if ticker != nil {
			<-ticker.C
		}

		err := scen(ctx, client, targetURL)
		atomic.AddInt64(&totalRequests, 1)
		if err != nil {
			atomic.AddInt64(&totalErrors, 1)
			// log.Printf("Error: %v", err) // Verbose
		}
	}
}
