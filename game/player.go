package game

import "github.com/mbaum0/euchrego/godeck"

type Player struct {
	hand         []godeck.Card
	tricksTaken  int
	name         string
	index        int
	playedCard   godeck.Card
	pointsEarned int
}

func InitPlayer(name string, index int) *Player {
	player := Player{}
	player.hand = make([]godeck.Card, 0)
	player.tricksTaken = 0
	player.pointsEarned = 0
	player.name = name
	player.index = index
	player.playedCard = godeck.EmptyCard()

	return &player
}

func (p *Player) GetName() string {
	return p.name
}

func (p *Player) GetTricksTaken() int {
	return p.tricksTaken
}

func (p *Player) GiveCards(cards []godeck.Card) {
	p.hand = append(p.hand, cards...)
}

func (p *Player) GiveCard(c godeck.Card) {
	p.hand = append(p.hand, c)
}

func (p *Player) ReturnCards() []godeck.Card {
	var cards = make([]godeck.Card, 0)
	cards = append(cards, p.hand...)
	p.hand = p.hand[:0]
	return cards
}

// removes the card from the players hand and returns it
func (p *Player) ReturnCard(rc godeck.Card) godeck.Card {
	for i, c := range p.hand {
		if c == rc {
			p.hand = append(p.hand[:i], p.hand[i+1:]...)
			return c
		}
	}
	return godeck.EmptyCard()
}
