package database

import (
	"learning/unit-testing/connections"
	"learning/unit-testing/models"
)

// Storage - Handle database functions
type Storage interface {
	GetUserConnectionGroupByName(userID, groupName string) (internal.UserConnectionGroupInfo, error)
	GetUserConnectionGroupByGroupID(userID, groupID string) (internal.UserConnectionGroupInfo, error)
	CreateUserConnectionGroup(userID string, group internal.UserConnectionGroupInfo) (string, error)
	GetPaginatedUserConnectionGroup(params connections.UsersConnectionsGroupsByUserIDGetParams) (groupsList []*models.Group, paginationMeta *models.PaginationData, err error)
	UpdateUserConnectionGroup(params connections.UsersConnectionsGroupsByUserIDAndGroupIDPatchParams) error
	DeleteUserConnectionGroup(params connections.UsersConnectionsGroupsByUserIDAndGroupIDDeleteParams) error
}
