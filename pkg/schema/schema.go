package schema

type ShortLink struct {
	ShortPath string `json:"short_path,omitempty" pg:",notnull,unique"`
	RealURL   string `json:"real_url,omitempty" pg:",notnull"`
}
