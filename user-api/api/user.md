### 1. "获取用户信息"

1. route definition

- Url: /user/login
- Method: POST
- Request: `LoginReq`
- Response: `LoginResp`

2. request definition



```golang
type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
```


3. response definition



```golang
type LoginResp struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Token string `json:"token"`
	ExpireAt string `json:"expireAt"`
}
```

