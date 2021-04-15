package message

import (
    "fmt"

/share/models"
)

type MessageRead struct {
	Model    `db:"-"`
	ID       int64
	RefID    int64
	MemberID int64
}

// --------------------------------------------------------------------------------

func NewMessageRead() *MessageRead {
	return &MessageRead{
		Model: Model{
			Db:        MyCat,
			TableName: TABLE_MESSAGE_READ,
		},
	}
}

func (this *MessageRead) InsertMessageRead(id int64, memberId int64) error {
	cmd := fmt.Sprintf(`INSERT INTO %v (RefID, MemberID) `+
		`SELECT %v, %v FROM dual WHERE not exists `+
		`(select * from %v where RefID=%v AND MemberID=%v)`,
		this.TableName, id, memberId, this.TableName, id, memberId)

	_, err := this.Db.Exec(cmd)
	return err
}
