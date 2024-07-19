package core

type Plugin interface {
	GetName() string
	GetDependencies() []string
	Install(*Bot) error
}
