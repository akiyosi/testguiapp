package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type appSignal struct {
	core.QObject
	_ func() `signal:"redrawSignal"`
}

type myapp struct {
	app *widgets.QApplication
	win *widgets.QMainWindow
	wid *widgets.QWidget

	signal   *appSignal
	queue    chan *core.QRect
	rect     *core.QRect
	keycount int
}

func (a *myapp) genRandRect() *core.QRect {
	x := rand.Intn(a.wid.Width()-1) + 1
	y := rand.Intn(a.wid.Height()-1) + 1
	width := int(a.wid.Width()) / (rand.Intn(64) + 1)
	height := int(a.wid.Height()) / (rand.Intn(64) + 1)

	rand.Seed(time.Now().UnixNano())

	return core.NewQRect4(
		x,
		y,
		width,
		height,
	)
}

func main() {
	app := widgets.NewQApplication(0, nil)
	win := widgets.NewQMainWindow(nil, 0)
	wid := widgets.NewQWidget(win, 0)
	win.SetCentralWidget(wid)

	a := &myapp{
		app:    app,
		win:    win,
		wid:    wid,
		signal: NewAppSignal(nil),
		queue:  make(chan *core.QRect, 400),
	}
	a.signal.ConnectRedrawSignal(func() {
		a.rect = <-a.queue
		a.keycount++
		// a.wid.Update()
		// a.wid.Update3(rect)
		a.wid.Update3(core.NewQRect4(0, 0, a.wid.Width()/2, a.wid.Height()/2))
		a.wid.Update3(core.NewQRect4(a.wid.Width()/2, a.wid.Height()/2, a.wid.Width(), a.wid.Height()))
	})

	rand.Seed(time.Now().UnixNano())

	// Draw random rectangle
	wid.ConnectPaintEvent(func(e *gui.QPaintEvent) {
		if a.rect == nil {
			return
		}
		p := gui.NewQPainter2(wid)

		/* Paint process 
		  This process is very fast and commenting out doesn't change update performance
		*/
		p.FillRect6(
			a.rect,
			gui.NewQColor3(
				rand.Intn(255),
				rand.Intn(255),
				rand.Intn(255),
				255,
			),
		)

		p.DestroyQPainter()
	})
	win.ConnectKeyPressEvent(func(e *gui.QKeyEvent) {
		a.queue <- a.genRandRect()
		a.signal.RedrawSignal()
	})

	// Notify the rate of key repeat
	go func() {
		for _ = range time.Tick(1000 * time.Millisecond) {
			fmt.Println(a.keycount)
			a.keycount = 0
		}
	}()

	win.Show()
	wid.SetFocus2()
	app.Exec()

}
