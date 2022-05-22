package components

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/memcache"
	"github.com/gofiber/storage/memory"
	"github.com/gofiber/storage/mysql"
	"github.com/gofiber/storage/redis"
	"github.com/gofiber/storage/ristretto"
	"github.com/gofiber/utils"
	"github.com/ohrimenko/sergo/config"
)

type storeStruct struct {
	store *session.Store
	err   error
	valid bool
}

var storeApp storeStruct

func (n *storeStruct) session() *session.Store {
	if !n.valid {
		expire, err := strconv.Atoi(config.Env("SESSION_EXPIRE"))
		if err != nil {
			expire = 60
		}

		if config.Env("SESSION_DRIVER") == "memcache" {
			n.store = session.New(session.Config{
				Expiration:   time.Duration(expire) * time.Minute,
				KeyLookup:    "cookie:session_id",
				KeyGenerator: utils.UUID,
				Storage: memcache.New(memcache.Config{
					Servers: config.Env("SESSION_HOST") + ":" + config.Env("SESSION_PORT"),
				}),
			})
		} else if config.Env("SESSION_DRIVER") == "redis" {
			port, err := strconv.Atoi(config.Env("SESSION_PORT"))
			if err != nil {
				port = 6379
			}

			n.store = session.New(session.Config{
				Expiration:   time.Duration(expire) * time.Minute,
				KeyLookup:    "cookie:session_id",
				KeyGenerator: utils.UUID,
				Storage: redis.New(redis.Config{
					Host:      config.Env("SESSION_HOST"),
					Port:      port,
					Username:  config.Env("SESSION_USERNAME"),
					Password:  config.Env("SESSION_PASSWORD"),
					Database:  0,
					Reset:     false,
					TLSConfig: nil,
				}),
			})
		} else if config.Env("SESSION_DRIVER") == "ristretto" {
			n.store = session.New(session.Config{
				Expiration:   time.Duration(expire) * time.Minute,
				KeyLookup:    "cookie:session_id",
				KeyGenerator: utils.UUID,
				Storage: ristretto.New(ristretto.Config{
					NumCounters: 1e7,     // number of keys to track frequency of (10M).
					MaxCost:     1 << 30, // maximum cost of cache (1GB).
					BufferItems: 64,      // number of keys per Get buffer.
				}),
			})
		} else if config.Env("SESSION_DRIVER") == "mysql" {
			port, err := strconv.Atoi(config.Env("SESSION_PORT"))
			if err != nil {
				port = 3306
			}

			n.store = session.New(session.Config{
				Expiration:   time.Duration(expire) * time.Minute,
				KeyLookup:    "cookie:session_id",
				KeyGenerator: utils.UUID,
				Storage: mysql.New(mysql.Config{
					Host:       config.Env("SESSION_HOST"),
					Port:       port,
					Username:   config.Env("SESSION_USERNAME"),
					Password:   config.Env("SESSION_PASSWORD"),
					Database:   config.Env("SESSION_DATABASE"),
					Table:      config.Env("SESSION_TABLE"),
					Reset:      false,
					GCInterval: 10 * time.Second,
				}),
			})
		} else {
			n.store = session.New(session.Config{
				Expiration:   time.Duration(expire) * time.Minute,
				KeyLookup:    "cookie:session_id",
				KeyGenerator: utils.UUID,
				Storage: memory.New(memory.Config{
					GCInterval: 10 * time.Second,
				}),
			})
		}

		n.valid = true
	}

	return n.store
}

func Session() *session.Store {
	return storeApp.session()
}
