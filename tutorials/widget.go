package tutorials

import (
	"encoding/base64"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	layout2 "gui/layout"
	"log"
)





func encodeTab(_ fyne.Window) fyne.CanvasObject {

	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("Please past text...")
	input.Wrapping = fyne.TextWrapBreak

	input.Resize(fyne.NewSize(512,200))


	output := widget.NewMultiLineEntry()
	output.Disable()
	enc := widget.NewButton("Encode", func() {
		output.SetText(base64.StdEncoding.EncodeToString([]byte(input.Text)))
	})

	dec := widget.NewButton("Decode", func() {
		text,err := base64.StdEncoding.DecodeString(input.Text)
		if err != nil {
			log.Println(err)
			output.SetText("decode error")
			return
		}
		output.SetText(string(text))
	})

	button := container.NewHBox(layout.NewSpacer(),enc,dec,layout.NewSpacer())
	output.Resize(fyne.NewSize(512,200))
	return container.New(layout2.NewVBoxLayout(),
		input,
		button,
		output,
	)
}

func signTab(win fyne.Window) fyne.CanvasObject {
	sign := widget.NewMultiLineEntry()
	sign.SetPlaceHolder("Please input signature. Support RSA and SM.")
	button := widget.NewButton("验证", func() {
		if sign.Text == "Success" {
			dialog.ShowInformation("验签结果","验签成功",win)
		}
	})

	return container.NewVBox(sign,button)
}