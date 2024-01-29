package main

import (
	"bufio"
	"errors"
	"math/rand"
	"os"
	"time"
)

type Words struct {
	wordList []string
}

func (w *Words) loadFrom(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		w.wordList = append(w.wordList, scanner.Text())
	}
	return scanner.Err()
}

func (w *Words) get(idx int) (string, error) {
	if idx >= len(w.wordList) {
		return "", errors.New("Out of bounds")
	}
	return w.wordList[idx], nil
}

func (w *Words) getLength() int {
	return len(w.wordList)
}

func (w *Words) getRandom() (string, error) {
	length := w.getLength()
	if length < 1 {
		return "", errors.New("Empty word list")
	}

	randSource := rand.NewSource(time.Now().UnixNano())
	r := rand.New(randSource)

	idx := r.Intn(length)
	return w.wordList[idx], nil
}
