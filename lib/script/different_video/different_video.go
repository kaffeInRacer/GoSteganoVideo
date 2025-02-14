package different_video

import (
	"fmt"
	"gocv.io/x/gocv"
)

func CallDifferent(originalSource, modifiedSource string) error {
	cv1, err := gocv.VideoCaptureFile(originalSource)
	if err != nil {
		return fmt.Errorf("failed to open original video: %w", err)
	}
	defer cv1.Close()

	cv2, err := gocv.VideoCaptureFile(modifiedSource)
	if err != nil {
		return fmt.Errorf("failed to open modified video: %w", err)
	}
	defer cv2.Close()

	cv1Width := int(cv1.Get(gocv.VideoCaptureFrameWidth))
	cv1Height := int(cv1.Get(gocv.VideoCaptureFrameHeight))
	cv2Width := int(cv2.Get(gocv.VideoCaptureFrameWidth))
	cv2Height := int(cv2.Get(gocv.VideoCaptureFrameHeight))

	if cv1Width != cv2Width || cv1Height != cv2Height {
		return fmt.Errorf("videos have different dimensions")
	}

	writer, err := gocv.VideoWriterFile("different_gap.avi", "FFV1", cv1.Get(gocv.VideoCaptureFPS), cv1Width, cv1Height, true)
	if err != nil {
		return err
	}
	defer writer.Close()

	frame1 := gocv.NewMat()
	defer frame1.Close()
	frame2 := gocv.NewMat()
	defer frame2.Close()

	for {
		if !cv1.Read(&frame1) || !cv2.Read(&frame2) {
			break
		}

		for y := 0; y < frame1.Rows(); y++ {
			for x := 0; x < frame1.Cols(); x++ {
				CV1pixel := frame1.GetVecbAt(y, x)
				b1, g1, r1 := CV1pixel[0], CV1pixel[1], CV1pixel[2]

				CV2pixel := frame2.GetVecbAt(y, x)
				b2, g2, r2 := CV2pixel[0], CV2pixel[1], CV2pixel[2]

				if b1 != b2 {
					b1 = 0
				}
				if g1 != g2 {
					g1 = 0
				}
				if r1 != r2 {
					r1 = 0
				}

				frame1.SetUCharAt(y, x*3, b1)
				frame1.SetUCharAt(y, x*3+1, g1)
				frame1.SetUCharAt(y, x*3+2, r1)
			}
		}

		err := writer.Write(frame1)
		if err != nil {
			return err
		}
	}

	return nil
}
