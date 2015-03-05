package slackerlib

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/url"
)

// Top-level exported Store interface for storage backends to implement
type Brain interface {
	Open() error
	Close() error
	Get(string) ([]byte, error)
	Set(key string, data []byte) error
	Delete(string) error
}

// NewStore returns an initialized store
func (b *Sbot) NewBrain() (*Brain, error) {
	var brain Brain
	var err error
	if b.Config.RedisURL != ``{
		Logger.Debug(`Brain:: setting up a Redis Brain to: `, b.Config.RedisURL)
		if brain, err =	newRedisBrain(b); err != nil{
			return &brain, err
		}
	}else{
		Logger.Debug(`Brain:: setting up an in-memory Brain`)
		if brain, err =	newRAMBrain(b); err != nil{
			return &brain, err
		}
	}
	return &brain, nil
}

//rambrain storage implementation
type ramBrain struct {
	data map[string][]byte
}

// New returns a new initialized ramBrain that implements Brain
func newRAMBrain(b *Sbot) (Brain, error) {
	rb := &ramBrain{
		data: map[string][]byte{},
	}
	return rb, nil
}

func (rb *ramBrain) Open() error {
	return nil
}

func (rb *ramBrain) Close() error {
	return nil
}

func (rb *ramBrain) Get(key string) ([]byte, error) {
	if val, ok := rb.data[key]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("key %s was not found", key)
}

func (rb *ramBrain) Set(key string, data []byte) error {
	rb.data[key] = data
	return nil
}

func (rb *ramBrain) Delete(key string) error {
	if _, ok := rb.data[key]; !ok {
		return fmt.Errorf("key %s was not found", key)
	}
	delete(rb.data, key)
	return nil
}

//redisbrain backend storage implementation
type redisBrain struct {
	url			string
	pw			   string
	nameSpace	string
	client redis.Conn
}

// New returns an new initialized store
func newRedisBrain(b *Sbot) (Brain, error) {
	s := &redisBrain{
		url: b.Config.RedisURL,
		nameSpace: b.Config.Name,
	}
	if b.Config.RedisPW != ``{
		s.pw = b.Config.RedisPW
	}
	return s, nil
}

func (rb *redisBrain) Open() error {
	uri, err := url.Parse(rb.url)
	if err != nil {
		Logger.Error(err)
	}

	conn, err := redis.Dial("tcp", uri.Host)
	if err != nil {
		Logger.Error(err)
		return err
	}
	
	rb.client = conn

	if rb.pw != ``{
		if _, err := rb.client.Do("AUTH", rb.pw); err != nil{
			return err
		}
	}

	return nil
}

func (rb *redisBrain) Close() error {
	if err := rb.client.Close(); err != nil {
		Logger.Error(err)
		return err
	}
	return nil
}

func (rb *redisBrain) Get(key string) ([]byte, error) {
	args := rb.namespace(key)
	data, err := rb.client.Do("GET", args)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return []byte{}, fmt.Errorf("%s not found", key)
	}
	return data.([]byte), nil
}

func (rb *redisBrain) Set(key string, data []byte) error {
	if _, err := rb.client.Do("SET", rb.namespace(key), data); err != nil {
		return err
	}
	return nil
}

func (rb *redisBrain) Delete(key string) error {
	res, err := rb.client.Do("DEL", rb.namespace(key))
	if err != nil {
		return err
	}
	if res.(int64) < 1 {
		return fmt.Errorf("%s not found", key)
	}
	return nil
}

func (rb *redisBrain) namespace(key string) string {
	return fmt.Sprintf("%s:%s", rb.nameSpace, key)
}

