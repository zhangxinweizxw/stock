package message

import (
    "fmt"
    "strconv"
    "time"

/share/models"

	"stock
/share/store/redis"
)

type MessageLive struct {
}

func NewMessageLive() *MessageLive {
	return &MessageLive{}
}

// 读取最后操作时间
func (this *MessageLive) GetLastDoingTime(id int64) (int64, error) {
	lastDongTime, err := redis.Get(
		fmt.Sprintf(REDIS_LIVE_LAST_DOING_TIME, id))
	result, _ := strconv.ParseInt(lastDongTime, 10, 64)

	return result, err
}

// 刷新最后操作时间
func (this *MessageLive) RefreshLastDoingTime(id int64) error {
	return redis.Set(
		fmt.Sprintf(REDIS_LIVE_LAST_DOING_TIME, id),
		[]byte(FormatInt(time.Now().Unix())))
}
