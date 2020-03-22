package cards

/*
card := Card{
		Id:              1,
		Name:            "foobar",
		CardType:        0,
		CharacterValues: CharacterValues{
			Velocity: 1,
			Attack:   2,
			Defense:  3,
			Power:    4,
		},
	}
	return card
*/

type CardType = int

const (
	CharacterCard = iota
	ActionCard    = iota
	TrapCard      = iota
	VehicleCard   = iota
)

type CharacterValues struct {
	Velocity int
	Attack   int
	Defense  int
	Power    int
}

type Card struct {
	Id              int64           `json:"id"`
	Name            string          `json:"name"`
	CardType        CardType        `json:"type"`
	CharacterValues CharacterValues `json:"values"`
}
