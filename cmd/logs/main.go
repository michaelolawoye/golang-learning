package main

import (
	"fmt"
	"bufio"
	"os"
	"flag"
	"unicode"
	"io"
	"errors"
)

func main() {

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
} 


func run() error {

	var stringFlag *string = flag.String("target", ".", "Character being counted in file")
	var inputFile *string = flag.String("in", "", "File being read from")
	var outputFile *string = flag.String("out", "", "File being written into")
	var caseFlag *bool = flag.Bool("casesen", false, "If true, target is case-sensitive")

	flag.Parse()

	if len(*stringFlag) == 0 {
		return errors.New("Length of target flag should not be zero")
	}

	runeArr := []rune(*stringFlag)
	var targetFlag rune = runeArr[0]

	rfile, err := os.Open(*inputFile)
	defer rfile.Close()
	if err != nil {
		return fmt.Errorf("Couldn't open file for reading: %w", err)
	}

	reader := bufio.NewReader(rfile)

	wfile, err := os.OpenFile(*outputFile, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	defer wfile.Close()
	if err != nil {
		return fmt.Errorf("Couldn't open file for writing: %w", err)
	}

	count, err := countChars(reader, targetFlag, *caseFlag)
	if err != nil {
		return fmt.Errorf("Failed to count chars: %w", err)
	}

	_, err = fmt.Fprintf(wfile, "There are %d occurences of %c in %s\n", count, targetFlag, *inputFile)
	if err != nil {
		return fmt.Errorf("Failed to write message to %v. Error: %w", *outputFile, err)
	}

	fmt.Printf("Char: %c. Count: %v\n", targetFlag, count)

	return nil
}

func countChars(reader io.Reader, letter rune, caseSensitive bool) (int, error){

	bufReader := bufio.NewReader(reader)
	var count int
	upperLetter := unicode.ToUpper(letter)

	for {
		char, _, err := bufReader.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, fmt.Errorf("Failed to read next rune in file: %w",err)
		}

		if char == letter {
			count++
		} else if !caseSensitive && unicode.ToUpper(char) == upperLetter {
			count++
		}
	}

	return count, nil
} 
