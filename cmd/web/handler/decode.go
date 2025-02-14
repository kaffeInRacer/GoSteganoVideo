package handler

import (
	"kaffein/lib/caesarcipher"
	"kaffein/lib/steganography"
	"kaffein/lib/transposecipher"
	"kaffein/utils"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

func IndexDecode(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"assets": utils.LoadAssets(nil),
	}

	switch r.Method {
	case http.MethodGet:
		home(w, r, data)
	case http.MethodPost:
		handlePost(w, r, data)
	default:
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
	}
}

func home(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	utils.LoadTemplate(w, "decode", data)
}

func handlePost(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		file = nil
	} else {
		defer file.Close()
	}

	keyShifter, _ := strconv.Atoi(r.FormValue("keyShifter"))
	keyTranspose, _ := strconv.Atoi(r.FormValue("keyTranspose"))

	RequestDecryptInput := &struct {
		Alphabet     string                `validate:"required,min=140" label:"K1"`
		KeyShifter   int                   `validate:"numeric,required,gte=1" label:"K2"`
		KeyTranspose int                   `validate:"numeric,required,gte=1" label:"K3"`
		File         *multipart.FileHeader `validate:"required" label:"Video"`
	}{
		Alphabet:     r.FormValue("keyAlphabet"),
		KeyShifter:   keyShifter,
		KeyTranspose: keyTranspose,
		File:         handler,
	}

	validator := utils.NewValidation()
	validatorInput := validator.ValidateInputVideo(RequestDecryptInput, handler)
	data["formValues"] = RequestDecryptInput

	if len(validatorInput) > 0 {
		data["validationError"] = validatorInput
		utils.LoadTemplate(w, "decode", data)
		return
	}

	path, err := utils.SaveFileToDirectory(file, baseDir, handler.Filename)
	if err != nil {
		utils.HandleError(w, data, "Failed to save file: "+err.Error(), "decode")
		return
	}

	videoEmbedded := steganography.NewVideoSteganoGraphy()
	cipherText, err := videoEmbedded.Decode(path)
	if err != nil {
		utils.HandleError(w, data, "Failed to decrypt message: "+err.Error(), "decode")
		return
	}

	decryptedTranspose := transposecipher.NewTranspose().Decrypt(
		cipherText,
		6,
		RequestDecryptInput.KeyTranspose,
	)

	plaintext, err := caesarcipher.NewCaesarCipher().Decrypt(
		decryptedTranspose,
		RequestDecryptInput.Alphabet,
		RequestDecryptInput.KeyShifter,
	)

	if err != nil {
		utils.HandleError(w, data, "Failed to decrypt message: "+err.Error(), "decode")
		return
	}

	data["plainText"] = plaintext
	delete(data, "formValues")
	utils.LoadTemplate(w, "decode", data)

	go func() {
		if err := os.Remove(path); err != nil {
			utils.HandleError(w, data, "Failed to decrypt message: "+err.Error(), "decode")
			return
		}
	}()
}
