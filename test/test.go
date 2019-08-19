package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jtguibas/cinema"
)

func main() {

	// loading the test video
	fmt.Println("Downloading Test Video...")
	video_url := "https://media.w3.org/2010/05/sintel/trailer.mp4"
	if err := DownloadFile("test_input.mp4", video_url); err != nil {
		panic(err)
	}

	// initializing the test video as a cinema video object
	v, err := cinema.MakeVideo("test_input.mp4")
	if err != nil {
		fmt.Println(err)
	}

	// testing all setters
	v.Trim(0, 10)               // trimming the video from t=0 seconds -> t=10 seconds
	v.SetStart(1)               //trimming only the start of the video from t=1
	v.SetEnd(9)                 // trimming only the end of the video at t=9
	v.SetSize(400, 400)         //resizing the video to 400x400
	v.Crop(200, 200, 0, 0)      //cropping the 400x400 video into a 200x200 video from position x=0,y=0
	v.SetSize(400, 400)         //resizing the cropped 200x200 video to a 400x400 video
	v.SetFps(48)                //set the output fps to 48
	v.Render("test_output.mov") // notice how format conversion is done with ease

	// testing all getters
	fmt.Println("Output Filepath", v.Filepath())
	fmt.Println("Start", v.Start())
	fmt.Println("End", v.End())
	fmt.Println("Width", v.Width())
	fmt.Println("Height", v.Height())
	fmt.Println("Duration", v.Duration())
	fmt.Println("FPS", v.Fps())
	fmt.Println("Filters", v.Filters())
	fmt.Println("FFMPEG Command", v.FFMPEG("render.mp4"))
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
