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
	anim   *gif.GIF
)

func main() {
	files := getFiles(path, name)
	for _, file := range files {
		img := decodeImage(file)
		addImage(img)
	}
	outputGif(output)
}

func init() {
	flag.StringVar(&path, "p", "", "png图片文件夹路径")
	flag.StringVar(&name, "n", "image", "png图片文件前缀")
	flag.StringVar(&output, "o", "output.gif", "生成gif的文件名")
	flag.IntVar(&delay, "d", 4, "每张图片的展示时间*15毫秒")
	flag.Parse()

	if path == "" {
		fmt.Println("请输入图片路径")
		flag.PrintDefaults()
		os.Exit(1)
	}
	anim = new(gif.GIF)
}

func getFiles(path string, name string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	numArr := []int{}
	for _, info := range files {
		numstr :=  strings.TrimSuffix(strings.TrimPrefix(info.Name(), name), ".png")
		num, err := strconv.Atoi(numstr)
		if err != nil {
			log.Fatal(err)
		}
		numArr = append(numArr, num)
	}
	sort.Ints(numArr)
	sortedFiles := []string{}
	for _, num := range numArr{
		sortedFile := fmt.Sprintf("%s%d.png", name, num)
		sortedFiles = append(sortedFiles, sortedFile)
	}
	return sortedFiles
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
	fmt.Println(f.Name())
	if err != nil {
		log.Fatalf("Could not open file %s. Error: %s\n", file, err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		log.Fatalf("Decode Error: %s\n", err)
	}
	return img
}

func addImage(img image.Image) {
	paletted := image.NewPaletted(img.Bounds(), palette.Plan9)
	draw.FloydSteinberg.Draw(paletted, img.Bounds(), img, image.ZP)
	anim.Image = append(anim.Image, paletted)
	anim.Delay = append(anim.Delay, delay*15)
}

func outputGif(output string) {
	f, _ := os.Create(output)
	defer f.Close()
	gif.EncodeAll(f, anim)
}
