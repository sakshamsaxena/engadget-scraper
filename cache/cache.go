package cache

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/sakshamsaxena/engadget-scraper/config"
)

const (
	topWordsSortedSet = "topWordsSortedSet"
)

var store *cache

type cache struct {
	redisClient *redis.Client
	wordBank    map[string]bool
}

func Initialize() {
	// Initialize the word bank file
	workBankFile, openErr := os.OpenFile(config.Get("static.wordBank").(string), os.O_RDONLY, 0444)
	if openErr != nil {
		panic(openErr)
	}

	// Initialize the redis cache
	redisHost := config.Get("redis.host").(string)
	redisPort := int(config.Get("redis.port").(float64))
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", redisHost, redisPort),
	})
	if redisErr := redisClient.Ping(context.Background()).Err(); redisErr != nil {
		panic(redisErr)
	}
	// Note: Assuming restart/fresh start here, the redis cache is also reset
	redisClient.Del(context.Background(), topWordsSortedSet)

	// Initialize the cache store
	store = &cache{
		wordBank:    make(map[string]bool),
		redisClient: redisClient,
	}

	// Load the wordBank to store
	buffer := strings.Builder{}
	for {
		singleByte := make([]byte, 1)
		bytesRead, readErr := workBankFile.Read(singleByte)
		if readErr != nil && readErr != io.EOF {
			panic(readErr)
		}
		if bytesRead == 0 || readErr == io.EOF {
			break
		}
		if singleByte[0] == '\n' {
			line := buffer.String()
			buffer.Reset()
			store.wordBank[line] = true
		} else {
			buffer.Write(singleByte)
		}
	}
	if buffer.Len() > 0 {
		line := buffer.String()
		store.wordBank[line] = true
		buffer.Reset()
	}
}

func GetTopNWords(n int) []map[string]any {
	response, err := store.redisClient.ZRevRangeWithScores(context.Background(), topWordsSortedSet,
		0,
		int64(n-1)).Result()
	if err != nil {
		panic(err)
	}
	results := make([]map[string]any, n)
	for index, set := range response {
		results[index] = map[string]any{
			"word": set.Member.(string),
			"freq": int(set.Score),
		}
	}
	return results
}

func SetWords(words []string) error {
	commandPipe := store.redisClient.Pipeline()
	for _, word := range words {
		commandPipe.ZIncrBy(context.Background(), topWordsSortedSet, 1, word)
	}
	_, err := commandPipe.Exec(context.Background())
	return err
}

func CheckWordBank(word string) bool {
	return store.wordBank[word]
}
