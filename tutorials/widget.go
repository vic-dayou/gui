package tutorials

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"gui/crypto/pfx"
	"gui/crypto/pkcs12"
	"gui/data"
	"gui/data/password"
	"gui/httpclient"
	layout2 "gui/layout"
	"log"
	"strings"
	"time"
)

func encodeTab(_ fyne.Window) fyne.CanvasObject {

	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("Please past text...")
	input.Wrapping = fyne.TextWrapBreak

	input.Resize(fyne.NewSize(512, 200))

	output := widget.NewMultiLineEntry()
	output.Disable()
	enc := widget.NewButton("编码", func() {
		output.SetText(base64.StdEncoding.EncodeToString([]byte(input.Text)))
	})

	dec := widget.NewButton("解码", func() {
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
	sign.SetPlaceHolder("请输入签名值，支持SM2和RSA签名")
	sign.Wrapping = fyne.TextWrapBreak
	sign.Resize(fyne.NewSize(512, 150))

	msg := widget.NewMultiLineEntry()
	msg.SetPlaceHolder("请输入Base64后的请求报文")
	msg.Wrapping = fyne.TextWrapBreak
	msg.Resize(fyne.NewSize(512, 150))
	radio := widget.NewRadioGroup([]string{"RSA", "SM"}, func(s string) {
		if s == "RSA" {

		} else if s == "SM" {

		} else {

		}
	})
	radio.Horizontal = true
	radio.SetSelected("SM")
	cradio := container.NewCenter(radio)
	var resMsg = "验签失败"
	var res = false

	button := widget.NewButton("验证", func() {
		defer showDialog(&res, &resMsg, win)
		if sign.Text != "" && msg.Text != "" {
			if radio.Selected == "RSA" {
				hash := sha1.New()
				s, err := hex.DecodeString(strings.Trim(strings.Trim(sign.Text, "\r\n"), " "))
				if err != nil {
					res = false
					resMsg = "signature is invalid."
					return
				}

				m, err := base64.StdEncoding.DecodeString(strings.Trim(msg.Text, " "))
				if err != nil {
					res = false
					resMsg = "message is invalid."
					return
				}
				hash.Write(m)
				for sn, key := range data.RSAPool {
					err = rsa.VerifyPKCS1v15(key, crypto.SHA1, hash.Sum(nil), s)
					if err == nil {
						res = true
						resMsg = fmt.Sprintf("使用SN:%s的证书验签成功.", sn)
						return
					}

				}

			} else if radio.Selected == "SM" {
				//log.Println(strings.Trim(strings.Trim(sign.Text,"\r\n")," "))

				s, err := hex.DecodeString(strings.Trim(strings.Trim(sign.Text, "\r\n"), " "))
				if err != nil {
					res = false
					resMsg = "signature is invalid."
					return
				}

				m, err := base64.StdEncoding.DecodeString(strings.Trim(msg.Text, " "))
				if err != nil {
					res = false
					resMsg = "message is invalid."
					return
				}

				for sn, key := range data.SM2Pool {
					verified := key.Verify(m, s)
					if verified {
						res = true
						resMsg = fmt.Sprintf("使用SN:%s的证书验签成功.", sn)
						return
					}

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

func sendMsg(win fyne.Window) fyne.CanvasObject {
	msg := widget.NewMultiLineEntry()
	msg.SetPlaceHolder("请输入XML格式的请求报文")
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

	button := widget.NewButton("发送", func() {
		if msg.Text == "" || input.Text == "" {
			return
		}
		p := ""
		message := base64.StdEncoding.EncodeToString([]byte(msg.Text))
		if p = password.Get(input.Text); p != "" {
			s, err := sign(input.Text, p, msg.Text)
			if err != nil {
				output.SetText(err.Error())
			}
			params := []httpclient.NameValuePair{{
				Key:   "message",
				Value: message,
			},
				{
					Key:   "signature",
					Value: s,
				},
			}

			body, err := httpclient.Post(params, "https://www.china-clearing.com/Gateway/InterfaceII")
			if err != nil {
				output.SetText(err.Error())
				return
			}
			output.SetText(string(body))
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

func showDialog(res *bool, msg *string, win fyne.Window) {
	if *res {
		dialog.ShowInformation("验签结果", *msg, win)
	} else {
		dialog.ShowError(errors.New(*msg), win)
	}
}
