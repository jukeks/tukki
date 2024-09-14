package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/jamiealquiza/tachymeter"
	kvv1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/kv/v1"
)

var (
	keyIndex = flag.Int("key", 6, "The index of the key column")
	port     = flag.Int("port", 50051, "The server port")
	hostname = flag.String("hostname", "localhost", "The server hostname")
)

func main() {
	r := csv.NewReader(os.Stdin)

	writeMeasurements := tachymeter.New(&tachymeter.Config{Size: 1000})

	records, err := r.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read CSV: %v\n", err)
		os.Exit(1)
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *hostname, *port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	client := kvv1.NewKvServiceClient(conn)

	ctx := context.Background()

	for _, record := range records[1:] {
		if len(record) > *keyIndex {
			start := time.Now()
			fmt.Printf("Key: %s\n", record[*keyIndex])

			resp, err := client.Set(ctx, &kvv1.SetRequest{
				Pair: &kvv1.KvPair{
					Key:   record[*keyIndex],
					Value: strings.Join(record, ","),
				},
			})
			writeMeasurements.AddTime(time.Since(start))

			if err != nil {
				log.Fatalf("can not set key-value pair %v", err)
			}

			fmt.Printf("Response: %v\n", resp)

		}
	}

	readMeasurements := tachymeter.New(&tachymeter.Config{Size: 1000})

	for _, record := range records[1:] {
		if len(record) > *keyIndex {
			start := time.Now()
			fmt.Printf("Key: %s\n", record[*keyIndex])

			_, err := client.Query(ctx, &kvv1.QueryRequest{
				Key: record[*keyIndex],
			})
			readMeasurements.AddTime(time.Since(start))

			if err != nil {
				log.Fatalf("can not get key-value pair %v", err)
			}
		}
	}

	fmt.Printf("Write results: %s\n", writeMeasurements.Calc().String())
	fmt.Printf("Read results: %s\n", readMeasurements.Calc().String())
}
