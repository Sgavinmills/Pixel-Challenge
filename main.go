package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/Sgavinmills/Pixel-Challenge/main/imagereader"
)

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
	
	
	// doesnt work atm, need to wait until finished reading before moving on... but...
	// could still bench mark / test the other function to check for race issues and improvements. 
	
	
	
	// resultCounter := imagereader.NewCounter(10)
	
	// type image struct {
	// 	name string
	// 	closeness float64
	// }

	// imagesReadCounter := make(chan struct{})
	// bytesComparedCounter := make(chan struct{})
	imagesProcessedCounter := make(chan imagereader.Image)
	//TODO FIND BETTER WAY THAN LEN OF FILES - 1, MIGHT NOT ALWAYS BE THE CASE
	// comparisonImages := imagereader.ReadFilesIntoBytes(files, referenceImage, directory)
	// imagereader.ReadFilesIntoBytes_WithChannels(files, referenceImage, directory, imagesReadCounter)


	//TODO IMPORTANT! STRIP OUT REFERENCE IMAGE FROM FILES LIST. Then dont need to pass in referenceImage
	// this function will read the bytes AND THEN compare to reference bytes and return closenss and name on image
	imagereader.ImageReader(files, referenceImage, referenceImageBytes, directory, imagesProcessedCounter)
	var comparedImages []imagereader.Image
	for i := 0; i < len(files)-1; i++ {
		// this will listen to replies on imagesProcessedCounter and add the results to a slice
		// in future will only add top 3 to save sorting time
		// var processedImage imagereader.Image
		processedImage := <- imagesProcessedCounter
		comparedImages = append(comparedImages, processedImage)

	}

	// for i := 0; i < len(files)-1; i++ {
	// 	// receive from the bytescomparedcounter with the name and closeness result
	// 	// add this to the list (later than only add if its in top 3 to save further time)
	// }
	// imagereader.StartResultsWatcher(resultCounter)
	// imagereader.CalculateClosenessForAllImages(referenceImageBytes, comparisonImages, testCounter)




	// sort comparisonImages to get winning order
	// might want to look into faster sorting algorithm OR only populate array with contenders in the first place, ie never allow more than 3 
	// basically check if there are more than 3 already stored, if so enter it in order and delete the slowest
	// time.Sleep(time.Second * 5)
	
	// for i := 0; i < len(comparisonImages); i++ {
	// 	<- testCounter
	// }
	// print top 3. Might need to add some checking to see that there are no ties.
	sort.Slice(comparedImages, func(i, j int) bool {
		return comparedImages[i].Closeness > comparedImages[j].Closeness
	})
	fmt.Printf("Image %v is 1st with a closeness score of: %v\n", comparedImages[0].Name, comparedImages[0].Closeness)
	fmt.Printf("Image %v is 2nd with a closeness score of: %v\n", comparedImages[1].Name, comparedImages[1].Closeness)
	fmt.Printf("Image %v is 3rd with a closeness score of: %v\n", comparedImages[2].Name, comparedImages[2].Closeness)
}



