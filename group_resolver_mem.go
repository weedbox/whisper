package whisper

type groupMemory struct {
	id      string
	members map[string]*Member
	muted   map[string]bool
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
		muted:   make(map[string]bool),
	}

	for _, mid := range members {
		g.members[mid] = &Member{
			ID: mid,
		}
	}

	gs.groups[groupID] = g

	return nil
}

func (gs *GroupResolverMemory) addMutedMembers(groupID string, members []string) error {

	g, ok := gs.groups[groupID]
	if !ok {
		return ErrGroupNotFound
	}

	for _, mid := range members {
		g.muted[mid] = true
	}

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

func (gs *GroupResolverMemory) IsMutedMember(groupID string, userID string) (bool, error) {

	g, ok := gs.groups[groupID]
	if !ok {
		return false, ErrGroupNotFound
	}

	if _, ok := g.muted[userID]; ok {
		return true, nil
	}

	return false, nil
}
