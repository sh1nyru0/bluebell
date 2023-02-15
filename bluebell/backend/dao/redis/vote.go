package redis

import (
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis"
)

/* 投票的集中情况：
direcation = 1时，有两种情况:
	1.之前没有投过票，现在投赞成票 		-->更新分数和投票记录 差值的绝对值：1 +432
	2.之前投反对票，现在改投赞成票		-->更新分数和投票记录 差值的绝对值：2 +432*2
direcation = 0时，有两中情况:
	1.之前投过赞成票，现在要取消投票	-->更新分数和投票记录 差值的绝对值：1 -432
	2.之前投过反对派，现在要取消投票	-->更新分数和投票记录 差值的绝对值：1 +432
direcation = -1时，有两种情况:
	1.之前没有投过票，现在投反对票		-->更新分数和投票记录 差值的绝对值：1 -432
	2.之前投赞成票，现在改投反对票		-->更新分数和投票记录 差值的绝对值：2 -432*2

投票的限制：
每个帖子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	1.到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
	2.到期之后删除那个 KeyPostVotedZSetPF
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote = 432 // 每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
)


func CreatePost(postID int64) error {
	pipeline := client.TxPipeline()
	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	
	_, err := pipeline.Exec()
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	// 1. 判断投票限制
	// 去redis取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZset), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2和3需要放到一个pipeline事务中操作

	// 2. 更新贴子的分数
	// 先查当前用户给当前帖子的投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedZsetPF+postID), userID).Val()
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) // 计算两次投票的差值
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZset), op*diff*scorePerVote, postID)

	// 3. 记录用户为该贴子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZsetPF+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZsetPF+postID), redis.Z{
			Score:  value, // 赞成票还是反对票
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}