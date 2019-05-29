// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//
// This is a simple http server which generates 500s randomly a percentage of
// the time.
package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	errorRatio     = stats.Float64("configured_error_ratio", "configured error ratio", stats.UnitDimensionless)
	errorRatioView = &view.View{
		Name:        "example/configured_error_ratio",
		Measure:     errorRatio,
		Description: "The current configured error ratio.",
		Aggregation: view.LastValue(),
	}
)

func init() {
	// Set a default error rate
	err := SetErrorRate(context.Background(), 0.001)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}
	log.Printf("Starting up on http://localhost:%s", port)

	pe, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		log.Fatalf("Failed to create Prometheus exporter: %v", err)
	}
	view.RegisterExporter(pe)

	err = view.Register(errorRatioView)
	if err != nil {
		log.Fatalf("Failed to register views: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		JSON(w, http.StatusNotFound, map[string]string{
			"error": "404: This page could not be found",
		})
	})

	r.Handle("/metrics", pe)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, http.StatusOK, map[string]string{
			"healthy": "true",
		})
	})

	r.Get("/quitquitquit", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("/quitquitquit called, exiting")
		os.Exit(1)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		rate, err := GetErrorRate(r.Context())
		if err != nil {
			log.Printf(err.Error())
			w.WriteHeader(500)
			return
		}

		if float64(rand.Intn(1000)) <= rate*1000 {
			w.WriteHeader(500)
			return
		}

		JSON(w, http.StatusOK, map[string]string{
			"Hello": "World",
		})
	})

	r.Get("/errors", func(w http.ResponseWriter, r *http.Request) {
		rate, err := GetErrorRate(r.Context())
		if err != nil {
			log.Printf(err.Error())
			w.WriteHeader(500)
			return
		}

		JSON(w, http.StatusOK, map[string]float64{
			"rate": rate,
		})
	})

	r.Get("/errors/{percent}", func(w http.ResponseWriter, r *http.Request) {
		rate, err := strconv.ParseFloat(chi.URLParam(r, "percent"), 64)
		if err != nil {
			log.Printf(err.Error())
			w.WriteHeader(500)
			return
		}

		if rate < 0 || rate > 100 {
			log.Printf("rate out of range")
			w.WriteHeader(500)
			return
		}

		err = SetErrorRate(r.Context(), rate)
		if err != nil {
			log.Printf(err.Error())
			w.WriteHeader(500)
			return
		}

		JSON(w, http.StatusOK, map[string]string{
			"status": "success",
		})
	})

	h := &ochttp.Handler{Handler: r}
	if err := view.Register(ochttp.DefaultServerViews...); err != nil {
		log.Fatal("Failed to register ochttp.DefaultServerViews")
	}

	log.Fatal(http.ListenAndServe(":"+port, h))
}

// JSON takes a piece of data and turns it into json and writes it out to the
// response with the correct headers.
func JSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	result, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	return nil
}

func SetErrorRate(ctx context.Context, rate float64) error {
	fp := filepath.Join(os.TempDir(), "rate.txt")
	content := []byte(strconv.FormatFloat(rate, 'E', -1, 64))

	err := ioutil.WriteFile(fp, content, 0644)
	if err != nil {
		return err
	}

	stats.Record(ctx, errorRatio.M(rate))

	return nil
}

func GetErrorRate(ctx context.Context) (float64, error) {
	fp := filepath.Join(os.TempDir(), "rate.txt")
	rateString, err := ioutil.ReadFile(fp)
	if err != nil {
		return 0, err
	}

	rate, err := strconv.ParseFloat(string(rateString), 64)
	if err != nil {
		return 0, err
	}

	stats.Record(ctx, errorRatio.M(rate))

	return rate, nil
}
