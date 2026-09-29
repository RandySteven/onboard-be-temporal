package enums

type RouterPrefix string

const (
	BasePrefix RouterPrefix = "/"
	AuthPrefix RouterPrefix = "/auth"
)

func (prefix RouterPrefix) ToString() string {
	return string(prefix)
}
