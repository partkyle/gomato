package main

import (
	"fmt"
	"log"
	"time"

	"github.com/getlantern/systray"
)

const DefaultTitle = "🍅"
const DefaultTimerTitle = "Toggle Timer"

type App struct {
	timer        *systray.MenuItem
	quit         *systray.MenuItem
	stopSignal   chan struct{}
	timerRunning bool
}

func NewApp() *App {
	app := &App{
		timer:      systray.AddMenuItem(DefaultTimerTitle, "Start a pomodoro"),
		quit:       systray.AddMenuItem("Quit", "No more tomatoes"),
		stopSignal: make(chan struct{}),
	}
	app.cleanup()
	return app
}

func (a *App) gomainloop() {
	go a.mainloop()
}

func (a *App) mainloop() {
	for {
		a.main()
	}
}

func (a *App) main() {
	select {
	case <-a.timer.ClickedCh:
		a.toggleTimer()
	case <-a.quit.ClickedCh:
		systray.Quit()
	}
}

func (a *App) updateTimeFromDuration(d time.Duration) {
	systray.SetTitle(fmt.Sprintf("%s %s", DefaultTitle, d.String()))
}

func (a *App) toggleTimer() {
	if a.timerRunning {
		a.stopTimer()
	} else {
		a.startTimer()
	}
}

func (a *App) startTimer() {

	a.timerRunning = true

	go func() {
		defer a.cleanup()

		ticker := time.NewTicker(time.Second)

		defer ticker.Stop()

		maxTime := 25 * 60
		for i := 0; i < maxTime; i = i + 1 {
			duration := time.Second * time.Duration(maxTime-i)
			a.updateTimeFromDuration(duration)
			select {
			case <-a.stopSignal:
				log.Println("got stop timers")
				return
			case <-ticker.C:
			}
		}
	}()
}

func (a *App) stopTimer() {
	a.stopSignal <- struct{}{}
}

func (a *App) cleanup() {
	a.timerRunning = false
	a.timer.SetTitle(DefaultTimerTitle)
	systray.SetTitle(DefaultTitle)
}

func onReady() {
	app := NewApp()
	app.gomainloop()
}

func onExit() {}

func main() {
	systray.Run(onReady, onExit)
}
