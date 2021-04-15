package common

import (
    "fmt"

/share/models"

	"stock
/share/store/redis"
)

type RankChart struct {
}

func NewRankChart() *RankChart {
	return &RankChart{}
}

type ChartType int

const (
	_                 ChartType = iota // 类型定义
	ChartTypeOpinion                   // 观点排行榜
	ChartTypeQuestion                  // 问答排行榜
)

const (
	INCREASE_CHART = 1 // INCR操作
	DECREASE_CHAR  = 2 // DECR操作
)

type ChartData struct {
	MemberID int64
	Count    int
}

func (this *RankChart) InDecrCount(memberId int64, chartType ChartType, operate int) error {
	var keyCache string
	var zsetKey string

	switch chartType {
	case ChartTypeOpinion:
		keyCache = fmt.Sprintf(REDIS_ADVISOR_OPINIONS, memberId)
		zsetKey = REDIS_WEBLIVE_CHARTS_OPINIONS
	case ChartTypeQuestion:
		keyCache = fmt.Sprintf(REDIS_ADVISOR_QUESTIONS, memberId)
		zsetKey = REDIS_WEBLIVE_CHARTS_QUESTIONS
	default:
		return ErrParameterError
	}

	var op string
	var increment int

	switch operate {
	case INCREASE_CHART:
		op = "INCR"
		increment = 1
	case DECREASE_CHAR:
		op = "DECR"
		increment = -1

	default:
		return ErrParameterError
	}

	// 调整单项数量(观点数或问答数)
	count, err := redis.Do(op, keyCache)
	if err != nil {
		return err
	}

	// 调整单项排行榜(观点排行或问答排行)
	exits, err := redis.Exists(zsetKey)
	if err != nil {
		return err
	}

	if exits {
		var score int
		score = int(count.(int64))

		if err := redis.Zadd(zsetKey, score, memberId); err != nil {
			return err
		}
	}

	// 调整综合排行榜
	exits, err = redis.Exists(REDIS_WEBLIVE_CHARTS_COMPREHENSIVE)
	if err != nil {
		return err
	}

	if exits {
		return redis.Zincrby(REDIS_WEBLIVE_CHARTS_COMPREHENSIVE, increment, memberId)
	}

	return nil
}
