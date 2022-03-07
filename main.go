package main

import (
	"flag"
	"fmt"
	"github.com/JamesHovious/w32"
	"github.com/go-vgo/robotgo"
	"log"
	"time"
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
		haybm := robotgo.ToBitmap(hay)

		if hay == nil {
			panic(fmt.Sprintf("File %s could not be loaded", *fname))
		}

		defer robotgo.FreeBitmap(hay)

		x, y := robotgo.FindBitmap(hay)
		w:=haybm.Width
		h:=haybm.Height
		t:=2
		log.Printf("FindBitmap------%d x %d : %d x %d ", x,y, w, h)

		top := &w32.RECT{
			Left:   int32(x),
			Top:    int32(y),
			Right:  int32(x + w),
			Bottom: int32(y + t),
		}
		bottom := &w32.RECT{
			Left:   int32(x),
			Top:    int32(y + h - t),
			Right:  int32(x + w),
			Bottom: int32(y + h),
		}
		left := &w32.RECT{
			Left:   int32(x),
			Top:    int32(y),
			Right:  int32(x + t),
			Bottom: int32(y + h),
		}
		right := &w32.RECT{
			Left:   int32(x + w - t),
			Top:    int32(y),
			Right:  int32(x + w),
			Bottom: int32(y + h),
		}
		hdc := w32.GetDC(0)
		lb := &w32.LOGBRUSH{
			LbStyle: w32.BS_SOLID,
			LbColor: 0x0000ff,
			LbHatch: 0,
		}
		brush := w32.CreateBrushIndirect(lb)
		w32.FillRect(hdc, top, brush)
		w32.FillRect(hdc, right, brush)
		w32.FillRect(hdc, bottom, brush)
		w32.FillRect(hdc, left, brush)

	default:
		log.Printf("Op: " + *op + " is not known")
	}

}
