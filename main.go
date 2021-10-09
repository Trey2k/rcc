package main

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if len(args) > 1 {
		dir = args[1]
	}

	err = remoteCompile(dir)
	if err != nil {
		switch err {
		case ErrNonDir:
			fmt.Println("Error: You must specify a directory.")
		case ErrNoConf:
			fmt.Println("The directory specified is missing a rcc.json file.")
		default:
			panic(err)
		}
	}
}

func remoteCompile(foldername string) error {

	fileInfo, err := os.Stat(foldername)
	if err != nil {
		return err
	}
	if !fileInfo.IsDir() {
		return ErrNonDir
	}

	_, err = os.Stat(fmt.Sprintf("%s/rcc.json", foldername))
	if err != nil {
		return ErrNoConf
	}

	filename := filepath.Base(foldername)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("tarball", filename)
	if err != nil {
		return err
	}

	err = Tar(foldername, part)
	if err != nil {
		return err
	}

	writer.Close()
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:%s/compile", Config.RCCServerIP, Config.RCCServerPort), body)
	if err != nil {
		return err
	}

	request.Header.Add("Auth", Config.RCCServerAuthToken)

	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	output := strings.Replace(resp.Header.Get("output"), "~n~", "\n", -1)

	fmt.Println(output)

	if resp.StatusCode == http.StatusOK {
		err = Untar(foldername, resp.Body)
		return err
	}
	return nil
}
