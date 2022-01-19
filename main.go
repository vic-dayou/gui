package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"gui/bundle"
	"gui/tutorials"
	"log"
	"os"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const preferenceCurrentTutorial = "currentTutorial"

var topWindow fyne.Window

func main() {

	_ = os.Setenv("FYNE_FONT", "C:\\Windows\\Fonts\\STKAITI.ttf")
	loadui()
	/*msg := "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9Im5vIj8+CjxSZXNwb25zZSB2ZXJzaW9uPSIyLjAiPgo8SGVhZD4KPENvZGU+MjQwMDAyPC9Db2RlPgo8TWVzc2FnZT7ns7vnu5/kuK3kuI3lrZjlnKjmjIflrprnmoTmnLrmnoTvvIzmn6XnnIvlj4LmlbBJbnN0aXR1dGlvbklEPC9NZXNzYWdlPgo8L0hlYWQ+CjwvUmVzcG9uc2U+,2A44B0F9563FF4306BA5C1EF4B3DD40F2B9FEB9A1012570A852CCF5502A23526C48BA302A3A767A4A1C1856D6C16C2A3866E7D6582F2CDB00DB8A664ABF586BF8B193740076602C229A08549718D2012093E776BF16D86C38F463C5639C884F4D51B7B509DA1756D5EEE6164E7EF4BD5EA304F0F77B840D3651A5370BB233FC0ED1FE4F0E256F65D42633AB9120767ED44C234DC46738E3FD4AFA3298BAC2BB2FE2A10B93BE8AB54A3F66460A50D3BF1D5FBF4EDB6DB9C3B676C6C7835645F6A3774BD98E1801DDDF50BA7C39A7740E7F1408C96D2AF97DEB1A158DD9D1C3C548341938CE3E83CA4715602E1EC5356F7B2489E6680DA7ED87AEF62987FB596B6"
	<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	<Request version="2.1">
	<Head>
	<TxCode>4616</TxCode>
	</Head>
	<Body>
	<InstitutionID>007507</InstitutionID>
	<SourceTxSN>202201061725150615775</SourceTxSN>
	<SourceTxCode>4611</SourceTxCode>
	<OperationFlag>10</OperationFlag>
	<UserType>11</UserType>
	</Body>
	</Request>
	bytes, err := base64.StdEncoding.DecodeString(msg)
	fmt.Println(len(bytes))
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(bytes))*/

}

func loadui() {
	a := app.NewWithID("cpcn.com")
	a.SetIcon(bundle.Cpcnlogo())
	logLifecycle(a)
	w := a.NewWindow("CPCN Test")
	topWindow = w

	//w.SetMainMenu(makeMenu(a, w))
	w.SetMaster()

	content := container.NewMax()
	//title := widget.NewLabel("Component name")
	//intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	//intro.Wrapping = fyne.TextWrapWord
	setTutorial := func(t tutorials.Tutorial) {
		if fyne.CurrentDevice().IsMobile() {
			child := a.NewWindow(t.Title)
			topWindow = child
			child.SetContent(t.View(topWindow))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = w
			})
			return
		}

		//title.SetText(t.Title)
		//intro.SetText(t.Intro)

		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	//tutorial := container.NewBorder(nil, nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(makeNav(setTutorial, false))
	} else {
		split := container.NewHSplit(makeNav(setTutorial, true), content)
		split.Offset = 0.2
		w.SetContent(split)
	}
	w.Resize(fyne.NewSize(640, 460))
	w.ShowAndRun()
}

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		log.Println("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
	})
}

func makeNav(setTutorial func(tutorial tutorials.Tutorial), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return tutorials.TutorialIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := tutorials.TutorialIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := tutorials.Tutorials[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := tutorials.Tutorials[uid]; ok {
				a.Preferences().SetString(preferenceCurrentTutorial, uid)
				setTutorial(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentTutorial, "welcome")
		tree.Select(currentPref)
	}

	themes := container.New(layout.NewGridLayout(2),
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}
