package cinema

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// Video contains information about a video file and all the operations that
// need to be applied to it. Call Load to initialize a Video from file. Call the
// transformation functions to generate the desired output. Then call Render to
// generate the final output video file.
type Video struct {
	filepath string
	width    int
	height   int
	fps      int
	start    float64
	end      float64
	duration float64
	filters  []string
}

// Load gives you a Video that can be operated on. Load does not open the file
// or load it into memory. Apply operations to the Video and call Render to
// generate the output video file.
func Load(path string) (*Video, error) {
	if _, err := exec.LookPath("ffprobe"); err != nil {
		return nil, errors.New("cinema.Load: ffprobe was not found in your PATH environment variable, make sure to install ffmpeg (https://ffmpeg.org/) and add ffmpeg, ffplay and ffprobe to your PATH")
	}

	if _, err := os.Stat(path); err != nil {
		return nil, errors.New("cinema.Load: unable to load file: " + err.Error())
	}

	cmd := exec.Command(
		"ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		path,
	)
	out, err := cmd.Output()

	if err != nil {
		return nil, errors.New("cinema.Load: ffprobe failed: " + err.Error())
	}

	json := string(out)
	width := int(gjson.Get(json, "streams.0.width").Int())
	height := int(gjson.Get(json, "streams.0.height").Int())
	duration := gjson.Get(json, "format.duration").Float()

	return &Video{
		filepath: path,
		width:    width,
		height:   height,
		fps:      30,
		start:    0,
		end:      duration,
		duration: duration,
	}, nil
}

// Render applies all operations to the Video and creates an output video file
// of the given name.
func (v *Video) Render(output string) error {
	filters := strings.Join(v.filters[:], ",")
	filters += ",setsar=1,fps=fps=" + strconv.Itoa(int(v.fps))

	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", v.filepath,
		"-ss", strconv.Itoa(int(v.start)),
		"-t", strconv.Itoa(int(v.end-v.start)),
		"-vf", filters,
		"-strict", "-2",
		output,
	)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return errors.New("cinema.Video.Render: ffmpeg failed: " + err.Error())
	}
	return nil
}

// Trim sets the start and end time of the output video in seconds. It is always
// relative to the original input video.
func (v *Video) Trim(start, end float64) {
	v.start = start
	v.end = end
}

// Start returns the start of the video in seconds.
func (v *Video) Start() float64 {
	return v.start
}

// SetStart sets the start time of the output video in seconds. It is always
// relative to the original input video.
func (v *Video) SetStart(start float64) {
	v.start = start
}

// End returns the end of the video in seconds.
func (v *Video) End() float64 {
	return v.end
}

// SetEnd sets the end time of the output video in seconds. It is always
// relative to the original input video.
func (v *Video) SetEnd(end float64) {
	v.end = end
}

// SetFPS sets the framerate (frames per second) of the output video.
func (v *Video) SetFPS(fps int) {
	v.fps = fps
}

// SetSize sets the width and height of the output video.
func (v *Video) SetSize(width int, height int) {
	v.width = width
	v.height = height
	v.filters = append(v.filters, "scale="+strconv.Itoa(width)+":"+strconv.Itoa(height))
}

// Width returns the width of the video in pixels.
func (v *Video) Width() int {
	return v.width
}

// Height returns the width of the video in pixels.
func (v *Video) Height() int {
	return v.height
}

// Crop makes the output video a sub-rectangle of the input video. (0,0) is the
// top-left of the video, x goes right, y goes down.
func (v *Video) Crop(x, y, width, height int) {
	v.width = width
	v.height = height
	v.filters = append(v.filters, "crop="+strconv.Itoa(width)+":"+strconv.Itoa(height)+":"+strconv.Itoa(int(x))+":"+strconv.Itoa(int(y)))
}

// Filepath returns the path of the input video.
func (v *Video) Filepath() string {
	return v.filepath
}

// Duration returns the duration of the video in seconds.
func (v *Video) Duration() float64 {
	return v.duration
}

// Get the set fps of the current video struct
func (v *Video) FPS() int {
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
