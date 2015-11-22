package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type wordMember struct {
	inWord  string
	outWord string
}

func transcribe(inFile, outFile string) error {

	var outWait sync.WaitGroup

	errCh := make(chan error)

	in := make(chan wordMember)

	outWait.Add(1)
	go readFile(inFile, in, errCh, &outWait)

	out := make(chan wordMember)
	for i := 0; i < 1000; i++ {
		outWait.Add(1)
		go func() {
			defer outWait.Done()
			for w := range in {
				w.process()
				out <- w
			}
		}()
	}

	go func() {
		outWait.Wait()
		close(out)
	}()

	go writeFile(outFile, out, errCh)

	// Get Read and Write errors
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil

}

func readFile(filename string, in chan wordMember, errCh chan error, wg *sync.WaitGroup) {

	defer wg.Done()
	defer close(in)

	from, err := os.Open(filename)
	defer from.Close()

	if err != nil {
		errCh <- fmt.Errorf("Error opening " + filename + ": " + err.Error())
		return
	}
	seen := make(map[string]bool)
	scanner := bufio.NewScanner(from)
	for scanner.Scan() {
		word := scanner.Text()
		if !seen[word] { // ignore duplicate entries
			seen[word] = true
			in <- wordMember{inWord: word}
		}
	}
	if err := scanner.Err(); err != nil {
		errCh <- fmt.Errorf("Error reading from " + filename + ": " + err.Error())
	}
}

func (w *wordMember) process() {
	if validWord(w.inWord) {
		w.outWord = "\"" + w.inWord + "\": true, "
	}
}

func validWord(candidate string) bool {
	return len(candidate) > 0 && candidate[0] != '#' && candidate[0] != ' '
}

func writeFile(filename string, out chan wordMember, errCh chan error) {
	// Write
	to, err := os.Create(filename)
	defer to.Close()
	if err != nil {
		errCh <- fmt.Errorf("Error creating " + filename + ": " + err.Error())
		return
	}

	// Write Header
	header := "package main\n\n"
	header += "type alphabet map[string]bool\n\n"
	header += "var englishWords = alphabet{"
	n, err := to.WriteString(header)
	if err != nil {
		errCh <- fmt.Errorf("Error writing header to " + filename + ". Wrote " + fmt.Sprint(n) + " bytes before error: " + err.Error())
		return
	}

	// Write Body
	for word := range out {
		n, err = to.WriteString(word.outWord)
		if err != nil {
			errCh <- fmt.Errorf("Error writing " + word.outWord + " to " + filename + ". Wrote " + fmt.Sprint(n) + " bytes before error: " + err.Error())
			return
		}
	}

	// Write Footer
	footer := "}"
	n, err = to.WriteString(footer)
	if err != nil {
		errCh <- fmt.Errorf("Error writing footer to " + filename + ". Wrote " + fmt.Sprint(n) + " bytes before error: " + err.Error())
		return
	}

	close(errCh)
}

func main() {
	err := transcribe("wiki-100k.txt", "words.go")
	if err != nil {
		fmt.Println("Error transcribing", err)
	}
}
