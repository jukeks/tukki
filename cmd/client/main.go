package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	kvv1 "github.com/jukeks/tukki/proto/gen/tukki/rpc/kv/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port     = flag.Int("port", 50051, "The server port")
	hostname = flag.String("hostname", "localhost", "The server hostname")
)

func main() {
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

func readAndParse(reader *bufio.Reader) (*Command, error) {
	input, err := reader.ReadString('\n')
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
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		cmd, err := readAndParse(reader)
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
			fmt.Printf("Key set\n")
		case "get":
			resp, err := client.Query(context.Background(), &kvv1.QueryRequest{
				Key: cmd.Key,
			})
			if err != nil {
				fmt.Printf("failed to get: %v\n", err)
				continue
			}
			if e := resp.GetError(); e != nil {
				fmt.Printf("key not found: %s\n", e.Message)
				continue
			}
			fmt.Printf("value: %s\n", resp.GetPair().Value)
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
