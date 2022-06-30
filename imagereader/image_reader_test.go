package imagereader

import (
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

func BenchmarkConvertBytesToPixels(b *testing.B) {
	testBytes, _ := GetImageBytes(sourceFiles[0])

	for i := 0; i < b.N; i++ {
		convertBytesToPixels(testBytes)
	}
}


// more tests for this one - inc edgecases n that. return errors if appropriate
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



// func TestCalculateClosenessFromRawBytesWithGoRoutines_DifferentSizeImages(t *testing.T) {
// 	testBytes1 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 0, 6, 54, 34}
// 	testBytes2 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 7}

// 	actualCloseness, _ := CalculateClosenessFromRawBytes(testBytes1, testBytes2)
// 	expectedCloseness := 0.0

// 	if actualCloseness != expectedCloseness {
// 		t.Error("Expected ", expectedCloseness, "Got ", actualCloseness)
// 	}

// }

// func TestCalculateClosenessFromRawBytesWithGoRoutines_ImageBytesLengthNotMultipleOfThree(t *testing.T) {
// 	testBytes1 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 0, 6, 54}
// 	testBytes2 := []byte{144, 144, 131, 255, 255, 255, 0, 0, 7, 6, 54}

// 	_, err := CalculateClosenessFromRawBytesWithGoRoutines(testBytes1, testBytes2)

// 	if err == nil {
// 		t.Error("Expected error Got nil")
// 	}

// }


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
//need the same tests for calculateclosenessfromrawbyteswithchannels ^

func BenchmarkCalculateClosenessFromRawBytes(b *testing.B) {
	testBytes1, _ := GetImageBytes(sourceFiles[0])
	testBytes2, _ := GetImageBytes(sourceFiles[1])

	for i := 0; i < b.N; i++ {
		CalculateClosenessFromRawBytes(testBytes1, testBytes2)
	}
}

func BenchmarkCalculateClosenessFromRawBytesWithGoRoutines(b *testing.B) {
	testBytes1, _ := GetImageBytes(sourceFiles[0])
	testBytes2, _ := GetImageBytes(sourceFiles[1])
	numberOfSlices := 1024
	for i := 0; i < b.N; i++ {
		CalculateClosenessFromRawBytesWithGoRoutines(testBytes1, testBytes2, numberOfSlices)
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



// func TestReadFilesIntoBytes(t *testing.T) {
// 	directory := "..\\Images\\Bronze\\"
// 	files, err := ioutil.ReadDir(directory)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	referenceImage := "1d25ea94-4562-4e19-848e-b60f1b58deee.raw"
// 	fmt.Println(files)
// 	actualResults := ReadFilesIntoBytes(files, referenceImage, directory)
// 	fmt.Println("xxxx")


// 	for i, result := range actualResults {
// 		if len(result.ImageBytes) != 3145728 {
// 			t.Errorf("Expected length of result %v's imageBytes to be 3145728, Got: %v", i, len(result.ImageBytes) )
// 		}
// 	}

// }




// func BenchmarkReadFilesIntoBytes(b *testing.B) {
// 	directory := "..\\Images\\Bronze\\"
// 	files, err := ioutil.ReadDir(directory)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	referenceImage := "1d25ea94-4562-4e19-848e-b60f1b58deee.raw"
// 	for i := 0; i < b.N; i++ {
// 		ReadFilesIntoBytes(files, referenceImage, directory)
// 	}

// }

// func BenchmarkCalculateClosenessForAllImages(b *testing.B) {
// 	directory := "..\\Images\\Bronze\\"
// 	files, err := ioutil.ReadDir(directory)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	referenceImage := "1d25ea94-4562-4e19-848e-b60f1b58deee.raw"
// 	referenceImageBytes, _ := GetImageBytes(directory+referenceImage)
// 	comparisonImages := ReadFilesIntoBytes(files, referenceImage, directory)

// 	// fmt.Println(comparisonImages)
// 	for i := 0; i < b.N; i++ {
// 		CalculateClosenessForAllImages(referenceImageBytes, comparisonImages)
// 	}
// }

//see what difference directorys make below
func BenchmarkProcessImageFiles(b *testing.B) {
	directory := "..\\Images\\Bronze\\"

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		b.Error("Problem reading files")
	}
	referenceImage := files[0].Name()
	referenceImageBytes, err := GetImageBytes(directory+referenceImage)

	referenceInfo := Reference{referenceImage, referenceImageBytes, directory}

	for i := 0; i < b.N; i++ {
		ProcessImageFiles(files, referenceInfo)
	}
}

func BenchmarkProcessImageFiles_OnlyCompareTop3(b *testing.B) {
	directory := "..\\Images\\Bronze\\"

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		b.Error("Problem reading files")
	}
	referenceImage := files[0].Name()
	referenceImageBytes, err := GetImageBytes(directory+referenceImage)

	referenceInfo := Reference{referenceImage, referenceImageBytes, directory}

	for i := 0; i < b.N; i++ {
		FindClosestImagesToReference(files, referenceInfo)
	}
}

func BenchmarkScanAndCompareImages(b *testing.B) {
	directory := "..\\Images\\Bronze\\"

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		b.Error("Problem reading files")
	}
	referenceImage := files[0].Name()
	referenceImageBytes, err := GetImageBytes(directory+referenceImage)

	referenceInfo := Reference{referenceImage, referenceImageBytes, directory}


	for i := 0; i < b.N; i++ {
		scanAndCompareImages(files, referenceInfo)
	}
}



func BenchmarkSortImagesByCloseness(b *testing.B) {
	directory := "..\\Images\\Bronze\\"

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		b.Error("Problem reading files")
	}
	referenceImage := files[0].Name()
	referenceImageBytes, err := GetImageBytes(directory+referenceImage)

	referenceInfo := Reference{referenceImage, referenceImageBytes, directory}
	comparedImages := scanAndCompareImages(files, referenceInfo)

	for i := 0; i < b.N; i++ {
		sortImagesByCloseness(comparedImages)
	}

}










