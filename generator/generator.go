package generator

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"time"
)

const minLength = 10

func GenerateFile(path string, linesCnt int64, maxLength int) error {
	f, err := os.Create(path)
	if err != nil {
		log.WithError(err).Errorf("Can't create by given path %q", path)
		return err
	}

	writer := bufio.NewWriterSize(f, 1*1024*1024)

	rand.Seed(time.Now().UTC().UnixNano())

	for i := int64(0); i < linesCnt; i++ {
		l := maxLength
		if maxLength > minLength {
			l = minLength + rand.Intn(maxLength-minLength)
		}

		line := make([]byte, l+1)
		for j := 0; j < l; j++ {
			line[j] = getRandomSymbol()
		}
		if i != linesCnt-1 {
			line[l] = '\n'
		}

		if _, err := writer.Write(line); err != nil {
			log.WithError(err).Error("Can't write to file")
			return err
		}
	}
	writer.Flush()
	return nil
}

func getRandomSymbol() byte {
	const (
		//asciiMin = 32 // space char
		//asciiMin = 48 // 0 char
		asciiMin = 65 // A char
		asciiMax = 90 // Z char
		//asciiMax = 126 // ~ char
	)
	return byte(asciiMin + rand.Intn(asciiMax-asciiMin))
}
