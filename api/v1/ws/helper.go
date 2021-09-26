package ws

import (
	logrus "cmdb/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"time"
)

func handleError(c *gin.Context, err error) bool {
	if err != nil {
		//logrus.WithError(err).Error("gin context http handler error")
		jsonError(c, err.Error())
		return true
	}
	return false
}
func jsonError(c *gin.Context, msg interface{}) {
	c.AbortWithStatusJSON(200, gin.H{"ok": false, "msg": msg})
}

func wshandleError(ws *websocket.Conn, err error) bool {
	if err != nil {
		logrus.SugarLogger.Error("handler ws ERROR:")
		dt := time.Now().Add(time.Second)
		if err := ws.WriteControl(websocket.CloseMessage, []byte(err.Error()), dt); err != nil {
			logrus.SugarLogger.Error("websocket writes control message failed:")
		}
		return true
	}
	return false
}
