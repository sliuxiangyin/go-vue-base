package ai

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"sync"
)

type Openai struct {
	key      string
	baseUrl  string
	once     sync.Once
	instance openai.Client
	err      error
}

func (f *Openai) Key() string {
	return f.key
}

func (f *Openai) BaseUrl() string {
	return f.baseUrl
}

func NewOpenai(key string, baseUrl string) *Openai {
	return &Openai{
		key:     key,
		baseUrl: baseUrl,
	}
}

func (f *Openai) Get() (openai.Client, error) {
	f.once.Do(func() {
		f.instance = openai.NewClient(
			option.WithAPIKey(f.key),
			option.WithBaseURL(f.baseUrl),
		)
	})
	return f.instance, f.err
}
