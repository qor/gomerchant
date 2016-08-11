package gomerchant

type Extra struct {
	params map[string]interface{}
}

func (extra *Extra) Get(key string) (interface{}, bool) {
	if extra.params == nil {
		return nil, false
	}

	value, ok := extra.params[key]
	return value, ok
}

func (extra *Extra) Set(key string, value interface{}) {
	if extra.params == nil {
		extra.params = map[string]interface{}{}
	}

	extra.params[key] = value
}

func (extra *Extra) Params() map[string]interface{} {
	return extra.params
}
