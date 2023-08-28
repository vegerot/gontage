package main

import (
	"fmt"
	"image"
	_ "image/png"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	sprites_folder, err := os.ReadDir(filepath.Join(pwd, "sprites"))
	sprites_folder = sprites_folder[0:59]
	if err != nil {
		log.Fatal(err)
	}

	var decoded_images_chunked [][]image.Image
	var chunked_sprite_dir_entries [][]fs.DirEntry
	chunked_sprite_dir_entries = chunkSpriteDirEntries(sprites_folder, 69)

	var chunk_images_waitgroup sync.WaitGroup
	for _, chunked_sprites_entry := range chunked_sprite_dir_entries {
		chunk_images_waitgroup.Add(1)
		go func(chunked_sprites_entry []fs.DirEntry) {
			one_chunk_of_decoded_images := decodeImages(chunked_sprites_entry, pwd, &chunk_images_waitgroup)
			decoded_images_chunked = append(decoded_images_chunked, one_chunk_of_decoded_images)
		}(chunked_sprites_entry)
	}
	start := time.Now()
	chunk_images_waitgroup.Wait()
	fmt.Println(time.Since(start))
}

func decodeImages(sprites_folder []fs.DirEntry, pwd string, wg *sync.WaitGroup) []image.Image {
	defer wg.Done()
	var sprites_array []image.Image
	for _, sprite := range sprites_folder {
		if reader, err := os.Open(filepath.Join(pwd, "sprites", sprite.Name())); err == nil {
			defer reader.Close()
			m, _, err := image.Decode(reader)
			if err != nil {
				panic(err)
			}
			sprites_array = append(sprites_array, m)
		}
	}
	return sprites_array
}

func chunkSpriteDirEntries(slice []fs.DirEntry, chunkSize int) [][]fs.DirEntry {
	var chunks [][]fs.DirEntry
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
