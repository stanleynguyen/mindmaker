package domain

// Bucket bucket of options
type Bucket struct {
	Name    string
	Options []Option `json:options`
}
