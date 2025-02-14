package MSE_PSNR

import (
	"errors"
	"fmt"
	"gocv.io/x/gocv"
	"math"
	"os"
)

func CallMSEandPSNR(originalVideo, modifiedVideo, saveResultAt string) error {
	cv1, err := gocv.VideoCaptureFile(originalVideo)
	if err != nil {
		return fmt.Errorf("gagal membuka video asli: %w", err)
	}
	defer cv1.Close()

	cv2, err := gocv.VideoCaptureFile(modifiedVideo)
	if err != nil {
		return fmt.Errorf("gagal membuka video hasil modifikasi: %w", err)
	}
	defer cv2.Close()

	cv1Count := int(cv1.Get(gocv.VideoCaptureFrameCount))
	cv2Count := int(cv2.Get(gocv.VideoCaptureFrameCount))
	if cv1Count != cv2Count {
		return errors.New("jumlah frame pada kedua video berbeda")
	}

	cv1Width, cv1Height := int(cv1.Get(gocv.VideoCaptureFrameWidth)), int(cv1.Get(gocv.VideoCaptureFrameHeight))
	cv2Width, cv2Height := int(cv2.Get(gocv.VideoCaptureFrameWidth)), int(cv2.Get(gocv.VideoCaptureFrameHeight))
	if cv1Width != cv2Width || cv1Height != cv2Height {
		return errors.New("resolusi video berbeda")
	}

	cv1Frame := gocv.NewMat()
	cv2Frame := gocv.NewMat()
	defer cv1Frame.Close()
	defer cv2Frame.Close()

	rows, cols := cv1Height, cv1Width
	var totalMSE, totalPSNR float64
	indexFrame := 0
loop:
	for i := 0; i < cv1Count; i++ {
		if ok := cv1.Read(&cv1Frame); !ok || cv1Frame.Empty() {
			break
		}
		if ok := cv2.Read(&cv2Frame); !ok || cv2Frame.Empty() {
			break
		}

		var frameSquaredError float64
		channels := cv1Frame.Channels()
		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				for c := 0; c < channels; c++ {
					val1 := cv1Frame.GetUCharAt(y, x*channels+c)
					val2 := cv2Frame.GetUCharAt(y, x*channels+c)

					diff := math.Pow(float64(val1)-float64(val2), 2)
					frameSquaredError += diff
				}
			}
		}

		// Menghitung MSE
		mse := frameSquaredError / float64(rows*cols)
		if mse == 0 {
			break loop
		} else {
			psnr := 10 * math.Log10(math.Pow(255, 2)/mse)
			totalPSNR += psnr
		}

		totalMSE += mse
		indexFrame++
	}

	averageMSE := totalMSE / float64(indexFrame)
	averagePSNR := totalPSNR / float64(indexFrame)
	fmt.Printf("Rata-rata MSE: %.4f\n", averageMSE)
	fmt.Printf("Rata-rata PSNR: %.4f dB\n", averagePSNR)
	fmt.Println()

	content := fmt.Sprintf("MSE: %.4f PSNR: %.4f\n", averageMSE, averagePSNR)
	path := fmt.Sprintf("%s/MSE_PSNR.txt", saveResultAt)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}
