package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type InvertedIndex map[string][]int

var (
	index InvertedIndex
	lines []string
	mutex sync.RWMutex
	// nGramSize is the size of the n-grams used to index the words a.k.a the accuracy of the search
	// The larger the n-gram size, the more accurate and faster the search will be at the cost of more memory
	// This is my favorite model to use with the nGram size of 5, I found it to be the best balance between speed and accuracy that I could find
	nGramSize = 5
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

func createNGrams(word string, n int) []string {
	var nGrams []string
	if len(word) < n {
		nGrams = append(nGrams, word)
	} else {
		for i := 0; i <= len(word)-n; i++ {
			nGrams = append(nGrams, word[i:i+n])
		}
	}
	return nGrams
}

func createIndex(lines []string, n int) InvertedIndex {
	index := make(InvertedIndex)
	for i, line := range lines {
		words := strings.Fields(line)
		for _, word := range words {
			word = strings.ToLower(word)
			nGrams := createNGrams(word, n)
			for _, nGram := range nGrams {
				index[nGram] = append(index[nGram], i)
			}
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

	results := search(query, nGramSize)
	jsonResponse, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func search(query string, n int) []string {
	query = strings.ToLower(query)
	queryNGrams := createNGrams(query, n)

	mutex.RLock()
	defer mutex.RUnlock()

	lineOccurrences := make(map[int]bool)
	var wg sync.WaitGroup
	resultsChan := make(chan int, len(queryNGrams))

	for _, nGram := range queryNGrams {
		wg.Add(1)
		go func(nGram string) {
			defer wg.Done()
			if lineNumbers, found := index[nGram]; found {
				for _, lineNumber := range lineNumbers {
					resultsChan <- lineNumber
				}
			}
		}(nGram)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for lineNumber := range resultsChan {
		lineOccurrences[lineNumber] = true
	}

	fullWordMatches := createNGrams(query, len(query))
	for _, nGram := range fullWordMatches {
		if lineNumbers, found := index[nGram]; found {
			for _, lineNumber := range lineNumbers {
				lineOccurrences[lineNumber] = true
			}
		}
	}

	var result []string
	for lineNumber := range lineOccurrences {
		result = append(result, lines[lineNumber])
	}
	return result
}

func main() {
	filePath := "../sample.txt"
	var err error
	lines, err = readFile(filePath)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	index = createIndex(lines, nGramSize)
	log.Printf("Created index with %d n-grams\n", len(index))

	http.HandleFunc("/search", searchHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
