// Copyright (C) 2020  Raziman

package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rivo/tview"
	"github.com/ztrue/tracerr"
)

type PlayingBar struct {
	*tview.Frame
	full      int
	progress  chan int
	_progress int
	skip      bool
	text      *tview.TextView
}

func (p *PlayingBar) help() []string {
	return []string{}
}

// Playing bar shows progress of the song and the title of the song
func newPlayingBar() *PlayingBar {

	textView := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	frame := tview.NewFrame(textView).SetBorders(1, 1, 1, 1, 1, 1)
	frame.SetBorder(true).SetTitle(" Now Playing ")

	p := &PlayingBar{
		Frame:    frame,
		text:     textView,
		progress: make(chan int),
	}

	return p
}

// Start processing progress bar
func (p *PlayingBar) run() error {

	for {

		// stop progressing if song ends or skipped
		if p._progress > p.full || p.skip {
			p.skip = false
			p._progress = 0
			break
		}

		p._progress += <-p.progress

		p.text.Clear()
		start, err := time.ParseDuration(strconv.Itoa(p._progress) + "s")

		if err != nil {
			return tracerr.Wrap(err)
		}

		end, err := time.ParseDuration(strconv.Itoa(p.full) + "s")

		if err != nil {
			return tracerr.Wrap(err)
		}

		_, _, width, _ := p.GetInnerRect()
		progressBar := progresStr(p._progress, p.full, width/2, "█", "━")
		// our progress bar
		p.text.SetText(fmt.Sprintf("%s ┃%s┫ %s",
			fmtDuration(start),
			progressBar,
			fmtDuration(end),
		))

	}

	return nil
}

// Updates song title
func (p *PlayingBar) setSongTitle(title string) {
	p.Clear()
	titleColor := gomu.colors.title
	p.AddText(title, true, tview.AlignCenter, titleColor)
}

// Resets progress bar, ready for execution
func (p *PlayingBar) newProgress(songTitle string, full int) {
	p.full = full
	p._progress = 0
	p.setSongTitle(songTitle)
}

// Sets default title and progress bar
func (p *PlayingBar) setDefault() {
	p.setSongTitle("---------:---------")
	_, _, width, _ := p.GetInnerRect()
	text := fmt.Sprintf(
		"%s ┣%s┫ %s", "00:00", strings.Repeat("━", width/2), "00:00",
	)
	p.text.SetText(text)
}

// Skips the current playing song
func (p *PlayingBar) stop() {
	p.skip = true
}
