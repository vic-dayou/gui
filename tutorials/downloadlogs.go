package tutorials

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gui/client/httpsclient"
	layout2 "gui/layout"
	"time"
)

func downloadLogFile(win fyne.Window) fyne.CanvasObject {
	/*label := widget.NewLabel("日志文件路径")
	path := widget.NewEntry()
	path.SetPlaceHolder("log.file.path")

	paths := strings.Split(path.Text,"/")
	if len(paths) != 5 {
		log.Println("日志文件格式不正确")
	}
	basePath := "https://10.136.6.12:8443/pay/"
	ip := paths[2]
	module := paths[3]
	httpsclient.Get()*/

	label := widget.NewLabel("URL:")
	url := widget.NewEntry()
	drv := fyne.CurrentApp().Driver()
	button := widget.NewButtonWithIcon("下载", theme.DownloadIcon(), func() {
		if url.Text == "" {
			return
		}
		err := httpsclient.Download(url.Text)
		if drv, ok := drv.(desktop.Driver); ok {
			w := drv.CreateSplashWindow()
			msg := "下载成功"
			if err != nil {
				msg = err.Error()
			}
			w.SetContent(widget.NewLabelWithStyle(msg,
				fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
			w.Show()
			go func() {
				time.Sleep(time.Second * 3)
				w.Close()
			}()

		}
	})
	box := container.New(layout2.NewHBoxLayout(), label, url)

	return container.NewVBox(box, button)

}
