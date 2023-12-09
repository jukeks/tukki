package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jukeks/tukki/internal/db"
)

func main() {
	tmpDir := os.TempDir()
	dbDir := filepath.Join(tmpDir, "tukki")
	err := os.MkdirAll(dbDir, 0755)
	if err != nil {
		log.Fatalf("failed to create db dir: %v", err)
	}

	db := db.NewDatabase(dbDir)
	repl(db)
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

	// Trim space and split the input
	parts := strings.Fields(strings.TrimSpace(input))

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid input, expected at least 2 parts")
	}

	cmd := &Command{
		Cmd: strings.ToLower(parts[0]),
		Key: parts[1],
	}

	// If there's a third part, treat it as the value
	if len(parts) > 2 {
		cmd.Value = strings.Join(parts[2:], " ")
	}

	return cmd, nil
}

func repl(db *db.Database) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		cmd, err := readAndParse(reader)
		if err != nil {
			log.Printf("failed to read command: %v", err)
			continue
		}

		switch cmd.Cmd {
		case "set":
			err = db.Set(cmd.Key, cmd.Value)
			if err != nil {
				log.Printf("failed to set: %v", err)
			}
		case "get":
			value, err := db.Get(cmd.Key)
			if err != nil {
				log.Printf("failed to get: %v", err)
			}
			log.Printf("value: %s", value)
		case "delete":
			err = db.Delete(cmd.Key)
			if err != nil {
				log.Printf("failed to delete: %v", err)
			}
		case "exit":
			return
		default:
			log.Printf("unknown command: %s", cmd.Cmd)
		}
	}
}
