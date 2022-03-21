package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

var asciiBrightnessTable string = `$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\|()1{}[]?-_+~<>i!lI;:,"^'.`
var asciiBrightnessTableLength = 68
var symbolsLength int
var inputFile string
var outputFile string
var cliOutput bool
var rescaleMult int

func main() {
	flag.StringVar(&inputFile, "i", "", "defines image path")
	flag.StringVar(&outputFile, "o", "ascii.txt", "defines output file path")
	flag.IntVar(&symbolsLength, "s", 119, "defines length of one row in symbols")
	flag.BoolVar(&cliOutput, "cli", false, "sets output to cli, will ignore row length change and output file name")
	flag.IntVar(&rescaleMult, "r", 10, "bigger - better quality of output")
	flag.Parse()
	if rescaleMult < 1 {
		fmt.Print("Rescale multiplier must be 1 or bigger")
		os.Exit(1)
	}
	if inputFile == "" {
		fmt.Print("Input image path must be specified")
		os.Exit(1)
	}
	if cliOutput {
		symbolsLength = 119
	}
	im, err := openImage(inputFile)
	if err != nil {
		fmt.Printf("Can't open input %s", err)
		os.Exit(1)
	}
	im = rescaleImage(im)
	grey := discolorImage(im)
	a := generateAsciiString(grey)
	if cliOutput {
		writeAsciiStringCli(a)
	} else {
		writeAsciiStringTxt(outputFile, a)
	}
}

func getMedian(x, y int, g *image.Gray) int {
	res := 0
	cnt := 0
	for i := x; i < x+4; i++ {
		for j := y; j < y+10; j++ {
			res += int(g.GrayAt(i, j).Y)
			cnt++
		}
		cnt++
	}
	return res / cnt
}

func openImage(name string) (image.Image, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func rescaleImage(i image.Image) image.Image {
	x1 := uint(symbolsLength) * uint(rescaleMult)
	i = resize.Resize(x1*2, x1*uint(i.Bounds().Dy())/uint(i.Bounds().Dx()), i, resize.Lanczos3)
	return i
}

func discolorImage(i image.Image) *image.Gray {
	gray := image.NewGray(i.Bounds())
	for x := 0; x < i.Bounds().Dx(); x++ {
		for y := 0; y < i.Bounds().Dy(); y++ {
			gray.Set(x, y, i.At(x, y))
		}
	}
	return gray
}

func generateAsciiString(g *image.Gray) string {
	s := strings.Builder{}
	for y := 0; y < g.Bounds().Dy(); y = y + rescaleMult*2 {
		for x := 0; x < g.Bounds().Dx(); x = x + rescaleMult*2 {
			s.WriteByte(asciiBrightnessTable[(asciiBrightnessTableLength*getMedian(x, y, g))/255])
		}
		s.WriteByte('\n')
	}
	return s.String()
}

func writeAsciiStringTxt(fileName string, asciiString string) {
	wr, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("Can't open output %s", err)
		os.Exit(1)
	}
	defer wr.Close()
	wr.WriteString(asciiString)
}

func writeAsciiStringCli(asciiString string) {
	fmt.Println(asciiString)
}
