package friedbot

type Plugin interface {
	Install(*Bot) error
}
