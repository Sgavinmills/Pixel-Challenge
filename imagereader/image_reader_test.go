package imagereader

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var sourceFiles []string = []string{
	
	"testdata\\Bronze\\1d25ea94-4562-4e19-848e-b60f1b58deee.raw",
	"testdate\\Bronze\\6c9952ef-e5bf-4de2-817b-fd0073be8449.raw",
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
	sourceFile := "filenotfound.inv"

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


func TestCreateSubSlices(t *testing.T) {
	testBytes1 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 0, 6, 54, 34}
	testBytes2 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 7, 6, 2, 34}

	result1, result2, _ := createSubSlices(testBytes1, testBytes2, 4)
	expectedResult1 := [][]byte{{144, 144, 131},{255,255,255},{0,0,0},{6,54,34}}
	expectedResult2 := [][]byte{{144, 144, 131},{255,255,255},{0,0,7},{6,2,34}}

	if !reflect.DeepEqual(result1, expectedResult1) {
		t.Error("Expected", expectedResult1, "Got ", result1)
	}
	if !reflect.DeepEqual(result2, expectedResult2) {
		t.Error("Expected", expectedResult2, "Got ", result2)
	}
}

func TestCreateSubSlices_InvalidSubSliceAmount(t *testing.T) {
testBytes1 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 0, 6, 54, 34}
testBytes2 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 7, 6, 2, 34}

_, _, err := createSubSlices(testBytes1, testBytes2, 3)
if err != ErrInvalidSubsliceValue {
	t.Error("Expected ", ErrInvalidSubsliceValue, "Got ", err)
}
}

func TestCalculateClosenessFromRawBytesWithGoRoutines(t *testing.T) {
	testBytes1 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 0, 6, 54, 34}
	testBytes2 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 7, 6, 2, 34}
	numberOfSlices := 4
	actualCloseness, err := CalculateClosenessFromRawBytesWithGoRoutines(testBytes1, testBytes2, numberOfSlices )
	expectedCloseness := 0.5
	if actualCloseness != expectedCloseness {
		t.Error("Expected ", expectedCloseness, "Got ", actualCloseness)
		if err != nil {
			t.Error(err.Error())
		}
	}
}


func TestCalculateClosenessFromRawBytesWithGoRoutines_DifferentSizeImages(t *testing.T) {
	testBytes1 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 0, 6, 54, 34}
	testBytes2 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 7}

	actualCloseness, _ := CalculateClosenessFromRawBytes(testBytes1, testBytes2)
	expectedCloseness := 0.0

	if actualCloseness != expectedCloseness {
		t.Error("Expected ", expectedCloseness, "Got ", actualCloseness)
	}

}


func TestCalculateClosenessFromRawBytes(t *testing.T) {
	testBytes1 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 0, 6, 54, 34}
	testBytes2 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 7, 6, 2, 34}

	actualCloseness, _ := CalculateClosenessFromRawBytes(testBytes1, testBytes2)
	expectedCloseness := 0.5

	if actualCloseness != expectedCloseness {
		t.Error("Expected ", expectedCloseness, "Got ", actualCloseness)
	}
}

func TestCalculateClosenessFromRawBytes_DifferentSizeImages(t *testing.T) {
	testBytes1 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 0, 6, 54, 34}
	testBytes2 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 7}

	actualCloseness, _ := CalculateClosenessFromRawBytes(testBytes1, testBytes2)
	expectedCloseness := 0.0

	if actualCloseness != expectedCloseness {
		t.Error("Expected ", expectedCloseness, "Got ", actualCloseness)
	}

}

func TestCalculateClosenessFromRawBytes_ImageBytesLengthNotMultipleOfThree(t *testing.T) {
	testBytes1 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 0, 6, 54}
	testBytes2 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 7, 6, 54}

	_, err := CalculateClosenessFromRawBytes(testBytes1, testBytes2)

	if err == nil {
		t.Error("Expected error Got nil")
	}

}

func TestFindClosestImagesToReference(t *testing.T) {

	referenceInfo := Reference{
		Directory: "testdata\\Bronze\\",
		RefImageName: "1d25ea94-4562-4e19-848e-b60f1b58deee.raw",
	}
	var err error
	referenceInfo.RefImageBytes, err = GetImageBytes(referenceInfo.Directory+referenceInfo.RefImageName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	files, err := ioutil.ReadDir(referenceInfo.Directory)
	results := FindClosestImagesToReference(files, referenceInfo)
	
	expectedResults := []string{
		"4424fa5a-c00d-4cd5-8525-fcf921b09ca8.raw",
		"e3c342e2-2429-4f47-8828-f7ee0703ad38.raw",
		"6c9952ef-e5bf-4de2-817b-fd0073be8449.raw",
	}

	for i, v := range expectedResults {
		if results[i].Name != v {
			t.Error("Expected result ", i, " to be ", expectedResults[i], "Got ", results[i])
		}
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

func BenchmarkCalculateClosenessFromPixels(b *testing.B) {
	testBytes1, _ := GetImageBytes(sourceFiles[0])
	testBytes2, _ := GetImageBytes(sourceFiles[1])
	testPixels1 := convertBytesToPixels(testBytes1)
	testPixels2 := convertBytesToPixels(testBytes2)

	for i := 0; i < b.N; i++ {
		calculateClosenessFromPixels(testPixels1, testPixels2)
	}
}


func BenchmarkConvertBytesToPixels(b *testing.B) {
	testBytes, _ := GetImageBytes(sourceFiles[0])

	for i := 0; i < b.N; i++ {
		convertBytesToPixels(testBytes)
	}
}

func BenchmarkCalculateClosenessFromRawBytes(b *testing.B) {
	testBytes1, err := GetImageBytes(sourceFiles[0])
	if err != nil {
		fmt.Println(err)
	}

	testBytes2, err := GetImageBytes(sourceFiles[1])
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < b.N; i++ {
		_, err := CalculateClosenessFromRawBytes(testBytes1, testBytes2)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkCalculateClosenessFromRawBytesWithGoRoutines(b *testing.B) {
	testBytes1, err := GetImageBytes(sourceFiles[0])
	if err != nil {
		fmt.Println(err)
	}

	testBytes2, err := GetImageBytes(sourceFiles[1])
	if err != nil {
		fmt.Println(err)
	}

	numberOfSlices := 1024
	for i := 0; i < b.N; i++ {
		_, err := CalculateClosenessFromRawBytesWithGoRoutines(testBytes1, testBytes2, numberOfSlices)
		if err != nil {
			fmt.Println(err)
		}
	}
}
