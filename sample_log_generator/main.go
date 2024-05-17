package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	hosts = []string{
		"example.com",
		"google.com",
		"stackoverflow.com",
		"github.com",
		"reddit.com",
		"amazon.com",
		"twitter.com",
		"linkedin.com",
		"facebook.com",
		"youtube.com",
	}

	paths = []string{
		"/path1",
		"/path2",
		"/path3",
		"/path4",
		"/path5",
		"/path6",
		"/path7",
		"/path8",
		"/path9",
		"/path10",
	}

	randomMessages = []string{
		"This is a random message.",
		"Here's another random message.",
		"Random message number three!",
		"Just adding some random text.",
		"Lorem ipsum dolor sit amet.",
		"Consectetur adipiscing elit.",
		"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	}
)

func main() {

	rand.Seed(time.Now().UnixNano())

	file, err := os.Create("../sample.txt")
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	logger := log.New(file, "", 0)

	numRequests := 1000000
	for i := 0; i < numRequests; i++ {
		logLine := generateLogLine()
		logger.Println(logLine)
	}

	fmt.Println("Mock request logs generated successfully!")
}

func generateLogLine() string {

	host := getRandomItem(hosts)
	path := getRandomItem(paths)
	message := getRandomItem(randomMessages)

	return fmt.Sprintf("%s %s %d %s %s %s", time.Now().Format("Jan 02 15:04:05.00"), "GET", 200, host, path, message)
}

func getRandomItem(items []string) string {
	return items[rand.Intn(len(items))]
}
