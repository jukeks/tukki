package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/jamiealquiza/tachymeter"
	kvv1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/kv/v1"
)

var (
	keyIndex    = flag.Int("key", 6, "The index of the key column")
	port        = flag.Int("port", 50051, "The server port")
	hostname    = flag.String("hostname", "localhost", "The server hostname")
	shouldWrite = flag.Bool("write", false, "Write data to the server")
	concurrency = flag.Int("concurrency", 1, "Number of concurrent readers")
)

func splitArray[T any](arr []T, n int) [][]T {
	var result [][]T
	chunkSize := len(arr) / n
	for i := 0; i < len(arr); i += chunkSize {
		end := i + chunkSize
		if end > len(arr) {
			end = len(arr)
		}
		chunkLen := end - i
		new := make([]T, chunkLen)
		copy(new, arr[i:end])
		result = append(result, new)
	}
	return result
}

func main() {
	flag.Parse()

	r := csv.NewReader(os.Stdin)
	records, err := r.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read CSV: %v\n", err)
		os.Exit(1)
	}
	recMap := make(map[string]string, len(records)-1)
	for _, record := range records[1:] {
		if len(record) > *keyIndex {
			recMap[record[*keyIndex]] = strings.Join(record, ",")
		}
	}
	keys := make([]string, 0, len(recMap))
	for key := range recMap {
		keys = append(keys, key)
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *hostname, *port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	client := kvv1.NewKvServiceClient(conn)
	ctx := context.Background()

	writeMeasurements := tachymeter.New(&tachymeter.Config{Size: 1000})
	if *shouldWrite {
		write(ctx, client, keys, recMap, writeMeasurements)
	}

	readMeasurements := tachymeter.New(&tachymeter.Config{Size: 1000})

	var wg sync.WaitGroup
	keySplits := splitArray(keys, *concurrency)
	fmt.Printf("Number of slices: %d\n", len(keySplits))
	start := time.Now()
	for i, split := range keySplits {
		wg.Add(1)
		go func(i int, keys []string) {
			defer wg.Done()
			read(ctx, client, keys, recMap, readMeasurements)
		}(i, split)
	}
	wg.Wait()
	fmt.Printf("Time taken: %v\n", time.Since(start))

	if *shouldWrite {
		fmt.Printf("Write results: %s\n", writeMeasurements.Calc().String())
	}
	fmt.Printf("Read results: %s\n", readMeasurements.Calc().String())
}

func read(ctx context.Context, client kvv1.KvServiceClient, keys []string, recMap map[string]string, readMeasurements *tachymeter.Tachymeter) {
	for _, key := range keys {
		start := time.Now()
		resp, err := client.Query(ctx, &kvv1.QueryRequest{
			Key: key,
		})
		readMeasurements.AddTime(time.Since(start))

		if err != nil {
			log.Fatalf("can not get key-value pair %v", err)
		}
		if resp.GetError() != nil {
			log.Fatalf("can not get key-value pair %v", resp.GetError().Message)
		}

		if resp.GetPair().Value != recMap[key] {
			log.Fatalf("expected\n%s, \n\n\ngot\n%s", recMap[key], resp.GetPair().Value)
		}

	}
}

func write(ctx context.Context, client kvv1.KvServiceClient, keys []string, recMap map[string]string, writeMeasurements *tachymeter.Tachymeter) {
	for _, key := range keys {
		start := time.Now()
		resp, err := client.Set(ctx, &kvv1.SetRequest{
			Pair: &kvv1.KvPair{
				Key:   key,
				Value: recMap[key],
			},
		})
		writeMeasurements.AddTime(time.Since(start))

		if err != nil {
			log.Fatalf("can not set key-value pair %v", err)
		}
		if resp.GetError() != nil {
			log.Fatalf("can not set key-value pair %v", resp.GetError().Message)
		}
	}
}
