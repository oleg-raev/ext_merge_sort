package ext_merge_sort

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"joom-sort/common"

	"github.com/sirupsen/logrus"
)

const (
	chunkSize  = 1000
	bufferSize = 64 * 1024 //128KB
)

type ExtMergeSort struct {
	inputPath string
	tempDir   string
	chunksCnt int
}

func New(path string) ExtMergeSort {
	return ExtMergeSort{
		inputPath: path,
	}
}

func (e *ExtMergeSort) Sort() error {
	if err := e.prepareTempDirectory(); err != nil {
		return err
	}

	if err := e.prepareChunks(); err != nil {
		return err
	}

	if err := e.mergeSort(); err != nil {
		return err
	}
	return nil
}

func (e *ExtMergeSort) mergeSort() error {
	for curSize := 1; curSize < e.chunksCnt; curSize *= 2 {
		for left := 0; left < e.chunksCnt; left += curSize * 2 {
			right := common.MinInt(left+curSize, e.chunksCnt-1)
			if right-left < curSize {
				continue
			}
			if err := e.merge(left, right); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *ExtMergeSort) getChunkScanner(num int) (*bufio.Scanner, error) {
	path := e.getChunkPath(num)
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		logrus.WithField("inputPath", path).WithError(err).Error("Can't open chunk")
		return nil, err
	}
	lreader := bufio.NewScanner(f)
	lreader.Split(bufio.ScanLines)

	return lreader, nil
}

func (e *ExtMergeSort) getOutputWriter(left, right int) (*bufio.Writer, error) {
	f, err := os.Create(e.getPairPath(left, right))
	if err != nil {
		logrus.WithError(err).Error("Error on creating output file for merging chunks %q and %q", left, right)
		return nil, err
	}
	return bufio.NewWriterSize(f, bufferSize), nil
}

func (e *ExtMergeSort) merge(left, right int) error {
	leftScan, err := e.getChunkScanner(left)
	if err != nil {
		return err
	}
	rightScan, err := e.getChunkScanner(right)
	if err != nil {
		return err
	}

	out, err := e.getOutputWriter(left, right)
	if err != nil {
		return err
	}

	leftScan.Scan()
	lineL := leftScan.Bytes()
	rightScan.Scan()
	lineR := rightScan.Bytes()
	for {
		//todo simplify it
		if bytes.Compare(lineL, lineR) == -1 {
			if err := e.addLine(out, lineL); err != nil {
				return err
			}
			if !leftScan.Scan() {
				if err := e.addLine(out, lineR); err != nil {
					logrus.WithError(err).Error("Can't write to lineR")
					return err
				}
				if err := e.copyRestRows(rightScan, out); err != nil {
					return err
				}
				break
			}
			lineL = append([]byte{'\n'}, leftScan.Bytes()...)
		} else {
			if err := e.addLine(out, lineR); err != nil {
				return err
			}
			if !rightScan.Scan() {
				if err := e.addLine(out, lineL); err != nil {
					logrus.WithError(err).Error("Can't write to lineR")
					return err
				}
				if err := e.copyRestRows(leftScan, out); err != nil {
					return err
				}
				break
			}
			lineR = append([]byte{'\n'}, rightScan.Bytes()...)
		}
	}
	if err := out.Flush(); err != nil {
		logrus.WithError(err).Error("Can't Flush writter")
		return err
	}

	if err := e.sanitizeMergedFiles(left, right); err != nil {
		return err
	}

	return nil
}

func (e *ExtMergeSort) getChunkPath(id int) string {
	return filepath.Join(e.tempDir, strconv.Itoa(id)+".txt")
}

func (e *ExtMergeSort) getPairPath(left, right int) string {
	return filepath.Join(e.tempDir, fmt.Sprintf("%d_%d.txt", left, right))
}

func (e *ExtMergeSort) removeChunk(id int) error {
	path := e.getChunkPath(id)
	if err := os.Remove(path); err != nil {
		logrus.WithError(err).Error("Can't remove chunk file")
		return err
	}
	return nil
}

func (e *ExtMergeSort) sanitizeMergedFiles(left, right int) error {
	if err := e.removeChunk(left); err != nil {
		return err
	}
	if err := e.removeChunk(right); err != nil {
		return err
	}
	path := e.getPairPath(left, right)
	if err := os.Rename(path, e.getChunkPath(left)); err != nil {
		logrus.WithError(err).WithField("inputPath", path).Error("Can't rename pair chunk")
	}
	return nil
}

func (e *ExtMergeSort) copyRestRows(from *bufio.Scanner, to *bufio.Writer) error {
	for from.Scan() {
		if _, err := to.Write(append(from.Bytes(), '\n')); err != nil {
			logrus.WithError(err).Error("Can't write bytes to out")
			return err
		}
	}
	return nil
}

func (e *ExtMergeSort) prepareTempDirectory() error {
	var err error
	e.tempDir, err = ioutil.TempDir("", "ext_merge_sort")
	if err != nil {
		logrus.WithError(err).Error("Can't create temporary directory")
		return err
	}
	logrus.Infof("Created temp directory: %s", e.tempDir)
	return nil
}

func (e *ExtMergeSort) prepareChunks() error {
	file, err := os.OpenFile(e.inputPath, os.O_RDONLY, 0)
	if err != nil {
		logrus.WithError(err).Error("Can't open file for sorting")
		return err
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, bufferSize)
	var (
		chunkNum  = -1
		lineNum   = 0
		chunkData = make([][]byte, 0)
	)

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.WithError(err).Error("Error given on reading from input")
			return err
		}

		// faced collision error when used line given by ReadLine because they may modify slice - it's linked object
		line = common.UnlinkSlice(line)

		chunkData = append(chunkData, line)
		lineNum++
		if lineNum >= chunkSize {
			chunkNum++
			if err := e.sortAndSaveChunk(chunkData, chunkNum); err != nil {
				return err
			}
			lineNum = 0
			chunkData = make([][]byte, 0)
		}
	}
	if len(chunkData) > 0 {
		chunkNum++
		if err := e.sortAndSaveChunk(chunkData, chunkNum); err != nil {
			return err
		}
	}
	e.chunksCnt = chunkNum + 1

	return nil
}

func (e *ExtMergeSort) addLine(writer *bufio.Writer, line []byte) error {
	if _, err := writer.Write(line); err != nil {
		logrus.WithError(err).Error("Can't write bytes to out")
		return err
	}
	return nil
}

func (e *ExtMergeSort) quickSort(data [][]byte) {
	l := len(data)
	for i := 0; i < l-1; i++ {
		minIndex := i
		for j := i + 1; j < l; j++ {
			if bytes.Compare(data[j], data[minIndex]) == -1 {
				data[j], data[minIndex] = data[minIndex], data[j]
			}
		}
	}
}

func (e *ExtMergeSort) selectSort(data [][]byte) {
	l := len(data)
	for i := 0; i < l-1; i++ {
		minIndex := i
		for j := i + 1; j < l; j++ {
			if bytes.Compare(data[j], data[minIndex]) == -1 {
				data[j], data[minIndex] = data[minIndex], data[j]
			}
		}
	}
}

func (e *ExtMergeSort) sortAndSaveChunk(data [][]byte, chunkNum int) error {
	e.selectSort(data)

	f, err := os.Create(e.getChunkPath(chunkNum))
	if err != nil {
		logrus.WithError(err).Error("Can't create chunk file")
		return err
	}
	defer f.Close()

	writer := bufio.NewWriterSize(f, bufferSize)

	if _, err := writer.Write(bytes.Join(data, []byte{'\n'})); err != nil {
		logrus.WithError(err).Error("Can't write to chunk file")
		return err
	}

	if err := writer.Flush(); err != nil {
		logrus.WithError(err).Error("Error on flushing writer buffer")
		return err
	}

	return nil
}
