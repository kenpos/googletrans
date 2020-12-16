package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func genURL(tl string, sl string, q string) string {
	baseURL, err := url.Parse("https://translate.google.com/m?")
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	params := url.Values{}
	params.Add("q", q)
	params.Add("tl", tl)
	params.Add("sl", sl)
	baseURL.RawQuery = params.Encode()
	return baseURL.String()
}

func getUrl2Html(url string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows Phone 10.0; Android 6.0.1; Microsoft; Kumia 640 LTE) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.99 Safari/537.36 Edge/14.14390")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return string(body)
}

func getTranslateResult(html string) string {
	req := regexp.MustCompile(`class="result-container">(.*?)</div>`)
	s := req.FindString(html)
	s = strings.Replace(s, "class=", "", 1)
	s = strings.Replace(s, "result-container", "", 1)
	s = strings.Replace(s, "</div>", "", 1)
	return s
}

// UTF-8 から ShiftJIS
func utf82sjis(str string) (string, error) {
	iostr := strings.NewReader(str)
	rio := transform.NewReader(iostr, japanese.ShiftJIS.NewEncoder())
	ret, err := ioutil.ReadAll(rio)
	if err != nil {
		return "", err
	}
	return string(ret), err
}

func main() {
	fontPath := "C:/WINDOW/FONTS/BIZ-UDMINCHOM.TTC"
	if _, err := os.Stat(fontPath); err == nil {
		os.Setenv("FYNE_FONT", fontPath)
	}

	for {
		text, _ := clipboard.ReadAll()
		flag.Parse()

		rptext := strings.Replace(text, "\n", "", -1)
		rptext = strings.Replace(text, ":", "", -1)

		gURL := genURL("ja", "en", rptext)
		html := getUrl2Html(gURL)
		str := getTranslateResult(html)

		fmt.Println(str)
		fmt.Println("----------------------------------")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
	}
}
