package imagereader

import (
	"reflect"
	"testing"
)

var sourceFiles []string = []string{
	
	"..\\Images\\Bronze\\1d25ea94-4562-4e19-848e-b60f1b58deee.raw",
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

func TestCalculateClosenessFromRawBytes(t *testing.T) {
	testBytes1 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 0, 6, 54, 34}
	testBytes2 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 7, 6, 2, 34}

	actualCloseness := CalculateClosenessFromRawBytes(testBytes1, testBytes2)
	expectedCloseness := 0.5

	if actualCloseness != expectedCloseness {
		t.Error("Expected ", expectedCloseness, "Got ", actualCloseness)
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










