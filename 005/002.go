package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

var account =make(map[string]string)

func main()  {
	b:=make([]byte,1024)
	f,err:=os.OpenFile("005/users.data",os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	defer f.Close()
	if err !=nil{
		fmt.Println("err:",err)
	}

	_,err=f.Read(b)
	index:=bytes.IndexByte(b,0)
	if err !=nil{
		fmt.Println("err:",err)
	}
	a:=b[0:index]
	err=json.Unmarshal(a,&account)
	if err !=nil{
		fmt.Println("err:",err)
	}


	r:=gin.Default()

	r.POST("/login",login)

	r.POST("/signup",func (c *gin.Context){
		username:=c.PostForm("username")
		password:=c.PostForm("password")
		passwordagain:=c.PostForm("passwordagain")
		_,exist:=account[username]
		if exist==true{
			c.JSON(http.StatusOK,"此账号已注册")
		}else if password!=passwordagain{
			c.JSON(http.StatusOK,"密码不一致")
		}else {
			account[username]=password
			os.Truncate("/005/users.data",0)
			k,_:=json.Marshal(account)
			_,err=f.Write(k)
			if err!=nil{
				c.JSON(http.StatusOK,"写入错误")
			}
			c.SetCookie("login_cookie", username, 3600, "/", "", false, true)
			c.JSON(http.StatusOK,"注册成功")
		}
	})

	r.GET("/source",cookie,source)

	r.Run(":6666")
}



func login(c *gin.Context){
	username:=c.PostForm("username")
	password:=c.PostForm("password")
	_,exist:=account[username]
	if exist==false {
		c.JSON(http.StatusOK,"无此账号")
	}
	if password==account[username] {
		c.SetCookie("login_cookie", username, 3600, "/", "", false, true)
		c.JSON(http.StatusOK,"登陆成功")
	}
}


func cookie(c *gin.Context){
	ck,err :=c.Cookie("login_cookie")
	if err!=nil{
		fmt.Println(err)
		c.JSON(403,"未登录")
		c.Abort()
	}else {
		c.Set("cookie",ck)
		c.Next()
		v,_:=c.Get("next")
		fmt.Println(v)
	}

}

func source(c *gin.Context){
	cookie,_:=c.Get("cookie")
	c.JSON(http.StatusOK,"哈啰，我是资源，你可以浏览我")
	s:=cookie.(string)
	c.JSON(http.StatusOK,s)
	fmt.Println(s)
}
