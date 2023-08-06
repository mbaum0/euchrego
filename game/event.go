package game

import "fmt"

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
	fmt.Println(event.name)
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
	fmt.Printf("%s: %s drew a %s\n", event.name, player.name, card.GetString())
	return event
}

type DrawnJackEvent struct {
	DefaultEvent
}

func NewDrawnJackEvent() DrawnJackEvent {
	event := DrawnJackEvent{}
	event.name = "DRAWN_JACK_EVENT"
	fmt.Printf("%s: Jack was drawn!\n", event.name)
	return event
}

type DealtCardsEvent struct {
	DefaultEvent
	DealtCards []*Card
}

func NewDealtCardsEvent(cards []*Card, player *Player) DealtCardsEvent {
	event := DealtCardsEvent{}
	event.name = "DEALT_CARDS_EVENT"
	event.DealtCards = cards
	event.player = player
	fmt.Printf("%s: %s was dealt %d cards\n", event.name, player.name, len(cards))
	return event
}

type FinishedDealingEvent struct {
	DefaultEvent
}

func NewFinishedDealingEvent() FinishedDealingEvent {
	event := FinishedDealingEvent{}
	event.name = "FINISHED_DEALING_EVENT"
	fmt.Println(event.name)

	return event
}
