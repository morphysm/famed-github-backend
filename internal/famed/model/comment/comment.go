package comment

type Comment interface {
	String() (string, error)
	Type() Type
}
