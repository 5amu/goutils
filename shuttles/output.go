package shuttles

import "encoding/json"

type ShuttleOutput struct {
	ID        int      `json:"id"`
	Output    string   `json:"output"`
	Arguments []string `json:"arguments"`
	Injected  []string `json:"injected"`
}

func (so *ShuttleOutput) String() string {
	out, _ := json.Marshal(so)
	return string(out)
}
