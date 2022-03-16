package main

import (
	"github.com/JamesHovious/w32"
)

func drawOnScreen(x, y, w, h int) {
	t := 2

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
}
