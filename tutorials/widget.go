package tutorials

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"gui/crypto/sm2"
	"gui/crypto/sm3"
	"gui/crypto/x509"
	"gui/data"
	layout2 "gui/layout"
	"log"
)

var rsaKey *rsa.PublicKey
var smKey *sm2.PublicKey

func init() {
	loadPublicKey()
}

func encodeTab(_ fyne.Window) fyne.CanvasObject {

	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("Please past text...")
	input.Wrapping = fyne.TextWrapBreak

	input.Resize(fyne.NewSize(512, 200))

	output := widget.NewMultiLineEntry()
	output.Disable()
	enc := widget.NewButton("Encode", func() {
		output.SetText(base64.StdEncoding.EncodeToString([]byte(input.Text)))
	})

	dec := widget.NewButton("Decode", func() {
		text, err := base64.StdEncoding.DecodeString(input.Text)
		if err != nil {
			log.Println(err)
			output.SetText("decode error")
			return
		}
		output.SetText(string(text))
	})

	button := container.NewHBox(layout.NewSpacer(), enc, dec, layout.NewSpacer())
	output.Resize(fyne.NewSize(512, 200))
	return container.New(layout2.NewVBoxLayout(),
		input,
		button,
		output,
	)
}

func signTab(win fyne.Window) fyne.CanvasObject {
	sign := widget.NewMultiLineEntry()
	sign.SetPlaceHolder("Please input signature. Support RSA and SM.")

	msg := widget.NewMultiLineEntry()
	radio := widget.NewRadioGroup([]string{"RSA", "SM"}, func(s string) {
		if s == "RSA" {

		} else if s == "SM" {

		} else {

		}
	})
	radio.Horizontal = true
	cradio := container.NewCenter(radio)
	button := widget.NewButton("验证", func() {
		var resMsg = "验签成功"
		var res = true
		defer showDialog(res, resMsg, win)
		if sign.Text != "" && msg.Text != "" {
			if radio.Selected == "RSA" {
				hash := sha1.New()
				m, err := base64.StdEncoding.DecodeString(msg.Text)
				if err != nil {
					res = false
					resMsg = "message is invalid."
					return
				}
				hash.Write(m)

				s, err := hex.DecodeString(sign.Text)
				if err != nil {
					res = false
					resMsg = "signature is invalid."
					return
				}
				err = rsa.VerifyPKCS1v15(rsaKey, crypto.SHA1, hash.Sum(nil), s)
				if err != nil {
					res = false
					resMsg = "验签失败"
				}
			} else if radio.Selected == "SM" {
				hash := sm3.New()
				m, err := base64.StdEncoding.DecodeString(msg.Text)
				if err != nil {
					res = false
					resMsg = "message is invalid."
					return
				}
				hash.Write(m)

				s, err := hex.DecodeString(sign.Text)
				if err != nil {
					res = false
					resMsg = "signature is invalid."
					return
				}
				verified := smKey.Verify(m, s)
				if !verified {
					res = false
					resMsg = "验签失败"
				}

			}
		}
	})

	return container.NewVBox(sign, msg, cradio, button)
}

func showDialog(res bool, msg string, win fyne.Window) {
	if res {
		dialog.ShowInformation("验签结果", msg, win)
	} else {
		dialog.ShowError(errors.New(msg), win)
	}
}

func loadPublicKey() {
	rsaCer, err := x509.ParseCertificate(data.GetPemBytes("RSA"))
	if err != nil {
		log.Println(err)
	}

	rsaKey = rsaCer.PublicKey.(*rsa.PublicKey)

	sm2Cer, err := x509.ParseCertificate(data.GetPemBytes("SM"))
	if err != nil {
		log.Println(err)
	}

	smKey = sm2Cer.PublicKey.(*sm2.PublicKey)

}
