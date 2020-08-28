package provider

type ProviderType string

const (
	ProviderLanguage  ProviderType = "language"
	ProviderFramework ProviderType = "framework"
	ProviderToolkit   ProviderType = "toolkit"
	ProviderIde       ProviderType = "ide"
)

type BaseModel struct {
	Id         string       `json:"id"`
	Name       string       `json:"name"`
	Type       ProviderType `json:"type"`
	Website    string       `json:"website"`
	Repository string       `json:"repository"`
}

type ScaleModel struct {
	Scale float32 `json:"scale"`
}

type PackageModel struct {
	Name       string `json:"name"`
	Repository string `json:"repository"`
	Website    string `json:"website"`
}

type ProductModel struct {
	Name    string `json:"name"`
	Website string `json:"website"`
}

type SnippetModel struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

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
