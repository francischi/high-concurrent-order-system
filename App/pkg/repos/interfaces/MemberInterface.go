package interfaces

import(
	"golang/pkg/repos/models"
)

type MemberRepo interface {
	
	Create(member models.Member)(err error)

	GetMember(memberId string)(member models.Member , err error)

	GetMemberByEmail(email string)(member models.Member , err error)

	ChangePwd(memberId string , newPwd string)(err error)

	IsEmailExist(email string) bool
}