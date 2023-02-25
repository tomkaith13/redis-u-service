package bf

type ReserveRequest struct {
	Name      string  `json:"name"`
	ErrorRate float64 `json:"errorRate"`
	Capacity  uint64  `json:"capacity"`
}
