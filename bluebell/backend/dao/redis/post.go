package redis

import "bluebell/models"

func GetPostIDsInOrder(p *models.ParamPostList)([]string,error) {
	// 从redis获取id
	// 1.根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZset)
	if p.Order == models.OrderScore{
		key = getRedisKey(KeyPostScoreZset)
	}
	// 2.确定查询的索引起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	// 3.ZREVRANGE 按分数从大到小的顺序查询指定数量的元素
	return client.ZRevRange(key,start,end).Result()
}