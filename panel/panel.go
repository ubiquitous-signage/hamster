package panel

import ("time")

type Panel struct {
	PanelHeader `bson:",inline"`
	Contents interface{} `json:"contents"`
}

type PanelHeader struct {
	Version  float64     `json:"version"`
	Type     string      `json:"type"`
	Title    string      `json:"title"`
	Category string      `json:"category"`
	Date	 time.Time	 `json:"date"`
 }

type Content struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}