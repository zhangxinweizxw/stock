package common

/share/models"
)

type AccessKeys struct {
	Model     `db:"-"`
	KeyID     string // Access Key ID
	KeySecret string // Access Key Secret
	AESKey    string // Encoding AES Key
	AppID     string // App ID
	Token     string // Access Token
	Level     int    // Privary Token
	Expires   int64  // Access Token Expires Time
}

func NewAccessKeys() *AccessKeys {
	return &AccessKeys{
		Model: Model{
			TableName: TABLE_ACCESS_KEYS,
			Db:        MyCat,
		},
	}
}

func (this *AccessKeys) Validate(id string, secret string) bool {
	exps := map[string]interface{}{
		"KeyID=?":     id,
		"KeySecret=?": secret,
		"Level=?":     ACCESS_KEYS_LEVEL_PRIVATE,
	}
	if ok, err := this.IsExist(exps, "ID", ""); err != nil || !ok {
		return false
	}
	return true
}
