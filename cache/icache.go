package cache

type ICache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
}
