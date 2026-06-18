package utils

import "regexp"

func IsEmailValido(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func IsSenhaValida(senha string) bool {
	return len(senha) >= 6
}
