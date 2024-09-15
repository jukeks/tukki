package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/chzyer/readline"
	kvv1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/kv/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port     = flag.Int("port", 50051, "The server port")
	hostname = flag.String("hostname", "localhost", "The server hostname")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *hostname, *port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	client := kvv1.NewKvServiceClient(conn)
	repl(client)
}

type Command struct {
	Cmd   string
	Key   string
	Value string
}

func readAndParse(rl *readline.Instance) (*Command, error) {
	input, err := rl.Readline()
	if err != nil {
		return nil, err
	}

	parts := strings.Fields(strings.TrimSpace(input))

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid input, expected at least 2 parts")
	}

	cmd := &Command{
		Cmd: strings.ToLower(parts[0]),
		Key: parts[1],
	}

	if len(parts) > 2 {
		cmd.Value = strings.Join(parts[2:], " ")
	}

	return cmd, nil
}

func repl(client kvv1.KvServiceClient) {
	rl, err := readline.New("> ")
	if err != nil {
		log.Fatalf("can not create readline: %v", err)
	}

	for {
		cmd, err := readAndParse(rl)
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				break
			}
			fmt.Printf("failed to read command: %v\n", err)
			continue
		}

		switch cmd.Cmd {
		case "set":
			start := time.Now()
			_, err := client.Set(context.Background(), &kvv1.SetRequest{
				Pair: &kvv1.KvPair{
					Key:   cmd.Key,
					Value: cmd.Value,
				},
			})
			if err != nil {
				fmt.Printf("failed to set: %v\n", err)
				continue
			}
			fmt.Printf("Key set (%v)\n", time.Since(start))
		case "get":
			start := time.Now()
			resp, err := client.Query(context.Background(), &kvv1.QueryRequest{
				Key: cmd.Key,
			})
			if err != nil {
				fmt.Printf("failed to get: %v (%v)\n", err, time.Since(start))
				continue
			}
			if e := resp.GetError(); e != nil {
				fmt.Printf("key not found: %s (%v)\n", e.Message, time.Since(start))
				continue
			}
			fmt.Printf("value: %s (%v)\n", resp.GetPair().Value, time.Since(start))
		case "delete":
			_, err := client.Delete(context.Background(), &kvv1.DeleteRequest{
				Key: cmd.Key,
			})
			if err != nil {
				fmt.Printf("failed to delete: %v\n", err)
			}
			fmt.Printf("Key deleted\n")
		case "exit":
			return
		default:
			fmt.Printf("unknown command: %s\n", cmd.Cmd)
		}
	}
}
