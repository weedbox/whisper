package whisper

type GroupManager interface {
	Init(w Whisper) error
	GetGroups() []Group
	GetGroupRule(groupID string) GroupRule
	AddGroup(groupID string, members []string) error
}

type groupManager struct {
	w      Whisper
	groups map[string]Group
}

func NewGroupManagerMemory() GroupManager {
	return &groupManager{
		groups: make(map[string]Group),
	}
}

func (gm *groupManager) Init(w Whisper) error {
	gm.w = w
	return nil
}

func (gm *groupManager) GetGroups() []Group {

	var groups []Group
	for _, g := range gm.groups {
		groups = append(groups, g)
	}

	return groups
}

func (gm *groupManager) GetGroupRule(groupID string) GroupRule {

	if g, ok := gm.groups[groupID]; ok {
		return g.Rule()
	}

	return nil
}

func (gm *groupManager) AddGroup(groupID string, members []string) error {

	gm.groups[groupID] = &group{
		id: groupID,
		rule: &groupRule{
			members: members,
		},
	}

	return nil
}
