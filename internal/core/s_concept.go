package core

import "github.com/wukong-app/ruyi/pkg/contract"

// 内置概念
var (
	png  = newConcept(contract.PNG, contract.File)
	jpeg = newConcept(contract.JPEG, contract.File, contract.JPG, contract.JPE)
)

// Concept 概念
type Concept struct {
	// name 名称
	name contract.ConceptName
	// kind 类型
	kind contract.Kind
	// aliases 别名
	aliases []contract.ConceptName
}

// newConcept 创建概念
func newConcept(name contract.ConceptName, kind contract.Kind, aliases ...contract.ConceptName) Concept {
	concept := Concept{
		name:    name,
		kind:    kind,
		aliases: aliases,
	}
	_conceptCache.put(concept) // 加入缓存
	return concept
}

// NormalizeConcept 根据name 或 alias 获取概念
func NormalizeConcept(name contract.ConceptName) (concept Concept, exist bool) {
	return _conceptCache.getFromByNameOrAliasesMap(name)
}

func (s Concept) Name() contract.ConceptName {
	return s.name
}

func (s Concept) Kind() contract.Kind {
	return s.kind
}

func (s Concept) Aliases() []contract.ConceptName {
	copyAliases := make([]contract.ConceptName, len(s.aliases))
	copy(copyAliases, s.aliases)
	return copyAliases
}

func PNG() Concept {
	return png
}

func JPEG() Concept {
	return jpeg
}
