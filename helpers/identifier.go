package helpers

import gonanoid "github.com/matoous/go-nanoid/v2"

func GenerateID(length int) string {
	id, err := gonanoid.New(length)
	if err != nil {
		panic(err)
	}
	return id
}
