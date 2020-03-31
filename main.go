package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/floj/mvweb-dl/conf"
	"github.com/floj/mvweb-dl/history"
	"github.com/floj/mvweb-dl/mvweb"
)

// "wurst":{
//	"queries":[{"fields":["channel"],"query":"ard"},{"fields":["topic"],"query":"sendung"}],
// "sortBy":"timestamp","sortOrder":"desc","future":false,"offset":0,"size":15},

func main() {
	confFile := flag.String("config", "", "Config file to use")
	skipDownload := flag.Bool("no-download", false, "Skip download, but add to history file")
	dryRun := flag.Bool("dry", false, "Don't download or add to history")

	flag.Parse()
	if *confFile == "" {
		log.Fatal("No config file specified")
	}

	configs, err := conf.Load(*confFile)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range configs {
		log.Printf("Running config '%s'", c.Name)
		req := c.ToMvwebRequest()
		resp, err := req.Run()
		if err != nil {
			log.Printf("Could not query config %s: %v", c.Name, err)
			continue
		}

		log.Printf("  results: %+v", resp.Result.QueryInfo)
		err = processResults(c, resp.Result.Results, *skipDownload, *dryRun)
		if err != nil {
			log.Printf("Could not process results: %v", err)
			continue
		}
	}
}
func processResults(c conf.Config, results []mvweb.Result, skipDownload, dryRun bool) error {
	hist, err := c.History()
	if err != nil {
		return err
	}
	defer hist.Close()
	for _, r := range results {
		processResult(c, r, hist, skipDownload, dryRun)
	}
	return nil
}

func processResult(c conf.Config, r mvweb.Result, hist history.History, skipDownload, dryRun bool) {
	log.Printf("  checking '%s' (ID: %s)", r.Title, r.ID)
	if match, filter := c.Matches(r); !match {
		log.Printf("    skipping - %s", filter)
		return
	}

	if hist.Exists(r.ID) {
		log.Printf("    skipping - found in history")
		return
	}

	filename := filepath.Join(c.DownloadDir, r.Filename())
	if existsFile(filename) {
		hist.Add(r.ID, r.Title)
		log.Printf("    skipping - file already exists")
		return
	}
	log.Printf("    downloading to '%s'", filename)
	if dryRun {
		log.Println("    DRY RUN - skipping history and download")
		return
	}

	if skipDownload {
		hist.Add(r.ID, r.Title)
		log.Println("    SKIP DOWNLOAD - adding to history but skipping download")
		return
	}

	bytes, duration, err := r.DownloadTo(filename)
	if err != nil {
		log.Printf("ERROR download failed %v", err)
		return
	}
	log.Printf("    finished after %s, %s", duration, formatBytes(bytes))
	hist.Add(r.ID, r.Title)

}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func existsFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
