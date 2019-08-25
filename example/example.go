package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jtguibas/cinema"
)

func main() {
	downloadTestVideo("example.mp4")

	video, err := cinema.Load("example.mp4")
	check(err)

	video.Trim(10, 20)              // trim video from 10 to 20 seconds
	video.SetStart(1)               // trim first second of the video
	video.SetEnd(9)                 // keep only up to 9 seconds
	video.SetSize(400, 300)         // resize video to 400x300
	video.Crop(0, 0, 200, 200)      // crop rectangle top-left (0,0) with size 200x200
	video.SetSize(400, 400)         // resize cropped 200x200 video to a 400x400
	video.SetFPS(48)                // set the output framerate to 48 frames per second
	video.Render("test_output.mov") // note format conversion by file extension

	// you can also generate the command line instead of applying it directly
	fmt.Println("FFMPEG Command", video.CommandLine("test_output.mov"))
}

func downloadTestVideo(to string) {
	const url = "https://media.w3.org/2010/05/sintel/trailer.mp4"

	fmt.Println("downloading test video...")
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()

	out, err := os.Create(to)
	check(err)
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
