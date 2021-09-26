package ws

import (
	"cmdb/model"
	"cmdb/terminal"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsSsh(c *gin.Context) {
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if handleError(c, err) {
		return
	}
	defer wsConn.Close()
	if handleError(c, err) {
		return
	}
	cols, err := strconv.Atoi(c.DefaultQuery("cols", "2000"))
	if wshandleError(wsConn, err) {
		return
	}
	rows, err := strconv.Atoi(c.DefaultQuery("rows", "500"))
	if wshandleError(wsConn, err) {
		return
	}
	user, _ := c.Get("wsusername")
	//通过对应的username以及选择的host，获取ssh认证后的对象
	id, _ := strconv.Atoi(c.Param("id"))
	mc, code := model.Permissions(id, user.(string))
	if code != 200 {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  1002,
				"message": mc,
			},
		)
		c.Abort()
	}
	client, err := terminal.NewSshClient(mc)
	if wshandleError(wsConn, err) {
		return
	}
	defer client.Close()
	//startTime := time.Now()

	sws, err := model.NewLogicSshWsSession(cols, rows, true, client, wsConn)
	if wshandleError(wsConn, err) {
		return
	}
	defer sws.Close()

	quitChan := make(chan bool, 3)
	sws.Start(quitChan)
	go sws.Wait(quitChan)

	<-quitChan
	//保存日志

	////write logs
	xtermLog := model.SshLog{
		Log: sws.LogString(),
		//ClientIp:  ,
		//StartTime:   startTime,
		//User: user.(string),
	}
	model.SshlogCreate(xtermLog)
	if wshandleError(wsConn, err) {
		return
	}
}
