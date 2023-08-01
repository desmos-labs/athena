package links

type Path struct {
	Path    string
	Actions []Action
}

func NewPath(path string, actions []Action) Path {
	return Path{
		Path:    path,
		Actions: actions,
	}
}

type Action struct {
	Name           string
	RequiredParams []string
}

func NewAction(name string, requiredArgs []string) Action {
	return Action{
		Name:           name,
		RequiredParams: requiredArgs,
	}
}

// -------------------------------------------------------------------------------------------------------------------

type GenerateLinkRequest struct {
	// URL represents the URL of the link that should be generated
	URL string `json:"url"`
}

type GenerateLinkResponse struct {
	Link string `json:"link"`
}

func NewGenerateLinkResponse(link string) *GenerateLinkResponse {
	return &GenerateLinkResponse{
		Link: link,
	}
}
