package routenames

import (
	"fmt"
)

const (
	Home                  = "home"
	Welcome               = "welcome"
	Dashboard             = "dashboard"
	AdminDashboard        = "admin_dashboard"
	AdminUserAdd          = "admin_user_add"
	AdminUserEdit         = "admin_user_edit"
	AdminUserDelete       = "admin_user_delete"
	About                 = "about"
	// [feature:mail] start
	Contact               = "contact"
	ContactSubmit         = "contact.submit"
	// [feature:mail] end
	Login                 = "login"
	LoginSubmit           = "login.submit"
	Register              = "register"
	RegisterSubmit        = "register.submit"
	ForgotPassword        = "forgot_password"
	ForgotPasswordSubmit  = "forgot_password.submit"
	Logout                = "logout"
	VerifyEmail           = "verify_email"
	ResetPassword         = "reset_password"
	ResetPasswordSubmit   = "reset_password.submit"
	Search                = "search"
	// [feature:tasks] start
	Task                  = "task"
	TaskSubmit            = "task.submit"
	// [feature:tasks] end
	Cache                 = "cache"
	CacheSubmit           = "cache.submit"
	// [feature:files] start
	Files                 = "files"
	FilesSubmit           = "files.submit"
	// [feature:files] end
	// [feature:tasks] start
	AdminTasks            = "admin:tasks"
	// [feature:tasks] end
	ProfileEdit           = "profile.edit"
	ProfileUpdate         = "profile.update"
	ProfileDestroy        = "profile.destroy"
	ProfileAppearance     = "profile.appearance"
	ProfilePassword       = "profile.password"
	ProfileUpdatePassword = "profile.update_password"
	// [feature:payment] start
	Plans                 = "plans"
	PlansSubscribe        = "plans.subscribe"
	Products              = "products"
	ProductsPurchase      = "products.purchase"
	Premium               = "premium"
	Billing               = "billing"
	BillingCancel         = "billing.cancel"
	// [feature:payment] end
	// [feature:chat] start
	ChatRooms             = "chat.rooms"
	ChatRoomCreate        = "chat.rooms.create"
	ChatRoom              = "chat.room"
	ChatWebSocket         = "chat.websocket"
	ChatBanUser           = "chat.ban"
	ChatUnbanUser         = "chat.unban"
	ChatDeleteRoom        = "chat.room.delete"
	// [feature:chat] end
)

func AdminEntityList(entityTypeName string) string {
	return fmt.Sprintf("admin:%s_list", entityTypeName)
}

func AdminEntityAdd(entityTypeName string) string {
	return fmt.Sprintf("admin:%s_add", entityTypeName)
}

func AdminEntityEdit(entityTypeName string) string {
	return fmt.Sprintf("admin:%s_edit", entityTypeName)
}

func AdminEntityDelete(entityTypeName string) string {
	return fmt.Sprintf("admin:%s_delete", entityTypeName)
}

func AdminEntityAddSubmit(entityTypeName string) string {
	return fmt.Sprintf("admin:%s_add.submit", entityTypeName)
}

func AdminEntityEditSubmit(entityTypeName string) string {
	return fmt.Sprintf("admin:%s_edit.submit", entityTypeName)
}

func AdminEntityDeleteSubmit(entityTypeName string) string {
	return fmt.Sprintf("admin:%s_delete.submit", entityTypeName)
}
