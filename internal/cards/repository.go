package cards

var (
	data = Table{cardsMap: map[int64]*Card{}}
)

type Table struct {
	cardsMap map[int64]*Card
}

func selectById(cardId int64) *Card {

	for id, card := range data.cardsMap {
		if id == cardId {
			return &Card{
				Id:              card.Id,
				Name:            card.Name,
				CardType:        card.CardType,
				CharacterValues: card.CharacterValues,
			}
		}
	}

	return nil
}

func selectAll() []*Card {

	values := make([]*Card, 0, len(data.cardsMap))
	for _, card := range data.cardsMap {
		values = append(values, &Card{
			Id:              card.Id,
			Name:            card.Name,
			CardType:        card.CardType,
			CharacterValues: card.CharacterValues,
		})
	}

	return values
}

func insert(card *Card) bool {
	data.cardsMap[card.Id] = card
	return true
}
