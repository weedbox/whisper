package whisper

type groupMemory struct {
	id      string
	members map[string]*Member
}

type GroupResolverMemory struct {
	w      Whisper
	groups map[string]*groupMemory
}

func NewGroupResolverMemory() *GroupResolverMemory {
	return &GroupResolverMemory{
		groups: make(map[string]*groupMemory),
	}
}

func (gs *GroupResolverMemory) Init(w Whisper) error {
	gs.w = w
	return nil
}

func (gs *GroupResolverMemory) getGroupIDs() []string {

	var groups []string
	for gid, _ := range gs.groups {
		groups = append(groups, gid)
	}

	return groups
}

func (gs *GroupResolverMemory) addGroup(groupID string, members []string) error {

	g := &groupMemory{
		id:      groupID,
		members: make(map[string]*Member),
	}

	for _, mid := range members {
		g.members[mid] = &Member{
			ID: mid,
		}
	}

	gs.groups[groupID] = g

	return nil
}

func (gs *GroupResolverMemory) GetMemberIDs(groupID string) ([]string, error) {

	var members []string

	g, ok := gs.groups[groupID]
	if !ok {
		return members, ErrGroupNotFound
	}

	for mid, _ := range g.members {
		members = append(members, mid)
	}

	return members, nil
}
