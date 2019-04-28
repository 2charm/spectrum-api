package sessions

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

//RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	//Redis client used to talk to redis server.
	Client *redis.Client
	//Used for key expiry time on redis.
	SessionDuration time.Duration
}

//NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	//initialize and return a new RedisStore struct
	return &RedisStore{
		Client:          client,
		SessionDuration: sessionDuration,
	}
}

//Store implementation

//Save saves the provided `sessionState` and associated SessionID to the store.
//The `sessionState` parameter is typically a pointer to a struct containing
//all the data you want to associated with the given SessionID.
func (rs *RedisStore) Save(sid SessionID, sessionState interface{}) error {
	j, err := json.Marshal(sessionState)
	if err != nil {
		return err
	}
	err = rs.Client.Set(sid.getRedisKey(), j, rs.SessionDuration).Err()
	if err != nil {
		return err
	}
	return nil
}

//Get populates `sessionState` with the data previously saved
//for the given SessionID
func (rs *RedisStore) Get(sid SessionID, sessionState interface{}) error {
	//TODO: get the previously-saved session state data from redis,
	//unmarshal it back into the `sessionState` parameter
	//and reset the expiry time, so that it doesn't get deleted until
	//the SessionDuration has elapsed.
	pipe := rs.Client.Pipeline()
	state := pipe.Get(sid.getRedisKey())
	boolCmd := pipe.Expire(sid.getRedisKey(), rs.SessionDuration)

	pipe.Exec()

	if state.Err() != nil {
		return ErrStateNotFound
	}
	stateBytes, errState := state.Bytes()
	if errState == redis.Nil {
		return ErrStateNotFound
	}
	if boolCmd.Err() != nil {
		return ErrStateNotFound
	}
	json.Unmarshal(stateBytes, sessionState)
	return nil
}

//Delete deletes all state data associated with the SessionID from the store.
func (rs *RedisStore) Delete(sid SessionID) error {
	rs.Client.Del(sid.getRedisKey())
	return nil
}

//getRedisKey() returns the redis key to use for the SessionID
func (sid SessionID) getRedisKey() string {
	//convert the SessionID to a string and add the prefix "sid:" to keep
	//SessionID keys separate from other keys that might end up in this
	//redis instance
	return "sid:" + sid.String()
}
