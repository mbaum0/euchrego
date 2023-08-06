package game

// an event is some input to the state machine

type Event interface {
	GetPlayer() *Player
}

type DefaultEvent struct {
	player *Player // the player who performed the event
	name   string
}

func (event DefaultEvent) GetPlayer() *Player {
	return event.player
}

type EndGameEvent struct {
	DefaultEvent
}

func NewEndGameEvent() EndGameEvent {
	event := EndGameEvent{}
	event.name = "END_GAME_EVENT"
	return event
}

type EmptyEvent struct {
	DefaultEvent
}

func NewEmptyEvent() EmptyEvent {
	event := EmptyEvent{}
	event.name = "EMPTY_EVENT"
	return event
}

type DrawnCardEvent struct {
	DefaultEvent
	DrawnCard *Card
}

func NewDrawnCardEvent(card *Card, player *Player) DrawnCardEvent {
	event := DrawnCardEvent{}
	event.name = "DRAWN_CARD_EVENT"
	event.DrawnCard = card
	event.player = player
	return event
}

type DrawnJackEvent struct {
	DefaultEvent
}

func NewDrawnJackEvent() DrawnJackEvent {
	event := DrawnJackEvent{}
	event.name = "DRAWN_JACK_EVENT"
	return event
}
