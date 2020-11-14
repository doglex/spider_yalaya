package main

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/levigross/grequests"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

const TemplateGallery = "https://www.yalayi.com/gallery/%v.html"
const TemplateGalleryPhoto = "https://img.yalayi.net/img/gallery/%v/z%v.jpg"
const Referer = "https://www.yalayi.com/"

func main() {
	for i := 1; i < 1024; i++ {
		FetchOne(i)
	}
}

func FetchOne(i int) {
	url := fmt.Sprintf(TemplateGallery, i)
	fmt.Println(url)
	ro := grequests.RequestOptions{
		RequestTimeout: time.Second * 90,
		Headers: map[string]string{
			"Referer":    Referer,
			"User-Agent": GetUA(),
		}}
	resp, err := grequests.Get(url, &ro)
	if err != nil {
		log.Println("Error", "get uri fail:", url)
		return
	}
	respStr := resp.String()
	if strings.Contains(respStr, "<title>404 Not Found") {
		fmt.Println(i,"=>",url, "404-NOT-FOUND")
		return
	}
	fmt.Println(i,"=>", url, "GOT")
	doc := soup.HTMLParse(respStr)
	for j := 1; j < 5; j++ {
		title := strings.ReplaceAll(doc.Find("title").FullText(), " ", "")
		title = strings.Split(title, "-")[0]
		title = fmt.Sprintf("%v.%v", i, title)
		down := fmt.Sprintf(TemplateGalleryPhoto, i, j)
		go DownloadFile(title, down, Referer)
	}
}

func GetUA() string {
	return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36 Edg/84.0.522.63"
}

func DownloadFile(folder string, url string, referer string) {
	defer DeferRecover()
	_ = os.MkdirAll(folder, os.ModePerm)
	fileName := LegalPathName(url)
	fullFileName := path.Join(folder, fileName)
	// 不要启用，有些图片没有下载完整
	//if FileExists(fullFileName) {
	//	return
	//}
	ro := grequests.RequestOptions{
		RequestTimeout: time.Second * 90,
		Headers: map[string]string{
			"Referer":    referer,
			"User-Agent": GetUA(),
		}}
	resp, err := grequests.Get(url, &ro)
	if err != nil {
		log.Println("Error", "get file fail:", url)
		return
	}
	if err := resp.DownloadToFile(fullFileName); err != nil {
		log.Println("Error", "save file fail:", err)
		return
	}
	fi, err := os.Stat(fullFileName)
	if err == nil {
		fileSize := fi.Size() / 1000
		if fileSize < 1 {
			fmt.Println("Warn", " file size too small: ", fullFileName, "|", fileSize, "kB", "|", url)
		}
		fmt.Println("GOT <- ", fullFileName, "|", fileSize, "kB", "|", url)
	} else {
		fmt.Println("Warn", " can not get file size : ", fullFileName, "|", url)
	}
}
