package links

const (
	ParamAddress = "address"
	ParamAction  = "action"
	ParamAmount  = "amount"

	EmptyAction = ""
	ActionSend  = "send"
	ActionView  = "view"
)

var (
	RegisteredPaths []Path
)

func RegisterPath(path string, actions []Action) {
	RegisteredPaths = append(RegisteredPaths, NewPath(path, actions))
}

func init() {
	RegisterPath("/user", []Action{
		NewAction(EmptyAction, []string{}),
		NewAction(ActionSend, []string{ParamAddress, ParamAmount}),
		NewAction(ActionView, []string{ParamAddress}),
	})
}
