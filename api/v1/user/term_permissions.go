package user

import (
	"cmdb/model"
	"cmdb/utils/errmsg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ServerPermissions(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	name := c.Query("name")
	private_ip_address := c.Query("private_ip_address")
	data, code := model.AllPermissions(id, name, private_ip_address)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"total":   len(data),
			"message": errmsg.GetErrMsg(code),
		},
	)
}

func AddPermissions(c *gin.Context) {
	var data model.AddPermissions
	_ = c.ShouldBindJSON(&data)
	var serverIds, userIds, termuserIds []int
	/*
		创建机器ssh权限
		查所有server，返回ip，name，id，批量选中ip
		查所有终端用户 返回 ,对应账号和id
		创建管理关系： serverid
		创建serverid=s.serverid 以及term_user_id= t.id
	*/
	servers, _ := model.GetServers(0, 0)
	for _, i := range data.Ips {
		for _, v := range servers {
			if i == "all" {
				serverIds = append(serverIds, v.ServerID)
			} else if v.PrivateIpAddress == i {
				serverIds = append(serverIds, v.ServerID)
				break
			}
		}
	}
	users, _ := model.GetUsers("", 0, 0)
	for _, i := range data.Users {
		for _, v := range users {
			if i == "all" {
				userIds = append(userIds, int(v.ID))
			} else if v.Username == i {
				userIds = append(userIds, int(v.ID))
				break
			}
		}
	}
	termUsers, _ := model.GetTermUsers("", 0, 0)

	for _, v := range termUsers {
		if v.Username == data.TermUsers {
			termuserIds = append(termuserIds, int(v.ID))
			break
		}
	}
	permissions := []model.UserPermissions{}
	for _, s := range serverIds {
		for _, u := range userIds {
			for _, t := range termuserIds {
				permissions = append(permissions, model.UserPermissions{
					Group:      data.Group,
					UserID:     u,
					ServerID:   s,
					TermUserID: t,
				})
			}
		}
	}
	code = model.CreatePermissions(&permissions)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// 编辑用户
func EditPermission(c *gin.Context) {

	c.JSON(
		http.StatusOK, gin.H{
			"status":  403,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// 删除用户
func DeletePermission(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code = model.DeletePermission(id)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}
