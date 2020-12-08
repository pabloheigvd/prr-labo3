package Interface

type Command interface {
	Match(string) bool
}
