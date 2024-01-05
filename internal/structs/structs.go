package structs

type DevEnvironment struct {
	Name         string       `json:"name"`
	Repositories []Repository `json:"repositories"`
}

type Repository struct {
	Path   string `json:"path"`
	Name   string `json:"name"`
	Action string `json:"action"`
}
