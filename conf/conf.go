package conf

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/floj/mvweb-dl/history"
	"github.com/floj/mvweb-dl/mvweb"
	"gopkg.in/yaml.v2"
)

type Query struct {
	Channel    string `yaml:"channel"`
	Topic      string `yaml:"topic"`
	MaxResults uint   `yaml:"max_results"`
}

type Skipper struct {
	Condition string `yaml:"condition"`
	Value     string `yaml:"value"`
}

type Config struct {
	Name        string    `yaml:"name"`
	DownloadDir string    `yaml:"download_to"`
	HistoryFile string    `yaml:"history_file"`
	Query       Query     `yaml:"query"`
	SkipIf      []Skipper `yaml:"skip_if"`
}

func (c *Config) History() (history.History, error) {
	return history.Load(c.HistoryFile)
}

func (q *Query) toMvwebQuery() []mvweb.Query {
	qq := make([]mvweb.Query, 0, 2)
	if q.Channel != "" {
		qq = append(qq, mvweb.NewQuery("channel", q.Channel))
	}
	if q.Topic != "" {
		qq = append(qq, mvweb.NewQuery("topic", q.Topic))
	}
	return qq
}

func (f *Skipper) Skip(r mvweb.Result) bool {
	switch f.Condition {
	case "title_contains":
		return strings.Contains(r.Title, f.Value)
	case "shorter_then":
		midDuration, err := time.ParseDuration(f.Value)
		if err != nil {
			panic(fmt.Errorf("Could not parse duration %s: %w", f.Value, err))
		}
		actDuration, err := time.ParseDuration(fmt.Sprintf("%ds", r.Duration))
		if err != nil {
			panic(fmt.Errorf("Could not parse duration %d: %w", r.Duration, err))
		}
		return actDuration < midDuration
	default:
		panic(fmt.Errorf("Unknown filter type: %s", f.Condition))
	}
}

func (f *Skipper) String() string {
	return fmt.Sprintf("%s(%s)", f.Condition, f.Value)
}

func (c *Config) Matches(r mvweb.Result) (bool, *Skipper) {
	for _, f := range c.SkipIf {
		if f.Skip(r) {
			return false, &f
		}
	}
	return true, nil

}

func (c *Config) ToMvwebRequest() mvweb.Request {
	return mvweb.NewRequest(c.Query.MaxResults, c.Query.toMvwebQuery()...)
}

func Load(path string) ([]Config, error) {
	var conf []Config
	f, err := os.Open(path)
	if err != nil {
		return conf, err
	}
	defer f.Close()
	dec, err := decoderFor(path, f)
	if err != nil {
		return conf, err
	}
	err = dec.Decode(&conf)
	if err != nil {
		return conf, err
	}
	for _, c := range conf {
		if c.Query.MaxResults == 0 {
			c.Query.MaxResults = 100
		}
	}
	return conf, nil
}

type confDecoder interface {
	Decode(v interface{}) error
}

func decoderFor(path string, r io.Reader) (confDecoder, error) {
	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		return json.NewDecoder(r), nil
	case ".yml":
		return yaml.NewDecoder(r), nil
	case ".yaml":
		return yaml.NewDecoder(r), nil
	}
	return nil, fmt.Errorf("No config decoder registered for %s", ext)
}
