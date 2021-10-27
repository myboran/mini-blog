package article_service

import (
	"gin-blog/pkg/qrcode"
	"gin-blog/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/file"
	"github.com/golang/freetype"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"os"
)

type ArticlePoster struct {
	PosterName string
	*Article
	Qr *qrcode.QrCode
}

func NewArticlePoster(posterName string, article *Article, qr *qrcode.QrCode) *ArticlePoster {
	return &ArticlePoster{
		PosterName: posterName,
		Article: article,
		Qr: qr,
	}
}

func GetPosterFlag() string {
	return "poster"
}

func (a *ArticlePoster) CheckMergedImage(path string) bool{
	if file.CheckNotExist(path + a.PosterName) == true {
		return false
	}
	return true
}

func (a *ArticlePoster) OpenMergedImage(path string) (*os.File, error) {
	f, err := file.MustOpen(a.PosterName, path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

type ArticlePosterBg struct {
	Name string
	*ArticlePoster
	*Rect
	*Pt
}

type Rect struct {
	Name string
	X0 int
	Y0 int
	X1 int
	Y1 int
}

type Pt struct {
	X int
	Y int
}

func NewArticlePosterBg(name string, ap *ArticlePoster, rect *Rect, pt *Pt) *ArticlePosterBg {
	return &ArticlePosterBg{
		Name: name,
		ArticlePoster: ap,
		Rect: rect,
		Pt: pt,
	}
}

func (a *ArticlePosterBg) Generate() (string, string, error) {
	// 获取二维码储存路径
	fullPath := qrcode.GetQrCodeFullPath()

	// 生成二维码图像
	fileName, path, err := a.Qr.Encode(fullPath)
	if err != nil {
		return "", "", err
	}

	// 检查合并后的二维码图像是否存在
	if !a.CheckMergedImage(path) {

		// 若不存在，则生成待合并的图像 mergedF
		mergedF, err := a.OpenMergedImage(path)
		if err != nil {
			return "", "", err
		}
		defer mergedF.Close()

		// 打开背景图bgF
		bgF, err := file.MustOpen(a.Name, path)
		if err != nil {
			return "", "", err
		}
		defer bgF.Close()

		// 打开生成的二维码图像qrf
		qrf, err := file.MustOpen(fileName, path)
		if err != nil {
			return "", "", err
		}
		defer qrf.Close()

		// 解码bgF和qrF返回 iamge.Image
		bgImage,err := jpeg.Decode(bgF)
		if err != nil {
			return "", "", err
		}

		qrImage, err := jpeg.Decode(qrf)
		if err != nil {
			return "", "", err
		}

		// 创建一个新的 RGBA 图像
		jpg := image.NewRGBA(image.Rect(a.Rect.X0, a.Rect.Y0, a.Rect.X1, a.Rect.Y1))

		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(a.Pt.X, a.Pt.Y)), draw.Over)

		err = a.DrawPoster(&DrawText{
			JPG:    jpg,
			Merged: mergedF,

			Title: "Golang Gin 系列文章",
			X0:    80,
			Y0:    160,
			Size0: 42,

			SubTitle: "---煎鱼",
			X1:       320,
			Y1:       220,
			Size1:    36,
		}, "msyhbd.ttc")
		if err != nil {
			return "", "", err
		}

		jpeg.Encode(mergedF, jpg, nil)
	}
	return fileName, path, nil
}

type DrawText struct {
	JPG draw.Image
	Merged *os.File

	Title string
	X0 int
	Y0 int
	Size0 float64

	SubTitle string
	X1 int
	Y1 int
	Size1 float64
}

func (a *ArticlePosterBg) DrawPoster(d *DrawText, fontName string) error {
	fontSource := setting.AppSetting.RuntimeRootPath + setting.AppSetting.FontSavePath + fontName
	fontSourceBytes, err := ioutil.ReadFile(fontSource)
	if err != nil {
		return err
	}

	trueTypeFont, err := freetype.ParseFont(fontSourceBytes)
	if err != nil {
		return err
	}

	fc := freetype.NewContext()
	fc.SetDPI(72)
	fc.SetFont(trueTypeFont)
	fc.SetFontSize(d.Size0)
	fc.SetClip(d.JPG.Bounds())
	fc.SetDst(d.JPG)
	fc.SetSrc(image.Black)

	pt := freetype.Pt(d.X0, d.Y0)
	_, err = fc.DrawString(d.Title, pt)
	if err != nil {
		return err
	}

	fc.SetFontSize(d.Size1)
	_, err = fc.DrawString(d.SubTitle, freetype.Pt(d.X1, d.Y1))
	if err != nil {
		return err
	}

	err = jpeg.Encode(d.Merged, d.JPG, nil)
	if err != nil {
		return err
	}

	return nil
}