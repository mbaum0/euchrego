package game

import (
	"fmt"
	"strings"
)

type Player struct {
	hand   []*Card
	tricks int
	score  int
	name   string
}

func InitPlayer(name string) Player {
	player := Player{}
	player.hand = make([]*Card, 0)
	player.score = 0
	player.tricks = 0
	player.name = name

	return player
}

func (p *Player) GetTricks() int {
	return p.tricks
}

func (p *Player) GetScore() int {
	return p.score
}

func (p *Player) GiveCards(cards []*Card) {
	p.hand = append(p.hand, cards...)
}

func (p *Player) ReturnCards() []*Card {
	var cards = make([]*Card, 0)
	cards = append(cards, p.hand...)
	p.hand = p.hand[:0]
	return cards
}

func (p *Player) PrintHand() {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s's Hand\n", p.name))
	for _, c := range p.hand {
		builder.WriteString(fmt.Sprintf("%s\t\n", c.Info()))
	}
	result := builder.String()
	fmt.Println(result)

}

func (p *Player) GetPlayableCards(trump Suite, lead Suite) []*Card {
	var cards = make([]*Card, 0)

	hasLeadCards := false
	// check if any cards match what was lead
	for _, c := range p.hand {
		if c.suite == lead {
			cards = append(cards, c)
			hasLeadCards = true
		}
	}

	hasTrumpCards := false
	// check if there trump cards
	if !hasLeadCards {
		for _, c := range p.hand {
			if c.suite == trump || (c.suite == LeftBauerSuite[trump] && c.rank == JACK) {
				cards = append(cards, c)
				hasTrumpCards = true
			}
		}
	}

	// if we don't have trump or lead cards, all cards are valid
	if !hasLeadCards && !hasTrumpCards {
		cards = append(cards, p.hand...)
	}

	return cards
}
