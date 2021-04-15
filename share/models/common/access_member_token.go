package common

/share/models"

	"stock
/share/gocraft/dbr"
	"stock
/share/logging"
	"stock
/share/store/redis"
)

type AccessMemberToken struct {
	Key       string // Access Token Key
	Expires   int64  // Access Token 生存周期，单位是秒
	RefMember string // 用于保存Token对应的用户GUID
}

type AccessMemberTokenJson struct {
	Key          string `json:"access_token"`
	Expires      int64  `json:"expires"`
	RefreshToken string `json:"refresh_token"`
}

func NewAccessMemberToken() *AccessMemberToken {
	return &AccessMemberToken{}
}

func (this *AccessMemberToken) GetMember(token string) (*Member, error) {
	refMember, err := redis.Get(REDIS_MAJOR_TOKEN + token)
	if err != nil || len(refMember) == 0 {
		return nil, dbr.ErrNotFound
	}
	m := NewMember()

	if err := m.GetSingle(IDDecrypt(refMember)); err != nil {
		logging.Debug("%v", err)
		return nil, err
	}
	return m, nil
}
