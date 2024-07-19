package core

type Plugin interface {
	Install(*Bot) error
}
