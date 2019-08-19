package cinema

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// The structure that defines the Video structure. All the values in this structure will be filled when MakeVideo is called.
type Video struct {
	filepath string
	width    int64
	height   int64
	fps      int64
	start    float64
	end      float64
	duration float64
	filters  []string
}

// MakeVideo takes in the filepath of any videofile supported by FFMPEG.
// It will return a Video structure with all of the values parsed from FFPROBE.
// Note: This will not load the video into memory.
func MakeVideo(filepath string) (Video, error) {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", filepath)
	out, err := cmd.Output()

	if err != nil {
		fmt.Println("Error: FFPROBE did not work. Check the filename and thatt ffmpeg is correctly installed.", filepath)
		return Video{}, err
	}

	json := string(out)
	width := gjson.Get(json, "streams.0.width").Int()
	height := gjson.Get(json, "streams.0.height").Int()
	duration := gjson.Get(json, "format.duration").Float()

	return Video{filepath: filepath, width: width, height: height, start: 0, end: duration, duration: duration}, nil
}

// This function will take the configuration set in the Video structure and properly apply it to the rendered video.
// Currently supports trimming, scaling, and video format conversion.
func (v *Video) Render(output string) {

	filter_chain := strings.Join(v.filters[:], ",") + ",setsar=1" + ",fps=fps=" + strconv.Itoa(int(v.fps))

	cmd := exec.Command("ffmpeg",
		"-y",
		"-i",
		v.filepath,
		"-ss",
		strconv.Itoa(int(v.start)),
		"-t",
		strconv.Itoa(int(v.end-v.start)),
		"-vf",
		filter_chain,
		"-strict",
		"-2",
		output)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

}

// Trim the video in seconds
func (v *Video) Trim(start float64, end float64) {
	v.start = start
	v.end = end
}

// Trim to the start of the video in seconds
func (v *Video) SetStart(start float64) {
	v.start = start
}

// Trim the end of the video using seconds
func (v *Video) SetEnd(end float64) {
	v.end = end
}

// Trim the end of the video using seconds
func (v *Video) SetFps(fps int64) {
	v.fps = fps
}

// Set the width and height of the video
func (v *Video) SetSize(width int64, height int64) {
	v.width = width
	v.height = height
	v.filters = append(v.filters, "scale="+strconv.Itoa(int(width))+":"+strconv.Itoa(int(height)))
}

// Crop the video based on width, height, x-coordinate, and y-coordinate (from top left)
func (v *Video) Crop(width int64, height int64, x int64, y int64) {
	v.width = width
	v.height = height
	v.filters = append(v.filters, "crop="+strconv.Itoa(int(width))+":"+strconv.Itoa(int(height))+":"+strconv.Itoa(int(x))+":"+strconv.Itoa(int(y)))
}

// Get the filepath of the current video struct
func (v *Video) Filepath() string {
	return v.filepath
}

// Get the start of the current video struct
func (v *Video) Start() float64 {
	return v.start
}

// Get the end of the current video struct
func (v *Video) End() float64 {
	return v.end
}

// Get the width of the current video struct
func (v *Video) Width() int64 {
	return v.width
}

// Get the height of the current video struct
func (v *Video) Height() int64 {
	return v.height
}

// Get the duration of the current video struct
func (v *Video) Duration() float64 {
	return v.duration
}

// Get the set fps of the current video struct
func (v *Video) Fps() int64 {
	return v.fps
}

// Get the current Filters using in the -vf flag of ffmpeg
func (v *Video) Filters() []string {
	return v.filters
}

// Get the output ffmpeg command cinema will run at .Render()
func (v *Video) FFMPEG(output string) string {
	filter_chain := strings.Join(v.filters[:], ",") + ",setsar=1" + ",fps=fps=" + strconv.Itoa(int(v.fps))
	cmd := "ffmpeg" + " " +
		"-y" + " " +
		"-i" + " " +
		v.filepath + " " +
		"-ss" + " " +
		strconv.Itoa(int(v.start)) + " " +
		"-t" + " " +
		strconv.Itoa(int(v.end-v.start)) + " " +
		"-filter:v" + " " +
		filter_chain + " " +
		"-strict" + " " +
		"-2" + " " +
		output
	return cmd
}
