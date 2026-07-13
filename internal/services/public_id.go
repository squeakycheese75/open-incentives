package services

import gonanoid "github.com/matoous/go-nanoid/v2"

type PublicIDGenerator interface {
	New(prefix string) (string, error)
}

type NanoIDGenerator struct{}

func (g NanoIDGenerator) New(prefix string) (string, error) {
	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz0123456789", 12)
	if err != nil {
		return "", err
	}

	return prefix + "_" + id, nil
}
