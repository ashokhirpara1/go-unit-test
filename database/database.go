package database

import (
	"log"

	"learning/unit-testing/models"
	"learning/unit-testing/restapi/operations/connections"

	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

// Connection Type representing the connection to a Firebase Database.
type Connection struct {
	Client  *firestore.Client
	Context context.Context
}

// NewConnection - Initialize new firestore connection
func NewConnection() (Storage, error) {

	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Printf("Error initializing Firebase app: %v\n", err)
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Connection{Client: client, Context: ctx}, nil
}

// GetUserConnectionGroupByName - function
func (c *Connection) GetUserConnectionGroupByName(userID string, groupName string) (internal.UserConnectionGroupInfo, error) {
	var groupinfoObj internal.UserConnectionGroupInfo

	existingGroupRef := c.Client.Collection(internal.GetGroupCollectionPath(userID)).Where("group_name", "==", groupName).Limit(1).Documents(c.Context)
	for {
		doc, err := existingGroupRef.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return groupinfoObj, err
		}

		if err := doc.DataTo(&groupinfoObj); err != nil {
			return groupinfoObj, err
		}
	}

	if groupinfoObj.GroupID == "" {
		return groupinfoObj, status.Error(codes.NotFound, "row does not found")
	}

	return groupinfoObj, nil
}

// GetUserConnectionGroupByGroupID - function
func (c *Connection) GetUserConnectionGroupByGroupID(userID string, groupID string) (internal.UserConnectionGroupInfo, error) {
	var groupinfoObj internal.UserConnectionGroupInfo

	// Get the connections group info.
	groupDoc, err := c.Client.Doc(internal.GetGroupDocPath(userID, groupID)).Get(c.Context)
	if err != nil {
		return groupinfoObj, err
	}

	if err := groupDoc.DataTo(&groupinfoObj); err != nil {
		return groupinfoObj, err
	}

	return groupinfoObj, nil
}

// CreateUserConnectionGroup - function
func (c *Connection) CreateUserConnectionGroup(userID string, group internal.UserConnectionGroupInfo) (string, error) {
	// Set the group into the database.
	groupRef := c.Client.Collection(internal.GetGroupCollectionPath(userID)).NewDoc()
	group.GroupID = groupRef.ID

	if _, groupSetErr := groupRef.Set(c.Context, group); groupSetErr != nil {
		return "", groupSetErr
	}
	return group.GroupID, nil
}

// GetPaginatedUserConnectionGroup - function
func (c *Connection) GetPaginatedUserConnectionGroup(params connections.UsersConnectionsGroupsByUserIDGetParams) (groupsList []*models.Group, paginationMeta *models.PaginationData, err error) {

	// Create the paginated query.
	var limit int32
	if internal.IsZeroOfUnderlyingType(params.Limit) {
		limit = internal.DefaultConnectionsQueryLimit
	} else {
		limit = *params.Limit
	}

	paginatedQuery := internal.NewPaginatedQuery(c.Client.Collection(internal.GetGroupCollectionPath(params.UserID)).Query, params.Offset, &limit, params.OrderBy, params.Order)

	// Set time interaction filter.
	if err := paginatedQuery.AddTimeSetFilterToQuery("latest_interaction_time", params.LatestInteractionTimeAfter, params.LatestInteractionTimeAfter); err != nil {
		return groupsList, paginationMeta, err
	}

	// Set filter for group name.
	if !internal.IsZeroOfUnderlyingType(params.GroupName) {
		*paginatedQuery.Query = paginatedQuery.Query.Where("group_name", "==", params.GroupName)
	}

	connectionGroupsDocs, err := paginatedQuery.Query.Documents(c.Context).GetAll()
	if err != nil {
		return groupsList, paginationMeta, err
	}

	dbConnection := internal.DataBaseConnection{Client: c.Client, Context: c.Context}
	// Get the pagination metadata.
	paginationMeta, err = paginatedQuery.GetPaginatedQueryMetadata(&dbConnection)
	if err != nil {
		return groupsList, paginationMeta, err
	}

	// If the query returned no results, its still a good query with no results.
	if len(connectionGroupsDocs) < 1 {
		return groupsList, paginationMeta, nil
	}

	for _, connectionGroupsDoc := range connectionGroupsDocs {
		var groupInfo internal.UserConnectionGroupInfo

		if err := connectionGroupsDoc.DataTo(&groupInfo); err != nil {
			return groupsList, paginationMeta, err
		}

		groupData := groupInfo.TransformToResponseGroup()

		groupsList = append(groupsList, groupData)
	}

	return groupsList, paginationMeta, nil
}

// UpdateUserConnectionGroup - function
func (c *Connection) UpdateUserConnectionGroup(params connections.UsersConnectionsGroupsByUserIDAndGroupIDPatchParams) error {

	groupObj, err := c.GetUserConnectionGroupByGroupID(params.UserID, params.GroupID)
	if err != nil {
		return err
	}

	var updates []firestore.Update

	if !internal.IsZeroOfUnderlyingType(params.Body.GroupName) {
		updates = append(updates, firestore.Update{
			Path:  "group_name",
			Value: params.Body.GroupName,
		})
	}

	changeConnectionUserIds := false
	connectionUserIds := groupObj.ConnectionUserIds

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
		updates = append(updates, firestore.Update{
			Path:  "connection_user_ids",
			Value: connectionUserIds,
		})
	}

	if !internal.IsZeroOfUnderlyingType(params.Body.GroupPic) {
		updates = append(updates, firestore.Update{
			Path:  "group_pic",
			Value: params.Body.GroupPic,
		})
	}

	if !internal.IsZeroOfUnderlyingType(updates) && len(updates) > 0 {
		if _, err := c.Client.Doc(internal.GetGroupDocPath(params.UserID, params.GroupID)).Update(c.Context, updates); err != nil {
			return err
		}
	}

	return nil
}

// DeleteUserConnectionGroup - function
func (c *Connection) DeleteUserConnectionGroup(params connections.UsersConnectionsGroupsByUserIDAndGroupIDDeleteParams) error {

	// Check Group exists before delete.
	_, err := c.GetUserConnectionGroupByGroupID(params.UserID, params.GroupID)
	if err != nil {
		return err
	}

	// Remove the connection group from the users groups database.
	if _, err := c.Client.Doc(internal.GetGroupDocPath(params.UserID, params.GroupID)).Delete(c.Context); err != nil {
		return err
	}

	return nil
}
