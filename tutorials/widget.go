package tutorials

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"gui/crypto/pfx"
	"gui/crypto/pkcs12"
	"gui/crypto/sm2"
	"gui/crypto/x509"
	"gui/data"
	"gui/data/password"
	layout2 "gui/layout"
	"log"
	"strings"
	"time"
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

func verifyTab(win fyne.Window) fyne.CanvasObject {
	sign := widget.NewMultiLineEntry()
	sign.SetPlaceHolder("Please input signature. Support RSA and SM.")

	msg := widget.NewMultiLineEntry()
	msg.SetPlaceHolder("Please input message.")
	radio := widget.NewRadioGroup([]string{"RSA", "SM"}, func(s string) {
		if s == "RSA" {

		} else if s == "SM" {

		} else {

		}
	})
	radio.Horizontal = true
	radio.SetSelected("SM")
	cradio := container.NewCenter(radio)
	var resMsg = "验签成功"
	var res = true

	button := widget.NewButton("验证", func() {
		defer showDialog(&res, &resMsg, win)
		if sign.Text != "" && msg.Text != "" {
			if radio.Selected == "RSA" {
				hash := sha1.New()

				s, err := hex.DecodeString(sign.Text)
				if err != nil {
					res = false
					resMsg = "signature is invalid."
					return
				}

				m, err := base64.StdEncoding.DecodeString(msg.Text)
				if err != nil {
					res = false
					resMsg = "message is invalid."
					return
				}
				hash.Write(m)

				err = rsa.VerifyPKCS1v15(rsaKey, crypto.SHA1, hash.Sum(nil), s)
				if err != nil {
					res = false
					resMsg = "验签失败"
				}
			} else if radio.Selected == "SM" {
				s, err := hex.DecodeString(sign.Text)
				if err != nil {
					res = false
					resMsg = "signature is invalid."
					return
				}

				m, err := base64.StdEncoding.DecodeString(msg.Text)
				if err != nil {
					res = false
					resMsg = "message is invalid."
					return
				}

				verified := smKey.Verify(m, s)
				if !verified {
					res = false
					resMsg = "验签失败"
				}

			}
		} else {
			res = false
			resMsg = "验签失败"
		}
	})

	return container.NewVBox(sign, msg, cradio, button)
}

func singTab(win fyne.Window) fyne.CanvasObject {
	msg := widget.NewMultiLineEntry()
	msg.SetPlaceHolder("Please input plain text.")
	msg.Resize(fyne.NewSize(512, 150))
	input := widget.NewEntry()
	input.Disable()
	selectFile := widget.NewButton("选择私钥", func() {
		f := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				input.SetText(reader.URI().Path())
			} else {
				input.SetText("")
			}
		}, win)
		f.SetFilter(storage.NewExtensionFileFilter([]string{".sm2", ".pfx"}))
		f.Show()
	})

	file := container.New(layout2.NewHBoxLayout(), selectFile, input)

	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("Output signature.")
	output.Resize(fyne.NewSize(512, 200))

	button := widget.NewButton("签名", func() {
		if msg.Text == "" || input.Text == "" {
			return
		}
		p := ""
		if p = password.Get(input.Text); p != "" {
			sign(input.Text, p, msg.Text)
		} else {
			passwordItem := widget.NewFormItem("密码", widget.NewPasswordEntry())
			passwordDialog := dialog.NewForm("请输入私钥密码", "确认", "取消", []*widget.FormItem{passwordItem}, func(b bool) {
				if !b {
					return
				}
				p = passwordItem.Widget.(*widget.Entry).Text
				s, err := sign(input.Text, p, msg.Text)
				if err != nil {
					output.SetText(err.Error())
				}
				output.SetText(s)
				password.Put(&password.Password{
					K:          input.Text,
					V:          p,
					ExpireTime: time.Now().Unix(),
				})

			}, win)
			passwordDialog.Resize(fyne.NewSize(250, 150))
			passwordDialog.Show()
		}

	})

	return container.NewVBox(msg, file, button, output)

}

func sign(path, password, msg string) (string, error) {
	index := strings.Index(path, ".")
	ext := path[index:]
	if ext == ".sm2" {
		privateKey, err := pkcs12.GetPrivateKeyFromSm2File(path, password)
		if err != nil {
			return "", err
		}
		signature, err := privateKey.Sign(rand.Reader, []byte(msg), nil)
		if err != nil {
			return "", err
		}
		return hex.EncodeToString(signature), nil

	} else if ext == ".pfx" {
		privatekey, err := pfx.GetPrivateKeyFromPfxFile(path, password)
		if err != nil {
			return "", err
		}
		hash := sha1.New()
		hash.Write([]byte(msg))

		signature, err := rsa.SignPKCS1v15(rand.Reader, privatekey, crypto.SHA1, hash.Sum(nil))
		if err != nil {
			return "", err
		}
		return hex.EncodeToString(signature), nil
	} else {
		return "", errors.New("私钥文件不正确")
	}
}

func showDialog(res *bool, msg *string, win fyne.Window) {
	if *res {
		dialog.ShowInformation("验签结果", *msg, win)
	} else {
		dialog.ShowError(errors.New(*msg), win)
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
