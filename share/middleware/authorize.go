package middleware

/share/models"

	"stock
/share/lib"
	"stock
/share/logging"
	"stock
/share/middleware/session"
	"stock
/share/models/common"

	"github.com/gin-gonic/gin"
)

type Authorize struct {
}

func NewAuthorize() *Authorize {
	return &Authorize{}
}

func (this *Authorize) GetClient(c *gin.Context) *Client {
	client := NewClient()

	// 获取 AccessToken
	token := c.Query("access_token")
	if len(token) > 0 {

		// 通过 AccessToken 获得会员信息
		t := common.NewAccessMemberToken()
		if m, err := t.GetMember(token); err != nil {
			logging.Debug("access_token: %v| Err:%v", token, err)
			return nil
		} else {
			client.Member.ID = m.ID
			client.Member.Name = m.FriendlyName
			client.Member.IsLogged = true
			client.Member.ClientIP = lib.TrimIPAddr(c.ClientIP())
		}
	} else {

		// 通过 Session 获取会员信息
		sess := session.Default(c)
		if sess == nil {
			return nil
		}

		value := sess.Get(SESSION_MEMBER_ID)
		if value == nil {
			return nil
		}

		client.Member.ID = value.(int64)
		client.Member.Name = sess.Get(SESSION_MEMBER_NAME).(string)
		client.Member.IsLogged = sess.Get(SESSION_MEMBER_LOGINED).(bool)
		client.Member.ClientIP = lib.TrimIPAddr(c.ClientIP())
	}

	return client
}

func (this *Authorize) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 如果会员未登录
		client := this.GetClient(c)
		if client == nil || !client.Member.IsLogged {
			logging.Debug("Middleware Error: 401")
			lib.WriteString(c, 401, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

// --------------------------------------------------------------------------------

type Client struct {
	Member *RefMember
}

type RefMember struct {
	ID       int64
	ClientIP string
	IsLogged bool
	Name     string
}

func NewClient() *Client {
	return &Client{
		Member: &RefMember{},
	}
}

// 验证会员特殊权限
func (this *Client) ValidateSpecialPermission(key string) bool {
	return true
}

// 验证会员权限
func (this *Client) ValidatePermission(key string, action int) bool {
	return true
}
