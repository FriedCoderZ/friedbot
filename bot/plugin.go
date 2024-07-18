package bot

type Plugin interface {
	GetName() string
	GetDependencies() []string
	Install(*Bot) error
}
