package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"os"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type Texture struct {
	rendererID         uint32
	filepath           string
	width, height, bpp int32
}

func CreateTexture(filepath string) (Texture, error) {
	texture := Texture{}
	texture.filepath = filepath

	imgFile, err := os.Open(filepath)
	if err != nil {
		return Texture{}, fmt.Errorf("texture %q not found on disk: %v", filepath, err)
	}
	defer imgFile.Close()
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return Texture{}, err
	}

	rgba := image.NewNRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return Texture{}, fmt.Errorf("Unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	rgba = imaging.FlipV(rgba)

	texture.width = int32(rgba.Rect.Size().X)
	texture.height = int32(rgba.Rect.Size().Y)

	gl.GenTextures(1, &texture.rendererID)
	texture.Bind(0, false)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, texture.width, texture.height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

	texture.Unbind()

	return texture, nil
}

func CreateEmptyTexture(width, height int32) Texture {
	texture := Texture{}
	texture.width = width
	texture.height = height

	gl.GenTextures(1, &texture.rendererID)
	texture.Bind(0, true)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, texture.width, texture.height, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)

	return texture
}

func (texture *Texture) Bind(slot uint32, write bool) {
	gl.ActiveTexture(gl.TEXTURE0 + slot)
	gl.BindTexture(gl.TEXTURE_2D, texture.rendererID)
	if write {
		gl.BindImageTexture(0, texture.rendererID, 0, false, 0, gl.READ_WRITE, gl.RGBA32F)
	}
}

func (texture *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (texture *Texture) Delete() {
	gl.DeleteTextures(1, &texture.rendererID)
}

func (texture *Texture) GetWidth() int32 {
	return texture.width
}
func (texture *Texture) GetHeight() int32 {
	return texture.height
}
