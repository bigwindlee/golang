package main

import (
	"fmt"
)

type Item struct {
	id       string
	price    float64
	quantity int
}

func (item *Item) Cost() float64 {
	return item.price * float64(item.quantity)
}

type SpecialItem struct {
	Item      // Anonymous field (embedding)
	catalogId int
}

type LuxuryItem struct {
	Item   // Anonymous field (embedding)
	markup float64
}

// Override the methods of an embedded field simply by creating a new method for
// the embedding struct that has the same name as one of the embedded fields’ methods.
func (item *LuxuryItem) Cost() float64 {
	return item.Item.Cost() * item.markup
}

func main() {
	special := SpecialItem{Item{"Green", 3, 5}, 207}
	fmt.Println(special.id, special.price, special.quantity, special.catalogId)
	// Green 3 5 207

	// When we call special.Cost(), since the SpecialItem type does not have its
	// own Cost() method, Go uses the Item.Cost() method——and passes it the embedded
	// Item value, not the entire SpecialItem that the method was originally called on.
	fmt.Println(special.Cost())
	// 15

	luxury := LuxuryItem{Item{"Gold", 30, 2}, 10}
	fmt.Println(luxury.Cost())
	// 600
}
