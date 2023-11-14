package errors

type errorJsonStruct struct {
	FrameworkName string `json:"framework"`
	Name          string `json:"errorName"`
	Msg           string `json:"Msg"`
}
