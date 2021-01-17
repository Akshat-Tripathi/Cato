package cate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/* This file contains functions related to downloading and parsing cate webpages
 */

func get(url string) ([]byte, error) {
	resp, err := login(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return html, nil
}

func download(url, location string) error {
	html, err := get(url)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(location, html, 0644)
}

func downloadHome() (*goquery.Document, error) {
	home, err := get(cateURL)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(bytes.NewBuffer(home))
}

//DownloadTimeTable needs info to be initialised before being called
func downloadTimeTable() (*goquery.Document, error) {
	currentYear := getAcademicYear()

	timetable, err := get(fmt.Sprintf(timeTableURL, currentYear, info.Term,
		info.Code, info.Shortcode))
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(bytes.NewBuffer(timetable))
}

//DownloadModule tries to download all tasks in a module in the appropriate folder
//It stops downloading as soon as it fails once
func downloadModule(module *Module) error {
	var err error
	location := "files/" + formatName(module.Name) + "/"
	for _, task := range module.Tasks {
		for _, file := range task.Files {
			err = download(cateURL+"/"+file, location+formatName(task.Name)+".pdf")
			if err != nil {
				log.Println("Error downloading module: " + module.Name)
				return err
			}
		}
	}
	return nil
}

func formatName(name string) string {
	return strings.ReplaceAll(name, ":", "")
}
