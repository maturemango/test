package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

var (
	fontkai *truetype.Font
)

func main() {
	templateFile, err := os.Open("D:\\project\\image\\d1ec268fd02b4800918c45e36fcfa3cf.png")
	if err != nil {
		panic(err)
	}
	defer templateFile.Close()

	templateFileImage, err := png.Decode(templateFile)
	if err != nil {
		panic(err)
	}

	newTempalteImage := image.NewRGBA(templateFileImage.Bounds())

	draw.Draw(newTempalteImage, templateFileImage.Bounds(), templateFileImage, templateFileImage.Bounds().Min, draw.Over)

	fontkai, err := loadFont("D:\\project\\ttf\\simkai.ttf")
	if err != nil {
		log.Panicln(err.Error())
		return
	}

	content := freetype.NewContext()
	content.SetClip(newTempalteImage.Bounds())
	content.SetDst(newTempalteImage)
	content.SetSrc(image.Black)
	content.SetDPI(72)
	content.SetFontSize(42)
	content.SetFont(fontkai)

	content.DrawString("董明宇同学:", freetype.Pt(160, 375))
	content.DrawString("您在2020年度表现突出,成绩优异、认真负责,", freetype.Pt(230, 450))
	content.DrawString("被评为", freetype.Pt(160, 520))
	content.DrawString("特发此证,以资鼓励。", freetype.Pt(520, 520))

	content.SetFontSize(42)
	content.SetSrc(image.NewUniform(color.RGBA{R: 237, G: 39, B: 90, A: 255}))
	content.DrawString("校级三等奖", freetype.Pt(300, 520))

	content.SetFontSize(32)
	content.SetSrc(image.Black)
	content.DrawString("软件信息工程学院", freetype.Pt(898, 660))
	content.DrawString("二零二零年十月", freetype.Pt(898, 726))

	saveFile(newTempalteImage)
}

func loadFont(path string) (font *truetype.Font, err error) {
	var fontBytes []byte
	fontBytes, err = ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("加载字体文件出错:%s", err.Error())
		return
	}
	font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		err = fmt.Errorf("解析字体文件出错,%s", err.Error())
		return
	}
	return
}

func saveFile(pic *image.RGBA) {
	dstFile, err := os.Create("D:\\project\\image\\2.png")
	if err != nil {
		fmt.Println(err)
	}
	defer dstFile.Close()
	png.Encode(dstFile, pic)
}
