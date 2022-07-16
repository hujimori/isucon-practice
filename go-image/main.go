package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	filename := "IMG_0902.JPG"
	imgfile := filepath.Join("/home/hujimori/isucon-prcatice/go-image", fmt.Sprintf("%s", filename))
	fmt.Print(imgfile)
	// f, err := os.Open(imgfile)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer f.Close()

	// img, err := io.ReadAll(f)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("%d", len(img))

	// dstfile := "/home/hujimori/isucon-prcatice/go-image/IMG_0904.JPG"
	// dst, err := os.Create(dstfile)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer dst.Close()

	// _, err = dst.Write(img)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
