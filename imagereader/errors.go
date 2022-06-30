package imagereader

import "errors"

var ErrInvalidLength = errors.New("Invalid length of imageBytes type []byte. All imageBytes must be length multiple of three")
var ErrInvalidSubsliceValue = errors.New("Invalid number of subslices requested. Subslices must contain multiples of 3 bytes")
