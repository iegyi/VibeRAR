package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
)

// Bilinear scale
func resizeImage(src *image.RGBA, targetW, targetH int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	srcW := src.Bounds().Dx()
	srcH := src.Bounds().Dy()

	xRatio := float64(srcW-1) / float64(targetW)
	yRatio := float64(srcH-1) / float64(targetH)

	for y := 0; y < targetH; y++ {
		for x := 0; x < targetW; x++ {
			gx := xRatio * float64(x)
			gy := yRatio * float64(y)
			gxi := int(gx)
			gyi := int(gy)
			xDiff := gx - float64(gxi)
			yDiff := gy - float64(gyi)

			c00 := src.RGBAAt(gxi, gyi)
			c10 := src.RGBAAt(gxi+1, gyi)
			c01 := src.RGBAAt(gxi, gyi+1)
			c11 := src.RGBAAt(gxi+1, gyi+1)

			blend := func(a, b, c, d uint8) uint8 {
				v := float64(a)*(1-xDiff)*(1-yDiff) +
					float64(b)*xDiff*(1-yDiff) +
					float64(c)*(1-xDiff)*yDiff +
					float64(d)*xDiff*yDiff
				if v < 0 {
					return 0
				}
				if v > 255 {
					return 255
				}
				return uint8(v + 0.5)
			}

			dst.SetRGBA(x, y, color.RGBA{
				R: blend(c00.R, c10.R, c01.R, c11.R),
				G: blend(c00.G, c10.G, c01.G, c11.G),
				B: blend(c00.B, c10.B, c01.B, c11.B),
				A: blend(c00.A, c10.A, c01.A, c11.A),
			})
		}
	}
	return dst
}

func main() {
	f, err := os.Open("Gemini_Generated_Image_f6xonif6xonif6xo.jpg")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer f.Close()

	src, _, err := image.Decode(f)
	if err != nil {
		fmt.Printf("Error decode: %v\n", err)
		return
	}

	w, h := src.Bounds().Dx(), src.Bounds().Dy()
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			rgba.Set(x, y, src.At(x, y))
		}
	}

	// BFS Flood Fill from outer edges
	visited := make([]bool, w*h)
	queue := make([][2]int, 0, w*h/4)

	// Outer border seeds
	for x := 0; x < w; x++ {
		queue = append(queue, [2]int{x, 0}, [2]int{x, h - 1})
		visited[x] = true
		visited[(h-1)*w+x] = true
	}
	for y := 0; y < h; y++ {
		queue = append(queue, [2]int{0, y}, [2]int{w - 1, y})
		visited[y*w] = true
		visited[y*w+(w-1)] = true
	}

	isBg := func(x, y int) bool {
		c := rgba.RGBAAt(x, y)
		lum := 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
		// Background or ground shadow (grayish/bluish-white)
		// Dark rim of the icon is < 75 luminance
		if lum > 140 {
			return true
		}
		// Shadow on bottom left
		if y > 1600 && x < 750 && lum > 35 {
			return true
		}
		if y > 1700 && lum > 45 {
			return true
		}
		return false
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		cx, cy := curr[0], curr[1]

		c := rgba.RGBAAt(cx, cy)
		lum := 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)

		if lum > 150 {
			rgba.SetRGBA(cx, cy, color.RGBA{0, 0, 0, 0})
		} else if lum > 90 {
			// Smooth edge falloff
			factor := (lum - 90) / (150 - 90)
			alpha := uint8((1.0 - factor) * 255.0)
			rgba.SetRGBA(cx, cy, color.RGBA{c.R, c.G, c.B, alpha})
		}

		neighbors := [][2]int{
			{cx + 1, cy}, {cx - 1, cy}, {cx, cy + 1}, {cx, cy - 1},
		}
		for _, nb := range neighbors {
			nx, ny := nb[0], nb[1]
			if nx >= 0 && nx < w && ny >= 0 && ny < h {
				idx := ny*w + nx
				if !visited[idx] && isBg(nx, ny) {
					visited[idx] = true
					queue = append(queue, [2]int{nx, ny})
				}
			}
		}
	}

	// Find bounding box
	minX, minY := w, h
	maxX, maxY := 0, 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if rgba.RGBAAt(x, y).A > 20 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	// Make square crop centered on content with 5% margin
	contentW := maxX - minX
	contentH := maxY - minY
	size := int(math.Max(float64(contentW), float64(contentH)) * 1.06)
	centerX := (minX + maxX) / 2
	centerY := (minY + maxY) / 2

	cropX := centerX - size/2
	cropY := centerY - size/2

	cropped := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			sx := cropX + x
			sy := cropY + y
			if sx >= 0 && sx < w && sy >= 0 && sy < h {
				cropped.Set(x, y, rgba.At(sx, sy))
			} else {
				cropped.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}

	// Save 512x512 transparent PNG
	png512 := resizeImage(cropped, 512, 512)
	outPng, _ := os.Create("viberar_icon.png")
	_ = png.Encode(outPng, png512)
	outPng.Close()
	fmt.Println("Saved transparent viberar_icon.png")

	// Create multi-resolution ICO file
	// Sizes: 256, 128, 64, 48, 32, 24, 16
	sizes := []int{16, 24, 32, 48, 64, 128, 256}
	type icoEntry struct {
		size int
		data []byte
	}
	var entries []icoEntry

	for _, sz := range sizes {
		resized := resizeImage(cropped, sz, sz)
		var buf bytes.Buffer
		_ = png.Encode(&buf, resized)
		entries = append(entries, icoEntry{size: sz, data: buf.Bytes()})
	}

	icoFile, err := os.Create("appicon.ico")
	if err != nil {
		fmt.Printf("Error creating ico: %v\n", err)
		return
	}
	defer icoFile.Close()

	// ICO Header: Reserved (0), Type (1), Count
	_ = binary.Write(icoFile, binary.LittleEndian, uint16(0))
	_ = binary.Write(icoFile, binary.LittleEndian, uint16(1))
	_ = binary.Write(icoFile, binary.LittleEndian, uint16(len(entries)))

	offset := uint32(6 + 16*len(entries))

	for _, e := range entries {
		var wByte, hByte byte
		if e.size >= 256 {
			wByte = 0
			hByte = 0
		} else {
			wByte = byte(e.size)
			hByte = byte(e.size)
		}
		_ = binary.Write(icoFile, binary.LittleEndian, wByte)
		_ = binary.Write(icoFile, binary.LittleEndian, hByte)
		_ = binary.Write(icoFile, binary.LittleEndian, byte(0)) // Color count
		_ = binary.Write(icoFile, binary.LittleEndian, byte(0)) // Reserved
		_ = binary.Write(icoFile, binary.LittleEndian, uint16(1)) // Planes
		_ = binary.Write(icoFile, binary.LittleEndian, uint16(32)) // BitCount
		_ = binary.Write(icoFile, binary.LittleEndian, uint32(len(e.data)))
		_ = binary.Write(icoFile, binary.LittleEndian, offset)
		offset += uint32(len(e.data))
	}

	for _, e := range entries {
		_, _ = icoFile.Write(e.data)
	}

	fmt.Println("Generated transparent multi-resolution appicon.ico successfully!")
}
