package dtos

// All not used since this program transfer
// raw msgs which are from Chomolungma to remote server
// Not parse them yet.
type HuoBiWs struct {
	Sub   string `json:"sub"`
	Unsub string `json:"unsub"`
	Req   string `json:"req"`
	From  string `json:"from"`
	To    string `json:"to"`
	ID    string `json:"id"`
}

type HuoBiWsV2 struct {
	Action string `json:"action"`
	Ch     string `json:"ch"`
	Cid    string `json:"cid"`
	Data   string `json:"data"`
}

type msgChomo struct {
	Close bool   `json:"close"`
	Data  string `json:"data"`
}
