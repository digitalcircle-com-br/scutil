package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"runtime"
	"time"

	"github.com/kbinani/screenshot"
)

//HAYSTACK
func compColor(a, b color.Color) bool {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	eq := ar == br && ag == bg && ab == bb && aa == ba
	//log.Printf("%v - R: %d:%d, G: %d:%d, B: %d:%d, A: %d,%d", eq, ar, br, ag, bg, ab, bb, aa, ba)
	return eq

}

type FindOpts struct {
	hay        image.Image
	nd         image.Image
	sx         int
	sy         int
	fx         int
	fy         int
	outChan    chan *image.Point
	cancelChan chan struct{}
}

func FindByCoord(o *FindOpts) error {
	proceed := true
	go func() {
		<-o.cancelChan
		close(o.cancelChan)
		proceed = false
	}()
	var p *image.Point

	for x := o.sx; x < o.fx && proceed; x++ {
		for y := o.sy; y < o.fy && proceed; y++ {

			for x1 := 0; x1 < o.nd.Bounds().Max.X && proceed; x1++ {
				for y1 := 0; y1 < o.nd.Bounds().Max.Y && proceed; y1++ {
					if !compColor(o.hay.At(x+x1, y+y1), o.nd.At(x1, y1)) {

						goto newloop
					}
				}

			}
			p = &image.Point{
				X: x,
				Y: y,
			}
			o.outChan <- p
			proceed = false
		newloop:
		}
	}
	if proceed {
		o.outChan <- nil
	}
	return nil
}

func FindMultiThread(hay, nd image.Image) (*image.Point, error) {
	cpus := runtime.NumCPU() - 2

	for math.Floor(float64(hay.Bounds().Max.Y)/float64(cpus)) < float64(nd.Bounds().Max.Y) {
		cpus--
	}

	if cpus < 1 {
		cpus = 1
	}

	jump := math.Floor(float64(hay.Bounds().Max.Y) / float64(cpus))
	outChan := make(chan *image.Point)
	cancelChans := make([]chan struct{}, 0)
	for i := 0; i < cpus; i++ {
		cancelChan := make(chan struct{})
		o := &FindOpts{
			hay:        hay,
			nd:         nd,
			sx:         0,
			fx:         hay.Bounds().Max.X,
			sy:         (i * int(jump)),
			fy:         ((i + 1) * int(jump)) + 1,
			outChan:    outChan,
			cancelChan: cancelChan,
		}
		cancelChans = append(cancelChans, cancelChan)
		go FindByCoord(o)
	}
	var ret *image.Point
	for i := 0; i < len(cancelChans); i++ {
		ret = <-outChan

		if ret != nil {
			for _, ch := range cancelChans {
				ch <- struct{}{}
			}
			break
		}
		if ret == nil && i == len(cancelChans)-1 {
			for _, ch := range cancelChans {
				ch <- struct{}{}
			}
		}
	}
	return ret, nil
}

func Find(hay, nd image.Image) (*image.Point, error) {

	for x := 0; x < hay.Bounds().Max.X-nd.Bounds().Max.X; x++ {
		for y := 0; y < hay.Bounds().Max.X-nd.Bounds().Max.Y; y++ {
			//log.Printf("Checking at: %d x %d", x, y)
			for x1 := 0; x1 < nd.Bounds().Max.X; x1++ {
				for y1 := 0; y1 < nd.Bounds().Max.Y; y1++ {
					if !compColor(hay.At(x+x1, y+y1), nd.At(x1, y1)) {
						//log.Printf("Not Found: %d x %d => %d x %d", x+x1, y+y1, x1, y1)
						goto newloop
					}
					//log.Printf("Found: %d x %d => %d x %d", x+x1, y+y1, x1, y1)
				}

			}
			log.Printf("Found at: %d x %d", x, y)
			return &image.Point{
				X: x,
				Y: y,
			}, nil
		newloop:
		}
	}

	return nil, nil
}

//END HAYSTACK

func main() {
	op := flag.String("op", "sc", "if tool should take screenshot")
	// coords := flag.String("coords", "0,0,0,0", "if only a piece of img should be saved, these coords will define it")
	fname := flag.String("fname", fmt.Sprintf("%d.png", time.Now().Unix()), "file name to save sc")
	flag.Parse()
	switch *op {
	case "sc":
		img, err := screenshot.CaptureDisplay(0)
		if err != nil {
			panic(err.Error())
		}
		f, err := os.OpenFile(*fname, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			panic(err.Error())
		}
		defer f.Close()
		err = png.Encode(f, img)
		if err != nil {
			panic(err.Error())
		}
	case "find":
		bs, err := os.ReadFile(*fname)
		if err != nil {
			panic(err.Error())
		}
		hay, err := png.Decode(bytes.NewReader(bs))
		if err != nil {
			panic(err.Error())
		}
		screen, err := screenshot.CaptureDisplay(0)
		if err != nil {
			panic(err.Error())
		}
		pt, err := Find(hay, screen)
		if err != nil {
			panic(err.Error())
		}
		log.Printf("Found image at: %#v", pt)
	default:
		log.Printf("Op: " + *op + " is not known")
	}

}
