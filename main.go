package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/go-vgo/robotgo"
)

func main() {
	op := flag.String("op", "sc", "if tool should take screenshot (sc) of find in screen (find)")
	// coords := flag.String("coords", "0,0,0,0", "if only a piece of img should be saved, these coords will define it")
	fname := flag.String("fname", fmt.Sprintf("%d.png", time.Now().Unix()), "file name to save sc")
	flag.Parse()
	switch *op {
	case "sc":
		rect := robotgo.GetScreenRect(0)
		bitmap := robotgo.CaptureScreen(rect.X, rect.Y, rect.W, rect.H)
		defer robotgo.FreeBitmap(bitmap)
		robotgo.SaveBitmap(bitmap, *fname)

	case "find":
		hay := robotgo.OpenBitmap(*fname)

		if hay == nil {
			panic(fmt.Sprintf("File %s could not be loaded", *fname))
		}

		defer robotgo.FreeBitmap(hay)

		fx, fy := robotgo.FindBitmap(hay)
		fmt.Println("FindBitmap------", fx, fy)

	default:
		log.Printf("Op: " + *op + " is not known")
	}

}
