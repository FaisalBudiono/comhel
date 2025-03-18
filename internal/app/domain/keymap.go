package domain

type Keymap struct {
	Keys        []string
	Description string
}

func NewKeymap(keys []string, description string) Keymap {
	return Keymap{
		Keys:        keys,
		Description: description,
	}
}
