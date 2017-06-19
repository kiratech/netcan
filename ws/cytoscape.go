package ws

type csdata struct {
	Id        string `json:"id"`
	Weight    int    `json:"weight"`
	Source    string `json:"source"`
	Target    string `json:"target"`
	FaveShape string `json:"faveShape"`
	Parent    string `json:"parent"`
}

type csedge struct {
	Data csdata `json:"data"`
}

type csnode struct {
	Data csdata `json:"data"`
}

type cselements struct {
	Nodes []csnode `json:"nodes"`
	Edges []csedge `json:"edges"`
}

type cytoscape struct {
	Elements cselements `json:"elements"`
}
