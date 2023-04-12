package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
)

var (
	fontkai *truetype.Font
)

func main() {
	r := gin.Default()

	r.POST("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"messag": "post",
		})
	})

	r.POST("/post", postPicture)

	r.Run(":8080")

}

func newpicture(name string, lox int, loy int, loz int, lon int) *image.RGBA {

	templateFile, err := os.Open("D:\\project\\image\\d1ec268fd02b4800918c45e36fcfa3cf.png")
	if err != nil {
		panic(err)
	}
	defer templateFile.Close()
	templateFileImage, err := png.Decode(templateFile)
	if err != nil {
		panic(err)
	}
	newTemplateImage := image.NewRGBA(templateFileImage.Bounds())
	draw.Draw(newTemplateImage, templateFileImage.Bounds(), templateFileImage, templateFileImage.Bounds().Min, draw.Over)

	fontkai, err := loadFont("./ttf/simkai.ttf")
	if err != nil {
		log.Panicln(err.Error())
	}

	content := freetype.NewContext()
	content.SetClip(newTemplateImage.Bounds())
	content.SetDst(newTemplateImage)
	content.SetSrc(image.Black)
	content.SetDPI(72)
	content.SetFontSize(42)
	content.SetFont(fontkai)
	content.DrawString(name, freetype.Pt(loz, lon))

	imageData, err := getDataByUrl("https://img-blog.csdnimg.cn/4e767dbcb43b447aba9b1539bbb8852c.png")
	if err != nil {
		fmt.Println("根据地址获取图片失败,err:", err.Error())
	}
	imageData = resize.Resize(387, 183, imageData, resize.Lanczos3)
	draw.Draw(newTemplateImage, newTemplateImage.Bounds().Add(image.Pt(lox, loy)), imageData, imageData.Bounds().Min, draw.Over)
	return newTemplateImage
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

func getDataByUrl(url string) (img image.Image, err error) {
	res, err := http.Get(url)
	if err != nil {
		err = fmt.Errorf("[%s]通过url获取数据失败,err:%s", url, err.Error())
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	// 读取获取的[]byte数据
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("读取数据失败,err:%s", err.Error())
		return
	}

	if !strings.HasSuffix(url, ".jpg") &&
		!strings.HasSuffix(url, ".jpeg") &&
		!strings.HasSuffix(url, ".png") {
		err = fmt.Errorf("[%s]不支持的图片类型,暂只支持.jpg、.png文件类型", url)
		return
	}

	// []byte 转 io.Reader
	reader := bytes.NewReader(data)

	if strings.HasSuffix(url, ".png") {
		img, err = png.Decode(reader)
		if err != nil {
			err = fmt.Errorf("png.Decode err:%s", err.Error())
			return
		}
	}

	return
}

type Picturemessage struct {
	Name string
	X    int
	Y    int
	Z    int
	N    int
}

func postPicture(c *gin.Context) {
	json := Picturemessage{}
	c.BindJSON(&json)
	pic := newpicture(json.Name, json.X, json.Y, json.Z, json.N)
	err := png.Encode(c.Writer, pic)
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(400, err)
		return
	}
}
