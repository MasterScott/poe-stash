package scraper

import (
	"encoding/json"
	"fmt"

	"gitlab.perso/poe-stash/inventory"
)

// parseCharacters parses a Path of Exile characters.
func parseCharacters(data []byte) ([]*inventory.Character, error) {
	characters := []*inventory.Character{}
	if err := json.Unmarshal(data, &characters); err != nil {
		return nil, err
	}

	return characters, nil
}

// ScrapCharacters scraps all characters own by a user.
func (s *Scraper) ScrapCharacters() ([]*inventory.Character, error) {
	url := fmt.Sprintf(ProfileCharactersURL, s.accountName)
	body, errRequest := s.CallAPI(url)
	if errRequest != nil {
		return nil, errRequest
	}
	characters, errCharacters := parseCharacters(body)
	if errCharacters != nil {
		return nil, errCharacters
	}

	return characters, nil
}

// parseInventory parses a Path of Exile character inventory.
func parseInventory(data []byte) (*inventory.CharacterInventory, error) {
	inventory := inventory.CharacterInventory{}
	if err := json.Unmarshal(data, &inventory); err != nil {
		return nil, err
	}

	return &inventory, nil
}

// ScrapCharacterInventory scraps the inventory of a given character.
func (s *Scraper) ScrapCharacterInventory(charName string) (*inventory.CharacterInventory, error) {
	url := fmt.Sprintf(ProfileCharacterItemsURL, s.accountName, s.realm, charName)
	body, errRequest := s.CallAPI(url)
	if errRequest != nil {
		return nil, errRequest
	}
	inventory, errInventory := parseInventory(body)
	if errInventory != nil {
		return nil, errInventory
	}

	return inventory, nil
}