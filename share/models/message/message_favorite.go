package message

/share/models"
)

type MessageFavorite struct {
	Model    `db:"-"`
	ID       int64
	RefID    int64 // 消息ID
	MemberID int64 // 成员ID
}

// --------------------------------------------------------------------------------

func NewMessageFavorite() *MessageFavorite {
	return &MessageFavorite{
		Model: Model{
			TableName: TABLE_MESSAGE_FAVORITES,
			Db:        MyCat,
		},
	}
}
