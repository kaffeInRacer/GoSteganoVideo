package steganography

import (
	"errors"
	"fmt"
	"gocv.io/x/gocv"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type videoSteganography struct {
	videoPath  string
	codec      string
	saveDir    string
	outputName string
	typeVideo  string
	message    string
	result
}

type result struct {
	outputNameVideo      string
	outputPathVideo      string
	outputNameAudio      string
	outputPathAudio      string
	outputFinalPathVideo string
}

func NewVideoSteganoGraphy() *videoSteganography {
	return &videoSteganography{}
}

// Encode returning PathFinal Video, Name Final Video, Error
func (v *videoSteganography) Encode(targetVideo, saveDir, outputName, codec, typeFile, message string) (string, string, error) {
	v.videoPath = targetVideo
	v.codec = codec
	v.saveDir = saveDir
	v.outputName = outputName
	v.typeVideo = typeFile
	v.message = message

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

func (v *videoSteganography) writeVideo() error {
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

	msg := stringToBinary(v.message)
	lengths2b := len(msg) / 8

	maxMessage := ((width * height * 3 * int(video.Get(gocv.VideoCaptureFrameCount))) / 8) - 1

	if maxMessage <= lengths2b {
		return errors.New("message too long")
	}
	msg = append(msg, stringToBinary(string(emmit))...)
	lengths2b = len(msg)

	msgIndex := 0
	for {
		if !video.Read(&frame) {
			break
		}
		if frame.Empty() {
			continue
		}

		for y := 0; y < height && msgIndex < lengths2b; y++ {
			for x := 0; x < width && msgIndex < lengths2b; x++ {
				pixel := frame.GetVecbAt(y, x)
				b, g, r := pixel[0], pixel[1], pixel[2]

				// Modify the least significant bit (LSB) of each color channel
				if msgIndex < lengths2b {
					r = (r & 0xFE) | msg[msgIndex]
					msgIndex++
				}
				if msgIndex < lengths2b {
					g = (g & 0xFE) | msg[msgIndex]
					msgIndex++
				}
				if msgIndex < lengths2b {
					b = (b & 0xFE) | msg[msgIndex]
					msgIndex++
				}

				// Set the modified pixel values
				frame.SetUCharAt(y, x*3, b)
				frame.SetUCharAt(y, x*3+1, g)
				frame.SetUCharAt(y, x*3+2, r)
			}
		}

		writer.Write(frame)
	}

	return nil
}

func (v *videoSteganography) extractAudio() (bool, error) {
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

func (v *videoSteganography) mergeAudioVideo() error {
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
