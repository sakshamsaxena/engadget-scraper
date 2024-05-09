package cache

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
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
	ww, err := os.OpenFile("words.txt", os.O_RDONLY, 0444)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if redisClient.Ping(context.Background()).Err() != nil {
		panic("redis err") // TODO: handle
	}
	dict := make(map[string]bool)
	buffer := strings.Builder{}
	for {
		singleByte := make([]byte, 1)
		bytesRead, readErr := ww.Read(singleByte)
		if readErr != nil && readErr != io.EOF {
			panic(readErr)
		}
		if bytesRead == 0 || readErr == io.EOF {
			break
		}
		if singleByte[0] == '\n' {
			line := buffer.String()
			buffer.Reset()
			dict[line] = true
		} else {
			buffer.Write(singleByte)
		}
	}
	if buffer.Len() > 0 {
		line := buffer.String()
		dict[line] = true
		buffer.Reset()
	}
	store = &cache{
		wordBank:    dict,
		redisClient: redisClient,
	}
}

func GetTopNWords() []string {
	val, err := store.redisClient.ZRevRange(context.Background(), topWordsSortedSet, 0, 9).Result()
	if err != nil {
		panic(err)
	}
	return val
}

func SetWords(words []string) error {
	pp := store.redisClient.Pipeline()
	for _, token := range words {
		pp.ZIncrBy(context.Background(), topWordsSortedSet, 1, token)
	}
	_, err := pp.Exec(context.Background())
	return err
}

func CheckBank(word string) bool {
	return store.wordBank[word]
}
