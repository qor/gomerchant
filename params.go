package gomerchant

type Params map[string]interface{}

func (params Params) Get(key string) (interface{}, bool) {
	if params == nil {
		return nil, false
	}

	value, ok := params[key]
	return value, ok
}

func (params Params) Set(key string, value interface{}) {
	if params == nil {
		params = map[string]interface{}{}
	}

	params[key] = value
}
