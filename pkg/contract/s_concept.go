package contract

// 内置概念
var (
	png  = newConcept(Png, File)
	jpeg = newConcept(Jpeg, File, Jpg, Jpe)
)

// Concept 概念
type Concept struct {
	// name 名称
	name ConceptName
	// kind 类型
	kind Kind
	// aliases 别名
	aliases []ConceptName
}

// newConcept 创建概念
func newConcept(name ConceptName, kind Kind, aliases ...ConceptName) Concept {
	concept := Concept{
		name:    name,
		kind:    kind,
		aliases: aliases,
	}
	_conceptCache.put(concept) // 加入缓存
	return concept
}

// NormalizeConcept 根据name 或 alias 获取概念
func NormalizeConcept(name ConceptName) (concept Concept, exist bool) {
	return _conceptCache.getFromByNameOrAliasesMap(name)
}

func (s Concept) Name() ConceptName {
	return s.name
}

func (s Concept) Kind() Kind {
	return s.kind
}

func (s Concept) Aliases() []ConceptName {
	copyAliases := make([]ConceptName, len(s.aliases))
	copy(copyAliases, s.aliases)
	return copyAliases
}

func PNG() Concept {
	return png
}

func JPEG() Concept {
	return jpeg
}
