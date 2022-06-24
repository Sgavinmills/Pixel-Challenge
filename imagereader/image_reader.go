package imagereader

import (
	"bytes"
	"io/ioutil"
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

