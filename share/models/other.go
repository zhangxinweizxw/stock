package models

type BirthdayJson struct {
	Date string `json:"date"`
	Type int    `json:"type"`
}

type LocationJson struct {
	City     string `json:"city"`
	Province string `json:"province"`
	District string `json:"district"`
}

type PageJson struct {
	Pagination PagePagination `json:"pagination"`
	Rows       []interface{}  `json:"rows"`
}

type PagePagination struct {
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

type ExchangeParams struct {
	Func   string                 `json:"func"`
	Params map[string]interface{} `json:"params"`
}

type PreferenceJson struct {
	NotifyDesktop      int                    `json:"notify_desktop"`       // 桌面通知
	NotifyMobile       int                    `json:"notify_mobile"`        // 移动端推送
	NotifyDesktopSound int                    `json:"notify_desktop_sound"` // 通知声音：0.静音、1.提示音一、2.提示音二
	NotifyAvoidDisturb NotifyAvoidDisturbJson `json:"notify_avoid_disturb"` // 勿扰模式设置
	TimeZone           string                 `json:"time_zone"`            // 时区设置
}

type NotifyAvoidDisturbJson struct {
	Enabled int `json:"enabled"`
	Start   int `json:"start"`
	End     int `json:"end"`
}

type EventPayload struct {
	Name     string   `json:"name,omitempty"`
	Channels []string `json:"channels"`
	Data     string   `json:"data"`
	SocketId *string  `json:"socket_id,omitempty"`
	State    string   `json:"state,omitempty"`
}

type ReceiveJsonWithQA struct {
	List []interface{} `json:"list"`
	Text string        `json:"text"`
	Type int           `json:"type"`
	Url  string        `json:"url"`
}
