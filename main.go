package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sprites_folder, err := os.ReadDir(filepath.Join(pwd, "sprites"))
	if err != nil {
		log.Fatal(err)
	}

	// sprite_height := 0
	// sprite_width := 0

	var sprites1 []image.Image
	var sprites2 []image.Image
	// var spritesd [][]image.Image
	var sprites []image.Image

	// chunked_images := chunk_images(sprites_folder, 2)
	// for _, chunk_image := range chunked_images {
	// 	wg.Add(1)
	// 	go func(chunk_image []fs.DirEntry) {
	// 		temp := decode_images(chunk_image, pwd, &wg)
	// 		spritesd = append(spritesd, temp)
	// 	}(chunk_image)

	// }

	// for _, sprite := range spritesd {
	// 	sprites = append(sprites, sprite...)

	// }

	mid := len(sprites_folder) / 2
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		sprites1 = decode_images(sprites_folder[:mid], pwd, &wg)
	}()
	wg.Add(1)
	go func() {
		sprites2 = decode_images(sprites_folder[mid:], pwd, &wg)
	}()
	wg.Wait()

	fmt.Println("sp1", sprites1)
	sprites = append(sprites, sprites1...)
	sprites = append(sprites, sprites2...)
	// fmt.Println(sprites)

	hframes := 8
	vframes := 3
	spritesheet_height := 128 * hframes
	spritesheet_width := 128 * vframes
	spritesheet := image.NewRGBA(image.Rect(0, 0, spritesheet_height, spritesheet_width))
	// transparent_background := color.RGBA{0, 0, 0, 0}
	draw.Draw(spritesheet, spritesheet.Bounds(), spritesheet, image.Point{}, draw.Src)

	// count_vertical_frames := -1
	// count_horizontal_frames := 0
	// var wg1 sync.WaitGroup

	// for i := 1; i < vframes+1; i++ {
	// for i, sprite_image := range hframes {
	// 	// for j := 0; j < hframes; j++ {
	// 	// 	if i*j < len(sprites_folder) {
	// 	// 		fmt.Println(i*j, "test")
	// 	// 		sprite_image := sprites[j*i]
	// 	// 		bounds := sprite_image.Bounds()
	// 	// 		w := bounds.Dx()
	// 	// 		h := bounds.Dy()
	// 	// 		fmt.Println(j+i, w, h)
	// 	// 		draw.Draw(spritesheet, image.Rect(j*w, (i-1)*h, w*8, h*3), sprite_image, image.Point{}, draw.Over)
	// 	// 	}
	// 	// }
	// 	// }
	// }
	sprites_chunked := chunkSlice(sprites, hframes)

	var wg1 sync.WaitGroup
	for count_vertical_frames, sprite_chunk := range sprites_chunked {
		wg1.Add(1)
		fmt.Println(count_vertical_frames)
		go func(count_vertical_frames int, sprite_chunk []image.Image) {
			defer wg1.Done()
			paint_spritesheet(sprite_chunk, hframes, count_vertical_frames, spritesheet)
		}(count_vertical_frames, sprite_chunk)
		// for sprite_index, sprite := range sprite_chunk {
		// }
	}
	wg1.Wait()

	// sprites_mid = len(sprites)
	// paint_spritesheet(sprites, hframes, count_vertical_frames, count_horizontal_frames, spritesheet)

	f, err := os.Create("spritesheet.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err = png.Encode(f, spritesheet); err != nil {
		log.Printf("failed to encode: %v", err)
	}
	// fmt.Println(len(sprites_folder))
	fmt.Println(time.Since(start))

}

func decode_images(sprites_folder []fs.DirEntry, pwd string, wg *sync.WaitGroup) []image.Image {
	defer wg.Done()
	var sprites_array []image.Image
	for _, sprite := range sprites_folder {
		// fmt.Println(sprite)
		if reader, err := os.Open(filepath.Join(pwd, "sprites", sprite.Name())); err == nil {
			defer reader.Close()
			m, _, err := image.Decode(reader)
			if err != nil {
				log.Fatal(err)
			}
			sprites_array = append(sprites_array, m)
			// if i == 0 {
			// 	bounds := m.Bounds()
			// 	w := bounds.Dx()
			// 	h := bounds.Dy()
			// 	sprite_height = h
			// 	sprite_width = w
			// }
			// fmt.Println(sprites_array, "arr")
		}
	}

	// fmt.Println(sprites_array, "tee")
	return sprites_array
}

func paint_spritesheet(sprites []image.Image, hframes int, count_vertical_frames int, spritesheet draw.Image) {
	for count_horizontal_frames, sprite_image := range sprites {
		bounds := sprite_image.Bounds()
		w := bounds.Dx()
		h := bounds.Dy()
		// if i%hframes == 0 {
		// 	count_vertical_frames++
		// }
		// if i%hframes != 0 {
		// 	count_horizontal_frames += 1
		// } else {
		// 	count_horizontal_frames = 0
		// }
		draw.Draw(spritesheet, image.Rect(count_horizontal_frames*h, count_vertical_frames*w, w*8, h*3), sprite_image, image.Point{}, draw.Over)
	}

}

func chunkSlice(slice []image.Image, chunkSize int) [][]image.Image {
	var chunks [][]image.Image
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
