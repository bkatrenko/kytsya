package main

// Example of PIPELINE multithreading pattern with kytsya.
// Here you could found a simple application that binarizzzze Hamlet
// by William Shakespeare: from text format to a set of 0's and 1's
import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bkatrenko/kytsya"
)

// go mod init {your module}
func main() {
	resCh := make(chan []byte) //<-- creating pipe channel
	now := time.Now()

	kytsya.NewBox().
		WithWaitGroup(). // <-- add a WaitGroup
		WithRecover().   // <-- add a recovery handler
		AddTask(func() { // <-- task for reading file byte by byte
			defer close(resCh)
			readBytes(resCh)
		}).
		AddTask(func() { // <-- task for write bytes in binary format into a file
			writeBytes(resCh)
		}).
		Run().
		Wait()

	fmt.Printf("done in: %v\n“To be, or not to be, that is the question.”\n", time.Since(now))
	// file content example: (@_@)
	// 1001000110011111011111111001011111110100110100111101001101100110010111111
	// 0101010101001110111110111110000010011011100101111001111100111100101110111
	// 0110011111001011110010111001110101010011110000111010011101100110111111100
	// 1011100111010100011111100101100001111011011001011100100110100111001111100
	// 1111100101111001010101000111111001011000011110110110010111001001101001110
	// 0111110011111001011110010100111111001110000011000111101111110110111100001
	// 1000011101110110100111011111101110101010001001101111110001111101001101111
	// 1110010100000110111111001101000001000100110100111101101101001110111011010
	// 0111101001111001101010000011110100111010011001011101110110010011000011101
	// 1101110100111001110110010000010011001101111111001011001001110011101100100
	// 0001000111111010111000011110010110010011100111011001000001001101111010111
	// 1001111010011100011110100111000011101110111001110110010000010011001100001
	// 1100101111001011101001100101111001110011111100111000001000110110111111011
	// 0011011001101111111011111001011110010111001110110010000010100111101111110
	// 1100110010011010011100101111001011100111011001000001001111110011011001101

	// “What a piece of work is a man! How noble in
	// reason, how infinite in faculty! In form and moving
	// how express and admirable! In action how like an Angel!
	// in apprehension how like a god! The beauty of the
	// world! The paragon of animals! And yet to me, what is
	// this quintessence of dust? Man delights not me; no,
	// nor Woman neither; though by your smiling you seem
	// to say so..”

	// (Hamlet, act 2 scene 2)
}

// readBytes will read file byte by byte
func readBytes(resCh chan []byte) {
	f, err := os.Open("hamlet.txt")
	if err != nil {
		panic(err)
	}

	b := make([]byte, 1024)

	for {
		_, err := f.Read(b)
		switch err {
		case io.EOF:
			return
		case nil:
		default:
			panic(err)
		}

		resCh <- b
	}
}

func writeBytes(resCh chan []byte) {
	f, err := os.Create("binary_hamlet.bin")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	kytsya.ForChan(resCh, func(val []byte) { // <-- read channel until its closed
		kytsya.ForEach(val, func(i int, val byte) { // <-- write every byte in binary
			f.Write([]byte(fmt.Sprintf("%b", val)))
		})
	})
}
