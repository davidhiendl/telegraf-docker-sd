package globalconfig

type GlobalConfigSpec struct {
	EnvMap map[string]string
	Backends []string
	Tags map[string]string
}
