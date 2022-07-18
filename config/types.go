package config

type Config struct {
	Paths []Path `json:"paths"`
}

type Path struct {
	Close Location `json:"close"`
	Far Location `json:"far"`
	Links []Link `json:"links"`
}

type Location struct {
	Name string `json:"name"`
}

type Link struct {
	Radio Radio `json:"radio"`
	Label uint32 `json:"label"`
}

type Radio struct {
	Ip string `json:"ip"`
	User string `json:"user"`
	Password string        `json:"password"`
	Polarity RadioPolarity `json:"polarity"`
	Type     RadioType     `json:"type"`
}

type RadioPolarity string
const PolarityClose = "close"
const PolarityFar = "far"

type RadioType string
const TypeAFLTU = "afltu"
const TypeAF = "af"