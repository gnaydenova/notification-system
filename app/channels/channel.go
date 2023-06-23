package channels

const (
	TypeEmail = "email"
	TypeSlack = "slack"
	TypeSMS   = "sms"
	// TypeLog is intended for testing purposes.
	TypeLog   = "log"
)

type Channel interface {
	Send(msg string) error
}

type Regisrty map[string]Channel

func NewRegisrty() Regisrty {
	return make(Regisrty)
}

func (r Regisrty) Add(name string, c Channel) {
	r[name] = c
}

func (r Regisrty) Exists(name string) bool {
	_, ok := r[name]

	return ok
}
