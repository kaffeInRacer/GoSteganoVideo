package handler

import (
	"fmt"
	"io/fs"
	"kaffein/lib/caesarcipher"
	"kaffein/lib/steganography"
	"kaffein/lib/transposecipher"
	"kaffein/utils"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	baseDir = "storage/stego_video"
)

func IndexEncode(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"assets": utils.LoadAssets(nil),
	}

	switch r.Method {
	case http.MethodGet:
		homeEncode(w, r, data)
	case http.MethodPost:
		handleEncode(w, r, data)
	default:
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
	}
}

func homeEncode(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	utils.LoadTemplate(w, "encode", data)
}

func handleEncode(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		file = nil
	} else {
		defer file.Close()
	}

	keyShifter, _ := strconv.Atoi(r.FormValue("keyShifter"))
	keyTranspose, _ := strconv.Atoi(r.FormValue("keyTranspose"))

	RequestEncryptInput := &struct {
		Alphabet     string                `validate:"required,min=140" label:"K1"`
		KeyShifter   int                   `validate:"numeric,required,gte=1" label:"K2"`
		KeyTranspose int                   `validate:"numeric,required,gte=1" label:"K3"`
		Message      string                `validate:"required" label:"Message"`
		File         *multipart.FileHeader `validate:"required" label:"Video"`
	}{
		Alphabet:     r.FormValue("keyAlphabet"),
		KeyShifter:   keyShifter,
		KeyTranspose: keyTranspose,
		Message:      r.FormValue("message"),
		File:         handler,
	}

	validator := utils.NewValidation()
	validatorInput := validator.ValidateInputVideo(RequestEncryptInput, handler)
	data["formValues"] = RequestEncryptInput

	if len(validatorInput) > 0 {
		data["validationError"] = validatorInput
		utils.LoadTemplate(w, "encode", data)
		return
	}

	path, err := utils.SaveFileToDirectory(file, baseDir, handler.Filename)
	if err != nil {
		utils.HandleError(w, data, "Failed to save file: "+err.Error(), "encode")
		return
	}

	encryptedCaesar, err := caesarcipher.NewCaesarCipher().Encrypt(
		RequestEncryptInput.Message,
		RequestEncryptInput.Alphabet,
		RequestEncryptInput.KeyShifter,
	)
	if err != nil {
		utils.HandleError(w, data, "Failed to encrypt message: "+err.Error(), "encode")
		return
	}

	encryptedTranspose := transposecipher.NewTranspose().Encrypt(
		encryptedCaesar,
		6,
		RequestEncryptInput.KeyTranspose,
	)

	videoEmbedded := steganography.NewVideoSteganoGraphy()
	_, filename, err := videoEmbedded.Encode(
		path,
		baseDir,
		handler.Filename,
		"FFV1",
		"avi",
		encryptedTranspose,
	)
	if err != nil {
		utils.HandleError(w, data, "Failed to encode video: "+err.Error(), "encode")
		return
	}

	data["success"] = fmt.Sprintf("/d/%s", filename)
	data["assets"] = utils.LoadAssets(map[string][]string{"js": {"/assets/js/blob.js"}})

	delete(data, "formValues")
	utils.LoadTemplate(w, "encode", data)
}

func DownloadEncryptFile(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/d/")
	filePath := filepath.Join(baseDir, "final", filename)
	log.Println(filePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, filePath)

	defer func() {
		err := filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())) == strings.TrimSuffix(filename, filepath.Ext(filename)) {
				if err := os.Remove(path); err != nil {
					fmt.Printf("Failed to delete %s: %v\n", path, err)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error walking the path %v: %v\n", baseDir, err)
		}
	}()
}
