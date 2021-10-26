package main

import (
	"os"
	"io"
	"log"
	"fmt"
	"sort"
	"time"
	"sync"
	"strconv"
	"crypto/rand"
	"encoding/binary"
	"github.com/nats-io/nats.go"
	hw "github.com/tylertreat/hdrhistogram-writer"
	hdrhistogram "github.com/HdrHistogram/hdrhistogram-go"
)

const PUB_RATE_PER_SEC = 10
const MSG_SIZE = 128
const NO_OF_MSGS = 1000
const FSECS = float64(time.Second)
var SEQUENTIAL_BATCH, _ = strconv.Atoi(os.Getenv("SEQUENTIAL_BATCH"))

func rps(count int, elapsed time.Duration) int {
	return int(float64(count) / (float64(elapsed) / FSECS))
}

func NatsLoadTest(natsConn *nats.Conn, subject string, wg *sync.WaitGroup) {
	defer wg.Done()

	var sub string
	var subWG sync.WaitGroup
	received := make([]int, SEQUENTIAL_BATCH)

	durations := make([]time.Duration, 0, SEQUENTIAL_BATCH*NO_OF_MSGS)
	maxLatency := time.Second * 0

	for si:=0; si<SEQUENTIAL_BATCH; si++ {
		sub = fmt.Sprintf("%s-%d", subject, si)
		go func(ind int) {
			subWG.Add(1)
			natsConn.Subscribe(sub, func(msg *nats.Msg) {
				sendTime := int64(binary.LittleEndian.Uint64(msg.Data))
				dur := time.Duration(time.Now().UnixNano()-sendTime).Truncate(time.Microsecond)
				if dur > maxLatency {
					maxLatency = dur
				}
				durations = append(durations, dur)
				received[ind]++
				if received[ind] % 300 == 0 {
					log.Println("Received N Msg", received[ind], sub)
				}
				if received[ind] == NO_OF_MSGS {
					log.Println("##########Sub Done#############", sub)
					subWG.Done()
				}
			})
		}(si)
	}

	var pubWG sync.WaitGroup
	pubWG.Add(1)
	data := make([]byte, MSG_SIZE)

	ticker := time.NewTicker(time.Second / time.Duration(PUB_RATE_PER_SEC))
	defer ticker.Stop()
	done := make(chan bool)

	start := time.Now()
	var msgsPublished = 0
	go func() {
		for {
			select {
			case <-done:
				log.Println("##########RECEIVED PUB DONE###########")
				pubWG.Done()
				ticker.Stop()
				return
			case t := <-ticker.C:
				//log.Println("Ticker Ticked !")
				go func () {
					for pi:=0; pi<SEQUENTIAL_BATCH && msgsPublished<NO_OF_MSGS; pi++ {
						sub = fmt.Sprintf("%s-%d", subject, pi)
						//log.Printf("Publishing to %s, %d", sub, msgsPublished, NO_OF_MSGS)
						io.ReadFull(rand.Reader, data)
						binary.LittleEndian.PutUint64(data[0:], uint64(t.UnixNano()))
						natsConn.Publish(sub, data)
					}
					natsConn.Flush()
					msgsPublished += 1
					if msgsPublished == NO_OF_MSGS {
						done <- true
					}
				}()
			}
		}
	}()
	log.Printf("######WAITING#####")
	pubWG.Wait()
	avgMsgRate := rps(NO_OF_MSGS, time.Since(start))
	subWG.Wait()
	//log.Println("Durations", durations)
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	h := hdrhistogram.New(1, int64(durations[len(durations)-1]), 5)
	for _, d := range durations {
		h.RecordValue(int64(d))
	}
	pctls := hw.Percentiles{10, 25, 50, 75, 90, 99, 99.9, 99.99, 99.999, 99.9999, 99.99999, 100.0}
	hw.WriteDistributionFile(h, pctls, 1.0/1000000.0, fmt.Sprintf("/tmp/Durations-%s.histogram", subject))
	log.Println("Avg Msg Rate Reached : ", avgMsgRate)
	log.Println("Max Latency Reached : ", maxLatency)
}