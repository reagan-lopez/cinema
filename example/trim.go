package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jtguibas/cinema"
)

const (
	DEFAULT_CLIP_DURATION   = 25
	DEFAULT_TRIM_DURATION   = 10
	DEFAULT_INPUT_FILENAME  = "Sham.mp4"
	DEFAULT_TRIM_DIR_PREFIX = "trimDir"
)

type VideoClip struct {
	path  string
	video cinema.Video
}

func main() {
	inputFilename := flag.String("input-file", DEFAULT_INPUT_FILENAME, "input file name")
	trimDuration := flag.Int("trim-duration", DEFAULT_TRIM_DURATION, "trim duration in seconds")

	trimBulk(*inputFilename, time.Duration(*trimDuration))

	// you can also generate the command line instead of applying it directly
	//fmt.Println("FFMPEG Command", video.CommandLine(outputFile))
}

func trimBulk(inputFile string, trimDuration time.Duration) {
	var err error
	var count int
	clipDuration := DEFAULT_CLIP_DURATION
	trimDirname := DEFAULT_TRIM_DIR_PREFIX + strings.TrimSuffix(inputFile, filepath.Ext(inputFile))

	log.Println("Reading input file: " + inputFile)
	video, err := cinema.Load(inputFile)
	check(err)

	inputDuration := int(video.Duration().Seconds())
	log.Println("Input duration: " + strconv.Itoa(inputDuration))

	mkdirHard(trimDirname) // Create output directory
	check(err)

	capacity := inputDuration/(clipDuration+int(trimDuration)) + 1
	var queue = make(chan VideoClip, capacity) // Create queue for trimmed videos

	for s := 0; s < int(video.Duration().Seconds()); s += clipDuration + int(trimDuration) {
		count = count + 1
		outputFilename := strings.Replace(inputFile, filepath.Ext(inputFile),
			"-"+strconv.Itoa(count)+filepath.Ext(inputFile), 1)
		outputFile := filepath.Join(trimDirname, outputFilename)

		log.Println("Trimming output file: " + outputFile)
		video.Trim(time.Duration(s)*time.Second, time.Duration(s+clipDuration)*time.Second)
		vc := VideoClip{outputFile, *video}

		queue <- vc // Adding to trim queue
	}
	close(queue)
	log.Println("Trimming complete")

	// Render video clips async
	/*
	timer := &stopWatch{w: os.Stdout}
	defer timer.WriteSummary(os.Stdout)
	timer.Start("Generating output..")
	*/

	var wg sync.WaitGroup
	for elem := range queue {
		wg.Add(1)
		go renderWorker(elem, &wg)
	}
	wg.Wait()

	log.Println("Complete!")
}

func renderWorker(vc VideoClip, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Writing output file: " + vc.path)
	vc.video.Render(vc.path)
	log.Println("Done: " + vc.path)
}

func mkdirHard(path string) {
	os.Remove(path)
	os.Mkdir(path, os.ModePerm)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

//ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 Sham4.mp4
