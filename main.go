package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

func find(fname string) {
	hay := robotgo.OpenBitmap(fname)
	haybm := robotgo.ToBitmap(hay)

	if hay == nil {
		panic(fmt.Sprintf("File %s could not be loaded", fname))
	}

	defer robotgo.FreeBitmap(hay)

	x, y := robotgo.FindBitmap(hay)
	w := haybm.Width
	h := haybm.Height

	// log.Printf("FindBitmap------%d x %d : %d x %d ", x, y, w, h)
	alert("Found %s: %d x %d : %d x %d ", fname, x, y, w, h)

}
func sc(fname string, to int) string {
	time.Sleep(time.Second * time.Duration(to))
	rect := robotgo.GetScreenRect(0)
	bitmap := robotgo.CaptureScreen(rect.X, rect.Y, rect.W, rect.H)
	defer robotgo.FreeBitmap(bitmap)
	ret := robotgo.SaveBitmap(bitmap, fname)
	alert("Screenshot saved: %s", fname)
	return ret

}

func findFiles() {
	ens, err := os.ReadDir(".")
	if err != nil {
		alert("Error reading dir: %s", err.Error())
		return
	}
	for _, e := range ens {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".png") {
			find(e.Name())
		}
	}

}

func finderDaemon() chan bool {
	alert("Starting Finder Daemon")
	ret := make(chan bool)
	go func() {
		for {
			select {
			case <-ret:
				alert("Stopping Finder Daemon")
				return
			case <-time.After(time.Second):
				findFiles()
			}
		}
	}()
	return ret
}

func alert(s string, p ...interface{}) {
	beeep.Alert("SCUtil", fmt.Sprintf(s, p...), "")
}

func daemon() {

	var chFinderDaemon chan bool
	var ch chan hook.Event

	robotgo.EventHook(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("ctrl-shift-q")
		robotgo.StopEvent()
		close(ch)
	})

	robotgo.EventHook(hook.KeyDown, []string{"w", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("ctrl-shift-w")
		fname := fmt.Sprintf("%d.png", time.Now().Unix())
		sc(fname, 0)
		alert("File saved: %s", fname)
	})

	robotgo.EventHook(hook.KeyDown, []string{"f", "ctrl", "shift"}, func(e hook.Event) {
		if chFinderDaemon == nil {
			chFinderDaemon = finderDaemon()

		} else {
			chFinderDaemon <- true
			close(chFinderDaemon)
			chFinderDaemon = nil
		}
	})

	robotgo.EventHook(hook.KeyDown, []string{"p", "ctrl", "shift", "alt"}, func(e hook.Event) {
		fmt.Println("ctrl-shift-p")
		//pfname := fmt.Sprintf("%d.png", time.Now().Unix())
		// sc(fname)

	})

	ch = robotgo.EventStart()
	<-robotgo.EventProcess(ch)
	alert("Exiting")
	os.Exit(0)

}
func main() {
	op := flag.String("op", "daemon", "[sc|find|daemon]")
	// coords := flag.String("coords", "0,0,0,0", "if only a piece of img should be saved, these coords will define it")
	fname := flag.String("fname", fmt.Sprintf("%d.png", time.Now().Unix()), "file name to save sc")
	to := flag.Int("to", 0, "Timeout for screenshot")
	flag.Parse()
	switch *op {
	case "sc":
		sc(*fname, *to)

	case "find":
		find(*fname)

	case "daemon":
		daemon()
	default:
		log.Printf("Op: " + *op + " is not known")
	}

}
