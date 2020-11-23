package swindow

import "github.com/nsf/termbox-go"

/*
 * Height     ウィンドウの高さ
 * Width      ウィンドウの幅
 */
type WindowInfo struct {
	Height    int
	Width     int
	TextColor termbox.Attribute
	BgColor   termbox.Attribute
}
