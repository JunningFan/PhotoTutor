package src

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	conn *gorm.DB
)

func Setup(DB_DSN string) {
	var err error
	if DB_DSN == "" {
		conn, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	} else {
		fmt.Println(DB_DSN)
		//	only support postgres connection
		conn, err = gorm.Open(postgres.Open(DB_DSN), &gorm.Config{})
	}
	if err != nil {
		panic(fmt.Sprintf("Fail to connect to database %v", err.Error()))
	}
	//conn = conn.LogMode(true).Set("gorm:auto_preload", true)

	//register objects
	err = conn.AutoMigrate(&User{})
	if err != nil {
		panic(fmt.Sprintf("Fail to migrate database %v", err.Error()))
	}
	//println("Finish set up databse conn dsn ", DB_DSN)
}

type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Username  string     `gorm:"unique"`
	Password  string     `json:"-"`
	Nickname  string
	Signature string
	ImgLoc    string `gorm:"-" json:"img"`
	Following []*User `gorm:"many2many:user_following"`
}

type UserRegisterInput struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
	Nickname string `binding:"required"`
}

type UserLoginInput struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

type UserChangePasswordInput struct {
	Password    string `binding:"required"`
	NewPassword string `binding:"required"`
}

type UserUpdateInput struct {
	Nickname  string `binding:"required"`
	Signature string `binding:"required"`
	Img       uint   `binding:"required"`
}

type UserFollowerInput struct {
	Following uint `binding:"required"`
}

type UserManager struct{}

func NewUserManager() UserManager {
	return UserManager{}
}

func (u *User) SetPassword(password string) error {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(encrypted)
	return nil
}

// CheckPassword Model to check the user password
func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid username / password")
	}
	return nil
}

// Create Create user with given user input
func (um *UserManager) Create(input *UserRegisterInput) (User, error) {
	user := User{Username: input.Username, Nickname: input.Nickname}
	if err := user.SetPassword(input.Password); err != nil {
		return User{}, err
	}
	res := conn.Create(&user)
	if err := res.Error; err != nil {
		return User{}, err
	}
	return user, nil
}

// Login perform the login operation
func (um *UserManager) Login(input *UserLoginInput) (User, error) {
	var user User
	res := conn.Where("Username = ?", input.Username).Find(&user)
	if err := res.Error; err != nil {
		return User{}, fmt.Errorf("user doesn't exist")
	}
	if err := user.CheckPassword(input.Password); err != nil {
		return User{}, err
	}
	return user, nil
}

// GetUser to get the user information by their id in model
func (um *UserManager) GetUser(uid uint) (User, error) {
	return GetUserByID(uid)
}

// GetUserByID helper function to get the user by its id id
func GetUserByID(uid uint) (User, error) {
	var ret User
	res := conn.First(&ret, uid)
	ret.Following = FollowerList(uid)
	return ret, res.Error
}

// Update User information
func (um *UserManager) Update(uid uint, input UserUpdateInput) (User, error) {
	user := User{}
	if img, err := GetImgInfo(input.Img, uid); err != nil {
		return User{}, err
	} else if res := conn.Find(&user, uid); res.Error != nil {
		return User{}, res.Error
	} else {
		user.Nickname = input.Nickname
		user.Signature = input.Signature
		user.ImgLoc = img.Small
		// user.ImgLoc = path.Join(util.ImgSmallPath, imgPath)

		if res := conn.Save(&user); res.Error != nil {
			return User{}, res.Error
		} else {
			return user, nil
		}
	}
}

// NicknameMap Only for mapping nicknames
type NicknameMap struct {
	ID       uint
	Nickname string
}

// ResolveNicknameByIds get a dist of id by an array
// the array must be increment by id
func (um *UserManager) ResolveNicknameByIds(ids []uint) ([]NicknameMap, error) {
	var ret []NicknameMap
	res := conn.Find(&User{}, ids).Order("id ASC").Pluck("nickname", &ret)
	return ret, res.Error
}

//Add user to following list
func (um *UserManager) AddFollower(uid uint, input UserFollowerInput) (User,error) {
	user := User{}
	followID,err := GetUserByID(input.Following)

	if err != nil{
		return User{}, err
	}
	if res := conn.Find(&user, uid); res.Error != nil {
		return User{}, res.Error
	} else if  res := conn.Find(&followID, followID.ID); res.Error != nil {
		return User{}, res.Error
	} else if  uid == followID.ID {
		return User{}, fmt.Errorf("Cannot follow self")
	} else {
		conn.Model(&user).Association("Following").Append(&followID)
		return user, nil
	}
}

//Remove user from following list
func (um *UserManager) Unfollow(uid uint, input UserFollowerInput) (User,error) {
	user := User{}
	followUser := User{}
	followID,err := GetUserByID(input.Following)

	if err != nil{
		return User{}, err
	}
	if res := conn.Find(&user, uid); res.Error != nil {
		return User{}, res.Error
	} else if  res := conn.Find(&followUser, followID); res.Error != nil {
		return User{}, res.Error
	} else {
		conn.Model(&user).Association("Following").Delete(&followUser)
		return user, nil
	}
}

//Get follower list
func FollowerList(uid uint) ([]*User) {
	user := User{}
	var userList []*User
	if res := conn.Find(&user, uid); res.Error != nil {
		return userList
	} else {
		conn.Model(&user).Association("Following").Find(&userList)
		return userList
	}
}