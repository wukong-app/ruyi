package core

// 内置概念
var (
	png  = newConcept("png", File)
	jpeg = newConcept("jpeg", File, "jpg", "jpe")
)

// Concept 概念
type Concept struct {
	// name 名称
	name string
	// kind 类型
	kind Kind
	// aliases 别名
	aliases []string
}

// newConcept 创建概念
func newConcept(name string, kind Kind, aliases ...string) Concept {
	concept := Concept{
		name:    name,
		kind:    kind,
		aliases: aliases,
	}
	_conceptCache.put(concept) // 加入缓存
	return concept
}

// NormalizeConcept 根据name 或 alias 获取概念
func NormalizeConcept(name string) (concept Concept, exist bool) {
	return _conceptCache.getFromByNameOrAliasesMap(name)
}

func (s Concept) Name() string {
	return s.name
}

func (s Concept) Kind() Kind {
	return s.kind
}

func (s Concept) Aliases() []string {
	copyAliases := make([]string, len(s.aliases))
	copy(copyAliases, s.aliases)
	return copyAliases
}

func PNG() Concept {
	return png
}

func JPEG() Concept {
	return jpeg
}
