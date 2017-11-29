package multiLanguageString

type MultiLanguageString struct {
	Ja string `json:"ja"`
	En string `json:"en"`
	Fr string `json:"fr"`
	Ru string `json:"ru"`
	Zh string `json:"zh"`
	Ko string `json:"ko"`
}

func NewMultiLanguageString(japanese string) *MultiLanguageString {
	return &MultiLanguageString{Ja: japanese}
}
