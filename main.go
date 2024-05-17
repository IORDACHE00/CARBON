package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type InvertedIndex map[string][]int

var (
	index InvertedIndex
	lines []string
	mutex sync.RWMutex
)

func readFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func createIndex(lines []string) InvertedIndex {
	index := make(InvertedIndex)
	for i, line := range lines {
		words := strings.Fields(line)
		for _, word := range words {
			word = strings.ToLower(word)
			index[word] = append(index[word], i)
		}
	}
	return index
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is missing", http.StatusBadRequest)
		return
	}

	results := search(query)
	jsonResponse, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func search(query string) []string {
	query = strings.ToLower(query)

	mutex.RLock()
	defer mutex.RUnlock()

	var result []string
	if lineNumbers, found := index[query]; found {
		for _, lineNumber := range lineNumbers {
			result = append(result, lines[lineNumber])
		}
	}
	return result
}

func GenerateRandomSentence(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ ")
	sentence := make([]rune, length)

	for i := range sentence {
		sentence[i] = letters[rand.Intn(len(letters))]
	}

	return string(sentence)
}

func generate() {
	rand.Seed(time.Now().UnixNano())

	numDocuments := 1000000
	avgDocumentLength := 100
	filePath := "sample.txt"

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	for i := 1; i <= numDocuments; i++ {
		documentLength := rand.Intn(avgDocumentLength*2) + avgDocumentLength/2
		randomSentence := GenerateRandomSentence(documentLength)
		_, err := file.WriteString(fmt.Sprintf("%d: %s\n", i, randomSentence))
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	fmt.Printf("Generated %d documents in %s\n", numDocuments, filePath)
}

func main() {
	// Uncomment to generate sample.txt - 1mil lines of words to mock the data for search
	generate()

	filePath := "./sample.txt"
	var err error
	lines, err = readFile(filePath)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	index = createIndex(lines)
	log.Printf("Created index with %d words\n", len(index))

	http.HandleFunc("/search", searchHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
