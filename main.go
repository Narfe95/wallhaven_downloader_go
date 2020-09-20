package main

import (
	"encoding/json"
	"flag"
	"github.com/lucasjones/reggen"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var ApiUrl = "https://wallhaven.cc/api/v1/search?"
var SearchQuery = "ratios=32x9&purity=110&categories=111"

type ApiResponse struct {
	Data []struct {
		Path       string
	}
}

func main() {
	local, dir := ParseFlags()
	FileName := ""
	if local {
		FileName = SelectLocal(dir)
	} else {
		FileName = DownloadRandomImage(dir)
	}
	SetWallpaper(FileName)
}

func ParseFlags() (bool, string) {
	var LocalSet bool
	var Directory string
	flag.BoolVar(&LocalSet, "l", false, "Dont download a new wallpaper but set one already downloaded.")
	flag.StringVar(&Directory, "d", "/home/jonas/Pictures/script_wallpapers/", "Chose a different directory for downloading to or selecting from.")
	flag.Parse()
	return LocalSet, Directory
}

func SelectLocal(directory string) string {
	fp, _ := os.Open(directory)

	files, _ := fp.Readdir(-1)
	fp.Close()

	// Set a random seed and randomize a number between 0 and WallhavenResponse length
	rand.Seed(time.Now().UnixNano())
	RandNum := rand.Intn(len(files))

	return files[RandNum].Name()
}

func SetWallpaper(filename string) {

}

func DownloadRandomImage(directory string) string {
	// Generate a random seed which is 6 (six) characters long with alphanumeric characters
	seed, _ := reggen.Generate("^&seed=[a-zA-Z0-9]{6}", 6)

	// Create the full url with all the properties supplied
	FullUrl := ApiUrl+SearchQuery+seed
	response, _ := http.Get(FullUrl)

	// Read out the response from the request
	body, _ := ioutil.ReadAll(response.Body)

	FileName := ""

	var WallhavenResponse ApiResponse
	if json.Valid(body) {
		// Decode the json from the response
		_ = json.Unmarshal(body, &WallhavenResponse)

		// Set a random seed and randomize a number between 0 and WallhavenResponse length
		rand.Seed(time.Now().UnixNano())
		RandNum := rand.Intn(len(WallhavenResponse.Data))

		ImageUrl := WallhavenResponse.Data[RandNum].Path

		// Split the url by / and pick out the last part as filename
		Splits := strings.Split(ImageUrl, "/")
		FileName = Splits[len(Splits)-1]

		// Download image to /home/jonas/Pictures/script_wallpapers/
		DownloadFile(ImageUrl, directory+FileName)
	}

	return FileName
}

func DownloadFile(URL, fileName string) {
	//Get the response bytes from the url
	response, _ := http.Get(URL)

	defer response.Body.Close()

	//Create a empty file
	file, _ := os.Create(fileName)

	defer file.Close()

	//Write the bytes to the file
	_, _ = io.Copy(file, response.Body)
}