package main

import (
	"testing"
)

func TestGetFilesFromDirectory(t *testing.T) {
	var testImageStorage imageStorage
	testImageStorage.comparisonImages = make(map[string]imageInfo)

	testImageStorage.getFilesFromDirectory("testdata\\testfiles\\", "file2.raw")
	
	if testImageStorage.comparisonImages["file1.raw"].name != "file1.raw" {
		t.Error("could not find file 1")
	}

	if testImageStorage.comparisonImages["file2.raw"].name == "file2.raw" {
		t.Error("File2 should have been excluded as the reference image")
	}


	

}