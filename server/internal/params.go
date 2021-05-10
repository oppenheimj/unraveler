package internal

type Params struct {
	N          int     `json:"n"`
	Kr         float64 `json:"kr"`
	Ka         float64 `json:"ka"`
	Kn         float64 `json:"kn"`
	MaxIters   int     `json:"maxIters"`
	MinError   float64 `json:"minError"`
	Theta      float64 `json:"theta"`
	NumThreads int     `json:"numThreads"`
}
