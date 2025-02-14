package convert_video

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ConvertVideo struct {
	videoPath  string
	codec      string
	saveDir    string
	outputName string
	typeVideo  string
	result
}

type result struct {
	outputNameVideo      string
	outputPathVideo      string
	outputNameAudio      string
	outputPathAudio      string
	outputFinalPathVideo string
}

func NewVideoConvert() *ConvertVideo {
	return &ConvertVideo{}
}

func (v *ConvertVideo) Encode(targetVideo, saveDir, outputName, codec, typeFile string) (string, string, error) {
	v.videoPath = targetVideo
	v.codec = codec
	v.saveDir = saveDir
	v.outputName = outputName
	v.typeVideo = typeFile

	if err := v.writeVideo(); err != nil {
		return "", "", err
	}
	hasAudio, err := v.extractAudio()
	if err != nil {
		return "", "", err
	}

	if hasAudio {
		if err = v.mergeAudioVideo(); err != nil {
			return "", "", err
		}
		return v.outputFinalPathVideo, v.outputNameVideo, nil
	}
	return v.outputPathVideo, v.outputNameVideo, nil
}

func (v *ConvertVideo) writeVideo() error {
	video, err := gocv.VideoCaptureFile(v.videoPath)
	if err != nil {
		return err
	}
	defer video.Close()
	frame := gocv.NewMat()
	defer frame.Close()

	if ok := video.Read(&frame); !ok {
		return fmt.Errorf("cannot read frame from %v", v.videoPath)
	}
	fps := video.Get(gocv.VideoCaptureFPS)
	width := int(video.Get(gocv.VideoCaptureFrameWidth))
	height := int(video.Get(gocv.VideoCaptureFrameHeight))

	v.outputNameVideo = strings.Split(v.outputName, ".")[0] + "." + v.typeVideo
	videoDir := filepath.Join(v.saveDir, "video")
	v.outputPathVideo = filepath.Join(videoDir, v.outputNameVideo)

	if err = os.MkdirAll(videoDir, os.ModePerm); err != nil {
		return err
	}

	writer, err := gocv.VideoWriterFile(v.outputPathVideo, v.codec, fps, width, height, true)
	if err != nil {
		return err
	}
	defer writer.Close()

	for {
		if !video.Read(&frame) {
			break
		}
		if frame.Empty() {
			continue
		}

		writer.Write(frame)
	}

	return nil
}

func (v *ConvertVideo) extractAudio() (bool, error) {
	v.outputNameAudio = strings.Split(v.outputName, ".")[0] + ".mp3"
	audioDir := filepath.Join(v.saveDir, "audio")
	v.outputPathAudio = filepath.Join(audioDir, v.outputNameAudio)

	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-select_streams", "a:0",
		"-show_entries", "stream=codec_type",
		"-of", "default=noprint_wrappers=1:nokey=1",
		v.videoPath,
	)
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	if strings.TrimSpace(string(output)) == "audio" {

		if err = os.MkdirAll(audioDir, os.ModePerm); err != nil {
			return false, err
		}

		cmd = exec.Command("ffmpeg",
			"-i", v.videoPath,
			"-q:a", "0",
			"-map", "a",
			v.outputPathAudio,
		)
		output, err = cmd.CombinedOutput()
		if err != nil {
			return false, fmt.Errorf("ffmpeg error: %v\nOutput: %s", err, string(output))
		}
		return true, nil
	}

	return false, nil
}

func (v *ConvertVideo) mergeAudioVideo() error {
	finalVideoDir := filepath.Join(v.saveDir, "final")
	v.outputFinalPathVideo = filepath.Join(finalVideoDir, v.outputNameVideo)
	if err := os.MkdirAll(finalVideoDir, os.ModePerm); err != nil {
		return err
	}

	cmd := exec.Command(
		"ffmpeg",
		"-i", v.outputPathVideo,
		"-i", v.outputPathAudio,
		"-c:v", "copy",
		"-c:a", "copy",
		"-y",
		v.outputFinalPathVideo,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v\nOutput: %s", err, string(output))
	}
	return nil
}
