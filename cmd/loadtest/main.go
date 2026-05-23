package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var queries = []string{
	"кроссовки", "платье", "наушники", "ноутбук", "телефон",
	"куртка", "джинсы", "кроссовки найк", "смартфон", "часы",
	"рюкзак", "кофта", "зарядка", "наушники беспроводные", "планшет",
}

var stopWords = []string{
	"казино", "ставки", "кредит", "займ", "18+",
}

type stats struct {
	total   atomic.Int64
	success atomic.Int64
	errors  atomic.Int64
	latency atomic.Int64 // суммарно в микросекундах
}

func (s *stats) record(start time.Time, ok bool) {
	s.total.Add(1)
	s.latency.Add(time.Since(start).Microseconds())
	if ok {
		s.success.Add(1)
	} else {
		s.errors.Add(1)
	}
}

func (s *stats) report(name string, elapsed time.Duration) {
	t := s.total.Load()
	if t == 0 {
		return
	}
	avgMs := float64(s.latency.Load()) / float64(t) / 1000.0
	fmt.Printf("  %-30s  запросов: %6d  RPS: %6.0f  avg: %5.1f ms  ошибок: %d\n",
		name, t, float64(t)/elapsed.Seconds(), avgMs, s.errors.Load())
}

func main() {
	addr := flag.String("addr", "localhost:8080", "адрес сервиса")
	dur := flag.Duration("duration", 15*time.Second, "длительность теста")
	readers := flag.Int("readers", 50, "горутин на GET /top")
	writers := flag.Int("writers", 5, "горутин на POST /event")
	stoplisters := flag.Int("stoplist", 2, "горутин на POST/DELETE /stoplist")
	flag.Parse()

	client := &http.Client{
		Transport: &http.Transport{MaxIdleConnsPerHost: 512},
		Timeout:   5 * time.Second,
	}

	var (
		readStats     stats
		writeStats    stats
		stoplistStats stats
	)

	deadline := time.Now().Add(*dur)
	var wg sync.WaitGroup
	start := time.Now()

	fmt.Printf("нагрузочный тест: %s, длительность %s\n", "http://"+*addr, *dur)
	fmt.Printf("горутин: readers=%d  writers=%d  stoplist=%d\n\n", *readers, *writers, *stoplisters)

	// GET /api/v1/top
	for range *readers {
		wg.Go(func() {
			url := "http://" + *addr + "/api/v1/top?n=10"
			for time.Now().Before(deadline) {
				t := time.Now()
				resp, err := client.Get(url)
				if err != nil {
					readStats.record(t, false)
					continue
				}
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				readStats.record(t, resp.StatusCode == http.StatusOK)
			}
		})
	}

	// POST /api/v1/event
	for range *writers {
		wg.Go(func() {
			url := "http://" + *addr + "/api/v1/event"
			r := rand.New(rand.NewSource(rand.Int63()))
			for time.Now().Before(deadline) {
				body, _ := json.Marshal(map[string]string{
					"query":   queries[r.Intn(len(queries))],
					"user_id": fmt.Sprintf("loadtest-user-%d", r.Intn(10000)),
				})
				t := time.Now()
				resp, err := client.Post(url, "application/json", bytes.NewReader(body))
				if err != nil {
					writeStats.record(t, false)
					continue
				}
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				writeStats.record(t, resp.StatusCode == http.StatusNoContent)
			}
		})
	}

	// POST/DELETE /api/v1/stoplist
	for range *stoplisters {
		wg.Go(func() {
			r := rand.New(rand.NewSource(rand.Int63()))
			for time.Now().Before(deadline) {
				word := stopWords[r.Intn(len(stopWords))]

				// добавляем
				body, _ := json.Marshal(map[string]string{"word": word})
				t := time.Now()
				resp, err := client.Post("http://"+*addr+"/api/v1/stoplist", "application/json", bytes.NewReader(body))
				if err != nil {
					stoplistStats.record(t, false)
				} else {
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
					stoplistStats.record(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusConflict)
				}

				// удаляем
				req, _ := http.NewRequest(http.MethodDelete, "http://"+*addr+"/api/v1/stoplist/"+word, nil)
				t = time.Now()
				resp, err = client.Do(req)
				if err != nil {
					stoplistStats.record(t, false)
				} else {
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
					stoplistStats.record(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound)
				}

				time.Sleep(100 * time.Millisecond)
			}
		})
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Println("результаты:")
	readStats.report("GET  /api/v1/top", elapsed)
	writeStats.report("POST /api/v1/event", elapsed)
	stoplistStats.report("POST+DELETE /api/v1/stoplist", elapsed)

	total := readStats.total.Load() + writeStats.total.Load() + stoplistStats.total.Load()
	fmt.Printf("\n  итого запросов: %d  за %s\n", total, elapsed.Round(time.Millisecond))
}
