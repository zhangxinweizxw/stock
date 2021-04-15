package graphics

import (
	"code.google.com/p/graphics-go/graphics"
    "errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

const (
	IMAGE_GIF = ".gif"
	IMAGE_PNG = ".png"
	IMAGE_JPG = ".jpg"
)

// 图片转换（jpg->png）
func Convert(file io.Reader, from, to string) (io.Reader, error) {
	return nil, nil
}

// 图片缩放
func Scale(file io.Reader, size int, ext string, filePath string) error {
	src, _, err := image.Decode(file)

	if err != nil {
		return errors.New(fmt.Sprintf("image decode err:%v", err))
	}
	bound := src.Bounds()
	dx := bound.Dx()
	dy := bound.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, size, size*dy/dx))
	err = graphics.Scale(dst, src)
	graphics.Thumbnail(dst, src)

	imgFile, err := os.Create(filePath)
	defer imgFile.Close()

	var headers = make(map[string]string)
	switch ext {
	case IMAGE_PNG:
		png.Encode(imgFile, dst)
		headers["Content-Type"] = "image/png"
		break
	case IMAGE_JPG:
		jpeg.Encode(imgFile, dst, &jpeg.Options{100})
		headers["Content-Type"] = "image/jpeg"
		break
	case IMAGE_GIF:
		gif.Encode(imgFile, dst, nil)
		headers["Content-Type"] = "image/gif"
		break
	default:
		break
	}

	return nil
}
