package service

type SemResult struct {
	ClassNames []string  `json:"classnames"`
	Scores     []int     `json:"scores"`
	Credits    []float32 `json:"credits"`
	Gpas       []float32 `json:"gpas"`
}

type QueryInput struct {
	Id     string `json:"id"`
	Passwd string `json:"passwd"`
}
