package controllers

import (
	"reflect"
	"testing"

	"learning/unit-testing/internal"
	"learning/unit-testing/models"
	"learning/unit-testing/restapi/operations/connections"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TestCaseCreateGroup struct {
	name                 string
	inputParams          connections.UsersConnectionsGroupsByUserIDPostParams
	inputPrincipal       *models.Principal
	expectedResponseType string
	expectedErrMsg       string
	expectedErr          error
}

var connectionUserIds []*models.UsersConnectionsGroupsPostRequestConnectionUserIdsItems0
var ctlr Ctlr

// init -
func init() {
	connectionUserId := &models.UsersConnectionsGroupsPostRequestConnectionUserIdsItems0{UserID: "1ca26428-98eb-4aa3-8943-5f459873ef85"}
	connectionUserIds = append(connectionUserIds, connectionUserId)

	ctlr = GetControllerMockDB()
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("%s != %s", a, b)
	}
}

// In Case if you need to add group manually into memory db then uncomment below code
func CreateConnectionGroup() {
	groupName := "Test Group Name"
	userID := "dc9dbe3e-60d5-4a07-8c9c-42027b555b01"

	ids := []internal.GroupConnectionUserID{}

	for _, CU := range connectionUserIds {
		id := internal.GroupConnectionUserID{UserID: CU.UserID}
		ids = append(ids, id)
	}

	group := internal.UserConnectionGroupInfo{
		GroupName:         groupName,
		ConnectionUserIds: ids,
		GroupPic:          "",
	}
	ctlr.DB.CreateUserConnectionGroup(userID, group)
}
func TestCreateConnectionsGroupsByUserID(t *testing.T) {

	groupName := "New Group Name"
	groupName1 := "New Group Name 1"
	testCases := []TestCaseCreateGroup{
		{
			name: "Created",
			inputParams: connections.UsersConnectionsGroupsByUserIDPostParams{
				UserID: "dc9dbe3e-60d5-4a07-8c9c-42027b555b01",
				Body: &models.UsersConnectionsGroupsPostRequest{
					GroupName:         &groupName,
					ConnectionUserIds: connectionUserIds,
					GroupPic:          "",
				},
			},
			inputPrincipal:       &models.Principal{},
			expectedResponseType: "Created",
			expectedErr:          nil,
			expectedErrMsg:       "",
		},
		{
			name: "GroupExists",
			inputParams: connections.UsersConnectionsGroupsByUserIDPostParams{
				UserID: "dc9dbe3e-60d5-4a07-8c9c-42027b555b01",
				Body: &models.UsersConnectionsGroupsPostRequest{
					GroupName:         &groupName,
					ConnectionUserIds: connectionUserIds,
					GroupPic:          "",
				},
			},
			inputPrincipal:       &models.Principal{},
			expectedResponseType: "OK",
			expectedErr:          nil,
			expectedErrMsg:       "Group Already Exists",
		},
		{
			name: "WithPicFailed",
			inputParams: connections.UsersConnectionsGroupsByUserIDPostParams{
				UserID: "dc9dbe3e-60d5-4a07-8c9c-42027b555b01",
				Body: &models.UsersConnectionsGroupsPostRequest{
					GroupName:         &groupName1,
					ConnectionUserIds: connectionUserIds,
					GroupPic:          "fake_image_base64",
				},
			},
			inputPrincipal:       &models.Principal{},
			expectedResponseType: "OK",
			expectedErr:          nil,
			expectedErrMsg:       "Failed to reduce image size",
		},
	}

	for _, test := range testCases {

		t.Run(test.name, func(t *testing.T) {
			res := ctlr.CreateConnectionsGroupsByUserID(
				test.inputParams,
				test.inputPrincipal,
			)

			t.Logf("Actual Response: %+v\n", res)
			assertEqual(t, res.resType, test.expectedResponseType)
			assertEqual(t, res.err, test.expectedErr)
			assertEqual(t, res.errMsg, test.expectedErrMsg)
		})
	}
}

type TestCaseGetGroup struct {
	name                 string
	inputParams          connections.UsersConnectionsGroupsByUserIDAndGroupIDGetParams
	inputPrincipal       *models.Principal
	expectedResponseType string
	expectedErrMsg       string
	expectedErr          error
}

func TestGetUsersConnectionsGroupsByUserIDAndGroupID(t *testing.T) {

	// In Case if you need to add group manually into memory db then uncomment below code
	// CreateConnectionGroup()

	testCases := []TestCaseGetGroup{
		{
			name: "OK",
			inputParams: connections.UsersConnectionsGroupsByUserIDAndGroupIDGetParams{
				UserID:  "dc9dbe3e-60d5-4a07-8c9c-42027b555b01",
				GroupID: "group_id_1",
			},
			inputPrincipal:       &models.Principal{},
			expectedResponseType: "OK",
			expectedErr:          nil,
			expectedErrMsg:       "",
		},
		{
			name: "NotFound",
			inputParams: connections.UsersConnectionsGroupsByUserIDAndGroupIDGetParams{
				UserID:  "dc9dbe3e-60d5-4a07-8c9c-42027b555b01",
				GroupID: "group_id_5",
			},
			inputPrincipal:       &models.Principal{},
			expectedResponseType: "errReturn404",
			expectedErr:          status.Error(codes.NotFound, "row does not found"),
			expectedErrMsg:       "record not found",
		},
		{
			name: "InternalError",
			inputParams: connections.UsersConnectionsGroupsByUserIDAndGroupIDGetParams{
				UserID:  "dc9dbe3e-60d5-4a07-8c9c-42027b55502",
				GroupID: "group_id_3",
			},
			inputPrincipal:       &models.Principal{},
			expectedResponseType: "errReturn500",
			expectedErr:          status.Error(codes.Internal, "something went wrong"),
			expectedErrMsg:       "failed to parse group from database",
		},
	}

	for _, test := range testCases {

		t.Run(test.name, func(t *testing.T) {
			res := ctlr.GetUsersConnectionsGroupsByUserIDAndGroupID(
				test.inputParams,
				test.inputPrincipal,
			)

			t.Logf("Actual Response: %+v\n", res)
			assertEqual(t, res.resType, test.expectedResponseType)
			assertEqual(t, res.err, test.expectedErr)
			assertEqual(t, res.errMsg, test.expectedErrMsg)
		})
	}
}

type TestCaseGetGroups struct {
	name                 string
	inputParams          connections.UsersConnectionsGroupsByUserIDGetParams
	inputPrincipal       *models.Principal
	expectedResponseType string
	expectedErrMsg       string
	expectedErr          error
}

func TestGetUsersConnectionsGroupsByUserID(t *testing.T) {

	// In Case if you need to add group manually into memory db then uncomment below code
	// CreateConnectionGroup()

	limit := int32(10)
	offset := int32(0)
	order := "asc"

	testCases := []TestCaseGetGroups{
		{
			name: "OK",
			inputParams: connections.UsersConnectionsGroupsByUserIDGetParams{
				UserID: "dc9dbe3e-60d5-4a07-8c9c-42027b555b01",
				Limit:  &limit,
				Offset: &offset,
				Order:  &order,
			},
			inputPrincipal:       &models.Principal{},
			expectedResponseType: "OK",
			expectedErr:          nil,
			expectedErrMsg:       "",
		},
		{
			name: "NotFound",
			inputParams: connections.UsersConnectionsGroupsByUserIDGetParams{
				UserID: "dc9dbe3e-60d5-4a07-8c9c-42027b555b02",
				Limit:  &limit,
				Offset: &offset,
				Order:  &order,
			},
			inputPrincipal:       &models.Principal{},
			expectedResponseType: "errReturn404",
			expectedErr:          status.Error(codes.NotFound, "row does not found"),
			expectedErrMsg:       "records not found",
		},
	}

	for _, test := range testCases {

		t.Run(test.name, func(t *testing.T) {
			res := ctlr.GetUsersConnectionsGroupsByUserID(
				test.inputParams,
				test.inputPrincipal,
			)

			t.Logf("Actual Response: %+v\n", res)
			assertEqual(t, res.resType, test.expectedResponseType)
			assertEqual(t, res.err, test.expectedErr)
			assertEqual(t, res.errMsg, test.expectedErrMsg)
		})
	}
}

type TestCaseUpdateGroup struct {
	name                 string
	inputParams          connections.UsersConnectionsGroupsByUserIDAndGroupIDPatchParams
	inputPrincipal       *models.Principal
	expectedResponseType string
	expectedErrMsg       string
	expectedErr          error
}

func TestUpdateUsersConnectionsGroupsByUserIDAndGroupID(t *testing.T) {

	// In Case if you need to add group manually into memory db then uncomment below code
	CreateConnectionGroup()

	testCases := []TestCaseUpdateGroup{
		{
			name: "Updated",
			inputParams: connections.UsersConnectionsGroupsByUserIDAndGroupIDPatchParams{
				UserID:  "dc9dbe3e-60d5-4a07-8c9c-42027b555b01",
				GroupID: "group_id_1",
				Body: &models.UsersConnectionsGroupsPatchRequest{
					GroupName:                "Update Group Name",
					ConnectionUserIDToAdd:    "1ca26428-98eb-4aa3-8943-5f459873ef85",
					ConnectionUserIDToRemove: "dc9dbe3e-60d5-4a07-8c9c-42027b555b01",
				},
			},
			inputPrincipal:       &models.Principal{},
			expectedResponseType: "Updated",
			expectedErr:          nil,
			expectedErrMsg:       "",
		},
		{
			name: "OKGroupExists",
			inputParams: connections.UsersConnectionsGroupsByUserIDAndGroupIDPatchParams{
				UserID:  "dc9dbe3e-60d5-4a07-8c9c-42027b555b01",
				GroupID: "group_id_2",
				Body: &models.UsersConnectionsGroupsPatchRequest{
					GroupName:                "Update Group Name",
					ConnectionUserIDToAdd:    "1ca26428-98eb-4aa3-8943-5f459873ef85",
					ConnectionUserIDToRemove: "dc9dbe3e-60d5-4a07-8c9c-42027b555b01",
				},
			},
			inputPrincipal:       &models.Principal{},
			expectedResponseType: "OK",
			expectedErr:          nil,
			expectedErrMsg:       "Group name is already in use, choose another group name.",
		},
		{
			name: "OKPicFailed",
			inputParams: connections.UsersConnectionsGroupsByUserIDAndGroupIDPatchParams{
				UserID:  "dc9dbe3e-60d5-4a07-8c9c-42027b555b01",
				GroupID: "group_id_2",
				Body: &models.UsersConnectionsGroupsPatchRequest{
					GroupName:                "Update Group Name 2",
					ConnectionUserIDToAdd:    "1ca26428-98eb-4aa3-8943-5f459873ef85",
					ConnectionUserIDToRemove: "dc9dbe3e-60d5-4a07-8c9c-42027b555b01",
					GroupPic:                 "fake_image",
				},
			},
			inputPrincipal:       &models.Principal{},
			expectedResponseType: "OK",
			expectedErr:          nil,
			expectedErrMsg:       "Failed to reduce image size",
		},
		{
			name: "InternalError",
			inputParams: connections.UsersConnectionsGroupsByUserIDAndGroupIDPatchParams{
				UserID:  "dc9dbe3e-60d5-4a07-8c9c-42027b555b03",
				GroupID: "group_id_2",
				Body: &models.UsersConnectionsGroupsPatchRequest{
					GroupName:                "Update Group Name 2",
					ConnectionUserIDToAdd:    "1ca26428-98eb-4aa3-8943-5f459873ef85",
					ConnectionUserIDToRemove: "dc9dbe3e-60d5-4a07-8c9c-42027b555b01",
				},
			},
			inputPrincipal:       &models.Principal{},
			expectedResponseType: "errReturn500",
			expectedErr:          status.Error(codes.Internal, "something went wrong"),
			expectedErrMsg:       "failed to parse group from database",
		},
	}

	for _, test := range testCases {

		t.Run(test.name, func(t *testing.T) {
			res := ctlr.UpdateUsersConnectionsGroupsByUserIDAndGroupID(
				test.inputParams,
				test.inputPrincipal,
			)

			t.Logf("Actual Response: %+v\n", res)
			assertEqual(t, res.resType, test.expectedResponseType)
			assertEqual(t, res.err, test.expectedErr)
			assertEqual(t, res.errMsg, test.expectedErrMsg)
		})
	}
}

type TestCaseDeleteGroup struct {
	name                 string
	inputParams          connections.UsersConnectionsGroupsByUserIDAndGroupIDDeleteParams
	inputPrincipal       *models.Principal
	expectedResponseType string
	expectedErrMsg       string
	expectedErr          error
}

func TestDeleteUsersConnectionsGroupsByUserIDAndGroupID(t *testing.T) {
	// In Case if you need to add group manually into memory db then uncomment below code
	CreateConnectionGroup()

	testCases := []TestCaseDeleteGroup{
		{
			name: "Deleted",
			inputParams: connections.UsersConnectionsGroupsByUserIDAndGroupIDDeleteParams{
				UserID:  "dc9dbe3e-60d5-4a07-8c9c-42027b555b01",
				GroupID: "group_id_1",
			},
			inputPrincipal:       &models.Principal{},
			expectedResponseType: "Deleted",
			expectedErr:          nil,
			expectedErrMsg:       "",
		},
		{
			name: "NotFound",
			inputParams: connections.UsersConnectionsGroupsByUserIDAndGroupIDDeleteParams{
				UserID:  "dc9dbe3e-60d5-4a07-8c9c-42027b555b01",
				GroupID: "group_id_1",
			},
			inputPrincipal:       &models.Principal{},
			expectedResponseType: "errReturn404",
			expectedErr:          status.Error(codes.NotFound, "row does not found"),
			expectedErrMsg:       "record not found",
		},
		{
			name: "InternalError",
			inputParams: connections.UsersConnectionsGroupsByUserIDAndGroupIDDeleteParams{
				UserID:  "dc9dbe3e-60d5-4a07-8c9c-42027b555b03",
				GroupID: "group_id_1",
			},
			inputPrincipal:       &models.Principal{},
			expectedResponseType: "errReturn500",
			expectedErr:          status.Error(codes.Internal, "something went wrong"),
			expectedErrMsg:       "failed to parse group from database",
		},
	}

	for _, test := range testCases {

		t.Run(test.name, func(t *testing.T) {
			res := ctlr.DeleteUsersConnectionsGroupsByUserIDAndGroupID(
				test.inputParams,
				test.inputPrincipal,
			)

			t.Logf("Actual Response: %+v\n", res)
			assertEqual(t, res.resType, test.expectedResponseType)
			assertEqual(t, res.err, test.expectedErr)
			assertEqual(t, res.errMsg, test.expectedErrMsg)
		})
	}
}
