package provider

type ProviderType string

const (
	ProviderLanguage  ProviderType = "language"
	ProviderFramework ProviderType = "framework"
	ProviderToolkit   ProviderType = "toolkit"
	ProviderIde       ProviderType = "ide"
)

// BaseModel structure which all provider models
// have in common.
//
// Properties of this struct should not be
// heritated to other provider structs because
// this may lead to some parsing inconsistencies
// for some de- and encoder modules.
type BaseModel struct {
	Id         string       `json:"id"`
	Name       string       `json:"name"`
	Type       ProviderType `json:"type"`
	Website    string       `json:"website"`
	Repository string       `json:"repository"`
}

// ScaleModel wraps a chart scale value.
type ScaleModel struct {
	Scale float32 `json:"scale"`
}

// PackageModel describes a langauge package/
// library/dependency object.
type PackageModel struct {
	Name       string `json:"name"`
	Repository string `json:"repository"`
	Website    string `json:"website"`
}

// ProductModel describes product information.
type ProductModel struct {
	Name    string `json:"name"`
	Website string `json:"website"`
}

// SnippetModel wraps information about
// a language snippet.
type SnippetModel struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

// LanguageModel wraps a language provider model.
type LanguageModel struct {
	Id         string       `json:"id"`
	Name       string       `json:"name"`
	Type       ProviderType `json:"type"`
	Website    string       `json:"website"`
	Repository string       `json:"repository"`

	Paradigm        []string        `json:"paradigm"`
	Difficulty      *ScaleModel     `json:"difficulty"`
	Popularity      *ScaleModel     `json:"popularity"`
	Targets         []string        `json:"targets"`
	PopularPackages []*PackageModel `json:"popular_packages"`
	PopularProducts []*ProductModel `json:"popular_products"`
	Snippets        []*SnippetModel `json:"snippets"`

	Community *struct {
		Scale      float32  `json:"scale"`
		Properties []string `json:"properties"`
	} `json:"community"`
}
