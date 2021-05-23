package session

import (
	"github.com/gin-gonic/gin"
	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/repositories"
)

// Accounts ...
var Accounts map[string]*model.User
var userRepository *repositories.UserRepository

func init() {
	Accounts = make(map[string]*model.User, 0)
}

func SetDb(repo *repositories.UserRepository) {
	userRepository = repo
}

// GetUser возращает текущего пользователя для запроса
// Сначало ищим в локальной мапе ,затем в базе
// если пользователя нету или кука не установдена возращаем пустую структуру model.User
func GetUser(c *gin.Context) *model.User {

	token, err := c.Cookie("token")

	if err != nil || len(token) == 0 {
		return nil
	}

	//Сначало ищим в колекции сессии
	if user, ok := Accounts[token]; ok {
		c.Set("user", user)
		return user
	}

	//Пытаемся найти в репозитории по токену
	if user, err := userRepository.FindByToken(token); user.ID != 0 && err == nil {
		Accounts[token] = user
		c.Set("user", user)
		return user
	}

	return nil
}
