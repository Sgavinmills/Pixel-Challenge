package imagereader

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

var sourceFiles []string = []string{
	
	"..\\Images\\Bronze\\1d25ea94-4562-4e19-848e-b60f1b58deee.raw",
	"..\\Images\\Bronze\\6c9952ef-e5bf-4de2-817b-fd0073be8449.raw",
}

func TestGetImageBytes(t *testing.T) {
	sourceFile := sourceFiles[0]

	imageBytes, _ := GetImageBytes(sourceFile) 
	expectedBytes := 3145728

	if len(imageBytes) != expectedBytes {
		t.Error("Expected ", expectedBytes, "Got", len(imageBytes))
	}
}

func TestGetImageBytes_FileError(t *testing.T) {
	sourceFile := "invalidsource.inv"

	_, err := GetImageBytes(sourceFile) 

	if err == nil {
		t.Error("Expected Error Got nil")
	}
}

func TestConvertBytesToPixels(t *testing.T) {
	testBytes := []byte{144, 155, 166, 188, 199, 220, 255, 255, 255}

	actualPixels := convertBytesToPixels(testBytes)
	expectedPixels := [][]byte{{144, 155, 166}, {188, 199, 220}, {255, 255, 255}}

	if !reflect.DeepEqual(actualPixels, expectedPixels) {
		t.Error("Expects ", expectedPixels, "Got", actualPixels)
	}
}

func BenchmarkConvertBytesToPixels(b *testing.B) {
	testBytes, _ := GetImageBytes(sourceFiles[0])

	for i := 0; i < b.N; i++ {
		convertBytesToPixels(testBytes)
	}
}

func TestCalculateClosenessFromRawBytes(t *testing.T) {
	testBytes1 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 0, 6, 54, 34}
	testBytes2 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 7, 6, 2, 34}

	actualCloseness := CalculateClosenessFromRawBytes(testBytes1, testBytes2)
	expectedCloseness := 0.5

	if actualCloseness != expectedCloseness {
		t.Error("Expected ", expectedCloseness, "Got ", actualCloseness)
	}
}

func BenchmarkCalculateClosenessFromPixels(b *testing.B) {
	testBytes1, _ := GetImageBytes(sourceFiles[0])
	testBytes2, _ := GetImageBytes(sourceFiles[1])
	testPixels1 := convertBytesToPixels(testBytes1)
	testPixels2 := convertBytesToPixels(testBytes2)

	for i := 0; i < b.N; i++ {
		calculateClosenessFromPixels(testPixels1, testPixels2)
	}
}



func TestCalculateClosenessFromPixels(t *testing.T) {
	testPixels1 := [][]byte{{144, 155, 166}, {188, 199, 220}, {255, 255, 255}, {0, 0, 0}}
	testPixels2 := [][]byte{{144, 155, 166}, {188, 199, 220}, {255, 255, 255}, {0, 0, 1}}

	actualCloseness := calculateClosenessFromPixels(testPixels1, testPixels2)
	expectedCloseness := 0.75

	if actualCloseness != expectedCloseness {
		t.Error("Expected ", expectedCloseness, "Got ", actualCloseness)
	}
}

func BenchmarkCalculateClosenessFromRawBytes(b *testing.B) {
	testBytes1, _ := GetImageBytes(sourceFiles[0])
	testBytes2, _ := GetImageBytes(sourceFiles[1])

	for i := 0; i < b.N; i++ {
		CalculateClosenessFromRawBytes(testBytes1, testBytes2)
	}
}

func TestReadFilesIntoBytes(t *testing.T) {
	directory := "..\\Images\\Bronze\\"
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println(err)
	}
	referenceImage := "1d25ea94-4562-4e19-848e-b60f1b58deee.raw"
	fmt.Println(files)
	actualResults := ReadFilesIntoBytes(files, referenceImage, directory)
	fmt.Println("xxxx")


	for i, result := range actualResults {
		if len(result.ImageBytes) != 3145728 {
			t.Errorf("Expected length of result %v's imageBytes to be 3145728, Got: %v", i, len(result.ImageBytes) )
		}
	}

}




func BenchmarkReadFilesIntoBytes(b *testing.B) {
	directory := "..\\Images\\Bronze\\"
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println(err)
	}
	referenceImage := "1d25ea94-4562-4e19-848e-b60f1b58deee.raw"
	for i := 0; i < b.N; i++ {
		ReadFilesIntoBytes(files, referenceImage, directory)
	}

}

func BenchmarkCalculateClosenessForAllImages(b *testing.B) {
	directory := "..\\Images\\Bronze\\"
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println(err)
	}
	referenceImage := "1d25ea94-4562-4e19-848e-b60f1b58deee.raw"
	referenceImageBytes, _ := GetImageBytes(directory+referenceImage)
	comparisonImages := ReadFilesIntoBytes(files, referenceImage, directory)

	// fmt.Println(comparisonImages)
	for i := 0; i < b.N; i++ {
		CalculateClosenessForAllImages(referenceImageBytes, comparisonImages)
	}
}










