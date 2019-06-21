package utils

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// ZipWriter is a function to zip directory
func ZipWriter(path string, filename string) {
	baseFolder := path

	// Get a Buffer to Write To
	outFile, err := os.Create(filename)
	CheckErr("", err)
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	AddFiles(w, baseFolder, "")

	CheckErr("", err)

	// Make sure to check the error on Close.
	err = w.Close()
	CheckErr("", err)
}

// AddFiles is a function to add files to zip file
func AddFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	CheckErr("", err)

	for _, file := range files {
		TimedPrintln("Adicionando arquivo: " + basePath + file.Name())
		if !file.IsDir() {
			log.Println("teste1")
			dat, err := ioutil.ReadFile(basePath + file.Name())
			log.Println("teste2")
			CheckErr("", err)
			log.Println("teste3")

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			CheckErr("", err)

			_, err = f.Write(dat)
			CheckErr("", err)
		} else if file.IsDir() {

			// Recurse
			newBase := basePath + file.Name() + "/"
			TimedPrintln("Recursing and Adding SubDir: " + file.Name())
			TimedPrintln("Recursing and Adding SubDir: " + newBase)

			AddFiles(w, newBase, file.Name()+"/")
		}
	}
}

// TimedPrintln is a function to print a message with the time
func TimedPrintln(message string) {
	t := time.Now()

	fmt.Printf("[%d/%d/%d às %02d:%02d:%02d] %s\n", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second(), message)
}

// CheckErr is a function to handle the erros of application
func CheckErr(message string, err error) {
	if err != nil {
		TimedPrintln(message)
		fmt.Printf("%+v", err)
		panic(err)
	}
}
