package provider

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Provider struct {
	languages map[string]*LanguageModel
}

func New() *Provider {
	return &Provider{
		languages: make(map[string]*LanguageModel),
	}
}

func (p *Provider) UpdateFromFile(path string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}

	base := new(BaseModel)
	if err = yaml.NewDecoder(f).Decode(base); err != nil {
		return
	}
	f.Seek(0, 0)

	switch base.Type {

	case ProviderLanguage:
		lang := new(LanguageModel)
		if err = yaml.NewDecoder(f).Decode(lang); err != nil {
			return
		}
		p.languages[lang.Id] = lang
	}

	return
}

func (p *Provider) GetLanguages() map[string]*LanguageModel {
	return p.languages
}
