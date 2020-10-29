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
