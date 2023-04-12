package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"messag": "pong",
		})
	})
	r.GET("/word", getWord)
	r.GET("/image", getImage)

	//请求静态图片资源
	r.Static("/img", "./image")

	r.Run(":8080")
}

func newImage(name string, reward string, x string, y string) *image.RGBA {
	//m := len([]rune(reward))
	z, _ := strconv.Atoi(x)
	n, _ := strconv.Atoi(y)
	if n == 0 || z == 0 {
		n = 520
		z = 300
	}
	picture, err := os.Open("./image/d1ec268fd02b4800918c45e36fcfa3cf.png")
	if err != nil {
		panic(err)
	}
	defer picture.Close()

	picturefile, err := png.Decode(picture)
	if err != nil {
		panic(err)
	}

	newpicture := image.NewRGBA(picturefile.Bounds())

	draw.Draw(newpicture, picturefile.Bounds(), picturefile, picturefile.Bounds().Min, draw.Over)

	fontkai, err := loadFont("./ttf/simkai.ttf")
	if err != nil {
		log.Panicln(err.Error())
	}

	content := freetype.NewContext()
	content.SetClip(newpicture.Bounds())
	content.SetDst(newpicture)
	content.SetSrc(image.Black)
	content.SetDPI(72)
	content.SetFontSize(42)
	content.SetFont(fontkai)
	content.DrawString(name+"同学:", freetype.Pt(160, 375))
	content.DrawString("您在2020年度表现突出,成绩优异、认真负责,", freetype.Pt(230, 450))
	content.DrawString("被评为", freetype.Pt(160, 520))
	//content.DrawString("特发此证,以资鼓励。", freetype.Pt((m+4)*42, n))
	content.DrawString("特发此证,以资鼓励。", freetype.Pt(520, 520))

	content.SetFontSize(42)
	content.SetSrc(image.NewUniform(color.RGBA{R: 237, G: 39, B: 90, A: 255}))
	content.DrawString(reward+"奖", freetype.Pt(z, n))

	content.SetFontSize(32)
	content.SetSrc(image.Black)
	content.DrawString("软件信息工程学院", freetype.Pt(898, 660))
	content.DrawString("二零二零年十月", freetype.Pt(898, 726))

	return newpicture
}

//图片动态位置
func newPicture(lox string, loy string) *image.RGBA {
	x, _ := strconv.Atoi(lox)
	y, _ := strconv.Atoi(loy)
	if x == 0 || y == 0 {
		x = 520
		y = 300
	}
	picture, err := os.Open("./image/d1ec268fd02b4800918c45e36fcfa3cf.png")
	if err != nil {
		panic(err)
	}
	defer picture.Close()
	picturefile, err := png.Decode(picture)
	if err != nil {
		panic(err)
	}
	newPicture := image.NewRGBA(picturefile.Bounds())
	draw.Draw(newPicture, picturefile.Bounds(), picturefile, picturefile.Bounds().Min, draw.Over)

	imageData, err := getDataByUrl("http://qiniu.yueda.vip/123.png")
	if err != nil {
		fmt.Println("根据地址获取图片失败,err:", err.Error())
	}
	imageData = resize.Resize(387, 183, imageData, resize.Lanczos3)
	draw.Draw(newPicture, newPicture.Bounds().Add(image.Pt(x, y)), imageData, imageData.Bounds().Min, draw.Over)
	return newPicture
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

func getWord(c *gin.Context) {
	name := c.Query("name")
	reward := c.Query("reward")
	x := c.Query("x")
	y := c.Query("y")
	image := newImage(name, reward, x, y)
	err := png.Encode(c.Writer, image)
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(400, err)
		return
	}
}

func getImage(c *gin.Context) {
	lox := c.Query("lox")
	loy := c.Query("loy")
	picture := newPicture(lox, loy)
	err := png.Encode(c.Writer, picture)
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(400, err)
		return
	}
}
