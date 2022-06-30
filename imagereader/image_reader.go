package imagereader

import (
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

type closenessResponse struct {
	closeness int
	err error
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

		if  fileName != referenceInfo.RefImageName {
			if !strings.HasSuffix(fileName, ".raw") {
				processedImage <- Image{Name: "invalid file", Closeness: 0}
			}
			
			go func(file, directory string){
				imageBytes, err := GetImageBytes(directory + file)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}

				imagesize1024x1024 := 3145728
				if len(imageBytes) != imagesize1024x1024 {
					processedImage <- Image{Name: "image wrong size", Closeness: 0}
					return
				}

				numberOfSlices := 16
				closeness, err := CalculateClosenessFromRawBytesWithGoRoutines(referenceInfo.RefImageBytes, imageBytes, numberOfSlices)
				// closeness, err := CalculateClosenessFromRawBytes(referenceInfo.RefImageBytes, imageBytes)
				if err != nil {
					log.Println(err.Error())
					os.Exit(1)
				}
				processedImage <- Image{Name: file, Closeness: closeness}

				
			}(fileName, referenceInfo.Directory)
		}
	}
}



func updateLeaderBoard(comparedImages []Image, processedImage Image, maxNumberOfResults int) []Image {

	if len(comparedImages) < maxNumberOfResults {
		comparedImages = append(comparedImages, processedImage)
		sort.Slice(comparedImages, func(i, j int) bool {
			return comparedImages[i].Closeness > comparedImages[j].Closeness
		})

		return comparedImages
	} 

	if comparedImages[maxNumberOfResults-1].Closeness < processedImage.Closeness {
		comparedImages[maxNumberOfResults-1] = processedImage
		sort.Slice(comparedImages, func(i, j int) bool {
			return comparedImages[i].Closeness > comparedImages[j].Closeness
		})
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

func CalculateClosenessFromRawBytes(imageBytes1, imageBytes2 []byte) (float64, error) {

	
	if len(imageBytes1) % 3 != 0 {
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

func CalculateClosenessFromRawBytesWithGoRoutines(imageBytes1, imageBytes2 []byte, numberOfSlices int) (float64, error) {

	subSlices1, subSlices2, err := createSubSlices(imageBytes1, imageBytes2, numberOfSlices)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	closenessCounter := make(chan closenessResponse)
	for i := 0; i < numberOfSlices; i++ {
		go func(i int){
			CalculatePixelClosenessFromSubSlices(subSlices1[i], subSlices2[i], closenessCounter)
		}(i)
	}

	totalMatchingPixels := 0
	for i := 0; i < numberOfSlices; i++ {
		returnedCloseness := <-closenessCounter
		if returnedCloseness.err != nil {
			return 0, returnedCloseness.err
		}
		totalMatchingPixels = totalMatchingPixels + returnedCloseness.closeness
	}

	closeness := float64(totalMatchingPixels) / float64(len(imageBytes1) / 3)

	return closeness, nil
}


func CalculatePixelClosenessFromSubSlices(imageBytes1, imageBytes2 []byte, closenessCounter chan closenessResponse) {

	var numberOfMatchingPixels int

	if len(imageBytes1) % 3 != 0 || len(imageBytes2) %3 != 0 {
		response := closenessResponse{err : ErrInvalidLength}
		closenessCounter <- response
		return
	}

	if len(imageBytes1) != len(imageBytes2) {
		response := closenessResponse{closeness : 0}
		closenessCounter <- response
		return
	} 

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


func createSubSlices(slice1, slice2 []byte, numberOfSlices int) ([][]byte, [][]byte, error) {

	lenAsFloat := float64(len(slice1))
	noOfSlicesAsFloat := float64(numberOfSlices)
	divided := lenAsFloat / noOfSlicesAsFloat

	// checks if its a whole number
	if divided != float64(len(slice1) / numberOfSlices) {
		return nil, nil, ErrInvalidSubsliceValue
	}

	if (len(slice1) / numberOfSlices) % 3 != 0 {
		return nil, nil, ErrInvalidSubsliceValue
	}

	var subSlices1 [][]byte
	var subSlices2 [][]byte
	itemsPerSlice := len(slice1) / numberOfSlices
	for i := 0; i < numberOfSlices; i++ {
		startOfSlice := i * itemsPerSlice
		subSlices1 = append(subSlices1, slice1[startOfSlice : startOfSlice + itemsPerSlice])
		subSlices2 = append(subSlices2, slice2[startOfSlice : startOfSlice + itemsPerSlice])
	}

	return subSlices1, subSlices2, nil
}
