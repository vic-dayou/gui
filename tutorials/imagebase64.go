package tutorials

import (
	"bytes"
	"encoding/base64"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	_ "golang.org/x/image/bmp"
	layout2 "gui/layout"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func imageBase64Tab(win fyne.Window) fyne.CanvasObject {
	intro := widget.NewLabel("对图片转为jpg并压缩进行base64，输出图片到D:/Users/目录下")
	output := widget.NewMultiLineEntry()
	output.Resize(fyne.NewSize(512, 250))
	output.Wrapping = fyne.TextTruncate
	label := widget.NewLabel("压缩比例")
	f := 0.75
	data := binding.BindFloat(&f)
	slide := widget.NewSliderWithData(0, 1, data)
	slide.Step = 0.01
	bar := widget.NewProgressBarWithData(data)
	slidebar := container.NewVBox(slide, bar)
	slidebarWithLabel := container.New(layout2.NewHBoxLayout(), label, slidebar)

	selectFile := widget.NewButton("选择图片", func() {
		f := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
			if r != nil {
				defer r.Close()
				if f == 1.0 {
					data, err := ioutil.ReadAll(r)
					if err != nil {
						output.SetText("读取文件失败")
					}
					bs := base64.StdEncoding.EncodeToString(data)
					output.SetText(bs)
					return
				}
				img, _, err := image.Decode(r)
				if err != nil || img == nil {
					output.SetText("图片格式不正确")
				}
				buf := bytes.Buffer{}
				ext := r.URI().Extension()
				log.Println("ext: ", ext)
				if ext == ".jpg" || ext == ".jpeg" {
					log.Println("quality ", int(f*100))
					err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: int(f * 100)})
				} else if ext == ".png" || ext == ".bmp" {
					newImage := image.NewRGBA(img.Bounds())
					draw.Draw(newImage, newImage.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
					draw.Draw(newImage, newImage.Bounds(), img, img.Bounds().Min, draw.Over)
					err = jpeg.Encode(&buf, newImage, &jpeg.Options{Quality: int(f * 100)})
				}
				if err != nil {
					log.Println(err.Error())
					output.SetText("图片解析错误")
					return
				}
				bs := base64.StdEncoding.EncodeToString(buf.Bytes())

				filename := strings.Split(r.URI().Name(), ".")
				_, err = os.Stat("D:/Users/")
				if os.IsNotExist(err) {
					os.Mkdir("D:/Users/", os.ModePerm)
				}
				f, err := os.Create("D:/Users/" + filename[0] + ".jpg")
				if err != nil {
					log.Println(err)
					output.SetText("创建文件失败")
					return
				}
				f.Write(buf.Bytes())
				f.Close()

				output.SetText(bs)
			}
		}, win)
		f.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".png", ".jpeg", ".bmp"}))
		f.Show()
	})
	return container.NewVBox(intro, selectFile, slidebarWithLabel, output)

}
