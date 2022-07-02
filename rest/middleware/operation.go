package middleware

//var bufferPool = sync.Pool{
//	New: func() interface{} {
//		return new(bytes.Buffer)
//	},
//}
//
//type bodyLogWriter struct {
//	gin.ResponseWriter
//	body *bytes.Buffer
//}

//
//func (w bodyLogWriter) Write(b []byte) (int, error) {
//	if n, err := w.body.Write(b); err != nil {
//		global.LOG.Error(err)
//		return n, err
//	}
//	return w.ResponseWriter.Write(b)
//}
//
//func OpRecord() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// Start timer
//		start := time.Now()
//		path := c.Request.URL.Path
//		if 0 <= strings.Index(path, `metrics`) {
//			c.Next()
//			return
//		}
//
//		var userId int
//		claims, _ := utils.GetClaims(c)
//		if nil != claims && claims.BaseClaims.ID != 0 {
//			userId = int(claims.BaseClaims.ID)
//		} else {
//			id, err := strconv.Atoi(c.Request.Header.Get("x-user-id"))
//			if err != nil {
//				userId = 0
//			}
//			userId = id
//		}
//
//		method := c.Request.Method
//		var btReq []byte
//		if http.MethodGet != method {
//			var err error
//			btReq, err = ioutil.ReadAll(c.Request.Body)
//			if err != nil {
//				global.LOG.Error("read body from request error:", err)
//			} else {
//				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(btReq))
//			}
//		} else {
//			query := c.Request.URL.RawQuery
//			query, _ = url.QueryUnescape(query)
//			split := strings.Split(query, "&")
//			m := make(map[string]string)
//			for _, v := range split {
//				kv := strings.Split(v, "=")
//				if len(kv) == 2 {
//					m[kv[0]] = kv[1]
//				}
//			}
//			btReq, _ = json.Marshal(&m)
//		}
//
//		//
//		blw := &bodyLogWriter{body: bufferPool.Get().(*bytes.Buffer), ResponseWriter: c.Writer}
//		c.Writer = blw
//
//		raw := c.Request.URL.RawQuery
//
//		deviceId := c.Request.Header.Get("deviceId")
//		version := c.Request.Header.Get("version")
//		debug := c.Request.Header.Get("debug")
//		clientTime := c.Request.Header.Get("timestamp")
//		channel := c.Request.Header.Get("channel")
//		appTag := c.Request.Header.Get("app_tag")
//
//		// Process request
//		c.Next()
//
//		//var resp string
//		btResp, _ := ioutil.ReadAll(blw.body)
//		//对象回收
//		blw.body.Reset()
//		bufferPool.Put(blw.body)
//
//		if raw != "" {
//			path = path + "?" + raw
//		}
//		// Stop timer
//		end := time.Now()
//
//		latency := fmt.Sprintf("%10v", end.Sub(start))
//		clientIp := c.ClientIP()
//		proto := c.Request.Proto
//		userAgent := c.Request.UserAgent()
//		keys := c.Keys
//		httpCode := c.Writer.Status()
//
//		sErr := ""
//		if nil != c.Errors {
//			sErr = c.Errors.String()
//		}
//		go doLog(userId, deviceId, version, debug, clientTime,
//			channel, appTag, clientIp, method, proto, userAgent,
//			path, latency, string(btReq), string(btResp), httpCode, keys, sErr)
//	}
//}
//func doLog(userId, deviceId, version, debug, clientTime, channel,
//	appTag, clientIp, method, proto, userAgent, path, latency,
//	req, resp, httpCode, keys, sErr interface{}) {
//
//	global.LOG.Info(log.Fields{
//		"user_id":     userId,
//		"device_id":   deviceId,
//		"version":     version,
//		"debug":       debug,
//		"client_time": clientTime,
//		"channel":     channel,
//		"app_tag":     appTag,
//		"client_ip":   clientIp,
//		"method":      method,
//		"proto":       proto,
//		"user_agent":  userAgent,
//		"path":        path,
//		"latency":     latency,
//		"keys":        keys,
//		"req":         req,
//		"resp":        resp,
//		"status_code": httpCode,
//	}, sErr)
//}
//
//func gzipDecode(in []byte) ([]byte, error) {
//	reader, err := gzip.NewReader(bytes.NewReader(in))
//	if err != nil {
//		var out []byte
//		return out, err
//	}
//	defer reader.Close()
//	return ioutil.ReadAll(reader)
//}
