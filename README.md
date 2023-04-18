# POST BACKEND USING GOLANG

## LANGUAGE: Golang
## DATABASE: MongoDB

## Các chức năng chính của dự án.
**1. API Chung cuả dự án**
``` Go
v1 := r.Group("api/v1")
	{
		NewUserApi(v1.Group("user"), userHandler)
		NewPostApi(v1.Group("post"), postHandler)
	}
	r.Use(static.Serve("/upload", static.LocalFile("./upload", true)))
```
**2. API cho USER**

``` Go
func NewUserApi(group *gin.RouterGroup, handler handler.UserHandler) *UserApi {
	s := &UserApi{
		RouterGroup: group,
		UserHandler: handler,
	}

	s.POST("/signup", s.UserHandler.Signup)
	s.POST("/login", s.UserHandler.Login)
	s.POST("/refresh", s.UserHandler.RefreshToken)
	s.Use(middleware.JwtAuthMiddleware(handler.Env.AccessTokenSecret))
	s.GET("/profile", s.UserHandler.Profile)
	s.PUT("/profile/change", s.UserHandler.ChangeProfile)
	s.PUT("/profile/change-password", s.UserHandler.ChangePassword)
	return s
}
```
API cho user Thực hiện các chức năng:
- "/signup": Đăng ký 1 tài khoản mới
- "/login": Đăng nhập 
- "/refresh": Refresh token khi token hết hạn.
- "/profile": Lấy thông tin chi tiết của user.
- "/profile/change": Update profile cho user.
- "/profile/change-password": Thay đổi password cho user.

**3. API cho POST**

``` Go
func NewPostApi(group *gin.RouterGroup, handler handler.PostHandler) *PostApi {
	s := &PostApi{
		RouterGroup: group,
		PostHandler: handler,
	}

	s.Use(middleware.JwtAuthMiddleware(handler.Env.AccessTokenSecret))
	s.GET("/all", s.PostHandler.ListAllPost)
	s.GET("/detail", s.PostHandler.GetPostByID)
	s.POST("/create", s.PostHandler.Create)
	s.PUT("/update", s.UpdatePost)
	s.DELETE("/delete", s.DeletePost)
	s.GET("/all/user", s.ListPostByUser)
	return s
}
```
API cho post thực hiện các chức năng:
- "/all" : lấy ra tất cả post theo skip và limit.
- "/detail" : Lấy ra thông tin chi tiết về 1 post theo id.
- "/create": Tạo mới 1 post.
- "/update": Update thông tin cho 1 post.
- "/delete": Xóa 1 post theo id của post.
- "/all/user": Lấy ra tất cả các post của 1 user.

## Data Của dự án
**1. USER**
``` GO
type UserReader struct {
	ID          string                `form:"_id,omitempty" json:"_id" bson:"_id,omitempty"`
	UserName    string                `form:"username" json:"username" bson:"username" binding:"required"`
	Email       string                `form:"email" json:"email" bson:"email" binding:"required,email"`
	FirstName   string                `form:"firstname" json:"firstname,omitempty" bson:"firstname,omitempty"`
	SurName     string                `form:"surname" json:"surname,omitempty" bson:"surname,omitempty"`
	MobilePhone string                `form:"phone" json:"phone,omitempty" bson:"phone,omitempty"`
	Address1    string                `form:"address1" json:"address1,omitempty" bson:"address1,omitempty"`
	Address2    string                `form:"address2" json:"address2,omitempty" bson:"address2,omitempty"`
	Education   string                `form:"education" json:"education,omitempty" bson:"education,omitempty"`
	Country     string                `form:"country" json:"country,omitempty" bson:"country,omitempty"`
	State       string                `form:"state" json:"state,omitempty" bson:"state,omitempty"`
	AvatarURL   string                `form:"avatar_url" json:"avatar_url,omitempty" bson:"avatar_url,omitempty"`
	AvatarFile  *multipart.FileHeader `form:"avatar_file" json:"-" bson:"avatar_file,omitempty"`
}
```
**2. POST **

``` GO
type Post struct {
	ID        string     `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string     `json:"title" bson:"title" binding:"required"`
	Content   string     `json:"content" bson:"content" binding:"required"`
	ImageURL  string     `json:"image,omitempty" bson:"image,omitempty"`
	Owner     UserReader `json:"owner" bson:"owner"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at,omitempty," bson:"updated_at,omitempty"`
	IsDeleted bool       `json:"is_deleted" bson:"is_deleted"`
}
```


## Framework và thư viện sử dụng.
**1. Framework**
- [https://github.com/gin-gonic/gin](https://github.com/gin-gonic/gin) Tạo router, binding data and response JSON.

**2. Thư viện**
- [https://github.com/mongodb/mongo-go-driver](https://github.com/mongodb/mongo-go-driver) cho kết nối golang and mongodb.
- [https://github.com/spf13/viper](https://github.com/spf13/viper) đọc file .env để lấy các biến môi trường.
- [github.com/google/uuid](github.com/google/uuid) Tạo uuid cho các đối tượng.
- [golang.org/x/crypto](golang.org/x/crypto) Generate hashPassword and compare hashPassword and Password.
- [github.com/golang-jwt/jwt/v4](github.com/golang-jwt/jwt/v4) Tạo jwt token cho authentication.

