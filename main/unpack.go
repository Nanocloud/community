package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func UnpackGo(sourcefile string) {

	file, err := os.Open(sourcefile)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	var fileReader io.ReadCloser = file

	// just in case we are reading a tar.gz file, add a filter to handle gzipped file
	if strings.HasSuffix(sourcefile, ".gz") {
		if fileReader, err = gzip.NewReader(file); err != nil {

			fmt.Println(err)
			os.Exit(1)
		}
		defer fileReader.Close()
	}

	tarBallReader := tar.NewReader(fileReader)

	// Extracting tarred files

	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}

		// get the individual filename and extract to the current directory
		filename := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
		case tar.TypeReg:
			// handle normal file
			fmt.Println("Untarring :", filename)
			if !strings.Contains(filename, "/") {
				writer, err := os.OpenFile("plugins/staging/"+filename, os.O_WRONLY|os.O_CREATE, os.FileMode(header.Mode))

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				io.Copy(writer, tarBallReader)

				//err = os.Chmod(filename, os.FileMode(header.Mode))

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				writer.Close()
				log.Println("filename: ", filename)
				log.Println("conf.StagDir: ", conf.StagDir)
				running_plugins = CreateEvent(running_plugins, filename, conf.StagDir+filename, sourcefile)
			}
		default:
			fmt.Printf("Unable to untar type : %c in file %s", header.Typeflag, filename)
		}
	}

}

func DeleteOldFront(sourcefile string) {
	file, err := os.Open(sourcefile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	var fileReader io.ReadCloser = file
	// just in case we are reading a tar.gz file, add a filter to handle gzipped file
	if strings.HasSuffix(sourcefile, ".gz") {
		if fileReader, err = gzip.NewReader(file); err != nil {

			fmt.Println(err)
			os.Exit(1)
		}
		defer fileReader.Close()
	}
	tarBallReader := tar.NewReader(fileReader)
	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}
		filename := header.Name
		switch header.Typeflag {
		case tar.TypeDir:
			// handle directory
			fmt.Println("Deleting directory :", "../front/"+filename)
			err = os.RemoveAll("../front/" + filename) // or use 0755 if you prefer
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case tar.TypeReg:
		default:
			fmt.Printf("Unable to delete type : %c in file %s", header.Typeflag, filename)
		}
	}
}

func UnpackFront(sourcefile string) {
	file, err := os.Open(sourcefile)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	var fileReader io.ReadCloser = file

	// just in case we are reading a tar.gz file, add a filter to handle gzipped file
	if strings.HasSuffix(sourcefile, ".gz") {
		if fileReader, err = gzip.NewReader(file); err != nil {

			fmt.Println(err)
			os.Exit(1)
		}
		defer fileReader.Close()
	}

	tarBallReader := tar.NewReader(fileReader)

	// Extracting tarred files

	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}

		// get the individual filename and extract to the current directory
		filename := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			// handle directory
			fmt.Println("Creating directory :", "../front/"+filename)
			err = os.MkdirAll("../front/"+filename, os.FileMode(header.Mode)) // or use 0755 if you prefer

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		case tar.TypeReg:
			// handle normal file
			fmt.Println("Untarring :", filename)
			writer, err := os.OpenFile("../front/"+filename, os.O_WRONLY|os.O_CREATE, os.FileMode(header.Mode))

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			io.Copy(writer, tarBallReader)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			writer.Close()
		default:
			fmt.Printf("Unable to untar type : %c in file %s", header.Typeflag, filename)
		}
	}

}
