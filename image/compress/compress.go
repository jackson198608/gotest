package compress

import (
	"github.com/donnie4w/go-logger/logger"
	"gopkg.in/gographics/imagick.v2/imagick"
	"strconv"
	"strings"
)

type Compress struct {
	imgaePath string
	width     int
	height    int
	filename  string
	suffix    string
}

func NewCompress(imgaePath string, width int, height int) *Compress {
	if imgaePath == "" || width == 0 || height == 0 {
		return nil
	}

	c := new(Compress)
	if c == nil {
		return nil
	}

	c.imgaePath = imgaePath
	c.width = width
	c.height = height

	c.parsePath()
	return c
}

func (c *Compress) Do() error {
	err := c.resizeImage(c.imgaePath, c.width, c.height)
	if err == nil {
		logger.Info("[sucess] compress image path is ", c.imgaePath, " width is ", c.width, " height is ", c.height)
	}
	return err
}

func (c *Compress) resizeImage(filename string, width int, height int) error {
	var err error

	mw := imagick.NewMagickWand()

	err = mw.ReadImage(filename)

	if err != nil {
		logger.Error(err)
		return err
	}
	hWidth := uint(width)
	hHeight := uint(height)

	err = mw.ResizeImage(hWidth, hHeight, imagick.FILTER_LANCZOS, 1)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = mw.SetImageCompressionQuality(80)
	if err != nil {
		logger.Error(err)
		return err
	}
	widthStr := strconv.Itoa(width)
	heightStr := strconv.Itoa(height)
	newimg := c.filename + "_" + widthStr + "_" + heightStr + "." + c.suffix
	err = mw.WriteImage(newimg)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (c *Compress) parsePath() error {
	rawSlice := []byte(c.imgaePath)
	rawLen := len(rawSlice)
	lastIndex := strings.LastIndex(c.imgaePath, ".")
	c.filename = string(rawSlice[0:lastIndex])
	c.suffix = string(rawSlice[lastIndex+1 : rawLen])

	return nil
}