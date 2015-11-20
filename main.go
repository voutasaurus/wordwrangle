package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func main() {
	orchestrateStreaming()
}

func orchestrateStreaming() {
	from := "wiki-100k.txt"
	to := "words.go"
	err := streamWords(from, to)
	if err != nil {
		fmt.Println("Stream failed:", err)
	}
}

func streamWords(fromFilename, toFilename string) error {

	// Read
	from, err := os.Open(fromFilename)
	defer from.Close()
	if err != nil {
		return errors.New("Error opening " + fromFilename + ": " + err.Error())
	}
	words := make(chan string)
	errCh := make(chan error)
	go func() {
		seen := make(map[string]bool)
		scanner := bufio.NewScanner(from)
		for scanner.Scan() {
			word := scanner.Text()
			if validWord(word) {
				if !seen[word] { // ignore duplicate entries
					seen[word] = true
					words <- word
				}
			}
		}
		close(words)
		if err := scanner.Err(); err != nil {
			errCh <- errors.New("Error reading from " + fromFilename + ": " + err.Error())
		}
		close(errCh)
	}()

	// Write
	to, err := os.Create(toFilename)
	defer to.Close()
	if err != nil {
		return errors.New("Error creating " + toFilename + ": " + err.Error())
	}

	// Write Header
	header := "package main\n\n"
	header += "type alphabet map[string]bool\n\n"
	header += "var englishWords = alphabet{"
	n, err := to.WriteString(header)
	if err != nil {
		return errors.New("Error writing header to " + toFilename + ". Wrote " + fmt.Sprint(n) + " bytes before error: " + err.Error())
	}

	// Write Body
	for {
		var open bool
		var word string

		select {
		case word, open = <-words: // includes word in map literal
			n, err = to.WriteString("\"" + word + "\": true, ")
			if err != nil {
				return errors.New("Error writing " + word + " to " + toFilename + ". Wrote " + fmt.Sprint(n) + " bytes before error: " + err.Error())
			}
			if !open {
				words = nil // mark for no additional recieves
			}
		case err, open = <-errCh: // interrupt
			if err != nil {
				break
			}
			if !open {
				errCh = nil // mark for no additional recieves
			}
		}

		if words == nil && errCh == nil {
			// both channels are nil (therefore they were closed)
			// must exit now before select blocks forever
			break
		}

	}

	if err != nil {
		return errors.New("Error reading from " + fromFilename + ": " + err.Error())
	}

	// Write Footer
	footer := "}"
	n, err = to.WriteString(footer)
	if err != nil {
		return errors.New("Error writing footer to " + toFilename + ". Wrote " + fmt.Sprint(n) + " bytes before error: " + err.Error())
	}

	return nil
}

func validWord(candidate string) bool {
	return len(candidate) > 0 && candidate[0] != '#' && candidate[0] != ' '
}
