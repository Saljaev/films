package models

import "time"

type Actors struct {
	Id          int
	FirstName   string
	LastName    string
	Gender      string
	DateOfBirth time.Time
	Films       []*Films
}

//func (a *Actors) IsValid() bool {
//	validGenders := map[string]struct{}{
//		"male":   {},
//		"female": {},
//		"other":  {},
//	}
//	_, isValid := validGenders[a.Gender]
//	return isValid && a.DateOfBirth.Year() >= 1800
//}
