package database

import (
	"fmt"
	"sync"

	"learning/unit-testing/models"
	"learning/unit-testing/restapi/operations/connections"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockConnection - handler
type MockConnection struct {
	userConnectionGroups map[string][]internal.UserConnectionGroupInfo
}

// NewMockConnection - Initialize Memory Storage
func NewMockConnection() Storage {
	return &MockConnection{userConnectionGroups: make(map[string][]internal.UserConnectionGroupInfo)}
}

// GenerateUUID -
func GenerateUUID() string {
	reqID := uuid.New()
	return reqID.String()
}

// GetUserConnectionGroupByName - function
func (m *MockConnection) GetUserConnectionGroupByName(userID string, groupName string) (internal.UserConnectionGroupInfo, error) {

	errNotFound := status.Error(codes.NotFound, "row does not found")

	var group internal.UserConnectionGroupInfo
	var mx sync.RWMutex

	mx.RLock()

	groups, ok := m.userConnectionGroups[userID]
	if !ok {
		return group, errNotFound
	}

	if len(groups) == 0 {
		return group, errNotFound
	}

	for _, g := range groups {
		if g.GroupName == groupName {
			group = g
			break
		}
	}

	if group.GroupID == "" {
		return group, errNotFound
	}

	mx.RUnlock()

	return group, nil
}

// GetUserConnectionGroupByGroupID - function
func (m *MockConnection) GetUserConnectionGroupByGroupID(userID string, groupID string) (internal.UserConnectionGroupInfo, error) {

	errNotFound := status.Error(codes.NotFound, "row does not found")
	errCollectionNotExists := status.Error(codes.Internal, "something went wrong")

	var group internal.UserConnectionGroupInfo
	var mx sync.RWMutex

	mx.RLock()

	groups, ok := m.userConnectionGroups[userID]
	if !ok {
		return group, errCollectionNotExists
	}

	if len(groups) == 0 {
		return group, errNotFound
	}

	for _, g := range groups {
		if g.GroupID == groupID {
			group = g
			break
		}
	}

	if group.GroupID == "" {
		return group, errNotFound
	}

	mx.RUnlock()

	return group, nil
}

// CreateUserConnectionGroup - function
func (m *MockConnection) CreateUserConnectionGroup(userID string, group internal.UserConnectionGroupInfo) (string, error) {

	var mx sync.RWMutex
	mx.Lock()

	// group.GroupID = GenerateUUID()

	if groups, ok := m.userConnectionGroups[userID]; ok {
		group.GroupID = fmt.Sprintf("group_id_%d", len(groups)+1)
		groups = append(groups, group)
		m.userConnectionGroups[userID] = groups
	} else {
		group.GroupID = fmt.Sprintf("group_id_%d", 1)
		connectionGroups := []internal.UserConnectionGroupInfo{}
		connectionGroups = append(connectionGroups, group)
		m.userConnectionGroups[userID] = connectionGroups
	}

	mx.Unlock()

	return group.GroupID, nil
}

// GetPaginatedUserConnectionGroup - function
func (m *MockConnection) GetPaginatedUserConnectionGroup(params connections.UsersConnectionsGroupsByUserIDGetParams) (groupsList []*models.Group, paginationMeta *models.PaginationData, err error) {

	errNotFound := status.Error(codes.NotFound, "row does not found")
	// errCollectionNotExists := status.Error(codes.Internal, "something went wrong")

	var mx sync.RWMutex
	mx.RLock()

	groups, ok := m.userConnectionGroups[params.UserID]
	if !ok {
		return groupsList, paginationMeta, errNotFound
	}

	// Create the paginated query.
	var limit int32
	if internal.IsZeroOfUnderlyingType(params.Limit) {
		limit = internal.DefaultConnectionsQueryLimit
	} else {
		limit = *params.Limit
	}

	if !internal.IsZeroOfUnderlyingType(params.GroupName) {
		var groupsF []internal.UserConnectionGroupInfo
		for _, g := range groups {
			if g.GroupName == *params.GroupName {
				groupsF = append(groupsF, g)
				break
			}
		}
		groups = groupsF
	}

	paginatedQuery := PaginatedQuery{CollectionName: "users_connections_groups", UserConnectionGroups: groups}

	// Set time interaction filter.
	if err := paginatedQuery.AddTimeSetFilterToQuery("latest_interaction_time", params.LatestInteractionTimeAfter, params.LatestInteractionTimeBefore); err != nil {
		return groupsList, paginationMeta, err
	}

	paginatedQuery.SetPaginatedQuery(params.Offset, &limit, params.OrderBy, params.Order)
	paginatedQuery.SortPaginatedQuery()
	paginatedQuery.LimitPaginatedQuery()

	paginationMeta = paginatedQuery.GetPaginatedQueryMetadata()

	if len(paginatedQuery.UserConnectionGroups) < 1 {
		return groupsList, paginationMeta, nil
	}

	for _, groupInfo := range paginatedQuery.UserConnectionGroups {
		groupData := groupInfo.TransformToResponseGroup()
		groupsList = append(groupsList, groupData)
	}

	mx.RUnlock()

	return groupsList, paginationMeta, nil
}

// UpdateUserConnectionGroup - function
func (m *MockConnection) UpdateUserConnectionGroup(params connections.UsersConnectionsGroupsByUserIDAndGroupIDPatchParams) error {

	var mx sync.RWMutex
	mx.RLock()
	group, err := m.GetUserConnectionGroupByGroupID(params.UserID, params.GroupID)
	if err != nil {
		return err
	}
	mx.RUnlock()

	mx.Lock()
	if !internal.IsZeroOfUnderlyingType(params.Body.GroupName) {
		group.GroupName = params.Body.GroupName
	}

	changeConnectionUserIds := false
	connectionUserIds := group.ConnectionUserIds

	if !internal.IsZeroOfUnderlyingType(params.Body.ConnectionUserIDToAdd) {
		cgUIDAdd := internal.GroupConnectionUserID{UserID: params.Body.ConnectionUserIDToAdd}
		connectionUserIds = append(connectionUserIds, cgUIDAdd)
		changeConnectionUserIds = true
	}

	if !internal.IsZeroOfUnderlyingType(params.Body.ConnectionUserIDToRemove) {

		for i, CU := range connectionUserIds {
			if CU.UserID == params.Body.ConnectionUserIDToRemove {
				connectionUserIds = append(connectionUserIds[:i], connectionUserIds[i+1:]...)
			}
		}

		changeConnectionUserIds = true
	}

	if changeConnectionUserIds {
		group.ConnectionUserIds = connectionUserIds
	}

	if !internal.IsZeroOfUnderlyingType(params.Body.GroupPic) {
		group.GroupPic = params.Body.GroupPic
	}

	for index, g := range m.userConnectionGroups[params.UserID] {
		if g.GroupID == params.GroupID {
			m.userConnectionGroups[params.UserID][index] = group
			break
		}
	}
	mx.Unlock()

	return nil
}

// DeleteUserConnectionGroup - function
func (m *MockConnection) DeleteUserConnectionGroup(params connections.UsersConnectionsGroupsByUserIDAndGroupIDDeleteParams) error {

	var mx sync.RWMutex
	mx.RLock()
	group, err := m.GetUserConnectionGroupByGroupID(params.UserID, params.GroupID)
	if err != nil {
		return err
	}
	mx.RUnlock()

	mx.Lock()
	for index, g := range m.userConnectionGroups[params.UserID] {
		if g.GroupID == group.GroupID {
			m.userConnectionGroups[params.UserID] = append(m.userConnectionGroups[params.UserID][:index], m.userConnectionGroups[params.UserID][index+1:]...)
			break
		}
	}
	mx.Unlock()

	return nil
}
