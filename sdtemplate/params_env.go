package sdtemplate

func (params *Params) EnvOrDefault(key string, def string) string {
	value, ok := params.EnvMap[key];
	if ok {
		return value
	} else {
		return def
	}
}

func (params *Params) EnvGet(key string) string {
	value, ok := params.EnvMap[key];
	if ok {
		return value
	} else {
		return ""
	}
}

func (params *Params) EnvHas(key string, def string) bool {
	_, ok := params.EnvMap[key];
	return ok
}

func (params *Params) EnvEquals(key string, value string) bool {
	return params.EnvMap[key] == value
}
