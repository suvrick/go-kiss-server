package session

import "github.com/suvrick/go-kiss-server/model"

// Accounts ...
var Accounts map[string]model.User

func init() {
	Accounts = make(map[string]model.User, 0)
}
