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
	"gui/crypto/sm2"
	"gui/data"
	"gui/data/password"
	"gui/httpclient"
	layout2 "gui/layout"
	"io/ioutil"
	"log"
	"net/url"
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
	var pwd *password.Password
	var reader fyne.URIReadCloser
	selectFile := widget.NewButton("选择私钥", func() {
		f := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
			if r != nil {
				input.SetText(r.URI().Path())
				reader = r
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
		// 如果密码记录为空，则提示输入密码
		pwd = password.Get(input.Text)
		index := strings.Index(input.Text, ".")
		ext := input.Text[index:]
		passwordItem := widget.NewFormItem("密码", widget.NewPasswordEntry())
		passwordDialog := dialog.NewForm("请输入私钥密码", "确认", "取消", []*widget.FormItem{passwordItem},
			func(b bool) {
				if !b {
					return
				}
				pwds := passwordItem.Widget.(*widget.Entry).Text
				bytes, err := ioutil.ReadAll(reader)
				if err != nil || len(bytes) == 0 {
					output.SetText("读取文件失败")
					return
				}
				privateKey, err := pkcs12.GetPrivateKeyFromBytes(bytes, ext, pwds)

				if err != nil {
					output.SetText(err.Error())
					return
				}
				signature, err := signByKey(ext, msg.Text, privateKey)
				if err != nil {
					output.SetText(err.Error())
					return
				}
				output.SetText(signature)
				password.Put(&password.Password{
					K:          input.Text,
					V:          pwds,
					Ext:        ext,
					Pk:         privateKey,
					ExpireTime: time.Now().Unix(),
				})
			}, win)
		passwordDialog.Resize(fyne.NewSize(250, 150))
		//passwordDialog.Show()
		if pwd == nil {
			passwordDialog.Show()
		} else {
			passwordDialog.Hide()
			signature, err := signByKey(ext, msg.Text, pwd.Pk)
			if err != nil {
				output.SetText(err.Error())
				return
			}
			output.SetText(signature)
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
func signByKey(ext, msg string, priv interface{}) (string, error) {
	if ext == ".sm2" {
		privateKey := priv.(*sm2.PrivateKey)
		signature, err := privateKey.Sign(rand.Reader, []byte(msg), nil)
		if err != nil {
			return "", err
		}
		return hex.EncodeToString(signature), nil
	} else if ext == ".pfx" {
		privateKey := priv.(*rsa.PrivateKey)
		hash := sha1.New()
		hash.Write([]byte(msg))

		signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hash.Sum(nil))
		if err != nil {
			return "", err
		}
		return hex.EncodeToString(signature), nil
	} else {
		return "", errors.New("签名失败")
	}
}

func sendMsgTab(win fyne.Window) fyne.CanvasObject {
	msg := widget.NewMultiLineEntry()
	msg.SetPlaceHolder("请输入XML格式的请求报文")
	msg.Resize(fyne.NewSize(512, 150))
	input := widget.NewEntry()
	input.Disable()
	var pwd *password.Password
	var data []byte
	selectFile := widget.NewButton("选择私钥", func() {
		f := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
			if r != nil {
				input.SetText(r.URI().Path())
				data, _ = ioutil.ReadAll(r)
			} else {
				input.SetText("")
			}
		}, win)
		f.SetFilter(storage.NewExtensionFileFilter([]string{".sm2", ".pfx"}))
		f.Show()
	})
	urlLabel := widget.NewLabel("请求地址")
	urlEntry := widget.NewEntry()
	urlEntry.SetText("https://test.cpcn.com.cn/Gateway/InterfaceII")
	urlEntry.OnChanged = func(s string) {
		if s == "" {
			urlEntry.SetText("https://test.cpcn.com.cn/Gateway/InterfaceII")
		}
	}

	urlContent := container.New(layout2.NewHBoxLayout(), urlLabel, urlEntry)

	file := container.New(layout2.NewHBoxLayout(), selectFile, input)

	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("响应报文.")
	output.Resize(fyne.NewSize(512, 200))

	button := widget.NewButton("发送", func() {
		if msg.Text == "" || input.Text == "" {
			return
		}
		// 如果密码记录为空，则提示输入密码
		pwd = password.Get(input.Text)
		index := strings.Index(input.Text, ".")
		ext := input.Text[index:]

		passwordItem := widget.NewFormItem("密码", widget.NewPasswordEntry())
		passwordDialog := dialog.NewForm("请输入私钥密码", "确认", "取消", []*widget.FormItem{passwordItem}, func(b bool) {
			if !b {
				return
			}
			pwds := passwordItem.Widget.(*widget.Entry).Text

			if len(data) == 0 {
				output.SetText("请重新选择私钥文件")
				return
			}
			privateKey, err := pkcs12.GetPrivateKeyFromBytes(data, ext, pwds)
			if err != nil {
				output.SetText(err.Error())
				return
			}
			signature, err := signByKey(ext, msg.Text, privateKey)
			if err != nil {
				output.SetText(err.Error())
				return
			}
			body, err := send(msg.Text, signature, urlEntry.Text)
			if err != nil {
				output.SetText(err.Error())
				return
			}
			output.SetText(body)
			password.Put(&password.Password{
				K:          input.Text,
				V:          pwds,
				Ext:        ext,
				Pk:         privateKey,
				ExpireTime: time.Now().Unix(),
			})

		}, win)
		passwordDialog.Resize(fyne.NewSize(250, 150))

		if pwd == nil {
			passwordDialog.Show()
		} else {
			signature, err := signByKey(ext, msg.Text, pwd.Pk)
			body, err := send(msg.Text, signature, urlEntry.Text)
			if err != nil {
				output.SetText(err.Error())
				return
			}
			output.SetText(body)
		}

	})

	return container.NewVBox(msg, file, urlContent, button, output)
}

func send(message, signature, link string) (string, error) {
	// 1. 对message进行base64编码
	m := base64.StdEncoding.EncodeToString([]byte(message))
	values := url.Values{}
	values.Set("message", m)
	values.Set("signature", signature)
	body, err := httpclient.Post(values, link)
	// resp[0]:响应信息，resp[1]:签名
	resp := strings.Split(string(body), ",")
	if len(resp) != 2 {
		return "", errors.New("响应信息非message,signature格式")
	}
	respText, _ := base64.StdEncoding.DecodeString(resp[0])
	return string(respText), err
}

func showDialog(res *bool, msg *string, win fyne.Window) {
	if *res {
		dialog.ShowInformation("验签结果", *msg, win)
	} else {
		dialog.ShowError(errors.New(*msg), win)
	}
}
