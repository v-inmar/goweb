package random_utils

import "github.com/google/uuid"

func GenerateString() string{
	generator := uuid.New()
	return generator.String()
}