package utils

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translator "github.com/go-playground/validator/v10/translations/en"
	"mime/multipart"
	"reflect"
)

type Validation struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewValidation() *Validation {
	translator := en.New()
	uni := ut.New(translator, translator)
	trans, _ := uni.GetTranslator("en")
	validate := validator.New()
	en_translator.RegisterDefaultTranslations(validate, trans)

	// Register custom tag name function
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("label")
		return name
	})
	//
	//// Register custom translations
	//validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
	//	return ut.Add("required", "{0} harus diisi", true)
	//}, func(ut ut.Translator, fe validator.FieldError) string {
	//	t, _ := ut.T("required", fe.Field())
	//	return t
	//})

	return &Validation{validate, trans}
}

func (v *Validation) Struct(s interface{}) map[string]string {
	errors := make(map[string]string)
	err := v.validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.StructField()] = err.Translate(v.trans)
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (v *Validation) ValidateInputVideo(reqInput interface{}, handler *multipart.FileHeader) map[string]string {
	allowedTypes := map[string]bool{
		"video/mp4":                true,
		"video/3gpp":               false,
		"video/avi":                true,
		"video/x-msvideo":          false,
		"video/quicktime":          false,
		"video/ogg":                false,
		"video/webm":               false,
		"video/x-ms-wmv":           false,
		"video/x-matroska":         false,
		"application/octet-stream": false,
	}

	vErrors := make(map[string]string)
	reqInputType := reflect.TypeOf(reqInput)
	reqInputValue := reflect.ValueOf(reqInput)

	if reqInputType.Kind() != reflect.Ptr || reqInputType.Elem().Kind() != reflect.Struct {
		panic("Input must be a pointer to a struct")
		return vErrors
	}

	// Validate the struct
	validationErrors := v.Struct(reqInput)
	for field, err := range validationErrors {
		vErrors[field] = err
	}

	// Check each field in the struct
	reqInputElem := reqInputValue.Elem()
	for i := 0; i < reqInputElem.NumField(); i++ {
		field := reqInputElem.Type().Field(i)
		if field.Type == reflect.TypeOf(&multipart.FileHeader{}) {
			if handler != nil && !allowedTypes[handler.Header.Get("Content-Type")] {
				vErrors[field.Name] = "Invalid file type"
			}
		}
	}

	return vErrors
}
