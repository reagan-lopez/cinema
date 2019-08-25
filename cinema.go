package cinema

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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
	start    time.Duration
	end      time.Duration
	duration time.Duration
	filters  []string
}

// Load gives you a Video that can be operated on. Load does not open the file
// or load it into memory. Apply operations to the Video and call Render to
// generate the output video file.
func Load(path string) (*Video, error) {
	if _, err := exec.LookPath("ffprobe"); err != nil {
		return nil, errors.New("cinema.Load: ffprobe was not found in your PATH " +
			"environment variable, make sure to install ffmpeg " +
			"(https://ffmpeg.org/) and add ffmpeg, ffplay and ffprobe to your " +
			"PATH")
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

	type description struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
		Format struct {
			DurationSec json.Number `json:"duration"`
		} `json:"format"`
	}
	var desc description
	if err := json.Unmarshal(out, &desc); err != nil {
		return nil, errors.New("cinema.Load: unable to parse JSON output " +
			"from ffprobe: " + err.Error())
	}
	if len(desc.Streams) == 0 {
		return nil, errors.New("cinema.Load: ffprobe does not contain stream " +
			"data, make sure the file " + path + " contains a valid video.")
	}

	secs, err := desc.Format.DurationSec.Float64()
	if err != nil {
		return nil, errors.New("cinema.Load: ffprobe returned invalid duration: " +
			err.Error())
	}
	duration := time.Duration(secs*float64(time.Second) + 0.5)

	return &Video{
		filepath: path,
		width:    desc.Streams[0].Width,
		height:   desc.Streams[0].Height,
		fps:      30,
		start:    0,
		end:      duration,
		duration: duration,
	}, nil
}

// Render applies all operations to the Video and creates an output video file
// of the given name.
func (v *Video) Render(output string) error {
	line := v.CommandLine(output)
	cmd := exec.Command(line[0], line[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return errors.New("cinema.Video.Render: ffmpeg failed: " + err.Error())
	}
	return nil
}

// CommandLine returns the command line that will be used to convert the Video
// if you were to call Render.
func (v *Video) CommandLine(output string) []string {
	var filters string
	if len(v.filters) > 0 {
		filters = strings.Join(v.filters, ",") + ","
	}
	filters += "setsar=1,fps=fps=" + strconv.Itoa(int(v.fps))

	return []string{
		"ffmpeg",
		"-y",
		"-i", v.filepath,
		"-ss", strconv.FormatFloat(v.start.Seconds(), 'f', -1, 64),
		"-t", strconv.FormatFloat((v.end - v.start).Seconds(), 'f', -1, 64),
		"-vf", filters,
		"-strict", "-2",
		output,
	}
}

// Trim sets the start and end time of the output video. It is always relative
// to the original input video.
func (v *Video) Trim(start, end time.Duration) {
	v.start = start
	v.end = end
}

// Start returns the start of the video .
func (v *Video) Start() time.Duration {
	return v.start
}

// SetStart sets the start time of the output video. It is always relative to
// the original input video.
func (v *Video) SetStart(start time.Duration) {
	v.start = start
}

// End returns the end of the video.
func (v *Video) End() time.Duration {
	return v.end
}

// SetEnd sets the end time of the output video. It is always relative to the
// original input video.
func (v *Video) SetEnd(end time.Duration) {
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
	v.filters = append(v.filters, fmt.Sprintf("scale=%d:%d", width, height))
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
	v.filters = append(
		v.filters,
		fmt.Sprintf("crop=%d:%d:%d:%d", width, height, x, y),
	)
}

// Filepath returns the path of the input video.
func (v *Video) Filepath() string {
	return v.filepath
}

// Duration returns the duration of the video in seconds.
func (v *Video) Duration() time.Duration {
	return v.duration
}

// Get the set fps of the current video struct
func (v *Video) FPS() int {
	return v.fps
}
