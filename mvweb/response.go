package mvweb

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/schollz/progressbar/v2"
)

type QueryInfo struct {
	FilmlisteTimestamp string `json:"filmlisteTimestamp"`
	ResultCount        uint   `json:"resultCount"`
	SearchEngineTime   string `json:"searchEngineTime"`
	TotalResults       uint   `json:"totalResults"`
}

type Result struct {
	Channel    string `json:"channel"`
	Topic      string `json:"topic"`
	Title      string `json:"title"`
	Duration   uint   `json:"duration"`
	ID         string `json:"id"`
	Size       uint   `json:"size"`
	Timestamp  uint   `json:"timestamp"`
	UrlVideo   string `json:"url_video"`
	UrlVideoHD string `json:"url_video_hd"`
	UrlVideoSD string `json:"url_video_low"`
}

type Response struct {
	Result struct {
		QueryInfo QueryInfo `json:"queryInfo"`
		Results   []Result  `json:"results"`
	} `json:"result"`
}

func (r *Result) DownloadTo(path string) (int64, time.Duration, error) {
	start := time.Now()
	n, err := r.download(r.UrlVideoHD, path)
	end := time.Now()
	return n, end.Sub(start), err
}

var filenameClean = regexp.MustCompile("[^a-zA-Z0-9äöüßÄÖÜ.()_+ -]")

func (r *Result) Filename() string {
	url := r.url()
	ext := filepath.Ext(url)
	name := filenameClean.ReplaceAllString(r.Title, "_")
	name = strings.Trim(name, " _-()+.")
	return name + ext

}

func (r *Result) url() string {
	return firstNonEmpty(r.UrlVideoHD, r.UrlVideoSD, r.UrlVideo)
}

func firstNonEmpty(s ...string) string {
	for _, e := range s {
		if e != "" {
			return e
		}
	}
	return ""
}

func (r *Result) download(url, path string) (int64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	f, err := os.Create(path)
	if err != nil {
		return -1, err
	}

	bar := progressbar.NewOptions(
		int(resp.ContentLength),
		progressbar.OptionSetBytes(int(resp.ContentLength)),
	)
	var out io.Writer = f
	if isatty.IsTerminal(os.Stdout.Fd()) {
		out = io.MultiWriter(f, bar)
	}

	n, err := io.Copy(out, resp.Body)
	f.Close()
	return n, err
}
