package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/logrusorgru/aurora"
)

var (
	typingSpeed int
	mu          sync.Mutex
	colorCode   int
	fileName    string
	pause       = make(chan os.Signal, 1)
	play        = make(chan os.Signal)
)

func init() {
	flag.IntVar(&typingSpeed, "s", 100, "The output speed in millisecond.")
	flag.IntVar(&colorCode, "c", 1, `Colors:
	1 - Grern
	2 - Red
	3 - Magenta
	4 - Blue
	0 - No color
 `)
	flag.Parse()
	// getting the filename
	args := flag.Args()
	if len(args) != 0 {
		fileName = args[0]
	}
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var (
		hole blackHole
	)

	stdin := os.Stdin
	fi, err := stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		f, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		go hole.suckToBuf(f)
	} else {
		go hole.suckToBuf(os.Stdin)
	}
	signal.Notify(pause, os.Interrupt)
	for {
		data := hole.flush()
		if len(data) != 0 {
			typingEffect(data)
		}
		if hole.done {
			os.Exit(0)
		}
	}
}

type blackHole struct {
	buff []byte
	done bool
}

// Take the whole things to buffer
func (hole *blackHole) suckToBuf(r io.Reader) {
	for {
		p := make([]byte, 4)
		_, err := r.Read(p)
		mu.Lock()
		hole.buff = append(hole.buff, p...)
		mu.Unlock()
		if err != io.EOF {
			hole.done = true
			return
		}
	}
}

// flush flushes the buffer
func (hole *blackHole) flush() []byte {
	mu.Lock()
	defer func() {
		hole.buff = nil
		mu.Unlock()
	}()
	return hole.buff
}

// colorFunc returns the best color function according to the provided flags.
func colorFunc() func(arg interface{}) aurora.Value {
	switch colorCode {
	case 1:
		return aurora.Green
	case 2:
		return aurora.Red
	case 3:
		return aurora.Magenta
	case 4:
		return aurora.Blue
	default:
		// no color
		return func(arg interface{}) aurora.Value {
			fmt.Print(arg)
			return aurora.Green("")
		}

	}
}

// typingEffect prints the file with a typing effect
func typingEffect(str []byte) {
	var colorizer = colorFunc()
	for _, v := range str {
		select {
		case <-pause:
			pauseFunc()
			fmt.Print(colorizer(string(v)))
		default:
			fmt.Print(colorizer(string(v)))
			// A random sleep for a realistic effect
			time.Sleep(
				time.Duration(rand.Intn(
					typingSpeed,
				)) * time.Millisecond,
			)
		}
	}
}

func pauseFunc() {
	<-pause
}
