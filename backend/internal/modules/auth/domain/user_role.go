package domain

type UserRole string

const (
	RoleRenter    UserRole = "RENTER"
	RoleOwner     UserRole = "OWNER"
	RoleOperator  UserRole = "OPERATOR"
	RoleLogistics UserRole = "LOGISTICS"
	RoleAdmin     UserRole = "ADMIN"
)
