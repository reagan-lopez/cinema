# cinema : a lightweight video editor for golang



![alt text](https://i.imgur.com/uYRpL29.jpg "github.com/jtguibas/cinema")

## Overview [![GoDoc](https://godoc.org/github.com/jtguibas/cinema?status.svg)](https://godoc.org/github.com/jtguibas/cinema)

cinema is a super simple video editor that supports video io, video trimming, and resizing. it is dependent on ffmpeg, an advanced command-line tool used for handling video, audio, and other multimedia files and streams. start programmatically editing videos with golang now!

## Install
You must have [FFMPEG](https://ffmpeg.org/download.html) installed on your machine! Make sure `ffmpeg` and `ffprobe` are available from the command line on your machine.
```
go get github.com/jtguibas/cinema
```

## Example Usage

```golang
func main() { // cinema/test/test.go

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
	fmt.Println("FFMPEG Command", v.FFMPEG(""))
}
```

## TODO

- [ ] add concatenation support
- [x] improve godoc documentation
- [x] add cropping support
- [ ] expand to audio
- [ ] test windows and ubuntu support 
- [x] implement fps support
- [ ] implement bitrate support

## Contact
[jtguibas](https://github.com/jtguibas)

