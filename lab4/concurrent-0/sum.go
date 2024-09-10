package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// read a file from a filepath and return a slice of bytes
func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return nil, err
	}
	return data, nil
}

// sum all bytes of a file
func sum(filePath string) (int, error) {
	data, err := readFile(filePath)
	if err != nil {
		return 0, err
	}

	_sum := 0
	for _, b := range data {
		_sum += int(b)
	}

	return _sum, nil
}

func processFile(path string, ch chan<- struct {
	sum  int
	file string
}, done chan<- bool) {
	_sum, err := sum(path)
	if err != nil {
		done <- true
		return
	}

	ch <- struct {
		sum  int
		file string
	}{_sum, path}

	done <- true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	var totalSum int64
	sums := make(map[int][]string)

	
	ch := make(chan struct {
		sum  int
		file string
	})

	done := make(chan bool)

	numFiles := len(os.Args) - 1

	for _, path := range os.Args[1:] {
		go processFile(path, ch, done)
	}

	go func() {
		for i := 0; i < numFiles; i++ {
			<-done 
		}
		close(ch) 
	}()

	for result := range ch {
		totalSum += int64(result.sum)
		sums[result.sum] = append(sums[result.sum], result.file)
	}

	fmt.Println("Total Sum:", totalSum)

	for sum, files := range sums {
		if len(files) > 1 {
			fmt.Printf("Sum %d: %v\n", sum, files)
		}
	}
}
