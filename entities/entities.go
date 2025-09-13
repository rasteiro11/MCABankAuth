package entities

import repository "github.com/rasteiro11/MCABankAuth/src/user/repository/models"

func GetEntities() []any {
	return []any{
		&repository.User{},
	}
}
