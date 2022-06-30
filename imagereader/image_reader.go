package imagereader

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

type Image struct {
	Name string
	Closeness float64
	ImageBytes []byte
}

type Reference struct {
	RefImageName string
	RefImageBytes []byte
	Directory string
}

func FindClosestImagesToReference(files []fs.FileInfo, referenceInfo Reference) []Image {
	
	processedImage := make(chan Image)

	calculateCloseness(files, referenceInfo, processedImage)
	var results []Image
	maxNumberOfResults := 3
	for i := 0; i < len(files)-1; i++ {
		processedImage := <- processedImage
		results = updateLeaderBoard(results, processedImage, maxNumberOfResults)
	}

	return results
}

func calculateCloseness(files []fs.FileInfo, referenceInfo Reference, processedImage chan Image) {
	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".raw") && fileName != referenceInfo.RefImageName {
			go func(file, directory string){
				imageBytes, err := GetImageBytes(directory + file)
				if err != nil {
					log.Println(err)
				}
				if len(imageBytes) == 3145728 {
					// numberOfSlices := 512
					// closeness, err := CalculateClosenessFromRawBytesWithGoRoutines(referenceInfo.RefImageBytes, imageBytes, numberOfSlices)
					closeness, err := CalculateClosenessFromRawBytes(referenceInfo.RefImageBytes, imageBytes)
					if err != nil {
						log.Println(err.Error())
						os.Exit(1)
					}
					processedImage <- Image{Name: file, Closeness: closeness}

				} else {
					processedImage <- Image{Name: file, Closeness: 0}
				}
				
			}(fileName, referenceInfo.Directory)
		}
	}
}

// this is dumb. just check if new closeness is bigger than smallest in eaderboard and
//if it is then just replace and resort the leaderboard. 
func updateLeaderBoard(comparedImages []Image, processedImage Image, maxNumberOfResults int) []Image {

	if len(comparedImages) >= maxNumberOfResults && processedImage.Closeness < comparedImages[len(comparedImages)-1].Closeness {
		return comparedImages
	}

	if len(comparedImages) == 0 {
		comparedImages = append(comparedImages, processedImage)
		return comparedImages
	}

	for i, v := range comparedImages {
		if processedImage.Closeness > v.Closeness {
			newResults := append(comparedImages[:i+1], comparedImages[i:]...)
			newResults[i] = processedImage
			newResults = newResults[:len(newResults)-1]
			return newResults
		}

		if i == len(comparedImages) - 1 && len(comparedImages) < maxNumberOfResults {
			comparedImages = append(comparedImages, processedImage)
			return comparedImages
		}
	}

	return comparedImages

}

func GetImageBytes(sourceFile string) ([]byte, error) {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}
var ErrInvalidLength = errors.New("Invalid length of imageBytes type []byte. All imageBytes must be length multiple of three")
var ErrInvalidSubsliceValue = errors.New("Invalid number of subslices requested. Subslices must contain multiples of 3 bytes")

func CalculateClosenessFromRawBytes(imageBytes1, imageBytes2 []byte) (float64, error) {

	
	if len(imageBytes1) % 3 != 0 || len(imageBytes2) %3 != 0 {
		return 0, ErrInvalidLength
	}
	if len(imageBytes1) != len(imageBytes2) {
		return 0, nil
	}
	
	
	var numberOfMatchingPixels int
		for i := 0; i < len(imageBytes1); i += 3 {
			if imageBytes1[i] == imageBytes2[i] &&
			   imageBytes1[i+1] == imageBytes2[i+1] && 
			   imageBytes1[i+2] == imageBytes2[i+2] {
				numberOfMatchingPixels++
			}
		}
	

	closeness := float64(numberOfMatchingPixels) / float64(len(imageBytes1) / 3)
	return closeness, nil

}

func CalculateMatchingPixelsFromRawBytesOnChannel(imageBytes1, imageBytes2 []byte, closenessCounter chan closenessResponse) {
	
	
	// issue seems to be imageBytes length not dividing by 3 :0

	var numberOfMatchingPixels int
	// TODO maybe check for byte len, if not equal return 0? (ie pictures not same size)
	// reject if len not divisible by 3? (not legit image in that case)

	if len(imageBytes1) % 3 != 0 || len(imageBytes2) %3 != 0 {
		response := closenessResponse{err : ErrInvalidLength}
		closenessCounter <- response
	}
	if len(imageBytes1) != len(imageBytes2) {

		response := closenessResponse{closeness : 0}
		closenessCounter <- response
	} else {

		for i := 0; i < len(imageBytes1); i += 3 {
			if imageBytes1[i] == imageBytes2[i] &&
			   imageBytes1[i+1] == imageBytes2[i+1] && 
			   imageBytes1[i+2] == imageBytes2[i+2] {
				numberOfMatchingPixels++
			}
		}
		response := closenessResponse{closeness : numberOfMatchingPixels}
		
		closenessCounter <- response

	}	


}
type closenessResponse struct {
	closeness int
	err error
}
func CalculateClosenessFromRawBytesWithGoRoutines(imageBytes1, imageBytes2 []byte, numberOfSlices int) (float64, error) {

	// var numberOfMatchingPixels int
	// TODO maybe check for byte len, if not equal return 0? (ie pictures not same size)
	// reject if len not divisible by 3? (not legit image in that case)

	
	// gunna split into 3 parts (to start)
	// so divide len by 3 and round up! thats how many items are included in each loop
	// so loop 5 times
	// each loop take the whole number of len / 3 rounded up and calc closeness on that slice
	// for the last loop just do the rest

	closenessCounterChan := make(chan closenessResponse)
	// numberOfSlices := 1024
	// itemsPerSlice := math.Ceil(float64(len(imageBytes1) / numberOfSlices))
	// for i := 0; i < numberOfSlices; i++ {
	// 	startOfSlice := i * int(itemsPerSlice)
	// 	tmpSlice1 := imageBytes1[startOfSlice:startOfSlice + int(itemsPerSlice)]
	// 	tmpSlice2 := imageBytes2[startOfSlice:startOfSlice + int(itemsPerSlice)]
	// 	go func(slice1, slice2 []byte){
	// 		CalculateMatchingPixelsFromRawBytesOnChannel(slice1, slice2, closenessCounterChan)
	// 	}(tmpSlice1, tmpSlice2)

	// }
	newSlices1, newSlices2, err := createSubSlices(imageBytes1, imageBytes2, numberOfSlices)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i := 0; i < numberOfSlices; i++ {
		go func(i int){
			CalculateMatchingPixelsFromRawBytesOnChannel(newSlices1[i], newSlices2[i], closenessCounterChan)
		}(i)
	}
	totalMatchingPixels := 0
	for i := 0; i < numberOfSlices; i++ {
		returnedCloseness := <-closenessCounterChan
		if returnedCloseness.err != nil {
			return 0, returnedCloseness.err
		}
		totalMatchingPixels = totalMatchingPixels + returnedCloseness.closeness
	}
	closeness := float64(totalMatchingPixels) / float64(len(imageBytes1) / 3)

	return closeness, nil
}

func createSubSlices(slice1, slice2 []byte, numberOfSlices int) ([][]byte, [][]byte, error) {
	// len(slice) / numberOfSlices must give an answer that divides by 3
	// fmt.Println(len(slice1))
	lenAsFloat := float64(len(slice1))
	noOfSlicesAsFloat := float64(numberOfSlices)
	divided := lenAsFloat / noOfSlicesAsFloat

	// checks if its a whole number
	if divided != float64(len(slice1) / numberOfSlices) {
		return nil, nil, ErrInvalidSubsliceValue

	}
	// fmt.Println(len(slice1) / numberOfSlices)
	if (len(slice1) / numberOfSlices) % 3 != 0 {
		return nil, nil, ErrInvalidSubsliceValue
	}
	// also error if they arent the same lengths etc
	var newSlices1 [][]byte
	var newSlices2 [][]byte
	itemsPerSlice := len(slice1) / numberOfSlices
	for i := 0; i < numberOfSlices; i++ {
		startOfSlice := i * itemsPerSlice
		tmpSlice1 := slice1[startOfSlice:startOfSlice + itemsPerSlice]
		tmpSlice2 := slice2[startOfSlice:startOfSlice + itemsPerSlice]
		newSlices1 = append(newSlices1, tmpSlice1)
		newSlices2 = append(newSlices2, tmpSlice2)
	}

	return newSlices1, newSlices2, nil
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



// func ReadFilesIntoBytes(files []fs.FileInfo, referenceImage string, directory string) []Image {
// 	var comparisonImages []Image
// 	for _, file := range files {
// 		fileName := file.Name()
// 		if strings.HasSuffix(fileName, ".raw") && fileName != referenceImage {
// 			imageBytes, err := GetImageBytes(directory + file.Name())
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 			image := Image{ImageBytes: imageBytes, Name: file.Name()}
// 			comparisonImages = append(comparisonImages, image)
// 		}
// 	}
// 	return comparisonImages
// }

func CalculateClosenessForAllImages(referenceImage []byte, comparisonImages []Image) {
	for i, v := range comparisonImages {
		go func(i int, v Image){
			closeness, _ := CalculateClosenessFromRawBytes(referenceImage, v.ImageBytes)
			comparisonImages[i].Closeness = closeness
		}(i, v)
	}
}



func ProcessImageFiles(files []fs.FileInfo, referenceInfo Reference) []Image {

	comparedImages := scanAndCompareImages(files, referenceInfo)
	sortImagesByCloseness(comparedImages)
	

	return comparedImages
}



func scanAndCompareImages(files []fs.FileInfo, referenceInfo Reference) []Image {
	imagesProcessedCounter := make(chan Image)

	calculateCloseness(files, referenceInfo, imagesProcessedCounter)
	var comparedImages []Image
	for i := 0; i < len(files)-1; i++ {
		processedImage := <- imagesProcessedCounter
		comparedImages = append(comparedImages, processedImage)
	}

	return comparedImages
}

func sortImagesByCloseness(comparedImages []Image) {
	sort.Slice(comparedImages, func(i, j int) bool {
		return comparedImages[i].Closeness > comparedImages[j].Closeness
	})
}


