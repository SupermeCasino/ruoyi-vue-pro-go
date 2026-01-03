package member

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	memberModel "github.com/wxlbd/ruoyi-mall-go/internal/model/member"
	memberSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type MemberUserHandler struct {
	userSvc  *memberSvc.MemberUserService
	levelSvc *memberSvc.MemberLevelService
	pointSvc *memberSvc.MemberPointRecordService
	groupSvc *memberSvc.MemberGroupService
	tagSvc   *memberSvc.MemberTagService
}

func NewMemberUserHandler(
	userSvc *memberSvc.MemberUserService,
	levelSvc *memberSvc.MemberLevelService,
	pointSvc *memberSvc.MemberPointRecordService,
	groupSvc *memberSvc.MemberGroupService,
	tagSvc *memberSvc.MemberTagService,
) *MemberUserHandler {
	return &MemberUserHandler{
		userSvc:  userSvc,
		levelSvc: levelSvc,
		pointSvc: pointSvc,
		groupSvc: groupSvc,
		tagSvc:   tagSvc,
	}
}

// UpdateUser 更新会员用户
func (h *MemberUserHandler) UpdateUser(c *gin.Context) {
	var r member.MemberUserUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	err := h.userSvc.AdminUpdateUser(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// UpdateUserLevel 更新会员等级
func (h *MemberUserHandler) UpdateUserLevel(c *gin.Context) {
	var r member.MemberUserUpdateLevelReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	err := h.levelSvc.UpdateUserLevel(c, r.ID, &r.LevelID, r.Reason)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// UpdateUserPoint 更新会员积分
func (h *MemberUserHandler) UpdateUserPoint(c *gin.Context) {
	var r member.MemberUserUpdatePointReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// 对应 Java: memberPointRecordService.createPointRecord(updateReqVO.getId(), updateReqVO.getPoint(),
	//           MemberPointBizTypeEnum.ADMIN, String.valueOf(getLoginUserId()));
	bizId := utils.ToString(context.GetLoginUserID(c))

	err := h.pointSvc.CreatePointRecord(c, r.ID, int(r.Point), consts.MemberPointBizTypeAdmin, bizId)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetUser 获得会员用户详情
func (h *MemberUserHandler) GetUser(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	user, err := h.userSvc.GetUser(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 关联查询
	var levelName, groupName string
	var tagNames []string

	if user.LevelID > 0 {
		if level, _ := h.levelSvc.GetLevel(c, user.LevelID); level != nil {
			levelName = level.Name
		}
	}
	if user.GroupID > 0 {
		if group, _ := h.groupSvc.GetGroup(c, user.GroupID); group != nil {
			groupName = group.Name
		}
	}
	if len(user.TagIds) > 0 {
		if tags, _ := h.tagSvc.GetTagListByIds(c, lo.Map(user.TagIds, func(id int, _ int) int64 { return int64(id) })); tags != nil {
			tagNames = lo.Map(tags, func(t *memberModel.MemberTag, _ int) string { return t.Name })
		}
	}

	response.WriteSuccess(c, h.convertRespWithExt(user, tagNames, levelName, groupName))
}

// GetUserPage 获得会员用户分页
func (h *MemberUserHandler) GetUserPage(c *gin.Context) {
	var r member.MemberUserPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	pageResult, err := h.userSvc.GetUserPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 批量获取关联数据
	levelIds := lo.Uniq(lo.FilterMap(pageResult.List, func(u *memberModel.MemberUser, _ int) (int64, bool) {
		return u.LevelID, u.LevelID > 0
	}))
	groupIds := lo.Uniq(lo.FilterMap(pageResult.List, func(u *memberModel.MemberUser, _ int) (int64, bool) {
		return u.GroupID, u.GroupID > 0
	}))
	tagIds := lo.Uniq(lo.Flatten(lo.Map(pageResult.List, func(u *memberModel.MemberUser, _ int) []int64 {
		return lo.Map(u.TagIds, func(id int, _ int) int64 { return int64(id) })
	})))

	levelMap := make(map[int64]string)
	groupMap := make(map[int64]string)
	tagMap := make(map[int64]string)

	if len(levelIds) > 0 {
		if levels, _ := h.levelSvc.GetLevelList(c, levelIds); levels != nil {
			for _, l := range levels {
				levelMap[l.ID] = l.Name
			}
		}
	}
	if len(groupIds) > 0 {
		if groups, _ := h.groupSvc.GetGroupListByIds(c, groupIds); groups != nil {
			for _, g := range groups {
				groupMap[g.ID] = g.Name
			}
		}
	}
	if len(tagIds) > 0 {
		if tags, _ := h.tagSvc.GetTagListByIds(c, tagIds); tags != nil {
			for _, t := range tags {
				tagMap[t.ID] = t.Name
			}
		}
	}

	respList := lo.Map(pageResult.List, func(user *memberModel.MemberUser, _ int) *member.MemberUserResp {
		var tagNames []string
		for _, tid := range user.TagIds {
			if name, ok := tagMap[int64(tid)]; ok {
				tagNames = append(tagNames, name)
			}
		}
		return h.convertRespWithExt(user, tagNames, levelMap[user.LevelID], groupMap[user.GroupID])
	})

	response.WritePage(c, pageResult.Total, respList)
}

func (h *MemberUserHandler) convertRespWithExt(user *memberModel.MemberUser, tagNames []string, levelName, groupName string) *member.MemberUserResp {
	if user == nil {
		return nil
	}
	return &member.MemberUserResp{
		ID:         user.ID,
		Mobile:     user.Mobile,
		Status:     user.Status,
		Nickname:   user.Nickname,
		Avatar:     user.Avatar,
		Name:       user.Name,
		Sex:        user.Sex,
		AreaID:     int64(user.AreaID),
		Birthday:   user.Birthday,
		Mark:       user.Mark,
		TagIDs:     lo.Map(user.TagIds, func(id int, _ int) int64 { return int64(id) }),
		LevelID:    user.LevelID,
		GroupID:    user.GroupID,
		RegisterIP: user.RegisterIP,
		LoginIP:    user.LoginIP,
		LoginDate:  user.LoginDate,
		CreateTime: user.CreateTime,
		Point:      user.Point,
		Experience: user.Experience,
		TagNames:   tagNames,
		LevelName:  levelName,
		GroupName:  groupName,
	}
}
