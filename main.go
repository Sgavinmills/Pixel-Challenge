package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type imageInfo struct {
	name        string
	imageBytes  []byte
	imagePixels [][]byte
	closeness   float64
}


type imageStorage struct {
	comparisonImages map[string]imageInfo
	referenceImage imageInfo
}

func main() {

	var imageStorage imageStorage
	imageStorage.comparisonImages = make(map[string]imageInfo)
	// TODO, get suffix from the reference image instead of hard coding .raw
	referenceImage := "0f0e5c84-3b99-4874-bb8c-e0228155c4b5.raw"
	directory := "Images\\Bronze\\"

	imageStorage.getFilesFromDirectory(directory, referenceImage)

	// take reference image and image directory as parameters

	// convert reference image to bytes
	// convert all the other images that aren't the reference images to bytes, store separately

	// loop over the other images calculating closeness score to reference image

	// return the images with the 3 highest closeness scores
}

func (store imageStorage) getFilesFromDirectory(directory string, referenceImage string) error {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".raw") && fileName != referenceImage {
			newImage := imageInfo{name: fileName}
			store.comparisonImages[fileName] = newImage
		}
	}

	return nil
}