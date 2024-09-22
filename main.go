package main

import (
	"strconv"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/JustTimmm/GoColor"
	"github.com/go-vgo/robotgo"
)

var (
	running bool
	mu      sync.Mutex
)

func AutoClicker(interval time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		mu.Lock()
		if !running {
			mu.Unlock()
			return
		}
		mu.Unlock()

		robotgo.Click()
		time.Sleep(interval)
	}
}

func main() {
	devMod := false

	GoColor.SuccessLog("GoClicker started ✨ \n")

	a := app.New()
	w := a.NewWindow("GoClicker")
	w.Resize(fyne.NewSize(600, 400))
	w.SetFixedSize(true)

	wSettings := a.NewWindow("GoClicker | Settings")
	wSettings.Resize(fyne.NewSize(400, 400))
	wSettings.SetFixedSize(true)

	cpsEntry := widget.NewEntry()
	cpsEntry.SetText("500")

	titleLabel := widget.NewLabel("GoClicker")
	titleContainer := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), titleLabel, layout.NewSpacer())

	startButton := widget.NewButton("START\nP", func() {
		if devMod {
			GoColor.DebugLog("START button clicked\n")
		}

		mu.Lock()
		if !running {
			GoColor.DebugLog("AutoClick Started !\n")
			running = true

			interval, err := strconv.Atoi(cpsEntry.Text)
			if err != nil {
				GoColor.ErrorLog("Invalid CPS value\n")
				mu.Unlock()
				return
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go AutoClicker(time.Duration(interval)*time.Millisecond, &wg)
		}
		mu.Unlock()
	})

	stopButton := widget.NewButton("STOP\nM", func() {
		if devMod {
			GoColor.DebugLog("STOP button clicked\n")
		}

		mu.Lock()
		if running {
			GoColor.DebugLog("AutoClick Stopped !\n")
			running = false
		}
		mu.Unlock()
	})

	settingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		wSettings.Show()
	})

	rightContainer := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), settingsButton)

	content := container.NewBorder(
		titleContainer,
		nil,
		nil,
		rightContainer,
		container.New(layout.NewGridLayout(1), startButton, stopButton),
	)

	w.SetContent(content)

	w.SetCloseIntercept(func() {
		if wSettings != nil {
			wSettings.Close()
		}
		a.Quit()
	})

	wSettings.SetContent(cpsEntry)

	w.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		switch key.Name {
		case fyne.KeyP:
			startButton.OnTapped()
		case fyne.KeyM:
			stopButton.OnTapped()
		case fyne.KeyB:
			if devMod {
				a.Quit()
			}
		}
	})

	w.ShowAndRun()
}
