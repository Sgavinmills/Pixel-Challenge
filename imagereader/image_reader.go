package imagereader

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"strings"
)

func GetImageBytes(sourceFile string) ([]byte, error) {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return []byte{}, err
	}
	// fmt.Println(input)
	return input, nil
}

func convertBytesToPixels(imageBytes []byte) [][]byte {

	var pixels [][]byte

	for i := 0; i < len(imageBytes); i +=3 {
		singlePixel := []byte{
			imageBytes[i],
			imageBytes[i+1],
			imageBytes[i+2],
		}
		pixels = append(pixels, singlePixel)
	}
	return pixels
}

func CalculateClosenessFromRawBytes(imageBytes1, imageBytes2 []byte) float64 {

	var numberOfMatchingPixels int
	// TODO maybe check for byte len, if not equal return 0? (ie pictures not same size)
	// reject if len not divisible by 3? (not legit image in that case)
	for i := 0; i < len(imageBytes1); i += 3 {
		// fmt.Println(imageBytes2[i])
		if imageBytes1[i] == imageBytes2[i] &&
		   imageBytes1[i+1] == imageBytes2[i+1] && 
		   imageBytes1[i+2] == imageBytes2[i+2] {
			numberOfMatchingPixels++
		}
	}

	closeness := float64(numberOfMatchingPixels) / float64(len(imageBytes1) / 3)
	return closeness

}

func calculateClosenessFromPixels(pixels1, pixels2 [][]byte) float64 {
	var numberOfMatchingPixels int

	for i := 0; i < len(pixels1); i++ {
		if bytes.Equal(pixels1[i], pixels2[i]) {
			numberOfMatchingPixels++
		}
	}

	closeness := float64(numberOfMatchingPixels) / float64(len(pixels1))

	return closeness
}

type Image struct {
	ImageBytes []byte
	Name string
	Closeness float64
}

// make a go routine version of this. 
func ReadFilesIntoBytes(files []fs.FileInfo, referenceImage string, directory string) []Image {
	var comparisonImages []Image
	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".raw") && fileName != referenceImage {
			imageBytes, err := GetImageBytes(directory + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			image := Image{ImageBytes: imageBytes, Name: file.Name()}
			comparisonImages = append(comparisonImages, image)
		}
	}
	return comparisonImages
}

func CalculateClosenessForAllImages(referenceImage []byte, comparisonImages []Image, counter chan struct{}) {
	for i, v := range comparisonImages {
		go func(i int, v Image){
			closeness := CalculateClosenessFromRawBytes(referenceImage, v.ImageBytes)
			comparisonImages[i].Closeness = closeness
			fmt.Println("Calculated a closeness")
			counter <- struct{}{}

		}(i, v)
	}
}


func StartResultsWatcher(counter *Counter) {

	go func(){
		fmt.Println("results watcher")
		for cmd := range counter.commandChan {
			fmt.Println("received a command")
			switch cmd {
			case "hello":
				fmt.Println("wor")

			}
		}
	}()

}

type Command struct {
	ty string
	replyChan chan Reply
}

type Reply struct {
	res string
}

type Counter struct {
	currentTotal int
	expectedTotal int
	commandChan chan string
}


func NewCounter(totalNumberOfImagesToCompare int) *Counter {
	return &Counter{expectedTotal : totalNumberOfImagesToCompare, commandChan : make(chan string)}
}