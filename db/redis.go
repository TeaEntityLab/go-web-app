package db

import (
	"strconv"
	"strings"
	"time"

	hashring "github.com/serialx/hashring"
	log "github.com/sirupsen/logrus"
	redis "gopkg.in/redis.v5"
)

func ParseRedisConnectionString(connStr string) ([]string, int) {
	slashIndex := strings.Index(connStr, "/")
	dbName := 0

	// strip off db name from the beginning
	if slashIndex != -1 {
		dbName, err := strconv.Atoi(connStr[:slashIndex])
		if err != nil {
			return []string{""}, 0
		}
		if slashIndex == len(connStr)-1 {
			return []string{""}, dbName
		}
		connStr = connStr[slashIndex+1:]
	}

	// split the hosts, and return them and the db name
	return strings.Split(connStr, ","), dbName
}

func NewRedisClient(endpoint string, workerPoolSize int, readFromSlave bool) (redis.Cmdable, error) {
	hosts, dbName := ParseRedisConnectionString(endpoint)
	var client redis.Cmdable
	if len(hosts) > 1 {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:          hosts,
			ReadOnly:       readFromSlave,
			RouteByLatency: readFromSlave,
			MaxRedirects:   5,

			// timeout setting
			DialTimeout:  5 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,

			// pool setting
			PoolSize:           workerPoolSize,
			PoolTimeout:        5 * time.Minute,
			IdleTimeout:        10 * time.Minute,
			IdleCheckFrequency: 10 * time.Minute,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr: hosts[0],
			DB:   dbName,
		})
	}
	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}
	return client, nil
}

func newRedisHashringClient(server string, workerPoolSize int) (redis.Cmdable, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     server,
		PoolSize: workerPoolSize,
	})
	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}
	return client, nil
}

// not a good way to do so, but it works
func CheckRedisNotExistError(err error) bool {
	if err != nil && err == redis.Nil {
		return true
	}
	return false
}

type RedisHashring struct {
	*hashring.HashRing
	redisServers map[string]redis.Cmdable
	keyPrefixLen int
}

func NewRedisHashring(endpoints string, keyPrefixLen int) (*RedisHashring, error) {
	servers, _ := ParseRedisConnectionString(endpoints)
	ring := hashring.New(servers)
	redisClients := make(map[string]redis.Cmdable)
	for _, server := range servers {
		redisClient, err := newRedisHashringClient(server, 100)
		if err != nil {
			return nil, err
		}
		redisClients[server] = redisClient
	}
	return &RedisHashring{
		HashRing:     ring,
		redisServers: redisClients,
		keyPrefixLen: keyPrefixLen,
	}, nil
}

func (rh *RedisHashring) GetNode(key string) redis.Cmdable {
	if rh.keyPrefixLen > 0 && rh.keyPrefixLen <= len(key) {
		key = key[:rh.keyPrefixLen]
	}
	nodeIP, _ := rh.HashRing.GetNode(key)
	log.WithFields(log.Fields{
		"NodeIP":       nodeIP,
		"KeyPrefixLen": rh.keyPrefixLen,
		"Key":          key,
	}).Debugf("Select redis client")
	return rh.redisServers[nodeIP]
}
