package dbs

import (
	"github.com/juggleim/imserver-console/commons/dbcommons"
)

type GroupMemberDao struct {
	ID             int64  `gorm:"primary_key"`
	GroupId        string `gorm:"group_id"`
	MemberId       string `gorm:"member_id"`
	MemberType     int    `gorm:"member_type"`
	AppKey         string `gorm:"app_key"`
	IsMute         int    `gorm:"is_mute"`
	IsAllow        int    `gorm:"is_allow"`
	MuteEndAt      int64  `gorm:"mute_end_at"`
	GrpDisplayName string `gorm:"grp_display_name"`
}

func (msg GroupMemberDao) TableName() string {
	return "groupmembers"
}

func (member GroupMemberDao) DeleteByGroupId(appkey, groupId string) error {
	return dbcommons.GetDb().Where("app_key=? and group_id=?", appkey, groupId).Delete(&GroupMemberDao{}).Error
}

func (member GroupMemberDao) CountByGroup(appkey, groupId string) int {
	var count int64
	err := dbcommons.GetDb().Model(&GroupMemberDao{}).Where("app_key=? and group_id=?", appkey, groupId).Count(&count).Error
	if err != nil {
		return 0
	}
	return int(count)
}
