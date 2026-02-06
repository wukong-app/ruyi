package core

import "github.com/wukong-app/ruyi/pkg/contract"

// 概念缓存
var _conceptCache = newConceptCache()

// conceptCache concept 索引
type conceptCache struct {
	// nameOrAliasesMap name or aliases -> concept
	nameOrAliasesMap map[contract.ConceptName]Concept

	// kindMap kind -> concepts
	kindMap map[contract.Kind][]Concept
}

func newConceptCache() *conceptCache {
	return &conceptCache{
		nameOrAliasesMap: make(map[contract.ConceptName]Concept),
		kindMap:          make(map[contract.Kind][]Concept),
	}
}

// put 添加概念， 若有重命直接覆盖
func (s *conceptCache) put(concept Concept) {
	s.putInNameOrAliasesMap(concept)
	s.putInKindMap(concept)
}

// putInByNameOrAliasesMap 添加概念到 nameOrAliasesMap
func (s *conceptCache) putInNameOrAliasesMap(concept Concept) {
	if s.nameOrAliasesMap == nil {
		s.nameOrAliasesMap = make(map[contract.ConceptName]Concept)
	}

	m := s.nameOrAliasesMap
	m[concept.name] = concept
	for _, alias := range concept.aliases {
		m[alias] = concept
	}
}

// getFromByNameOrAliasesMap 从 nameOrAliasesMap 获取概念
// @param name 概念 name or alias
func (s *conceptCache) getFromByNameOrAliasesMap(name contract.ConceptName) (concept Concept, exist bool) {
	concept, exist = s.nameOrAliasesMap[name]
	return
}

// putInKindMap 添加概念到 kindMap
func (s *conceptCache) putInKindMap(concept Concept) {
	if s.kindMap == nil {
		s.kindMap = make(map[contract.Kind][]Concept)
	}

	kind := concept.kind
	s.kindMap[kind] = append(s.kindMap[kind], concept)
}

// getFromByKindMap 从 kindMap 获取概念
func (s *conceptCache) getFromByKindMap(kind contract.Kind) (concepts []Concept) {
	return s.kindMap[kind]
}
