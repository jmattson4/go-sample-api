package accounts

type Reader interface {
	GetUser(u uint) *Account
}
type Writer interface {
	Create(account *Account) error
}

type Constructor interface {
	Construct(interface{}) interface{}
}

type BaseRepo interface {
	Reader
	Writer
	Constructor
}

type DbRepo interface {
	BaseRepo
	Create(account *Account) error
}

type CacheRepo interface {
	BaseRepo
	CreateToken(accountID uint) (*TokenDetails, error)
	CreateAuth(userid uint, td *TokenDetails) error
	DeleteAuth(givenUuid string) (int64, error)
}

type BaseService interface {
	Login(email, password string) map[string]interface{}
}
