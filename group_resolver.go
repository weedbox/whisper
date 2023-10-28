package whisper

type GroupResolver interface {
	Init(w Whisper) error
	GetGroupRule(groupID string) GroupRule
}

type GroupResolverMemory struct {
	w      Whisper
	groups map[string]Group
}

func NewGroupResolverMemory() *GroupResolverMemory {
	return &GroupResolverMemory{
		groups: make(map[string]Group),
	}
}

func (gs *GroupResolverMemory) Init(w Whisper) error {
	gs.w = w
	return nil
}

func (gs *GroupResolverMemory) GetGroups() []Group {

	var groups []Group
	for _, g := range gs.groups {
		groups = append(groups, g)
	}

	return groups
}

func (gs *GroupResolverMemory) GetGroupRule(groupID string) GroupRule {

	if g, ok := gs.groups[groupID]; ok {
		return g.Rule()
	}

	return nil
}

func (gs *GroupResolverMemory) AddGroup(groupID string, members []string) error {

	gs.groups[groupID] = &group{
		id: groupID,
		rule: &groupRule{
			members: members,
		},
	}

	return nil
}
