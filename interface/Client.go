package _interface

type Client interface {
	election()
	getElu() int
}
