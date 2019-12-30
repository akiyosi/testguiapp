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
	queue    chan bool
	keycount int
}

func (a *myapp) genRandRect() *core.QRectF {
	x := rand.Intn(a.wid.Width()-1) + 1
	y := rand.Intn(a.wid.Height()-1) + 1
	width := int(a.wid.Width()) / (rand.Intn(64) + 1)
	height := int(a.wid.Height()) / (rand.Intn(64) + 1)

	rand.Seed(time.Now().UnixNano())

	return core.NewQRectF4(
		float64(x),
		float64(y),
		float64(width),
		float64(height),
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
		queue:  make(chan bool, 1000),
	}
	a.signal.ConnectRedrawSignal(func() {
		<-a.queue
		a.keycount++
		a.wid.Update()
	})

	rand.Seed(time.Now().UnixNano())

	// Draw random rectangle
	wid.ConnectPaintEvent(func(e *gui.QPaintEvent) {
		p := gui.NewQPainter2(wid)

		/* Paint process 
		  This process is very fast and commenting out doesn't change update performance
		*/
		p.FillRect4(
			a.genRandRect(),
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
		a.queue <- true
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
