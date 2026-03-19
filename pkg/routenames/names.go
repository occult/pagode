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
	Contact               = "contact"
	ContactSubmit         = "contact.submit"
	Login           = "login"
	Register        = "register"
	Logout          = "logout"
	Search                = "search"
	Task                  = "task"
	TaskSubmit            = "task.submit"
	Cache                 = "cache"
	CacheSubmit           = "cache.submit"
	Files                 = "files"
	FilesSubmit           = "files.submit"
	AdminTasks            = "admin:tasks"
	ProfileEdit           = "profile.edit"
	ProfileUpdate         = "profile.update"
	ProfileDestroy        = "profile.destroy"
	ProfileAppearance     = "profile.appearance"
	ProfilePassword       = "profile.password"
	Plans                 = "plans"
	PlansSubscribe        = "plans.subscribe"
	Products              = "products"
	ProductsPurchase      = "products.purchase"
	Premium               = "premium"
	Billing               = "billing"
	BillingCancel         = "billing.cancel"
	ChatRooms             = "chat.rooms"
	ChatRoomCreate        = "chat.rooms.create"
	ChatRoom              = "chat.room"
	ChatWebSocket         = "chat.websocket"
	ChatBanUser           = "chat.ban"
	ChatUnbanUser         = "chat.unban"
	ChatDeleteRoom        = "chat.room.delete"
	CasdoorCallback       = "casdoor.callback"
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
