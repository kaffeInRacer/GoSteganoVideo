package utils

import (
	"fmt"
	"html/template"
	"io"
	"kaffein/assets/templates"
	"kaffein/config"
	"net/http"
	"os"
	"path/filepath"
)

func LoadTemplate(w http.ResponseWriter, nameFile string, data map[string]interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmplBytes, err := templates.TemplatesFiles.ReadFile(nameFile + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New(nameFile).Parse(string(tmplBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleError(w http.ResponseWriter, data map[string]interface{}, message string, view string) {
	vErrors := make(map[string]string)
	vErrors["alert"] = message
	data["validationError"] = vErrors
	LoadTemplate(w, view, data)
}

func SaveFileToDirectory(file io.Reader, dirName, fileName string) (string, error) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if err = os.MkdirAll(dirName, os.ModePerm); err != nil {
			return "", err
		}
	}

	filePath := filepath.Join(dirName, fileName)
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, file); err != nil {
		return "", err
	}

	return filePath, nil
}

func LoadAssets(st map[string][]string) map[string][]string {
	cfg := config.Config{}

	based := map[string][]string{
		"css": {
			cfg.Base + "/assets/css/bootstrap.min.css",
			cfg.Base + "/assets/css/main.css",
			cfg.Base + "/assets/css/font-awesome/all.min.css",
		},
		"js": {
			cfg.Base + "/assets/js/bootstrap.bundle.min.js",
			cfg.Base + "/assets/js/main.js",
		},
	}

	for key, paths := range st {
		for _, path := range paths {
			based[key] = append(based[key], cfg.Base+path)
		}
	}

	return based
}

func CleanTempDir(directory string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return nil
	}

	err := os.RemoveAll(directory)
	if err != nil {
		return fmt.Errorf("failed to delete temp directory: %w", err)
	}

	return nil
}
