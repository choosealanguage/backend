// Package provider provides models and
// functionalities to parse provider
// information from file and update them
// to internal maps.
package provider

import (
	"os"

	"gopkg.in/yaml.v2"
)

// ProviderHost wraps compiled maps of provider
// model instances.
type ProviderHost struct {
	languages map[string]*LanguageModel
}

// New returns a new instance of Provider.
func New() *ProviderHost {
	return &ProviderHost{
		languages: make(map[string]*LanguageModel),
	}
}

// UpdateFromFile tries to read the passed file
// and pases it using a YAML decoder. Depending
// on the provider type, the provider is then
// updated to the maps of the ProviderHost.
func (p *ProviderHost) UpdateFromFile(path string) (err error) {
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

// GetLanguages returns the provider map of
// registered language models.
func (p *ProviderHost) GetLanguages() map[string]*LanguageModel {
	return p.languages
}
