package services

import (
	"context"

	apimodels "github.com/juggleim/imserver-console/apis/models"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/imsdk"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/dbs"
)

func QryGroups(ctx context.Context, appkey, groupId, name string, offset string, limit int64, isPositive bool) (errs.AdminErrorCode, *apimodels.Groups) {
	var startId int64
	if offset != "" {
		startId, _ = tools.DecodeInt(offset)
	}
	ret := &apimodels.Groups{
		Items: []*apimodels.Group{},
	}
	storage := dbs.GroupDao{}
	memberStorage := dbs.GroupMemberDao{}
	if groupId != "" {
		grp, err := storage.FindById(appkey, groupId)
		if err == nil && grp != nil {
			ret.Items = append(ret.Items, groupDaoToAPI(appkey, grp, memberStorage))
		}
	} else {
		grps, err := storage.QryGroups(appkey, name, startId, limit, isPositive)
		if err == nil {
			for _, grp := range grps {
				ret.Offset, _ = tools.EncodeInt(grp.ID)
				ret.Items = append(ret.Items, groupDaoToAPI(appkey, grp, memberStorage))
			}
		}
	}
	return errs.AdminErrorCode_Success, ret
}

func QryGroupInfo(appkey, groupId string) *apimodels.Group {
	storage := dbs.GroupDao{}
	grp, err := storage.FindById(appkey, groupId)
	if err != nil || grp == nil {
		return &apimodels.Group{GroupId: groupId}
	}
	return &apimodels.Group{
		GroupId:       grp.GroupId,
		GroupName:     grp.GroupName,
		GroupPortrait: grp.GroupPortrait,
	}
}

func groupDaoToAPI(appkey string, grp *dbs.GroupDao, memberStorage dbs.GroupMemberDao) *apimodels.Group {
	return &apimodels.Group{
		GroupId:       grp.GroupId,
		GroupName:     grp.GroupName,
		GroupPortrait: grp.GroupPortrait,
		Owner:         QryUserInfo(appkey, grp.CreatorId),
		CreatedTime:   grp.CreatedTime.UnixMilli(),
		MemberCount:   memberStorage.CountByGroup(appkey, grp.GroupId),
	}
}

func DissolveGroups(ctx context.Context, req *apimodels.GroupIds) errs.AdminErrorCode {
	appkey := req.AppKey
	sdk := imsdk.GetImSdk(appkey)
	storage := dbs.GroupDao{}
	memberStorage := dbs.GroupMemberDao{}
	for _, groupId := range req.GroupIds {
		_ = storage.Delete(appkey, groupId)
		_ = memberStorage.DeleteByGroupId(appkey, groupId)
		if sdk != nil {
			_, _, _ = sdk.DissolveGroup(groupId)
		}
	}
	return errs.AdminErrorCode_Success
}
