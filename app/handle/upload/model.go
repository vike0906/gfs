package upload

type SmallFUResult struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
	Url  string `json:"url"`
	Size int64  `json:"size"`
}

func NewSmallFUResult(name, hash, url string, size int64) *SmallFUResult {
	return &SmallFUResult{Name: name, Hash: hash, Url: url, Size: size}
}
