package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/mapping"
	"strconv"
	"time"
)

type (
	Client interface {
		redis.Cmdable
	}

	clientMode func(conf Config) Client

	Pair struct {
		Key   string
		Score int64
	}

	Z              = redis.Z
	ZStore         = redis.ZStore
	GeoLocation    = redis.GeoLocation
	GeoRadiusQuery = redis.GeoRadiusQuery
	GeoPos         = redis.GeoPos
)

const (
	blockingQueryTimeout = 5 * time.Second
)

func Single() clientMode {
	return func(conf Config) Client {
		return redis.NewClient(&redis.Options{
			Addr:     conf.Addrs[0],
			Password: conf.Password,
			DB:       conf.DB,
		})
	}
}

func Sentinel() clientMode {
	return func(conf Config) Client {
		return redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    conf.MasterName,
			SentinelAddrs: conf.Addrs,
			Password:      conf.Password,
		})
	}
}

func Cluster() clientMode {
	return func(conf Config) Client {
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    conf.Addrs,
			Password: conf.Password,
		})
	}
}

type RedisConn interface {
	// Set args: 0->mode
	Set(key string, val interface{}, args ...int) error
	// SetNX args: 0->mode
	SetNX(key, val string, args ...int) (bool, error)
	Get(key string) (string, error)
	GetSet(key, val string) (string, error)
	Incr(key string) (int64, error)
	Incrby(key string, increment int64) (int64, error)
	IncrByFloat(key string, increment float64) (float64, error)
	Decr(key string) (int64, error)
	DecrBy(key string, decrement int64) (int64, error)
	MGet(keys ...string) ([]string, error)

	HSet(key, field, val string) error
	HSetNX(key, field, val string) (bool, error)
	HGet(key, field string) (string, error)
	HExists(key, field string) (bool, error)
	HDel(key string, fields ...string) (bool, error)
	HLen(key string) (int, error)
	HIncrBy(key, field string, increment int) (int, error)
	HIncrByFloat(key, field string, increment float64) (float64, error)
	HMset(key string, fieldsAndValues map[string]string) error
	HMget(key string, fields ...string) ([]string, error)
	HKeys(key string) ([]string, error)
	HVals(key string) ([]string, error)
	HGetAll(key string) (map[string]string, error)
	HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error)

	LPush(key string, values ...interface{}) (int, error)
	RPush(key string, values ...interface{}) (int, error)
	LPop(key string) (string, error)
	RPop(key string) (string, error)
	LRem(key string, count int, val string) (int, error)
	LLen(key string) (int, error)
	LIndex(key string, index int64) (string, error)
	LRange(key string, start, stop int64) ([]string, error)
	LTrim(key string, start, stop int64) error
	BLPop(node Client, key string) (string, error)

	SAdd(key string, values ...interface{}) (int, error)
	SIsMember(key string, val interface{}) (bool, error)
	SPop(key string) (string, error)
	SRandMember(key string, count int) ([]string, error)
	SRem(key string, values ...interface{}) (int, error)
	SCard(key string) (int64, error)
	SMembers(key string) ([]string, error)
	SScan(key string, cursor uint64, match string, count int64) (keys []string, cur uint64, err error)
	SInter(keys ...string) ([]string, error)
	SInterStore(destination string, keys ...string) (int, error)
	SUnion(keys ...string) ([]string, error)
	SUnionStore(destination string, keys ...string) (int, error)
	SDiff(keys ...string) ([]string, error)
	SDiffStore(destination string, keys ...string) (int, error)

	ZAdd(key string, score int64, val string) (bool, error)
	ZScore(key, val string) (int64, error)
	ZIncrBy(key string, increment int64, field string) (int64, error)
	ZCard(key string) (int, error)
	ZCount(key string, start, stop int64) (int, error)
	ZRange(key string, start, stop int64) ([]string, error)
	ZRevRangeWithScores(key string, start, stop int64) ([]Pair, error)
	ZRangeByScoreWithScores(key string, start, stop int64) ([]Pair, error)
	ZRevRangeByScoreWithScores(key string, start, stop int64) ([]Pair, error)
	ZRank(key, field string) (int64, error)
	ZRevRank(key, field string) (int64, error)
	ZRem(key string, values ...interface{}) (int, error)
	ZRemRangeByRank(key string, start, stop int64) (int, error)
	ZRemRangeByScore(key string, start, stop int64) (int, error)
	ZUnionStore(dest string, store *ZStore) (int64, error)

	GeoAdd(key string, geoLocation ...*GeoLocation) (int64, error)
	GeoPos(key string, members ...string) ([]*GeoPos, error)
	GeoDist(key, member1, member2, unit string) (float64, error)
	GeoRadius(key string, longitude, latitude float64, query *GeoRadiusQuery) ([]GeoLocation, error)
	GeoRadiusByMember(key, member string, query *GeoRadiusQuery) ([]GeoLocation, error)
	GeoHash(key string, members ...string) ([]string, error)

	SetBit(key string, offset int64, val int) (int, error)
	GetBit(key string, offset int64) (int, error)
	BitCount(key string, start, end int64) (int64, error)
	BitPos(key string, bit, start, end int64) (int64, error)
	BitOpAnd(destKey string, keys ...string) (int64, error)

	Exists(keys ...string) (bool, error)
	Del(keys ...string) error
	Scan(cursor uint64, match string, count int64) (keys []string, cur uint64, err error)

	Expire(key string, expiration time.Duration) error
	TTL(key string) (int, error)
	Persist(key string) (bool, error)
}

func toPairs(vals []redis.Z) []Pair {
	pairs := make([]Pair, len(vals))
	for i, val := range vals {
		switch member := val.Member.(type) {
		case string:
			pairs[i] = Pair{
				Key:   member,
				Score: int64(val.Score),
			}
		default:
			pairs[i] = Pair{
				Key:   mapping.Repr(val.Member),
				Score: int64(val.Score),
			}
		}
	}
	return pairs
}

func NewConn(conf Config, client clientMode) RedisConn {
	return &redisConn{
		client: client(conf),
		brk:    breaker.NewBreaker(),
	}
}

func acceptable(err error) bool {
	return err == nil || err == redis.Nil || err == context.Canceled
}

func toStrings(vals []interface{}) []string {
	ret := make([]string, len(vals))
	for i, val := range vals {
		if val == nil {
			ret[i] = ""
		} else {
			switch val := val.(type) {
			case string:
				ret[i] = val
			default:
				ret[i] = mapping.Repr(val)
			}
		}
	}
	return ret
}

type redisConn struct {
	client Client
	brk    breaker.Breaker
}

func (r redisConn) Set(key string, val interface{}, args ...int) (err error) {
	seconds := 0
	if 0 < len(args) {
		seconds = args[0]
	}
	err = r.brk.DoWithAcceptable(func() error {
		err = r.client.Set(context.Background(), key, val, time.Duration(seconds)*time.Second).Err()
		return err
	}, acceptable)
	return
}

func (r redisConn) SetNX(key, val string, args ...int) (b bool, err error) {
	seconds := 0
	if 0 < len(args) {
		seconds = args[0]
	}

	err = r.brk.DoWithAcceptable(func() error {
		b, err = r.client.SetNX(context.Background(), key, val, time.Duration(seconds)*time.Second).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) Get(key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.Get(context.Background(), key).Result()
		if err == redis.Nil {
			return nil
		}
		return err
	}, acceptable)
	return
}

func (r redisConn) GetSet(key, inVal string) (outVal string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		if outVal, err = r.client.GetSet(context.Background(), key, outVal).Result(); err == redis.Nil {
			return nil
		}

		return err
	}, acceptable)

	return
}

func (r redisConn) Incr(key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.Incr(context.Background(), key).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) Incrby(key string, increment int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.IncrBy(context.Background(), key, increment).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) IncrByFloat(key string, increment float64) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.IncrByFloat(context.Background(), key, increment).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) Decr(key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.Decr(context.Background(), key).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) DecrBy(key string, decrement int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.DecrBy(context.Background(), key, decrement).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) MGet(keys ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.MGet(context.Background(), keys...).Result()

		val = toStrings(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) HSet(key, field, val string) (err error) {
	err = r.brk.DoWithAcceptable(func() error {
		return r.client.HSet(context.Background(), key, field, val).Err()
	}, acceptable)
	return
}

func (r redisConn) HSetNX(key, field, value string) (b bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		b, err = r.client.HSetNX(context.Background(), key, field, value).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) HGet(key, field string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.HGet(context.Background(), key, field).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) HExists(key, field string) (b bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		b, err = r.client.HExists(context.Background(), key, field).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) HDel(key string, fields ...string) (b bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.HDel(context.Background(), key, fields...).Result()
		b = v >= 1
		return err
	}, acceptable)
	return
}

func (r redisConn) HLen(key string) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.HLen(context.Background(), key).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) HIncrBy(key, field string, increment int) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.HIncrBy(context.Background(), key, field, int64(increment)).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) HIncrByFloat(key, field string, increment float64) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.HIncrByFloat(context.Background(), key, field, increment).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) HMset(key string, values map[string]string) (err error) {
	err = r.brk.DoWithAcceptable(func() error {
		//vals := make(map[string]interface{}, len(values))
		//for k, v := range values {
		//	vals[k] = v
		//}
		return r.client.HMSet(context.Background(), key, values).Err()
	}, acceptable)
	return
}

func (r redisConn) HMget(key string, fields ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.HMGet(context.Background(), key, fields...).Result()

		val = toStrings(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) HKeys(key string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.HKeys(context.Background(), key).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) HVals(key string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.HVals(context.Background(), key).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) HGetAll(key string) (val map[string]string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.HGetAll(context.Background(), key).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) HScan(key string, cursor uint64, match string, count int64) (
	keys []string, cur uint64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		keys, cur, err = r.client.HScan(context.Background(), key, cursor, match, count).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) LPush(key string, values ...interface{}) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.LPush(context.Background(), key, values...).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) RPush(key string, values ...interface{}) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.RPush(context.Background(), key, values...).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) LPop(key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.LPop(context.Background(), key).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) RPop(key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.RPop(context.Background(), key).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) LRem(key string, count int, value string) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.LRem(context.Background(), key, int64(count), value).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) LLen(key string) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.LLen(context.Background(), key).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) LIndex(key string, index int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.LIndex(context.Background(), key, index).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) LRange(key string, start, stop int64) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.LRange(context.Background(), key, start, stop).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) LTrim(key string, start, stop int64) (err error) {
	err = r.brk.DoWithAcceptable(func() error {
		return r.client.LTrim(context.Background(), key, start, stop).Err()
	}, acceptable)
	return
}

func (r redisConn) BLPop(node Client, key string) (string, error) {
	vals, err := node.BLPop(context.Background(), blockingQueryTimeout, key).Result()
	if err != nil {
		return "", err
	}
	if len(vals) < 2 {
		return "", fmt.Errorf("no value on key: %s", key)
	}

	return vals[0], nil
}

func (r redisConn) SAdd(key string, values ...interface{}) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.SAdd(context.Background(), key, values...).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) SIsMember(key string, value interface{}) (b bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		b, err = r.client.SIsMember(context.Background(), key, value).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) SPop(key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.SPop(context.Background(), key).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) SRandMember(key string, count int) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.SRandMemberN(context.Background(), key, int64(count)).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) SRem(key string, values ...interface{}) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.SRem(context.Background(), key, values).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) SCard(key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.SCard(context.Background(), key).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) SMembers(key string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.SMembers(context.Background(), key).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) SScan(key string, cursor uint64, match string, count int64) (
	keys []string, cur uint64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		keys, cur, err = r.client.SScan(context.Background(), key, cursor, match, count).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) SInter(keys ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.SInter(context.Background(), keys...).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) SInterStore(destination string, keys ...string) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.SInterStore(context.Background(), destination, keys...).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) SUnion(keys ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.SUnion(context.Background(), keys...).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) SUnionStore(destination string, keys ...string) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.SUnionStore(context.Background(), destination, keys...).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) SDiff(keys ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.SDiff(context.Background(), keys...).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) SDiffStore(destination string, keys ...string) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.SDiffStore(context.Background(), destination, keys...).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) ZAdd(key string, score int64, value string) (b bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.ZAdd(context.Background(), key, &redis.Z{
			Score:  float64(score),
			Member: value,
		}).Result()
		if err != nil {
			return err
		}
		b = v == 1

		return err
	}, acceptable)
	return
}

func (r redisConn) ZScore(key, value string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.ZScore(context.Background(), key, value).Result()
		val = int64(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) ZIncrBy(key string, increment int64, field string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.ZIncrBy(context.Background(), key, float64(increment), field).Result()
		val = int64(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) ZCard(key string) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.ZCard(context.Background(), key).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) ZCount(key string, start, stop int64) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.ZCount(context.Background(), key,
			strconv.FormatInt(start, 10),
			strconv.FormatInt(stop, 10)).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) ZRange(key string, start, stop int64) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.ZRange(context.Background(), key, start, stop).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) ZRevRangeWithScores(key string, start, stop int64) (val []Pair, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.ZRevRangeWithScores(context.Background(), key, start, stop).Result()
		if err != nil {
			return err
		}

		val = toPairs(v)
		return nil
	}, acceptable)
	return
}

func (r redisConn) ZRangeByScoreWithScores(key string, start, stop int64) (val []Pair, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.ZRangeByScoreWithScores(context.Background(), key,
			&redis.ZRangeBy{
				Min: strconv.FormatInt(start, 10),
				Max: strconv.FormatInt(stop, 10),
			}).Result()
		if err != nil {
			return err
		}

		val = toPairs(v)
		return nil
	}, acceptable)
	return
}

func (r redisConn) ZRevRangeByScoreWithScores(key string, start, stop int64) (val []Pair, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.ZRevRangeByScoreWithScores(context.Background(), key,
			&redis.ZRangeBy{
				Min: strconv.FormatInt(start, 10),
				Max: strconv.FormatInt(stop, 10),
			}).Result()
		if err != nil {
			return err
		}

		val = toPairs(v)
		return nil
	}, acceptable)
	return
}

func (r redisConn) ZRank(key, field string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.ZRank(context.Background(), key, field).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) ZRevRank(key, field string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.ZRevRank(context.Background(), key, field).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) ZRem(key string, values ...interface{}) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.ZRem(context.Background(), key, values...).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) ZRemRangeByRank(key string, start, stop int64) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.ZRemRangeByRank(context.Background(), key, start, stop).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) ZRemRangeByScore(key string, start, stop int64) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.ZRemRangeByScore(context.Background(), key,
			strconv.FormatInt(start, 10),
			strconv.FormatInt(stop, 10)).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) ZUnionStore(dest string, store *ZStore) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.ZUnionStore(context.Background(), dest, store).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) GeoAdd(key string, geoLocation ...*GeoLocation) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.GeoAdd(context.Background(), key, geoLocation...).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) GeoPos(key string, members ...string) (val []*GeoPos, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.GeoPos(context.Background(), key, members...).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) GeoDist(key, member1, member2, unit string) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.GeoDist(context.Background(), key, member1, member2, unit).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) GeoRadius(key string, longitude, latitude float64, query *GeoRadiusQuery) (
	val []GeoLocation, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.GeoRadius(context.Background(), key, longitude, latitude, query).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) GeoRadiusByMember(key, member string, query *GeoRadiusQuery) (
	val []GeoLocation, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.GeoRadiusByMember(context.Background(), key, member, query).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) GeoHash(key string, members ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.GeoHash(context.Background(), key, members...).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) SetBit(key string, offset int64, value int) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.SetBit(context.Background(), key, offset, val).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) GetBit(key string, offset int64) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.GetBit(context.Background(), key, offset).Result()
		val = int(v)
		return err
	}, acceptable)
	return
}

func (r redisConn) BitCount(key string, start, end int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.BitCount(context.Background(), key, &redis.BitCount{
			Start: start,
			End:   end,
		}).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) BitPos(key string, bit, start, end int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.BitPos(context.Background(), key, bit, start, end).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) BitOpAnd(destKey string, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		val, err = r.client.BitOpAnd(context.Background(), destKey, keys...).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) Exists(keys ...string) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		v, err := r.client.Exists(context.Background(), keys...).Result()
		val = v == 1
		return err
	}, acceptable)
	return
}

func (r redisConn) Del(keys ...string) (err error) {
	err = r.brk.DoWithAcceptable(func() error {
		err = r.client.Del(context.Background(), keys...).Err()
		return err
	}, acceptable)
	return
}

func (r redisConn) Scan(cursor uint64, match string, count int64) (keys []string, cur uint64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		keys, cur, err = r.client.Scan(context.Background(), cursor, match, count).Result()
		return err
	}, acceptable)
	return
}

func (r redisConn) Expire(key string, expiration time.Duration) (err error) {
	err = r.brk.DoWithAcceptable(func() error {
		err = r.client.Expire(context.Background(), key, expiration).Err()
		return err
	}, acceptable)
	return
}

func (r redisConn) TTL(key string) (val int, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		duration, err := r.client.TTL(context.Background(), key).Result()

		val = int(duration / time.Second)
		return err
	}, acceptable)
	return
}

func (r redisConn) Persist(key string) (b bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		b, err = r.client.Persist(context.Background(), key).Result()
		return err
	}, acceptable)
	return
}
