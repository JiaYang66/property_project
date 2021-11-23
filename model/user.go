package model

import (
	"relocate/util/logging"
	"relocate/util/times"
)

type UserStatus int

const (
	NotVerifiedUser UserStatus = iota
	RejectedUser
	AutomaticMatchingUser
	VerifiedUser
)

var AllUserStatus map[UserStatus]string

func (us UserStatus) String() string {
	switch us {
	case NotVerifiedUser: // 0
		return "未核验"
	case RejectedUser: // 1
		return "已驳回"
	case AutomaticMatchingUser: // 2
		return "自动匹配"
	case VerifiedUser:
		return "已核验" // 3
	default:
		return "unknown"
	}
}

type IdNumberType int

const (
	MainlandId IdNumberType = iota
	HongKongId
	HuZhaoId
)

func (it IdNumberType) String() string {
	switch it {
	case MainlandId:
		return "大陆"
	case HongKongId:
		return "港澳证件"
	case HuZhaoId:
		return "护照"
	default:
		return "unknown"
	}
}

// 用户
// 通过手机号码和密码注册登录，注册成功后状态默认为未核验
// 用户进行身份核验，填写真实姓名，身份证号码，正反图片，若完全和合同表一致，自动匹配关联合同表，状态为自动匹配
// 管理员可以收到通过核验或驳回
// 只要核验过后的状态才可以进行申报操作
type User struct {
	PhoneNumber   string          `json:"phone_number" gorm:"primary_key;comment:'手机号码-主键'"`
	Password      string          `json:"-" gorm:"not null;comment:'密码'"`
	Salt          string          `json:"-" gorm:"not null;comment:'混淆盐'"`
	Name          string          `json:"name" gorm:"null;comment:'真实姓名'"`
	IdNumberType  IdNumberType    `json:"id_number_type" gorm:"null;comment:'身份证类型'"`
	IdNumber      string          `json:"id_number" gorm:"null;comment:'身份证号码'"`
	PositiveImage string          `json:"positive_image" gorm:"null;comment:'身份证正面图片'"`
	NegativeImage string          `json:"negative_image" gorm:"null;comment:'身份证反面图片'"`
	UserStatus    UserStatus      `json:"user_status" gorm:"not null;comment:'用户状态'"`
	CreatedAt     times.JsonTime  `json:"created_at" gorm:"not null;comment:'创建时间'"`
	UpdatedAt     times.JsonTime  `json:"-" gorm:"comment:'更新时间'"`
	DeletedAt     *times.JsonTime `json:"-" gorm:"comment:'删除时间'"`
}

// 表名
func (u User) TableName() string {
	return "user"
}

func initUser() {
	if !db.HasTable(&User{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&User{}).Error; err != nil {
			panic(err)
		}
		user := User{
			PhoneNumber: "13600479145",
			Password:    "e7b3e907df88943aff31c0007feecdb6", //123456
			UserStatus:  NotVerifiedUser,
			Salt:        "ABCDEF",
		}
		user.Create()
	}
	AllUserStatus = make(map[UserStatus]string, 0)
	AllUserStatus[NotVerifiedUser] = NotVerifiedUser.String()
	AllUserStatus[RejectedUser] = RejectedUser.String()
	AllUserStatus[AutomaticMatchingUser] = AutomaticMatchingUser.String()
	AllUserStatus[VerifiedUser] = VerifiedUser.String()
}

func (u *User) Create() error {
	return db.Create(&u).Error
}

// 通过用户获取信息
func GetUserInfo(userName string) (*User, error) {
	var user User
	if err := db.Where("phone_number = ?", userName).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func AutoMatchUser(phoneNumber, name, idCard string) bool {
	query := "(phone_number1 LIKE ? OR phone_number2 LIKE ?) AND peoples LIKE ? AND card_number LIKE ?"
	if err := db.Where(query, "%"+phoneNumber+"%", "%"+phoneNumber+"%", "%"+name+"%", "%"+idCard+"%").First(&Contract{}).Error; err != nil {
		return false
	}
	return true
}

func (u *User) UpdateUser() error {
	sql := db.Model(u).Where("phone_number = ?", u.PhoneNumber).Updates(&u)
	rowsAffected := sql.RowsAffected
	logging.Infof("更新影响的记录数%d", rowsAffected)
	logging.Infoln(sql.Error)
	return sql.Error
}

func QueryUserByPhone(phone string) (*User, error) {
	var user User
	if err := db.Where("phone_number = ?", phone, AutomaticMatchingUser, VerifiedUser).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUserByPhone(phone string) (*User, error) {
	var user User
	if err := db.Where("phone_number = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func ExistsUserByIdNumber(idNumber string) bool {
	var user User
	db.Where("id_number = ?", idNumber).First(&user)
	return user.PhoneNumber != ""
}

func ExistsUserByPhoneNumber(phoneNumber string) bool {
	var user User
	db.Where("phone_number = ?", phoneNumber).First(&user)
	return user.PhoneNumber != ""
}

func UserQueryByStatusOrFilterName(filterStatus string, filterName string, page uint, pageSize uint) (data *PaginationQ, err error) {
	q := PaginationQ{
		Data:     &[]User{},
		PageSize: pageSize,
		Page:     page,
	}

	filterName = "%" + filterName + "%"

	// 状态，姓名，身份证或电话号码
	//select * from user where phone_number like '%测试1%' or id_number like '%测试1%' or name like '%测试1%' and user_status = 1
	result := db.Model(&User{})
	if filterStatus != "" {
		result = result.Where("user_status = ?", filterStatus)
	}
	if filterName != "" {
		result = result.Where("name LIKE ? or id_number LIKE ? or phone_number LIKE ?", filterName, filterName, filterName)
	}
	data, err = q.SearchAll(result)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func (user *User) Update() error {
	UpUser := User{
		UserStatus: user.UserStatus,
	}

	return db.Model(&user).Updates(UpUser).Error
}
