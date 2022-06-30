# Pixel-Challenge

## Usage

Images should be placed in a sub folder of the Images directory.
The reference image should also be placed in the same directory.

Use the -dir and -ref flags to specify a sub directory and reference image.

If no -ref image is set then the first image in the sub directory alphabetically will be used. 

If no -dir is set then the application will look for a Bronze directory as default.

Example usage:
To compare images in a sub directory Silver against a default reference image then ensure there is a folder named Silver containing the images in the Image folder and run:

    go run main.go -dir Silver

To compare images in a sub directory Gold/GoldImages against a speciric image then use

    go run main.go -dir Gold -ref 0aeb3950-8d36-4c29-be82-8bcdc82eb216.raw

