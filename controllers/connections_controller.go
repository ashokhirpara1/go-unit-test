package controllers

import (
	"log"
	"strings"
	"time"
)

// CreateConnectionsGroupsByUserIDResponse - Holding reponse for CreateConnectionsGroupsByUserID()
type CreateConnectionsGroupsByUserIDResponse struct {
	existsPayload  models.UsersConnectionsGroupsExistsPostResponse
	createdPayload models.UsersConnectionsGroupsPostResponse
	resType        string
	errMsg         string
	err            error
}

// GetUsersConnectionsGroupsByUserIDAndGroupIDResponse - Holding reponse for GetUsersConnectionsGroupsByUserIDAndGroupID()
type GetUsersConnectionsGroupsByUserIDAndGroupIDResponse struct {
	payload models.UsersConnectionsGroupsResponse
	resType string
	errMsg  string
	err     error
}

// GetUsersConnectionsGroupsByUserIDResponse - Holding reponse for GetUsersConnectionsGroupsByUserID()
type GetUsersConnectionsGroupsByUserIDResponse struct {
	payload models.UsersConnectionsGroupsGetResponse
	resType string
	errMsg  string
	err     error
}

// UpdateUsersConnectionsGroupsByUserIDAndGroupIDResponse - Holding reponse for UpdateUsersConnectionsGroupsByUserIDAndGroupID()
type UpdateUsersConnectionsGroupsByUserIDAndGroupIDResponse struct {
	existsPayload models.UsersConnectionsGroupsExistsPostResponse
	resType       string
	errMsg        string
	err           error
}

// DeleteUsersConnectionsGroupsByUserIDAndGroupIDResponse - Holding reponse for DeleteUsersConnectionsGroupsByUserIDAndGroupID()
type DeleteUsersConnectionsGroupsByUserIDAndGroupIDResponse struct {
	resType string
	errMsg  string
	err     error
}

// CreateConnectionsGroupsByUserIDController -
func CreateConnectionsGroupsByUserIDController(params connections.UsersConnectionsGroupsByUserIDPostParams, principal *models.Principal) middleware.Responder {

	ctlr, err := GetControllerDB()
	if err != nil {
		return connections.NewUsersConnectionsGroupsByUserIDPostInternalServerError()
	}

	response := ctlr.CreateConnectionsGroupsByUserID(params, principal)
	if response.err != nil {
		switch response.resType {
		case "errReturn500":
			return connections.NewUsersConnectionsGroupsByUserIDPostInternalServerError()
		}
	}

	if response.resType == "OK" {
		return connections.NewUsersConnectionsGroupsByUserIDPostOK().WithPayload(&response.existsPayload)
	}

	return connections.NewUsersConnectionsGroupsByUserIDPostCreated().WithPayload(&response.createdPayload)
}

// UsersConnectionsGroupsByUserIDAndGroupIDGetController - Get an individual Connections Group.
func UsersConnectionsGroupsByUserIDAndGroupIDGetController(params connections.UsersConnectionsGroupsByUserIDAndGroupIDGetParams, principal *models.Principal) middleware.Responder {

	ctlr, err := GetControllerDB()
	if err != nil {
		return connections.NewUsersConnectionsGroupsByUserIDAndGroupIDGetInternalServerError()
	}

	response := ctlr.GetUsersConnectionsGroupsByUserIDAndGroupID(params, principal)
	if response.err != nil {
		switch response.resType {
		case "errReturn404":
			return connections.NewUsersConnectionsGroupsByUserIDAndGroupIDGetNotFound()
		case "errReturn500":
			return connections.NewUsersConnectionsGroupsByUserIDAndGroupIDGetInternalServerError()
		}
	}

	return connections.NewUsersConnectionsGroupsByUserIDAndGroupIDGetOK().WithPayload(&response.payload)
}

// UsersConnectionsGroupsByUserIDGetController - Get a batch of Users Connections Groups.
func UsersConnectionsGroupsByUserIDGetController(params connections.UsersConnectionsGroupsByUserIDGetParams, principal *models.Principal) middleware.Responder {

	ctlr, err := GetControllerDB()
	if err != nil {
		return connections.NewUsersConnectionsGroupsByUserIDGetInternalServerError()
	}

	response := ctlr.GetUsersConnectionsGroupsByUserID(params, principal)
	if response.err != nil {
		switch response.resType {
		case "errReturn404":
			return connections.NewUsersConnectionsGroupsByUserIDGetNotFound()
		case "errReturn500":
			return connections.NewUsersConnectionsGroupsByUserIDGetInternalServerError()
		}
	}

	return connections.NewUsersConnectionsGroupsByUserIDGetOK().WithPayload(&response.payload)
}

// UsersConnectionsGroupsByUserIDAndGroupIDPatchController - Updates a specific user's group.
func UsersConnectionsGroupsByUserIDAndGroupIDPatchController(params connections.UsersConnectionsGroupsByUserIDAndGroupIDPatchParams, principal *models.Principal) middleware.Responder {

	ctlr, err := GetControllerDB()
	if err != nil {
		return connections.NewUsersConnectionsGroupsByUserIDAndGroupIDPatchInternalServerError()
	}

	response := ctlr.UpdateUsersConnectionsGroupsByUserIDAndGroupID(params, principal)
	if response.err != nil {
		switch response.resType {
		case "errReturn500":
			return connections.NewUsersConnectionsGroupsByUserIDAndGroupIDPatchInternalServerError()
		}
	}

	if response.errMsg != "" {
		return connections.NewUsersConnectionsGroupsByUserIDPostOK().WithPayload(&response.existsPayload)
	}

	return connections.NewUsersConnectionsGroupsByUserIDAndGroupIDPatchOK()
}

// UsersConnectionsGroupsByUserIDAndGroupIDDeleteController - Delete an individual Connections Group.
func UsersConnectionsGroupsByUserIDAndGroupIDDeleteController(params connections.UsersConnectionsGroupsByUserIDAndGroupIDDeleteParams, principal *models.Principal) middleware.Responder {

	ctlr, err := GetControllerDB()
	if err != nil {
		return connections.NewUsersConnectionsGroupsByUserIDAndGroupIDDeleteInternalServerError()
	}

	response := ctlr.DeleteUsersConnectionsGroupsByUserIDAndGroupID(params, principal)
	if response.err != nil {
		switch response.resType {
		case "errReturn404":
			return connections.NewUsersConnectionsGroupsByUserIDAndGroupIDDeleteNotFound()
		case "errReturn500":
			return connections.NewUsersConnectionsGroupsByUserIDAndGroupIDDeleteInternalServerError()
		}
	}

	return connections.NewUsersConnectionsGroupsByUserIDAndGroupIDDeleteNoContent()
}

// CreateConnectionsGroupsByUserID -
func (c Ctlr) CreateConnectionsGroupsByUserID(params connections.UsersConnectionsGroupsByUserIDPostParams, principal *models.Principal) CreateConnectionsGroupsByUserIDResponse {

	groupinfoObj, err := c.DB.GetUserConnectionGroupByName(params.UserID, *params.Body.GroupName)
	if err != nil && status.Code(err) != codes.NotFound {
		return CreateConnectionsGroupsByUserIDResponse{resType: "errReturn500", errMsg: "failed to parse group from database", err: err}
	}

	if !IsZeroOfUnderlyingType(groupinfoObj.GroupID) {
		var msg string = "Group Already Exists"
		responseExistsPayload := models.UsersConnectionsGroupsExistsPostResponse{
			ErrorMessage: &msg,
			GroupID:      &groupinfoObj.GroupID,
			GroupName:    &groupinfoObj.GroupName,
		}

		return CreateConnectionsGroupsByUserIDResponse{resType: "OK", errMsg: msg, existsPayload: responseExistsPayload}
	}

	var reducedGroupPic *string
	if !IsZeroOfUnderlyingType(params.Body.GroupPic) {

		start := time.Now()

		if strings.Contains(params.Body.GroupPic, "base64,") {
			params.Body.GroupPic = params.Body.GroupPic[strings.IndexByte(params.Body.GroupPic, ',')+1:]
		}
		max := GetGroupPicMaxSizeBytes()
		min := int(float64(max) * ProfileGraphicRatioMinThresholdDefault)
		sizeSpecs := ImageOptions{
			MaxImageSizeBytes: &max,
			MinImageSizeBytes: &min,
		}
		reducedBgImageData, err := ReduceBase64EncodedImage(params.Body.GroupPic, &sizeSpecs)
		if err != nil {
			log.Printf("failed to save group picture for user (%s) (%s)", params.UserID, err.Error())
			msg := err.Error()
			responseExistsPayload := models.UsersConnectionsGroupsExistsPostResponse{
				ErrorMessage: &msg,
			}
			return CreateConnectionsGroupsByUserIDResponse{resType: "OK", errMsg: "Failed to reduce image size", existsPayload: responseExistsPayload}
		}

		reducedGroupPic = reducedBgImageData

		elapsed := time.Since(start)

		log.Printf("Image Processing took %s", elapsed)
	}

	ids := []GroupConnectionUserID{}

	for _, CU := range params.Body.ConnectionUserIds {
		id := GroupConnectionUserID{UserID: CU.UserID}
		ids = append(ids, id)
	}

	groupPicStr := ""
	if reducedGroupPic != nil {
		groupPicStr = *reducedGroupPic
	}
	group := UserConnectionGroupInfo{
		GroupName:             *params.Body.GroupName,
		ConnectionUserIds:     ids,
		GroupPic:              groupPicStr,
		LatestInteractionTime: time.Now(),
	}

	// Set the group into the database.
	groupID, err := c.DB.CreateUserConnectionGroup(params.UserID, group)

	if err != nil {
		return CreateConnectionsGroupsByUserIDResponse{resType: "errReturn500", errMsg: "failed to create new Group entry in database", err: err}
	}

	responsePayload := models.UsersConnectionsGroupsPostResponse{
		GroupID: &groupID,
	}

	return CreateConnectionsGroupsByUserIDResponse{resType: "Created", createdPayload: responsePayload}
}

// GetUsersConnectionsGroupsByUserIDAndGroupID -
func (c Ctlr) GetUsersConnectionsGroupsByUserIDAndGroupID(params connections.UsersConnectionsGroupsByUserIDAndGroupIDGetParams, principal *models.Principal) GetUsersConnectionsGroupsByUserIDAndGroupIDResponse {

	groupInfo, err := c.DB.GetUserConnectionGroupByGroupID(params.UserID, params.GroupID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return GetUsersConnectionsGroupsByUserIDAndGroupIDResponse{resType: "errReturn404", errMsg: "record not found", err: err}
		}
		return GetUsersConnectionsGroupsByUserIDAndGroupIDResponse{resType: "errReturn500", errMsg: "failed to parse group from database", err: err}
	}

	groupData := groupInfo.TransformToResponseGroup()

	payload := models.UsersConnectionsGroupsResponse{
		Group: groupData,
	}

	return GetUsersConnectionsGroupsByUserIDAndGroupIDResponse{resType: "OK", payload: payload}
}

// GetUsersConnectionsGroupsByUserID - Get a batch of Users Connections Groups.
func (c Ctlr) GetUsersConnectionsGroupsByUserID(params connections.UsersConnectionsGroupsByUserIDGetParams, principal *models.Principal) GetUsersConnectionsGroupsByUserIDResponse {

	groupsList, paginationMeta, err := c.DB.GetPaginatedUserConnectionGroup(params)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return GetUsersConnectionsGroupsByUserIDResponse{resType: "errReturn404", errMsg: "records not found", err: err}
		}
		return GetUsersConnectionsGroupsByUserIDResponse{resType: "errReturn500", errMsg: "failed to parse groups from database", err: err}
	}

	payload := models.UsersConnectionsGroupsGetResponse{
		Groups:             groupsList,
		PaginationMetadata: paginationMeta,
	}

	return GetUsersConnectionsGroupsByUserIDResponse{resType: "OK", payload: payload}
}

// UpdateUsersConnectionsGroupsByUserIDAndGroupID -
func (c Ctlr) UpdateUsersConnectionsGroupsByUserIDAndGroupID(params connections.UsersConnectionsGroupsByUserIDAndGroupIDPatchParams, principal *models.Principal) UpdateUsersConnectionsGroupsByUserIDAndGroupIDResponse {

	groupinfoObj, err := c.DB.GetUserConnectionGroupByName(params.UserID, params.Body.GroupName)
	if err != nil && status.Code(err) != codes.NotFound {
		return UpdateUsersConnectionsGroupsByUserIDAndGroupIDResponse{resType: "errReturn500", errMsg: "failed to parse group from database", err: err}
	}

	if !IsZeroOfUnderlyingType(groupinfoObj.GroupID) {
		var msg string = "Group name is already in use, choose another group name."
		responseExistsPayload := models.UsersConnectionsGroupsExistsPostResponse{
			ErrorMessage: &msg,
			GroupID:      &groupinfoObj.GroupID,
			GroupName:    &groupinfoObj.GroupName,
		}

		return UpdateUsersConnectionsGroupsByUserIDAndGroupIDResponse{resType: "OK", errMsg: msg, existsPayload: responseExistsPayload}
	}

	var reducedGroupPic *string
	if !IsZeroOfUnderlyingType(params.Body.GroupPic) {

		start := time.Now()

		if strings.Contains(params.Body.GroupPic, "base64,") {
			params.Body.GroupPic = params.Body.GroupPic[strings.IndexByte(params.Body.GroupPic, ',')+1:]
		}
		max := GetGroupPicMaxSizeBytes()
		min := int(float64(max) * ProfileGraphicRatioMinThresholdDefault)
		sizeSpecs := ImageOptions{
			MaxImageSizeBytes: &max,
			MinImageSizeBytes: &min,
		}
		reducedBgImageData, err := ReduceBase64EncodedImage(params.Body.GroupPic, &sizeSpecs)
		if err != nil {
			log.Printf("failed to save group picture for user (%s) (%s)", params.UserID, err.Error())
			msg := err.Error()
			responseExistsPayload := models.UsersConnectionsGroupsExistsPostResponse{
				ErrorMessage: &msg,
			}

			return UpdateUsersConnectionsGroupsByUserIDAndGroupIDResponse{resType: "OK", errMsg: "Failed to reduce image size", existsPayload: responseExistsPayload}
		}

		reducedGroupPic = reducedBgImageData

		elapsed := time.Since(start)

		log.Printf("Image Processing took %s", elapsed)
	}

	params.Body.GroupPic = ""
	if !IsZeroOfUnderlyingType(reducedGroupPic) {
		params.Body.GroupPic = *reducedGroupPic
	}

	err = c.DB.UpdateUserConnectionGroup(params)
	if err != nil {
		return UpdateUsersConnectionsGroupsByUserIDAndGroupIDResponse{resType: "errReturn500", errMsg: "failed to parse group from database", err: err}
	}

	return UpdateUsersConnectionsGroupsByUserIDAndGroupIDResponse{resType: "Updated"}
}

// DeleteUsersConnectionsGroupsByUserIDAndGroupID -
func (c Ctlr) DeleteUsersConnectionsGroupsByUserIDAndGroupID(params connections.UsersConnectionsGroupsByUserIDAndGroupIDDeleteParams, principal *models.Principal) DeleteUsersConnectionsGroupsByUserIDAndGroupIDResponse {

	err := c.DB.DeleteUserConnectionGroup(params)
	if err != nil {

		if status.Code(err) == codes.NotFound {
			return DeleteUsersConnectionsGroupsByUserIDAndGroupIDResponse{resType: "errReturn404", errMsg: "record not found", err: err}
		}

		return DeleteUsersConnectionsGroupsByUserIDAndGroupIDResponse{resType: "errReturn500", errMsg: "failed to parse group from database", err: err}
	}

	return DeleteUsersConnectionsGroupsByUserIDAndGroupIDResponse{resType: "Deleted"}
}
