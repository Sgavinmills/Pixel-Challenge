package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/Sgavinmills/Pixel-Challenge/main/imagereader"
)

// type image struct {
// 	imageBytes []byte
// 	name string
// 	closeness float64
// }

func main() {

	refImgPtr := flag.String("ref", "", "Name of image")
	dirPtr := flag.String("dir", "Bronze", "Which image level to use")
	flag.Parse()

	directory := "Images\\"+*dirPtr+"\\"
	referenceImage := *refImgPtr
	
	files, err := ioutil.ReadDir(directory)
	
	// if not ref image provided we take first one from directory instead
	if referenceImage == "" {
		referenceImage = files[0].Name()
	}
	
	fmt.Printf("Comparing image %v against images in %v\n", referenceImage, directory)
	
	// get referenceimage bytes
	referenceImageBytes, err := imagereader.GetImageBytes(directory+referenceImage)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	
	comparisonImages := imagereader.ReadFilesIntoBytes(files, referenceImage, directory)
	// doesnt work atm, need to wait until finished reading before moving on... but...
	// could still bench mark / test the other function to check for race issues and improvements. 

	
	
	// resultCounter := imagereader.NewCounter(10)


	testCounter := make(chan struct{})
	// imagereader.StartResultsWatcher(resultCounter)
	imagereader.CalculateClosenessForAllImages(referenceImageBytes, comparisonImages, testCounter)




	// sort comparisonImages to get winning order
	// might want to look into faster sorting algorithm OR only populate array with contenders in the first place, ie never allow more than 3 
	// basically check if there are more than 3 already stored, if so enter it in order and delete the slowest
	// time.Sleep(time.Second * 5)
	
	for i := 0; i < 9; i++ {
		<- testCounter
	}
	// print top 3. Might need to add some checking to see that there are no ties.
	sort.Slice(comparisonImages, func(i, j int) bool {
		return comparisonImages[i].Closeness > comparisonImages[j].Closeness
	})
	fmt.Printf("Image %v is 1st with a closeness score of: %v\n", comparisonImages[0].Name, comparisonImages[0].Closeness)
	fmt.Printf("Image %v is 2nd with a closeness score of: %v\n", comparisonImages[1].Name, comparisonImages[1].Closeness)
	fmt.Printf("Image %v is 3rd with a closeness score of: %v\n", comparisonImages[2].Name, comparisonImages[2].Closeness)
}



