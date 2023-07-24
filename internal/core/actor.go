package core

type Actor int

const (
	ActorAnonymous Actor = iota
	ActorUser
	ActorSystem
)
