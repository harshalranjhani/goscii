package cmd

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strconv"
	"strings"

	creator "github.com/faakern/ascii-creator"
	"github.com/nfnt/resize"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate <filepath> <width> <height>",
	Short: "Generate the ASCII art from the image. Adding the width and height is optional. The default width is 120 and the height is calculated based on the aspect ratio of the image.",
	Long:  `Generate the ASCII art from the image and print it to the console. The image path is passed as an argument to the command. The ASCII art is generated using the image.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: generate <filepath> [<width> <height>]")
			os.Exit(1)
		}
		// Register supported image types
		image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
		image.RegisterFormat("jpg", "jpg", jpeg.Decode, jpeg.DecodeConfig)
		fmt.Printf("Converting image '%s'...\n", args[0])

		// Open the input file and decode the image
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Printf("Could not find or open file: %s\n", err)
			os.Exit(1)
		}
		img, _, err := image.Decode(file)

		if err != nil {
			fmt.Printf("Could not find or open file: %s\n", err)
			os.Exit(1)
		}

		// Determine the desired width and height for the ASCII art
		var desiredWidth, desiredHeight int

		if len(args) > 1 {
			desiredWidth, err = strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Width value is invalid")
				os.Exit(1)
			}
		} else {
			desiredWidth = 120
		}

		if len(args) > 2 {
			desiredHeight, err = strconv.Atoi(args[2])
			if err != nil {
				fmt.Println("Height value is invalid")
				os.Exit(1)
			}
		}

		// Calculate the scaling factor for width and height
		scaleFactorWidth := float64(desiredWidth) / float64(img.Bounds().Dx())
		if desiredHeight == 0 {
			desiredHeight = int(float64(img.Bounds().Dy()) * scaleFactorWidth)
		}
		scaleFactorHeight := float64(desiredHeight) / float64(img.Bounds().Dy())

		// Resize the image
		resizedImg := resize.Resize(uint(float64(img.Bounds().Dx())*scaleFactorWidth), uint(float64(img.Bounds().Dy())*scaleFactorHeight), img, resize.Lanczos3)

		// Create a builder for ASCII conversion/creation
		builder := creator.NewBuilder()

		// Provide a list of characters - these should be arranged from 'darker' to 'lighter' values,
		// an input image, and build a generator to be used as the basis for the conversion.
		gamma, err := strconv.ParseFloat("1", 32)
		if err != nil {
			fmt.Println("Gamma correction value is invalid")
			os.Exit(1)
		}

		generator := builder.WithCharSet(creator.CharSet{
			Characters: []byte{' ', '.', ',', ':', ';', '+', '*', '?', '%', '&', '#', '@'},
		}).WithGammaCorrection(float32(gamma)).WithInput().Image(resizedImg).Build()

		// Do the actual conversion/ASCII generation
		var out creator.Result
		err = generator.Generate(&out)
		if err != nil {
			fmt.Printf("Error converting image: %s\n", err)
			os.Exit(1)
		}

		// Write the result to an output file
		file, err = os.Create(fmt.Sprintf("%s.txt", strings.Split(args[0], ".")[0]))
		if err != nil {
			fmt.Println("Could not create output file")
			os.Exit(1)
		}

		size, err := file.Write(out.Ascii)
		if err != nil {
			fmt.Println("Could not write output to file")
			os.Exit(1)
		}

		fmt.Printf("Wrote %d bytes to %s\n", size, file.Name())

		err = file.Close()
		if err != nil {
			os.Exit(1)
		}

		fmt.Println("Displaying ASCII art:")
		fmt.Println(string(out.Ascii))
	},
}
