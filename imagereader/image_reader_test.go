package imagereader

import (
	"reflect"
	"testing"
)

var sourceFiles []string = []string{
	"..\\Images\\Bronze\\0f0e5c84-3b99-4874-bb8c-e0228155c4b5.raw",
	"..\\Images\\Bronze\\2ee00cf5-62b0-4808-a97d-b5fc33c5d81d.raw",
	"..\\Images\\Bronze\\3a09cc9f-0f41-4880-af85-7b256a8b4b3f.raw",
	"..\\Images\\Bronze\\6a9dbaf2-fae6-4ca2-bf88-43283d0b5534.raw",
	"..\\Images\\Bronze\\95fe8559-a9dc-4bb1-9286-5737f51cdec8.raw",
	"..\\Images\\Bronze\\172a9dbb-85af-4708-abbd-4ebfd5b84b52.raw",
	"..\\Images\\Bronze\\699eb551-5ab0-4102-8af9-32d69c3d03fc.raw",
	"..\\Images\\Bronze\\824603d7-6d45-4373-9906-e004dcaea97d.raw",
	"..\\Images\\Bronze\\af814aac-eb49-4cfb-a44f-42f500daf75f.raw",
	"..\\Images\\Bronze\\dcb76a5d-8c03-4054-b94d-117e24a8cba5.raw",
}

func TestGetImageBytes(t *testing.T) {
	sourceFile := sourceFiles[0]

	imageBytes, _ := getImageBytes(sourceFile) 
	expectedBytes := 3145728

	if len(imageBytes) != expectedBytes {
		t.Error("Expected ", expectedBytes, "Got", len(imageBytes))
	}
}

func TestGetImageBytes_FileError(t *testing.T) {
	sourceFile := "invalidsource.inv"

	_, err := getImageBytes(sourceFile) 

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

	actualCloseness := calculateClosenessFromRawBytes(testBytes1, testBytes2)
	expectedCloseness := 0.5

	if actualCloseness != expectedCloseness {
		t.Error("Expected ", expectedCloseness, "Got ", actualCloseness)
	}
}

func TestCalculateClosenessFromRawBytes_RoundsTo2dp(t *testing.T) {
	testBytes1 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 0}
	testBytes2 := []byte{144, 144, 131, 255, 255, 254, 0, 0, 7}

	actualCloseness := calculateClosenessFromRawBytes(testBytes1, testBytes2)
	expectedCloseness := 0.33

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

func TestCalculateClosenessFromPixels_RoundsTo2dp(t *testing.T) {
	testPixels1 := [][]byte{{144, 155, 166}, {188, 199, 220}, {0, 0, 0}}
	testPixels2 := [][]byte{{144, 155, 166}, {188, 199, 220}, {0, 0, 1}}

	actualCloseness := calculateClosenessFromPixels(testPixels1, testPixels2)
	expectedCloseness := 0.66

	if actualCloseness != expectedCloseness {
		t.Error("Expected ", expectedCloseness, "Got ", actualCloseness)
	}
}








