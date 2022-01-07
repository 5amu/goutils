package filequeue

import (
	"bufio"
	"os"
)

const BufferSize = 10

type FileQueue struct {
	next    chan string
	empty   chan struct{}
	counter int
}

func NewFileQueue() *FileQueue {
	next := make(chan string, 1)
	empty := make(chan struct{}, 1)
	return &FileQueue{next: next, empty: empty}
}

func (fq *FileQueue) ScanFile(fname string) error {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make(chan string, BufferSize)
	go func() {
		for {
			b := <-buffer
			fq.next <- b
		}
	}()

	fq.counter = 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		buffer <- scanner.Text()
		fq.counter++
	}

	for {
		if len(buffer) == 0 && fq.counter == 0 {
			fq.empty <- struct{}{}
			return nil
		}
	}
}

func (fq *FileQueue) IsEmpty() <-chan struct{} {
	return fq.empty
}

func (fq *FileQueue) Pop() <-chan string {
	fq.counter--
	return fq.next
}
