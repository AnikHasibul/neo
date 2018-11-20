package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/logrusorgru/aurora"
)

var (
	typingSpeed int
	colorCode   int
	fileName    string
	pause       = make(chan os.Signal, 1)
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
		err  error
		data []byte
	)

	stdin := os.Stdin
	fi, err := stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		data, err = ioutil.ReadFile(fileName)
	} else {
		data, err = ioutil.ReadAll(os.Stdin)
	}
	if err != nil {
		log.Fatal(err)
	}
	signal.Notify(pause, os.Interrupt)
	typingEffect(data)
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

func pauseFunc() {
	defer signal.Notify(pause, os.Interrupt)
	signal.Stop(pause)
	fmt.Print(aurora.Red(":"))
	char := ""
	fmt.Scanf("%s", &char)
	if char == "q" {
		os.Exit(0)
	}
}
