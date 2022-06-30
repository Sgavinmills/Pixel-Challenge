package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/Sgavinmills/Pixel-Challenge/main/imagereader"

	"net/http"
	_ "net/http/pprof"
)

func main() {
	start := time.Now()
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
		
	referenceInfo, files := setup()

	setupDuration := time.Since(start)
	fmt.Printf("Setup duration: %v\n", setupDuration)
	fmt.Printf("Comparing image %v against images in %v\n", referenceInfo.RefImageName, referenceInfo.Directory)
	
	imageResults := imagereader.FindClosestImagesToReference(files, referenceInfo)
	fmt.Printf("Image %v is 1st with a closeness score of: %v\n", imageResults[0].Name, imageResults[0].Closeness)
	fmt.Printf("Image %v is 2nd with a closeness score of: %v\n", imageResults[1].Name, imageResults[1].Closeness)
	fmt.Printf("Image %v is 3rd with a closeness score of: %v\n", imageResults[2].Name, imageResults[2].Closeness)

	fullExecutionDuration := time.Since(start)
	// time.Sleep(time.Second * 1000)
	fmt.Printf("Program run in %v ", fullExecutionDuration)

}

func setup() (imagereader.Reference, []fs.FileInfo) {

	var referenceInfo imagereader.Reference

	refImgPtr := flag.String("ref", "", "Name of image")
	dirPtr := flag.String("dir", "Bronze", "Which image level to use")
	flag.Parse()

	referenceInfo.Directory = "Images\\"+*dirPtr+"\\"
	referenceInfo.RefImageName = *refImgPtr
	
	files, err := ioutil.ReadDir(referenceInfo.Directory)

	if referenceInfo.RefImageName == "" {
		referenceInfo.RefImageName = files[0].Name()
	}
	referenceInfo.RefImageBytes, err = imagereader.GetImageBytes(referenceInfo.Directory+referenceInfo.RefImageName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return referenceInfo, files
}





