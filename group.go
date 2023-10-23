package main

type Group interface {
	ID() string
	Rule() GroupRule
}

type group struct {
	id   string
	rule GroupRule
}

func (g *group) ID() string {
	return g.id
}

func (g *group) Rule() GroupRule {
	return g.rule
}

type GroupRule interface {
	GetMembers() []string
}

type groupRule struct {
	members []string
}

func NewGroupRule() GroupRule {
	return &groupRule{
		members: make([]string, 0),
	}
}

func (gr *groupRule) GetMembers() []string {
	return gr.members
}