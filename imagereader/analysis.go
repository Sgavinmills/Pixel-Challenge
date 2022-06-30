package imagereader

import (
	"bytes"
)

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

func CalculateClosenessForAllImages(referenceImage []byte, comparisonImages []Image) {
	for i, v := range comparisonImages {
		go func(i int, v Image){
			closeness, _ := CalculateClosenessFromRawBytes(referenceImage, v.ImageBytes)
			comparisonImages[i].Closeness = closeness
		}(i, v)
	}
}

// func updateLeaderBoardWithoutSort(comparedImages []Image, processedImage Image, maxNumberOfResults int) []Image {

// 	if len(comparedImages) >= maxNumberOfResults && processedImage.Closeness < comparedImages[len(comparedImages)-1].Closeness {
// 		return comparedImages
// 	}

// 	if len(comparedImages) == 0 {
// 		comparedImages = append(comparedImages, processedImage)
// 		return comparedImages
// 	}

// 	for i, v := range comparedImages {
// 		if processedImage.Closeness > v.Closeness {
// 			newResults := append(comparedImages[:i+1], comparedImages[i:]...)
// 			newResults[i] = processedImage
// 			newResults = newResults[:len(newResults)-1]
// 			return newResults
// 		}

// 		if i == len(comparedImages) - 1 && len(comparedImages) < maxNumberOfResults {
// 			comparedImages = append(comparedImages, processedImage)
// 			return comparedImages
// 		}
// 	}

// 	return comparedImages

// }

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