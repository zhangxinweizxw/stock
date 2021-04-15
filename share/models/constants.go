package models

import (
    "errors"
)

var (
	ERROR_OBJECT_IS_NIL = errors.New("object is nil!")
)

const (
	PRODUCT_NAME = "海纳牛牛"
)

// Base URL
// --------------------------------------------------------------------------------

const (
	BASE_URL_AVATAR = "http://avatar.0606.com.cn/"
	BASE_URL_ATTACH = "http://attach.0606.com.cn/"
)

const (
	BASE_WEBLIVE_IMAGE = "weblive/u0.png"
	DEFAULT_AVATAR_URI = "u0.png"
)

// Notify Content
// --------------------------------------------------------------------------------

const (
	NOTIFY_MODEL_CONTENT_WEBLIVE     = `您关注的［%v］已经开启，今日主题“%v”，快去围观吧！`
	NOTIFY_MODEL_CONTENT_BUY_SUCCESS = `您已经成功订阅［%v］，产品服务周期%v天，如有疑问可拨打客服电话：400-156-6699`
	NOTIFY_MODEL_EXPIRING_OTHER      = `您订阅的［%v］，还有5天就要到期了哦！%v最新%v［%v］热卖中！`
	NOTIFY_MODEL_EXPIRING            = `您订阅的［%v］，还有5天就要到期了哦！`
	NOTIFY_MODEL_EXPIRED_OTHER       = `您订阅的［%v］已经到期，想要继续享受老师的服务，%v最新%v［%v］热卖中！`
	NOTIFY_MODEL_EXPIRED             = `您订阅的［%v］已经到期。`
	NOTIFY_MODEL_PRODUCT_UPDATE      = `您订阅的［%v%v-%v］有新内容了，快去看看吧！`
	NOTIFY_MODEL_STRATEGY_TRANSFER   = `策略"%v" %v %v %v%v`
)

// Access Keys Level
// --------------------------------------------------------------------------------

const (
	ACCESS_KEYS_LEVEL_PUBLIC  = 0 // 公有
	ACCESS_KEYS_LEVEL_PRIVATE = 1 // 私有
)

// Affix Assort Categorys
// --------------------------------------------------------------------------------

const (
	AFFIX_ASSORT_NORMAL  = 0 // 普通
	AFFIX_ASSORT_PICTURE = 1 // 图片
	AFFIX_ASSORT_VOICE   = 2 // 音频
	AFFIX_ASSORT_VIDEO   = 3 // 视频
)

// Context and Session Categorys
// --------------------------------------------------------------------------------

const (
	CONTEXT_ADVISOR_GUID        = "advisor_guid"
	CONTEXT_AFFIX_GUID          = "affix_guid"
	CONTEXT_ARTICLES_GUID       = "articles_guid"
	CONTEXT_CODE                = "code"
	CONTEXT_DATE                = "date"
	CONTEXT_DIRECTION           = "direction"
	CONTEXT_FILTER              = "filter"
	CONTEXT_GUID                = "guid"
	CONTEXT_ID                  = "id"
	CONTEXT_JOB_NUMBER          = "job_number_guid"
	CONTEXT_LATEST_STAMP        = "latest_stamp"
	CONTEXT_LIMIT               = "limit"
	CONTEXT_MEMBER_GUID         = "member_guid"
	CONTEXT_MEMBER_ID           = "member_id"
	CONTEXT_MESSAGE_GUID        = "message_guid"
	CONTEXT_MONTH               = "month"
	CONTEXT_OPINION_GUID        = "opinion_guid"
	CONTEXT_ORDER_GUID          = "order_guid"
	CONTEXT_PAGE                = "page"
	CONTEXT_QUESTION_GUID       = "question_guid"
	CONTEXT_REF_ID              = "ref_id"
	CONTEXT_REF_TYPE            = "ref_type"
	CONTEXT_REPORT_GUID         = "report_guid"
	CONTEXT_ROOM                = "room"
	CONTEXT_ROOM_GUID           = "room_guid"
	CONTEXT_ROOM_NUMBER         = "room_number"
	CONTEXT_SESSION_GUID        = "session_guid"
	CONTEXT_SIZE                = "size"
	CONTEXT_SNS_TYPE            = "sns_type"
	CONTEXT_TACTICS_ADVISES_ID  = "advise_guid"
	CONTEXT_TACTICS_ID          = "tactic_guid"
	CONTEXT_TOPIC_GUID          = "topic_guid"
	CONTEXT_TYPE                = "type"
	CONTEXT_UCODE_GUID          = "ucode_guid"
	CONTEXT_YEAR                = "year"
	CONTEXT_CHANNEL_GUID        = "channel_guid"
	CONTEXT_TITLE               = "title"
	CONTEXT_PRODUCT_GUID        = "product_guid"
	CONTEXT_PROPOSAL_GUID       = "proposal_guid"
	CONTEXT_SYSTEM_MESSAGE_GUID = "system_message_guid"
	CONTEXT_ACCESS_TOKEN        = "access_token"
	CONTEXT_STRATEGY_GUID       = "strategy_guid"
	CONTEXT_COURSE_GUID         = "course_guid"
	CONTEXT_LESSON_GUID         = "lesson_guid"
	CONTEXT_TIMESTAMP           = "timestamp"
	CONTEXT_RECOMMEND           = "recommend"
	CONTEXT_SIGNATURE           = "signature"
	CONTEXT_ASSEMBLE_GUID       = "assemble_guid" // add by yh 20170808 (复合产品guid)
)

const (
	SESSION_MEMBER_ID          = "user_id"
	SESSION_MEMBER_GROUP       = "group"
	SESSION_MEMBER_NAME        = "user_name"
	SESSION_MEMBER_LOGINED     = "user_is_logged"
	SESSION_MEMBER_DEVICETOKEN = "device_token"
	SESSION_MEMBER_AVATAR      = "avatar"
)

// Entity RefType Categorys
// --------------------------------------------------------------------------------
const (
	REFTYPE_TACTIC         = 1  // 锦囊
	REFTYPE_WEBLIVE        = 2  // 聊天直播
	REFTYPE_TEXTLIVE       = 3  // 图文直播
	REFTYPE_REPORT         = 4  // 内参
	REFTYPE_GROUP_MEMBER   = 5  // 群组
	REFTYPE_MEMBER_UPDATE  = 7  // 会员升级
	REFTYPE_NEWS           = 8  // 资讯
	REFTYPE_ROOM           = 9  // 房间消息频道
	REFTYPE_SYSTEM         = 10 // 系统消息
	REFTYPE_STRATEGY       = 11 // 量化策略
	REFTYPE_COURSE         = 12 // 课堂
	REFTYPE_PRIVATE_NOTIFY = 13 // 私信提醒
	REFTYPE_MEMBER_FOLLOW  = 14 // 用户关注（提醒等）
	REFTYPE_COURSE_LESSON  = 15 // 课堂 精品课直播 wdk 20170815 add
	REFTYPE_ASSEMBLE       = 99 // 组合产品
)

// Member Status Categorys
// --------------------------------------------------------------------------------

const (
	MEMBER_STATUS_DELETED      = 0 // 删除
	MEMBER_STATUS_NORMAL       = 1 // 正常
	MEMBER_STATUS_LOCK         = 2 // 锁定
	MEMBER_STATUS_LEAVE_OFFICE = 3 // 离职
)

// Weblive Status
// --------------------------------------------------------------------------------

const (
	WEBLIVE_ROOM_TYPE_TEXT   = 0 // 图文直播
	WEBLIVE_ROOM_TYPE_SINGLE = 1 // 单人视频直播
	WEBLIVE_ROOM_TYPE_MULTI  = 2 // 多人视频直播

	WEBLIVE_MESSAGE_TYPE_TEXT = 1 // 图文直播
	WEBLIVE_MESSAGE_TYPE_CHAT = 2 // 聊天

	WEBLIVE_ROOM_START = 1 // 开启
	WEBLIVE_ROOM_END   = 0 // 关闭

	WEBLIVE_ROOM_RECOMMEND = 1 // 精选

	WEBLIVE_ADVISORS_ASSORT_ADVISOR  = 0 // 投顾
	WEBLIVE_ADVISORS_ASSORT_CUSTOMER = 1 // 嘉宾

	WEBLIVE_ROOM_ENABLED = 1 // 直播启用 wdk 20170807 add
)

// Message target types
// --------------------------------------------------------------------------------

const (
	MESSAGE_TARGET_TYPE_CHANNEL = 1 // 群组
	MESSAGE_TARGET_TYPE_SESSION = 2 // 会话
	MESSAGE_TARGET_TYPE_LIVE    = 3 // 直播
	MESSAGE_TARGET_TYPE_HELPER  = 4 // 阿牛智能助手
	MESSAGE_TARGET_TYPE_MESSAGE = 9 // 消息
)

// Message body types
// --------------------------------------------------------------------------------

const (
	MESSAGE_BODY_TYPE_SYSTEM     = 1 // 系统
	MESSAGE_BODY_TYPE_TEXT       = 2 // 文本
	MESSAGE_BODY_TYPE_ATTACHMENT = 3 // 附件
	MESSAGE_BODY_TYPE_INLINE     = 4 // 内联
	MESSAGE_BODY_TYPE_UNDO       = 9 // 已撤回
)

// Request URL query Keys
// --------------------------------------------------------------------------------

const (
	REQUEST_QUERY_DOWNLOAD        = "download"
	REQUEST_QUERY_FILENAME        = "filename"
	REQUEST_QUERY_LATEST_STAMP    = "latest_stamp"
	REQUEST_QUERY_LIMIT           = "limit"
	REQUEST_QUERY_MEDIUM          = "medium"
	REQUEST_QUERY_REF_ID          = "ref_id"
	REQUEST_QUERY_QUALITY         = "quality"
	REQUEST_QUERY_LESSON          = "lesson"
	REQUEST_QUERY_REF_TYPE        = "ref_type"
	REQUEST_QUERY_PRIORITY        = "tf" // 排序优先级
	REQUEST_QUERY_ID              = "id"
	REQUEST_QUERY_REFRESH_TOKEN   = "refresh_token"
	REQUEST_QUERY_PAGE            = "page"
	REQUEST_QUERY_DEVICE_IDENTITY = "device_identity"
	REQUEST_QUERY_MOBILE          = "mobile"
	REQUEST_QUERY_CODE            = "code"
	REQUEST_QUERY_STATE           = "state"
	REQUEST_QUERY_TYPE            = "type"
	REQUEST_ROOM_GUID             = "room"
	REQUEST_YEAR                  = "year"
	REQUEST_MONTH                 = "month"
	REQUEST_QUERY_CATEGORY_ID     = "category_id"
	REQUEST_QUERY_STATUS          = "status"
	REQUEST_QUERY_BEGIN_STAMP     = "begin_stamp"
	REQUEST_QUERY_END_STAMP       = "end_stamp"
	REQUEST_QUERY_IS_PAY          = "is_pay"
	REQUEST_QUERY_POSITION        = "position"
	REQUEST_QUERY_DATE            = "date"
	REQUEST_QUERY_COUNTRY         = "country"
	REQUEST_QUERY_START_TIME      = "start_time"
	REQUEST_QUERY_END_TIME        = "end_time"
	REQUEST_QUERY_RUNSTATUS       = "run_status"
	REQUEST_QUERY_IS_RECOMMENED   = "is_recommened"
)

// Redis Prefix for Other
// --------------------------------------------------------------------------------

const (
	// member
	REDIS_MAJOR_ACCESS_KEYS         = "m:access:keys:"         // Access Keys（HASH）
	REDIS_MAJOR_ACCESS_TOKEN        = "m:access:token:"        // Access Token
	REDIS_MAJOR_ACCESS_MEMBER_TOKEN = "m:access:member:token:" // Access Member Token
	REDIS_MAJOR_CAPTCHA_SESSION     = "m:captcha:session:"     // Session -> 图片验证码Key
	REDIS_MAJOR_CAPTCHA_CODE        = "m:captcha:code:"        // 图片验证码Key -> Session
	REDIS_MAJOR_CAPTCHA_RATE        = "m:rate:%v:captcha"      // 图片验证码请求频率
	REDIS_MAJOR_SIGNUP_INDEX        = "m:signup:today"         // 当日注册顺序
	REDIS_MAJOR_PROTOCOL_INDEX      = "m:protocol:today"       // 当日购买协议顺序

	REDIS_MAJOR_MAIL               = "m:mail:"                    // 邮箱验证码信息（HASH）
	REDIS_MAJOR_MOBILE             = "m:mobile:"                  // 手机验证码信息（HASH）
	REDIS_MAJOR_MOBILE_TOKEN       = "m:mobile:token:"            // 手机验证码（Token -> Mobile）
	REDIS_MAJOR_SEQ_SESSION        = "m:seq:session"              // 会话序列号
	REDIS_MAJOR_TOKEN              = "m:token:"                   // Token
	REDIS_CARD_NO_AUTH             = "m:cardno:auth:%v"           // 实名认证 [HASH]
	REDIS_NOTIFY_SUBSCRIBE_PRODUCT = "m:%v:notify:product"        // 消息产品订阅（memberId） [SET]
	REDIS_NOTIFY_SUBSCRIBE_CHANNEL = "m:%v:notify:channel"        // 消息频道订阅（memberId）[SET]
	REDIS_NOTIFY_MESSAGE_LASTTIME  = "m:%v:message:lasttime"      // 用户最后拉取订阅消息时间（memberId）[STRING]
	REDIS_NOTIFY_CMS_LASTTIME      = "m:%v:cms:lasttime"          // 用户最后拉取资讯时间（memberId）[STRING]
	REDIS_NOTIFY_SYSTEM_LASTTIME   = "m:%v:system:lasttime"       // 用户最后拉取系统时间（memberId）[STRING]
	REDIS_PROPOSAL_TOTAL_NUMBER    = "m:%v:proposal:totalnumber"  // 用户投诉建议总数（memberId）[STRING]
	REDIS_PROPOSAL_SOLVED_NUMBER   = "m:%v:proposal:solvednumber" // 用户投诉建议已处理数（memberId）[STRING]
	REDIS_MEMBER_AUTH              = "m:%v:auth"                  // 用户实名认证（memberId）[STRING]

	REDIS_ADVISOR_MESSAGE_INC_ID = "a:lastid:message"   // 消息最后的插入ID  [INT]
	REDIS_ADVISOR_CHANNEL_INC_ID = "a:lastid:channel"   // 频道最后的ID [INT]
	REDIS_MAJOR_ADVISOR_GROUP    = "m:advisor:group:%v" // 投顾会员组 (group.weight) [SETS]

	REDIS_MEMBERS_FOLLOW = "m:%v:follow" // 会员关注列表（memberId）[Set]
	REDIS_MEMBERS_FANS   = "m:%v:fan"    // 会员粉丝数量（memberId）[INT]

	REDIS_MEMBERS_CHANNELS_UNREAD  = "m:%v:unread:channels:%v" // 频道未读数 （memberId,channelId）[INT]
	REDIS_MEMBERS_SESSIONS_UNREAD  = "m:%v:unread:sessions:%v" // 消息未读数 （memberId,channelId） INT]
	REIDS_MEMBERS_UNREAD           = "m:%v:unread:*"           // 所有未读数 (memberId)
	REDIS_UPLOAD_AFFIXS            = "m:upload:affixs:%v:info" // 附件信息 (attachmentId) [HASH]
	REDIS_SESSIONS_LAST_DOING_TIME = "m:sessions:%v:ldt"       // 会话最后操作时间 (sessionId) [STRING]
	REDIS_SESSIONS_LAST_MESSAGE    = "m:sessions:%v:msg"       // 会话最后消息内容 (sessionId) [STRING]

	REDIS_BUY_TIME_PERMIT = "m:%v:buy:permit:%v" // 用户购买产品时间许可（memberId,productId） [STRING]

	// a 代表 advisor
	REDIS_ADVISOR_BOXES     = "a:%v:boxes"     // 投顾宝箱数量（memberId）[INT]
	REDIS_ADVISOR_OPINIONS  = "a:%v:opinions"  // 投顾观点数量（memberId）[INT]
	REDIS_ADVISOR_QUESTIONS = "a:%v:questions" // 投顾问答数量（memberId）[INT]

	REDIS_ADVISOR_OPINION_COMMENTS = "a:opinion:%v:comments" // 投顾观点评论数（opinionId）[INT]
	REDIS_ADVISOR_OPINION_READS    = "a:opinion:%v:reads"    // 投顾观点阅读数（opinionId）[INT]
	REDIS_ADVISOR_OPINION_PRAISES  = "a:opinion:%v:praises"  // 投顾观点点赞数（opinionId）[INT]

	//ch 代表 channel
	REDIS_CHANNELS                 = "ch:%v:info"       // 频道信息  （channelId）[HASH]
	REDIS_CHANNELS_LAST_DOING_TIME = "ch:%v:ldt"        // 频道最后操作时间 （channelId）[STRING]
	REDIS_CHANNELS_LAST_MESSAGE    = "ch:%v:msg"        // 频道最后消息内容 (channelId) [STRING]
	REDIS_CHANNELS_MEMBERS         = "ch:%v:members:%v" // 频道成员信息 （channelId,memberId）[HASH]

	REDIS_PUSHER_CHANNEL_OF_SCHEMAS = "p:channel:%v"
	REDIS_MEMBERS                   = "m:%v"        // 会员信息（HASH）
	REDIS_SIMPLE_MEMBERS            = "m:simple:%v" // 简要会员信息（HASH）
	REDIS_MEMBER_CODE               = "m:code"      // 会员编码 [INT]

	//wl 代表 weblive
	REDIS_WEBLIVE_ROOM_VISITOR         = "wl:visitor:%v"         // 参与人数 [STRING]
	REDIS_WEBLIVE_ROOM_QUESTION        = "wl:question:%v"        // 问答人数 [SET]
	REDIS_WEBLIVE_ROOM_FOCUS           = "wl:focus:%v"           // 直播聚焦 (roomId) [HASH]
	REDIS_WEBLIVE_ROOM_INFO            = "wl:room:%v:info"       // 直播室的信息 (roomId) [HASH]
	REDIS_WEBLIVE_ROOM_CHANNEL         = "wl:room:%v:channel:%v" // 直播室订阅的频道 (roomId,assort) [STRING]
	REDIS_WEBLIVE_ROOM_BAN             = "wl:room:%v:ban:%v"     // 直播室禁言列表 (roomId,memberId) [STRING]
	REDIS_LIVE_LAST_DOING_TIME         = "wl:topic:%v:ldt"       // 直播消息最后操作时间 （topicId) [STRING]
	REDIS_WEBLIVE_CHARTS_OPINIONS      = "wl:charts:opinions"    // 观点排行榜  [ZSet]
	REDIS_WEBLIVE_CHARTS_QUESTIONS     = "wl:charts:questions"   // 问答排行榜  [ZSet]
	REDIS_WEBLIVE_CHARTS_COMPREHENSIVE = "wl:charts:compres"     // 综合排行榜  [ZSet]
	REDIS_WEBLIVE_TOTAL_OPINION        = "wl:total:opinion"      // 观点总数  [STRING]
	REDIS_WEBLIVE_TOTAL_WEBLIVE        = "wl:total:weblive"      // 直播总数  [STRING]
	REDIS_WEBLIVE_TOTAL_TREASURE       = "wl:total:treasure"     // 百宝箱总数  [STRING]
	REDIS_WEBLIVE_TOTAL_QUESTION       = "wl:total:question"     // 问答总数  [STRING]
	REDIS_WEBLIVE_TOTAL_VISITOR        = "wl:total:visitor"      // 影响人数  [STRING]
	REDIS_WEBLIVE_TOTAL_UPDATE_TIME    = "wl:total:updatetime"   // 海纳牛牛更新时间  [STRING]

	// c 代表 cms
	REDIS_CMS_CATEGORY_MAGAZINE = "c:category:magazine" // 首证期刊分类 [HASH]
	REDIS_CMS_CATEGORY_FINANCE  = "c:category:finance"  // 财经头条分类 [HASH]

	REDIS_CMS_ARTICLES_VIEW_COUNT = "c:articles:viewcount:%v" // 咨询访问量总数 [STRING] wdk 20170731 add

	// order
	REDIS_MEMBER_ORDER_TREASURE_BOX = "m:%v:order:treasurebox" // 用户百宝箱购买记录 （memberID）[SETS]
	REDIS_MEMBER_ORDER_ACTIVITY     = "m:%v:order:activity"    // 用户活动课堂购买记录 （memberID）[SETS]
	REDIS_MEMBER_ORDER_ASSEMBLE     = "m:%v:order:assemble"    // 用户复合产品购买记录 （memberID）[SETS]
	REDIS_MEMBER_ORDER_STRATEGY     = "m:%v:order:strategy"    // 用户策略购买记录 （memberID）[SETS]

	// mj 代表 major tb 代表 treasure_box
	REDIS_TREASURE_BOX_SUBSCRIBES = "mj:tb:%v:subscribes" // 宝箱订阅数量 （treasureBoxId）[STRING]

	// gp 代表group
	REDIS_GROUPS = "gp:%v" // 会员组信息（groupId）（HASH）

	// s 代表 strategy
	REDIS_STRATEGY_BUY_COUNT       = "s:%v:buy:count"       // 策略订阅数（strategyID）[STRING]
	REDIS_STRATEGY_BUY_COUNT_CRIME = "s:%v:buy:count:crime" // 策略订阅数(假)（strategyID）[STRING]

	// sw 代表 sensitive word
	REDIS_SENSITIVE_WORD             = "sw:word"             // 敏感词 [SET]
	REDIS_SENSITIVE_WORD_REG         = "sw:word:reg"         // 敏感词正则表达式 [STRING]
	REDIS_SENSITIVE_WORD_MSG_VERSION = "sw:word:msg:version" // 敏感词更新标志 [STRING] wdk 20170713 modify （Andy说以后会有四种敏感词库，现在这里分开）
	REDIS_SENSITIVE_WORD_NAME        = "sw:word:name"        // 用户名/昵称 敏感词 [SET]
	REDIS_SENSITIVE_WORD_NAME_REG    = "sw:word:name:reg"    // 用户名/昵称 敏感词正则表达式 [STRING]

	// co 代表 course
	REDIS_COURSE_BUY_COUNT          = "co:%v:buy:count"          // 课堂购买数量 （courseID） [STRING]
	REDIS_COURSE_BUY_COUNT_CRIME    = "co:%v:buy:count:crime"    // 课堂购买数量(假) （courseID） [STRING]
	REDIS_COURSE_ONLINE_MEMBER      = "co:%v:online:member"      // 课堂直播在线用户 （lessonID） [SET]
	REDIS_COURSE_ONLINE_COUNT       = "co:%v:online:count"       // 课堂直播在线人数 （lessonID） [STRING]
	REDIS_COURSE_ONLINE_COUNT_CRIME = "co:%v:online:count:crime" // 课堂直播在线人数(假) （lessonID） [STRING]
	REDIS_COURSE_VIDEO_COUNT        = "co:%v:video:count"        // 课堂视频播放人数 （lessonID） [STRING]
	REDIS_COURSE_VIDEO_COUNT_CRIME  = "co:%v:video:count:crime"  // 课堂视频播放人数(假) （lessonID） [STRING]

	// bl 代表 blacklist
	REDIS_BLACKLIST_IP = "bl:ip:set" // IP黑名单 (SET)
)

// MNS Exchange
// --------------------------------------------------------------------------------

const (
	MNS_EXCHANGE_ADVISER = "exchange-advisor"
	MNS_EXCHANGE_BOX     = "exchange-box"
	MNS_PUSHER_QUEUE     = "pusher-queue-%v"
)

const (
	FUNC_DELETE_ATTACHMENT  = "DELETE_ATTACHMENT"
	FUNC_DELETE_MESSAGE     = "DELETE_MESSAGE"
	FUNC_SAVE_MESSAGE       = "SAVE_MESSAGE"
	FUNC_SAVE_VOICE_MESSAGE = "SAVE_VOICE_MESSAGE"
)

// Database Schema Name
// --------------------------------------------------------------------------------

const (
	SCHEMA_MAJOR = "haina_major."
)

// Database Table Name
// --------------------------------------------------------------------------------

// MAJOR SCHEMA

const (
	TABLE_ACCESS_KEYS             = "hn_accesskeys"
	TABLE_ADVISOR_SIGNED          = "hn_advisor_signed"
	TABLE_ADVISOR_OPINION         = "hn_advisor_opinions"
	TABLE_ADVISOR_OPINION_READER  = "hn_advisor_opinion_reader"
	TABLE_ADVISOR_OPINION_PRAISE  = "hn_advisor_opinion_praise"
	TABLE_ADVISOR_OPINION_COMMENT = "hn_advisor_opinion_comments"
	TABLE_ARTICLE_CATEGORYS       = "hn_article_categorys"
	TABLE_ARTICLES                = "hn_articles"
	TABLE_CHANNEL_MEMBERS         = "hn_channel_members"
	TABLE_CHANNEL_MESSAGES        = "hn_channel_messages"
	TABLE_CHANNELS                = "hn_channels"
	TABLE_FOLLOWS                 = "hn_follows"
	TABLE_GROUPS                  = "hn_groups"
	TABLE_IDCARD_CHECK            = "hn_idcard_check"
	TABLE_INVITE_EMAIL            = "hn_invite_email"
	TABLE_LAYOUT_COLUMN           = "hn_layout_column"
	TABLE_LAYOUT_FIELDS           = "hn_layout_fields"
	TABLE_MEMBER_ADVISORS         = "hn_member_advisors"
	TABLE_MEMBER_AUTH             = "hn_member_auth"
	TABLE_MEMBER_CODE             = "hn_member_code"
	TABLE_MEMBER_FOLLOW           = "hn_member_follows"
	TABLE_MEMBER_LOGGING          = "hn_member_logging"
	TABLE_MEMBER_RISK_TESTING     = "hn_member_risk_testing"
	TABLE_MEMBER_SNS              = "hn_member_sns"
	TABLE_MEMBERS                 = "hn_members"
	TABLE_MESSAGE_FAVORITES       = "hn_message_favorites"
	TABLE_MESSAGE_READ            = "hn_message_read"
	TABLE_MESSAGES                = "hn_messages"
	TABLE_MOBILE_AUTHCODE         = "hn_mobile_authcode"
	TABLE_ORDERS                  = "hn_orders"
	TABLE_PRODUCTS                = "hn_products"
	TABLE_PRODUCT_MEMBER_BUY      = "hn_product_member_buy"
	TABLE_RESEARCH_REPORT_AFFIXS  = "hn_research_report_affixs"
	TABLE_RESEARCH_REPORTS        = "hn_research_reports"
	TABLE_RISK_LEVEL              = "hn_risk_level"
	TABLE_RISK_QUESTION           = "hn_risk_question"
	TABLE_RISK_QUESTION_OPTION    = "hn_risk_question_option"
	TABLE_RISK_TESTING_DETAIL     = "hn_risk_testing_details"
	TABLE_SESSIONS                = "hn_sessions"
	TABLE_TACTICS_ADVISES         = "hn_tactics_advises"
	TABLE_TACTICS_ORDER           = "hn_tactic_order"
	TABLE_TOPIC_COMMENT           = "hn_topic_comment"
	TABLE_UPLOAD_AFFIXS           = "hn_upload_affixs"
	TABLE_UPLOAD_RELEVANCE        = "hn_upload_relevance"
	TABLE_WEBLIVE_ROOMS           = "hn_weblive_rooms"
	TABLE_WEBLIVE_ROOMS_LOGS      = "hn_weblive_room_logs"
	TABLE_WEBLIVE_VISITORS        = "hn_weblive_visitors"
	TABLE_WEBLIVE_ADVISORS        = "hn_weblive_advisors"
	TABLE_ADVISOR_QUESTIONS       = "hn_advisor_questions"
	TABLE_CATEGORYS               = "hn_categorys"
	TABLE_TREASURE_BOX            = "hn_treasure_box"
	TABLE_CONSUME_DETAILS         = "hn_consume_details"
	TABLE_ADVERTS                 = "hn_adverts"
	TABLE_NOTIFY                  = "hn_notify"
	TABLE_NOTIFY_SYSTEM           = "hn_notify_system"
	TABLE_NOTIFY_SUBSCRIBE        = "hn_notify_subscribe"
	TABLE_PAY_OFF                 = "hn_pay_off"
	TABLE_PROPOSALS               = "hn_proposals"
	TABLE_PROPOSAL_DETAILS        = "hn_proposal_details"
	TABLE_STRATEGYS               = "hn_strategys"
	TABLE_STRATEGY_STOCK_PROFIT   = "hn_strategy_stock_profit"
	TABLE_SENSITIVE_WORD          = "hn_sensitive_word"
	TABLE_SENSITIVE_WORD_MSG      = "hn_sensitive_word_msg" //wdk 20170713 modify （Andy说以后会有四种敏感词库，现在这里分开）
	TABLE_SENSITIVE_WORD_NAME     = "hn_sensitive_word_name"
	TABLE_COURSES                 = "hn_courses"
	TABLE_COURSE_LESSON           = "hn_course_lessons"
	TABLE_COURSE_DEMO             = "hn_course_demo"          //课堂改版添加 zxw 20170814 add
	TABLE_COURSE_GOLDEN_STOCKS    = "hn_course_golden_stocks" //课堂改版添加 zxw 20170812 add
	TABLE_ASSEMBLE                = "hn_assemble"
	TABLE_BLACKLIST_IP            = "hn_blacklist_ip"
)

// Visibility Range Level
// --------------------------------------------------------------------------------

const (
	VISIBILITY_PUBLIC  = 0 // 公开
	VISIBILITY_PRIVATE = 1 // 私有
)

// Mobile MD5 suffix
//---------------------------------------------------------------------------------

const (
	MOBILE_ENCRYPT_SUFFIX = "7dbssea4OBILE"
)

//
//---------------------------------------------------------------------------------

const (
	SNS_TYPE_QQ     = 1 // qq
	SNS_TYPE_SINA   = 2 // sina
	SNS_TYPE_WECHAT = 3 // wechat
	SNS_TYPE_JD     = 4 // jd

	SNS_TYPE_SINA_STR   = "sina"
	SNS_TYPE_QQ_STR     = "qq"
	SNS_TYPE_WECHAT_STR = "wechat"
	SNS_TYPE_JD_STR     = "jd"
)

// Advisor Type
//---------------------------------------------------------------------------------

const (
	ADVISOR_TYPE_STAFF    = 0 // 普通会员
	ADVISOR_TYPE_INTERNAL = 1 // 内部员工
	ADVISOR_TYPE_ADVISOR  = 2 // 持牌投顾
)

// Tactic Status
//---------------------------------------------------------------------------------

const (
	TACTIC_STATUS_DISABLE = 0 // 关闭
	TACTIC_STATUS_ENABLED = 1 // 开启
)

// Research Report Status
//---------------------------------------------------------------------------------

const (
	REPORT_STATUS_DISABLE = 0 // 禁用
	REPORT_STATUS_ENABLED = 1 // 启用
)

// Notice Event Categorys
// --------------------------------------------------------------------------------

const (
	NOTICE_EVENT_DELETE_MESSAGE = "delete-message"
	NOTICE_EVENT_JOIN_GROUP     = "join-group"
	NOTICE_EVENT_LEAVE_GROUP    = "leave-group"
	NOTICE_EVENT_PULL_MESSAGE   = "pull-message"
	NOTICE_EVENT_UNDO_MESSAGE   = "undo-message"
	NOTICE_EVENT_START_WEBLIVE  = "start-weblive"
	NOTICE_EVENT_END_WEBLIVE    = "end-weblive"
	NOTICE_EVENT_NOTIFY         = "pull-notify"
	NOTICE_EVENT_NEWS           = "pull-news"
	NOTICE_EVENT_SYSTEM         = "pull-system"
)

// Category Assort
// --------------------------------------------------------------------------------

const (
	CATEGORY_ASSORT_FINANCE      = 0 // 财经头条
	CATEGORY_ASSORT_MAGAZINE     = 1 // 首证期刊
	CATEGORY_ASSORT_LIVING       = 2 // 直播
	CATEGORY_ASSORT_TREASURE_BOX = 3 // 百宝箱
	CATEGORY_ASSORT_PROPOSAL     = 4 // 投诉建议
)

//
//---------------------------------------------------------------------------------

const OPINION_INTRO_LENGTH = 200  // Opinion Intro Length
const NOTIFY_CONTENT_LENGTH = 100 // Notify Content Length
const CMS_INTRO_LENGTH = 200      // CMS Intro Length

// Risk Level
//---------------------------------------------------------------------------------

const (
	RISK_LEVEL_HATE_RISK    = 1 // 厌恶风险型
	RISK_LEVEL_CONSERVATIVE = 2 // 保守型
	RISK_LEVEL_STEADY       = 3 // 稳健型
	RISK_LEVEL_RIDACAL      = 4 // 激进型
)

// Order Status
//---------------------------------------------------------------------------------

const (
	ORDER_STATUS_UNPAID    = 1 // 待支付
	ORDER_STATUS_PAID      = 2 // 已支付
	ORDER_STATUS_CANCELED  = 3 // 已取消
	ORDER_STATUS_TO_REFUND = 4 // 提交退款
	ORDER_STATUS_REFUNDING = 5 // 退款中
	ORDER_STATUS_REFUNDED  = 6 // 已退款
	ORDER_STATUS_IN_REVIEW = 7 // 审核中
)

// Consume Origin
//---------------------------------------------------------------------------------

const (
	CONSUME_ORIGIN_PC     = 1 // PC端
	CONSUME_ORIGIN_WEB    = 2 // Web端
	CONSUME_ORIGIN_MOBILE = 3 // 移动端
)

// Treasure State
//---------------------------------------------------------------------------------

const (
	TREASURE_STATE_NOT_RUNNING = 1 // 待运行
	TREASURE_STATE_RUNNING     = 2 // 运行中
	TREASURE_STATE_ENDED       = 3 // 已结束
	TREASURE_STATE_STOP_SELL   = 4 // 停售
)

/// Treasure Ruuning State Version 2
//---------------------------------------------------------------------------------

const (
	TREASURE_RUN_STATE_BEFORE_SELLING = 1 // 预售中
	TREASURE_RUN_STATE_RUNNING        = 2 // 运行中
	TREASURE_RUN_STATE_STOP_SELL      = 3 // 已停售
	TREASURE_RUN_STATE_ENDED          = 4 // 已结束
)

// SNS Operate Type
//---------------------------------------------------------------------------------

const (
	SNS_OPERATE_TYPE_BIND  = "binduser" // 账号绑定
	SNS_OPERATE_TYPE_LOGIN = "login"    // 账号登陆
)

// Group Assort
//---------------------------------------------------------------------------------

const (
	GROUP_ASSORT_MEMBER_GROUP  = 0 // 会员组
	GROUP_ASSORT_ADVISOR_LEVEL = 1 // 投顾等级
)

// JPush Alias and Tags
// --------------------------------------------------------------------------------

const (
	JPUSH_TAG_MEMBER = "member_%v"
	JPUSH_TAG_GROUP  = "group_%v"
)

const (
	REFTYPE_JPUSH_TAG_NOTIFY = 1 // 提醒
	REFTYPE_JPUSH_TAG_NEWS   = 8 // 资讯
)

// Affix Url
//---------------------------------------------------------------------------------

const (
	AFFIX_URL = "http://attach.0606.com.cn/" // 附件URL
)

// Date
//---------------------------------------------------------------------------------

const (
	ONE_YEAR_DAYS   float32 = 365
	ONE_DAY_SECONDS int64   = 24 * 3600
)

// Number
//---------------------------------------------------------------------------------

const (
	ONE_BILLION = 1000000000
)

// Group Weight
//---------------------------------------------------------------------------------

const (
	MEMBER_GROUP_WEIGHT_VIP0 = 1
	MEMBER_GROUP_WEIGHT_VIP1 = 2
	MEMBER_GROUP_WEIGHT_VIP2 = 3
	MEMBER_GROUP_WEIGHT_VIP3 = 4
)

// Proposal Status
//---------------------------------------------------------------------------------

const (
	PROPOSAL_STATUS_UNSOLVE = 1 // 已受理
	PROPOSAL_STATUS_SOLVING = 2 // 处理中
	PROPOSAL_STATUS_SOLVED  = 0 // 已解决
)

// Advert Status
//---------------------------------------------------------------------------------

const (
	ADVERT_STATUS_DISABLE = 0 // 禁用
	ADVERT_STATUS_ENABLE  = 1 // 启用
)

// Strategy Status
//---------------------------------------------------------------------------------

const (
	STRATEGY_STATUS_DISABLE = 0 // 禁用
	STRATEGY_STATUS_ENABLE  = 1 // 启用
)

// TreasureBox Status
//---------------------------------------------------------------------------------

const (
	TREASUREBOX_STATUS_DISABLE   = 0 // 禁用
	TREASUREBOX_STATUS_ENABLE    = 1 // 启用
	TREASUREBOX_STATUS_STOP_SELL = 2 // 停售
)

// Message Channel Categorys
// --------------------------------------------------------------------------------

const (
	CHANNEL_PRIVATE_GROUP  = "private-channel-%s"
	CHANNEL_PRIVATE_MEMBER = "private-member-%s"
	CHANNEL_PRIVATE_LIVE   = "private-live-%s"
)

// QA Types
// --------------------------------------------------------------------------------

const (
	LIB_QA_TYPE_TEXT = 1 // 文本类
	LIB_QA_TYPE_URL  = 2 // 链接类
	LIB_QA_TYPE_NEWS = 3 // 新闻类
	LIB_QA_TYPE_MENU = 4 // 菜谱类
)

// Course Lesson Assort Categorys
// --------------------------------------------------------------------------------

const (
	LESSON_ASSORT_NORMAL = 0 // 其他
	LESSON_ASSORT_SYSTEM = 1 // 系统课
	LESSON_ASSORT_ACTION = 2 // 实战课
)

// Course Lesson LiveStatus Categorys
// --------------------------------------------------------------------------------

const (
	LESSON_LIVE_STATUS_OTHER = 0 // 其他
	LESSON_LIVE_STATUS_START = 1 // 直播开始
	LESSON_LIVE_STATUS_END   = 2 // 直播结束
)

// Course Type

const (
	COURSE_TYPE_PUBLIC  = 1 // 公开课
	COURSE_TYPE_PRIVATE = 2 // 精品课
)

// Assemble Status
//---------------------------------------------------------------------------------

const (
	ASSEMBLE_STATUS_DISABLE   = 0 // 禁用
	ASSEMBLE_STATUS_ENABLE    = 1 // 启用
	ASSEMBLE_STATUS_STOP_SELL = 2 // 停售
)
