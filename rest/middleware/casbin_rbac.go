package middleware

//// 拦截器 回调服务内逻辑?
//func CasbinHandler(env string) gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		waitUse, _ := utils.GetClaims(ctx)
//		// 获取请求的PATH
//		obj := ctx.Request.URL.Path
//		// 获取请求方法
//		act := ctx.Request.Method
//		// 获取用户的角色
//		sub := waitUse.RoleId
//		e := service.Group.Sys.Casbin()
//		// 判断策略中是否存在
//		success, _ := e.Enforce(sub, obj, act)
//		if env == "dev" || success {
//			ctx.Next()
//		} else {
//			common.RespFailed(ctx, "权限不足")
//			ctx.Abort()
//			return
//		}
//	}
//}
