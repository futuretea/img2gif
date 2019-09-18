package main

import (
	"flag"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

var (
	path   string
	name   string
	output string
	delay  int
	bound  int
	anim   *gif.GIF
)

func getFiles(path string, name string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalf("Could not open dir %s. Error: %s\n", path, err)
	}
	numArr := []int{}
	for _, info := range files {
		numstr := strings.TrimSuffix(strings.TrimPrefix(info.Name(), name), ".png")
		num, err := strconv.Atoi(numstr)
		if err != nil {
			log.Fatalf("Could not read image name %s. Error: %s\n", info.Name(), err)
		}
		numArr = append(numArr, num)
	}
	sort.Ints(numArr)
	sortedFiles := []string{}
	for _, num := range numArr {
		sortedFile := fmt.Sprintf("%s%d.png", name, num)
		sortedFiles = append(sortedFiles, sortedFile)
	}
	return sortedFiles
}

func getBounds(files []string, bound int) image.Rectangle {
	return decodeImage(files[bound-1]).Bounds()
}

func decodeImage(file string) image.Image {
	sysType := runtime.GOOS
	fpath := ""
	if sysType == "linux" {
		fpath = path + "/" + file
	}
	if sysType == "windows" {
		fpath = path + "\\" + file
	}
	f, err := os.Open(fpath)
	if err != nil {
		log.Fatalf("Could not open file %s. Error: %s\n", file, err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		log.Fatalf("Could not Decode file %s. Error: %s\n", f.Name(), err)
	}
	return img
}

func addImage(img image.Image, bounds image.Rectangle) {
	//bounds := img.Bounds()
	paletted := image.NewPaletted(bounds, palette.Plan9)
	draw.FloydSteinberg.Draw(paletted, bounds, img, image.ZP)
	anim.Image = append(anim.Image, paletted)
	anim.Delay = append(anim.Delay, delay*15)
}

func outputGif(output string) {
	f, _ := os.Create(output)
	defer f.Close()
	err := gif.EncodeAll(f, anim)
	if err != nil {
		log.Fatalf("Could not Encode gif %s. Error: %s\n", output, err)
	}
}

func init() {
	flag.StringVar(&path, "p", "", "png图片文件夹路径")
	flag.StringVar(&name, "n", "image", "png图片文件前缀")
	flag.StringVar(&output, "o", "output.gif", "生成gif的文件名")
	flag.IntVar(&delay, "d", 8, "每张图片的展示时间*15毫秒")
	flag.IntVar(&bound, "b", 1, "gif边界参照")
	flag.Parse()

	if path == "" {
		fmt.Println("请输入图片路径")
		flag.PrintDefaults()
		os.Exit(1)
	}
	anim = new(gif.GIF)
}

func main() {
	files := getFiles(path, name)
	bounds := getBounds(files, bound)
	fmt.Printf("gif边界: %d -> %d\n", bounds.Min, bounds.Max)
	for _, file := range files {
		img := decodeImage(file)
		addImage(img, bounds)
		fmt.Println(file)
	}
	outputGif(output)
}