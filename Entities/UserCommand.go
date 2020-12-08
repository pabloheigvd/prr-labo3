package Entities

type ElectionCmd struct {}

func (c ElectionCmd) Match(cmd string) bool {
	return "e" == cmd
}

type GetEluCmd struct {}

func (c GetEluCmd) Match(cmd string) bool {
	return "g" == cmd
}
