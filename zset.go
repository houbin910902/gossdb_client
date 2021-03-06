package gossdb_client

import (
	"fmt"
	"strconv"
	"github.com/houbin910902/to"
)

//  设置 zset 中指定 key 对应的权重值.
//  setName zset 名称
//  key zset 中的 key.
//  score 整数, key 对应的权重值
//  返回 err, 可能的错误, 操作成功返回 nil
func (c *DbClient) ZSet(setName, key string, score int64) (err error) {
	resp, err := c.Client.Do("zset", setName, key, score)
	if err != nil {
		return fmt.Errorf("ZSet %s %s error: %s", setName, key, err.Error())
	}
	if len(resp) > 0 && resp[0] == "ok" {
		return nil
	}
	return handError(resp, setName, key)
}


//  获取zset中指定 key 对应的权重值.
//
//  setName zset名称
//  key zset 中的 key.
//  返回 score 整数, key 对应的权重值
//  返回 err, 可能的错误, 操作成功返回 nil
func (c *DbClient) ZGet(setName, key string) (score int64, err error) {
	resp, err := c.Client.Do("zget", setName, key)
	if err != nil {
		return 0, fmt.Errorf("ZGet %s/%s error: %s", setName, key, err.Error())
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return strconv.ParseInt(resp[1], 10, 64)
	}
	return 0, handError(resp, setName, key)
}

//删除 zset 中指定 key
//
//  setName zset名称
//  key zset 中的 key.
//  返回 err, 可能的错误, 操作成功返回 nil
func (c *DbClient) ZDel(setName, key string) (err error) {
	resp, err := c.Client.Do("zdel", setName, key)
	if err != nil {
		return fmt.Errorf("ZDel %s/%s error: %s", setName, key, err.Error())
	}
	if len(resp) > 0 && resp[0] == "ok" {
		return nil
	}
	return handError(resp, setName, key)
}

//判断指定的 key 是否存在于 zset 中.
//  setName zset名称
//  key zset 中的 key.
//  返回 re 如果存在, 返回 true, 否则返回 false.
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZExists(setName, key string) (re bool, err error) {
	resp, err := c.Client.Do("zexists", setName, key)
	if err != nil {
		return false, fmt.Errorf("ZExists %s/%s error: %s", setName, key, err.Error())
	}

	if len(resp) == 2 && resp[0] == "ok" {
		return resp[1] == "1", nil
	}
	return false, handError(resp, setName, key)
}

//  返回处于区间 [start,end] key 数量.
//  setName zset名称
//  start key 的最小权重值(包含), 空字符串表示 -inf.
//  end key 的最大权重值(包含), 空字符串表示 +inf.
//  返回 count 返回符合条件的 key 的数量.
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZCount(setName string, start, end interface{}) (count int64, err error) {
	resp, err := c.Client.Do("zcount", setName, start, end)
	if err != nil {
		return -1, fmt.Errorf("ZCount %s %v %v error: %s", setName, start, end, err)
	}

	if len(resp) == 2 && resp[0] == "ok" {
		return strconv.ParseInt(resp[1], 10, 64)
	}
	return -1, handError(resp, setName)
}

//  删除 zset 中的所有 key.
//  setName zset名称
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZClear(setName string) (err error) {
	resp, err := c.Client.Do("zclear", setName)
	if err != nil {
		return fmt.Errorf("%s ZClear %s error: %s", err, setName, err.Error())
	}

	if len(resp) > 0 && resp[0] == "ok" {
		return nil
	}
	return handError(resp, setName)
}


//  列出 zset 中处于区间 (key_start+score_start, score_end] 的 key-score 列表.
//  如果 key_start 为空, 那么对应权重值大于或者等于 score_start 的 key 将被返回. 如果 key_start 不为空, 那么对应权重值大于 score_start 的 key, 或者大于 key_start 且对应权重值等于 score_start 的 key 将被返回.
//  也就是说, 返回的 key 在 (key.score == score_start && key > key_start || key.score > score_start), 并且 key.score <= score_end 区间. 先判断 score_start, score_end, 然后判断 key_start.
//
//  setName zset名称
//  keyStart score_start 对应的 key.
//  scoreStart 返回 key 的最小权重值(可能不包含, 依赖 key_start), 空字符串表示 -inf.
//  scoreEnd 返回 key 的最大权重值(包含), 空字符串表示 +inf.
//  limit  最多返回这么多个元素.
//  返回 keys 返回符合条件的 key 的数组.
//  返回 scores 返回符合条件的 key 对应的权重.
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZScan(setName string, keyStart string, scoreStart, scoreEnd interface{}, limit int64) (keys []string, scores []int64, err error) {
	resp, err := c.Client.Do("zscan", setName, keyStart, scoreStart, scoreEnd, limit)

	if err != nil {
		return nil, nil, fmt.Errorf("ZScan %s %v %v %v %v error: %s", setName, keyStart, scoreStart, scoreEnd, limit, err.Error())
	}
	if len(resp) > 0 && resp[0] == "ok" {
		size := len(resp)
		keys := make([]string, 0, (size-1)/2)
		scores := make([]int64, 0, (size-1)/2)

		for i := 1; i < size-1; i += 2 {
			keys = append(keys, resp[i])
			sco, _:= strconv.ParseInt(resp[1], 10, 64)
			scores = append(scores, sco)
		}
		return keys, scores, nil
	}
	return nil, nil, handError(resp, setName, keyStart, scoreStart, scoreEnd, limit)
}

//  列出 zset 中的 key-score 列表, 反向顺序
//  setName zset名称
//  keyStart score_start 对应的 key.
//  scoreStart 返回 key 的最小权重值(可能不包含, 依赖 key_start), 空字符串表示 -inf.
//  scoreEnd 返回 key 的最大权重值(包含), 空字符串表示 +inf.
//  limit  最多返回这么多个元素.
//  返回 keys 返回符合条件的 key 的数组.
//  返回 scores 返回符合条件的 key 对应的权重.
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZrScan(setName string, keyStart string, scoreStart, scoreEnd interface{}, limit int64) (keys []string, scores []int64, err error) {
	resp, err := c.Client.Do("zrscan", setName, keyStart, scoreStart, scoreEnd, limit)

	if err != nil {
		return nil, nil, fmt.Errorf("ZrScan %s %v %v %v %v error: %s", setName, keyStart, scoreStart, scoreEnd, limit, err.Error())
	}

	if len(resp) > 0 && resp[0] == "ok" {
		size := len(resp)
		keys := make([]string, 0, (size-1)/2)
		scores := make([]int64, 0, (size-1)/2)

		for i := 1; i < size-1; i += 2 {
			keys = append(keys, resp[i])
			sco, _:= strconv.ParseInt(resp[1], 10, 64)
			scores = append(scores, sco)
		}
		return keys, scores, nil
	}
	return nil, nil, handError(resp, setName, keyStart, scoreStart, scoreEnd, limit)
}

//批量设置 zset 中的 key-score.
//
//  setName zset名称
//  kvs 包含 key-score 的map
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) MultiZSet(setName string, kvs map[string]int64) (err error) {

	var args []interface{}
	for k, v := range kvs {
		args = append(args, k)
		args = append(args, v)
	}
	resp, err := c.Client.Do("multi_zset", setName, args)

	if err != nil {
		return fmt.Errorf("MultiZset %s %v error: %s", setName, kvs, err.Error())
	}

	if len(resp) > 0 && resp[0] == "ok" {
		return nil
	}
	return handError(resp, setName, kvs)
}

//批量获取 zset 中的 key-score.
//
//  setName zset名称
//  key 要获取key的列表，支持多个key
//  返回 val 包含 key-score 的map
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) MultiZGet(setName string, key ...string) (val map[string]int64, err error) {
	if len(key) == 0 {
		return make(map[string]int64), nil
	}
	resp, err := c.Client.Do("multi_zget", setName, key)

	if err != nil {
		return nil, fmt.Errorf("MultiZget %s %s error: %s", setName, key, err.Error())
	}
	size := len(resp)
	if size > 0 && resp[0] == "ok" {
		val = make(map[string]int64)
		for i := 1; i < size && i+1 < size; i += 2 {
			sco, _:= strconv.ParseInt(resp[i+1], 10, 64)
			val[resp[i]] = sco
		}
		return val, nil
	}
	return nil, handError(resp, key)
}

//批量获取 zset 中的 key-score.
//
//  setName zset名称
//  key 要获取key的列表，支持多个key
//  返回 keys 包含 key的slice
//  返回 scores 包含 key对应权重的slice
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) MultiZGetSlice(setName string, key ...string) (keys []string, scores []int64, err error) {
	if len(key) == 0 {
		return []string{}, []int64{}, nil
	}
	resp, err := c.Client.Do("multi_zget", setName, key)

	if err != nil {
		return nil, nil, fmt.Errorf("MultiZGetSlice %s %s error: %s", setName, key, err.Error())
	}

	size := len(resp)
	if size > 0 && resp[0] == "ok" {

		keys := make([]string, (size-1)/2)
		scores := make([]int64, (size-1)/2)

		for i := 1; i < size && i+1 < size; i += 2 {
			keys = append(keys, resp[i])
			scores = append(scores, to.Int64(resp[i+1]))
		}
		return keys, scores, nil
	}
	return nil, nil, handError(resp, setName, key)
}

//  批量获取 zset 中的 key-score.
//  setName zset名称
//  key 要获取key的slice
//  返回 val 包含 key-score 的map
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) MultiZGetArray(setName string, key []string) (val map[string]int64, err error) {
	if len(key) == 0 {
		return make(map[string]int64), nil
	}
	resp, err := c.Client.Do("multi_zget", setName, key)

	if err != nil {
		return nil, fmt.Errorf("MultiZGetArray %s %s error: %s", setName, key, err.Error())
	}
	size := len(resp)
	if size > 0 && resp[0] == "ok" {
		val = make(map[string]int64)
		for i := 1; i < size && i+1 < size; i += 2 {
			val[resp[i]] = to.Int64(resp[i+1])
		}
		return val, nil
	}
	return nil, handError(resp, key)
}

//批量获取 zset 中的 key-score.
//
//  setName zset名称
//  key 要获取key的slice
//  返回 keys 包含 key的slice
//  返回 scores 包含 key对应权重的slice
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) MultiZGetSliceArray(setName string, key []string) (keys []string, scores []int64, err error) {
	if len(key) == 0 {
		return []string{}, []int64{}, nil
	}
	resp, err := c.Client.Do("multi_zget", setName, key)

	if err != nil {
		return nil, nil, fmt.Errorf("MultiZGetSliceArray %s %s error: %s", setName, key, err.Error())
	}

	size := len(resp)
	if size > 0 && resp[0] == "ok" {

		keys := make([]string, (size-1)/2)
		scores := make([]int64, (size-1)/2)

		for i := 1; i < size && i+1 < size; i += 2 {
			keys = append(keys, resp[i])
			scores = append(scores, to.Int64(resp[i+1]))
		}
		return keys, scores, nil
	}
	return nil, nil, handError(resp, setName, key)
}


//批量删除 zset 中的 key-score.
//
//  setName zset名称
//  key 要删除key的列表，支持多个key
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) MultiZDel(setName string, key ...string) (err error) {
	if len(key) == 0 {
		return nil
	}
	resp, err := c.Client.Do("multi_zdel", key)

	if err != nil {
		return fmt.Errorf("MultiZDel %s %s error: %s", setName, key, err.Error())
	}

	if len(resp) > 0 && resp[0] == "ok" {
		return nil
	}
	return handError(resp, setName, key)
}

//使 zset 中的 key 对应的值增加 num. 参数 num 可以为负数.
//
//  setName zset名称
//  key 要增加权重的key
//  num 要增加权重值
//  返回 int64 增加后的新权重值
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZIncR(setName string, key string, num int64) (int64, error) {
	if len(key) == 0 {
		return 0, nil
	}
	resp, err := c.Client.Do("zincr", setName, key, num)
	if err != nil {
		return 0, fmt.Errorf("ZIncR %s %s %v error: %s", setName, key, num, err.Error())
	}

	if len(resp) > 1 && resp[0] == "ok" {
		return to.Int64(resp[1]), nil
	}
	return 0, handError(resp, setName, key)
}

//列出名字处于区间 (name_start, name_end] 的 zset.
//
//  name_start - 返回的起始名字(不包含), 空字符串表示 -inf.
//  name_end - 返回的结束名字(包含), 空字符串表示 +inf.
//  limit  最多返回这么多个元素.
//  返回 []string 返回包含名字的slice.
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZList(nameStart, nameEnd string, limit int64) ([]string, error) {
	resp, err := c.Client.Do("zlist", nameStart, nameEnd, limit)
	if err != nil {
		return nil, fmt.Errorf("ZList %s %s %v error: %s", nameStart, nameEnd, limit, err.Error())
	}
	if len(resp) > 0 && resp[0] == "ok" {
		size := len(resp)
		keyList := make([]string, 0, size-1)

		for i := 1; i < size; i += 1 {
			keyList = append(keyList, resp[i])
		}
		return keyList, nil
	}
	return nil, handError(resp, nameStart, nameEnd, limit)
}

//返回 zset 中的元素个数.
//
//  name zset的名称.
//  返回 val 返回包含名字元素的个数.
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZSize(name string) (val int64, err error) {
	resp, err := c.Client.Do("zsize", name)
	if err != nil {
		return 0, fmt.Errorf("ZSize %s  error: %s", name, err.Error())
	}

	if len(resp) > 0 && resp[0] == "ok" {
		val = to.Int64(resp[1])
		return val, nil
	}
	return 0, handError(resp, name)
}

//列出 zset 中的 key 列表. 参见 zscan().
//
//  setName zset名称
//  keyStart score_start 对应的 key.
//  scoreStart 返回 key 的最小权重值(可能不包含, 依赖 key_start), 空字符串表示 -inf.
//  scoreEnd 返回 key 的最大权重值(包含), 空字符串表示 +inf.
//  limit  最多返回这么多个元素.
//  返回 keys 返回符合条件的 key 的数组.
//  返回 scores 返回符合条件的 key 对应的权重.
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZKeys(setName string, keyStart string, scoreStart, scoreEnd interface{}, limit int64) (keys []string, err error) {
	resp, err := c.Client.Do("zkeys", setName, keyStart, scoreStart, scoreEnd, limit)

	if err != nil {
		return nil, fmt.Errorf("ZKeys %s %v %v %v %v error: %s", setName, keyStart, scoreStart, scoreEnd, limit, err.Error())
	}
	if len(resp) > 0 && resp[0] == "ok" {
		size := len(resp)
		keys := []string{}

		for i := 1; i < size; i++ {
			keys = append(keys, resp[i])
		}
		return keys, nil
	}
	return nil, handError(resp, setName, keyStart, scoreStart, scoreEnd, limit)
}

//  返回 key 处于区间 [start,end] 的 score 的和.
//  setName zset名称
//  scoreStart  key 的最小权重值(可能不包含, 依赖 key_start), 空字符串表示 -inf.
//  scoreEnd  key 的最大权重值(包含), 空字符串表示 +inf.
//  返回 val 符合条件的 score 的求和
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZSum(setName string, scoreStart, scoreEnd interface{}) (val int64, err error) {
	resp, err := c.Client.Do("zsum", setName, scoreStart, scoreEnd)

	if err != nil {
		return 0, fmt.Errorf("ZSum %s %v %v  error: %s", setName, scoreStart, scoreEnd, err.Error())
	}
	if len(resp) > 0 && resp[0] == "ok" {
		val = to.Int64(resp[1])
		return val, nil
	}
	return 0, handError(resp, setName, scoreStart, scoreEnd)
}

//返回 key 处于区间 [start,end] 的 score 的平均值.
//
//  setName zset名称
//  scoreStart  key 的最小权重值(可能不包含, 依赖 key_start), 空字符串表示 -inf.
//  scoreEnd  key 的最大权重值(包含), 空字符串表示 +inf.
//  返回 val 符合条件的 score 的平均值
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZAvg(setName string, scoreStart, scoreEnd interface{}) (val int64, err error) {
	resp, err := c.Client.Do("zavg", setName, scoreStart, scoreEnd)

	if err != nil {
		return 0, fmt.Errorf("ZAvg %s %v %v  error: %s", setName, scoreStart, scoreEnd, err.Error())
	}
	if len(resp) > 0 && resp[0] == "ok" {
		val = to.Int64(resp[1])
		return val, nil
	}
	return 0, handError(resp, setName, scoreStart, scoreEnd)
}

//返回指定 key 在 zset 中的排序位置(排名), 排名从 0 开始. 注意! 本方法可能会非常慢! 请在离线环境中使用.
//
//  setName zset名称
//  key 指定key名
//  返回 val 排名
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZRank(setName, key string) (val int64, err error) {
	resp, err := c.Client.Do("zrank", setName, key)

	if err != nil {
		return 0, fmt.Errorf("ZRank %s %s  error: %s", setName, key, err.Error())
	}
	if len(resp) > 0 && resp[0] == "ok" {
		val = to.Int64(resp[1])
		return val, nil
	}
	return 0, handError(resp, setName, key)
}

//返回指定 key 在 zset 中的倒序排名.注意! 本方法可能会非常慢! 请在离线环境中使用.
//
//  setName zset名称
//  key 指定key名
//  返回 val 排名
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZRRank(setName, key string) (val int64, err error) {
	resp, err := c.Client.Do("zrrank", setName, key)

	if err != nil {
		return 0, fmt.Errorf("ZRRank %s %s  error: %s", setName, key, err.Error())
	}
	if len(resp) > 0 && resp[0] == "ok" {
		val = to.Int64(resp[1])
		return val, nil
	}
	return 0, handError(resp, setName, key)
}

//根据下标索引区间 [offset, offset + limit) 获取 key-score 对, 下标从 0 开始.注意! 本方法在 offset 越来越大时, 会越慢!
//
//  setName zset名称
//  offset 从此下标处开始返回. 从 0 开始.
//  limit  最多返回这么多个 key-score 对.
//  返回 val 排名
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZRange(setName string, offset, limit int64) (val map[string]int64, err error) {
	resp, err := c.Client.Do("zrange", setName, offset, limit)

	if err != nil {
		return nil, fmt.Errorf("ZRange %s %v %v  error: %s", setName, offset, limit, err.Error())
	}
	if len(resp) > 0 && resp[0] == "ok" {
		val = make(map[string]int64)
		size := len(resp)
		for i := 1; i < size-1; i += 2 {
			val[resp[i]] = to.Int64(resp[i+1])
		}
		return val, nil
	}
	return nil, handError(resp, setName, offset, limit)
}

//根据下标索引区间 [offset, offset + limit) 获取 获取 key和score 数组对, 下标从 0 开始.注意! 本方法在 offset 越来越大时, 会越慢!
//
//  setName zset名称
//  offset 从此下标处开始返回. 从 0 开始.
//  limit  最多返回这么多个 key-score 对.
//  返回 val 排名
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZRangeSlice(setName string, offset, limit int64) (key []string, val []int64, err error) {
	resp, err := c.Client.Do("zrange", setName, offset, limit)

	if err != nil {
		return nil, nil, fmt.Errorf("ZRangeSlice %s %v %v  error: %s", setName, offset, limit, err.Error())
	}
	if len(resp) > 0 && resp[0] == "ok" {
		val = []int64{}
		key = []string{}
		size := len(resp)
		for i := 1; i < size-1; i += 2 {
			key = append(key, resp[i])
			val = append(val, to.Int64(resp[i+1]))
		}
		return key, val, nil
	}
	return nil, nil, handError(resp, setName, offset, limit)
}

//根据下标索引区间 [offset, offset + limit) 获取 key-score 对, 反向顺序获取.注意! 本方法在 offset 越来越大时, 会越慢!
//
//  setName zset名称
//  offset 从此下标处开始返回. 从 0 开始.
//  limit  最多返回这么多个 key-score 对.
//  返回 val 排名
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZRRange(setName string, offset, limit int64) (val map[string]int64, err error) {
	resp, err := c.Client.Do("zrrange", setName, offset, limit)

	if err != nil {
		return nil, fmt.Errorf("ZRRange %s %v %v error: %s", setName, offset, limit, err.Error())
	}
	if len(resp) > 0 && resp[0] == "ok" {
		val = make(map[string]int64)
		size := len(resp)

		for i := 1; i < size-1; i += 2 {
			val[resp[i]] = to.Int64(resp[i+1])
		}
		return val, nil
	}
	return nil, handError(resp, setName, offset, limit)
}


//  根据下标索引区间 [offset, offset + limit) 获取 key和score 数组对, 反向顺序获取.注意! 本方法在 offset 越来越大时, 会越慢!
//  setName zset名称
//  offset 从此下标处开始返回. 从 0 开始.
//  limit  最多返回这么多个 key-score 对.
//  返回 val 排名
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZRRangeSlice(setName string, offset, limit int64) (key []string, val []int64, err error) {
	resp, err := c.Client.Do("zrrange", setName, offset, limit)

	if err != nil {
		return nil, nil, fmt.Errorf("ZRRangeSlice %s %v %v error: %s", setName, offset, limit, err.Error())
	}
	if len(resp) > 0 && resp[0] == "ok" {
		val = []int64{}
		key = []string{}
		size := len(resp)

		for i := 1; i < size-1; i += 2 {
			key = append(key, resp[i])
			val = append(val, to.Int64(resp[i+1]))
		}
		return key, val, nil
	}
	return nil, nil, handError(resp, setName, offset, limit)
}

//  删除位置处于区间 [start,end] 的元素.
//  setName zset名称
//  start 区间开始，包含start值
//  end  区间结束，包含end值
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZRemRangeByRank(setName string, start, end int64) (err error) {
	resp, err := c.Client.Do("zremrangebyrank", setName, start, end)

	if err != nil {
		return fmt.Errorf("ZRemRangeByRank %s %v %v error: %s", setName, start, end, err.Error())
	}
	if len(resp) > 0 && resp[0] == "ok" {
		return nil
	}
	return handError(resp, setName, start, end)
}

//  删除权重处于区间 [start,end] 的元素.
//  setName zset名称
//  start 区间开始，包含start值
//  end  区间结束，包含end值
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZRemRangeByScore(setName string, start, end int64) (err error) {
	resp, err := c.Client.Do("zremrangebyscore", setName, start, end)

	if err != nil {
		return fmt.Errorf("ZRemRangeByScore %s %v %v  error: %s", setName, start, end, err.Error())
	}
	if len(resp) > 0 && resp[0] == "ok" {
		return nil
	}
	return handError(resp, setName, start, end)
}

//  从 zset 首部删除并返回 `limit` 个元素.
//  setName zset名称
//  limit 最多要删除并返回这么多个 key-score 对.
//  返回 包含 key-score 的map
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZPopFront(setName string, limit int64) (val map[string]int64, err error) {
	resp, err := c.Client.Do("zpop_front", setName, limit)

	if err != nil {
		return nil, fmt.Errorf("ZPopFront %s %v  error: %s", setName, limit, err.Error())
	}
	size := len(resp)
	if size > 0 && resp[0] == "ok" {
		val = make(map[string]int64)
		for i := 1; i < size && i+1 < size; i += 2 {
			val[resp[i]] = to.Int64(resp[i+1])
		}
		return val, nil
	}
	return nil, handError(resp, setName, limit)
}

//  从 zset 尾部删除并返回 `limit` 个元素.
//  setName zset名称
//  limit 最多要删除并返回这么多个 key-score 对.
//  返回 包含 key-score 的map
//  返回 err，可能的错误，操作成功返回 nil
func (c *DbClient) ZPopBack(setName string, limit int64) (val map[string]int64, err error) {
	resp, err := c.Client.Do("zpop_back", setName, limit)

	if err != nil {
		return nil, fmt.Errorf("ZPopBack %s %v  error: %s", setName, limit, err.Error())
	}
	size := len(resp)
	if size > 0 && resp[0] == "ok" {
		val = make(map[string]int64)
		for i := 1; i < size && i+1 < size; i += 2 {
			val[resp[i]] = to.Int64(resp[i+1])
		}
		return val, nil
	}
	return nil, handError(resp, setName, limit)
}
