export enum BasicStatus {
	DISABLE = 2,
	ENABLE = 1,
}

export enum ResultEnum {
	SUCCESS = 0,
	ERROR = -1,
	TIMEOUT = 401,
	TOKEN_ERROR = 3001,
}

export enum StorageEnum {
	UserInfo = "userInfo",
	UserToken = "userToken",
	Settings = "settings",
	Nav = "nav",
	I18N = "i18nextLng",
	Menu = "menu",
	Role = "role",
	STSToken = "stsToken",
	SysConfig = "sysConfig",
	UserStore = "userStore",
	DeviceId = "deviceId",
}

export enum ThemeMode {
	Light = "light",
	Dark = "dark",
}

export enum ThemeLayout {
	Vertical = "vertical",
	Horizontal = "horizontal",
	Mini = "mini",
}

export enum ThemeColorPresets {
	Default = "default",
	Cyan = "cyan",
	Purple = "purple",
	Blue = "blue",
	Orange = "orange",
	Red = "red",
}

export enum LocalEnum {
	en_US = "en_US",
	zh_CN = "zh_CN",
}

export enum MultiTabOperation {
	FULLSCREEN = "fullscreen",
	REFRESH = "refresh",
	CLOSE = "close",
	CLOSEOTHERS = "closeOthers",
	CLOSEALL = "closeAll",
	CLOSELEFT = "closeLeft",
	CLOSERIGHT = "closeRight",
}

export enum PermissionType {
	CATALOGUE = 0,
	MENU = 1,
	BUTTON = 2,
}

export enum HtmlDataAttribute {
	ColorPalette = "data-color-palette",
	ThemeMode = "data-theme-mode",
}

export enum SortDirection {
	SortASC = "asc",
	SortDesc = "desc",
}

export enum PagePath {
	Login = "/auth/login",
}

export const MessageType = {
	"You have been logged in from another device": {
		title: "账号已在别处登录",
		content: "您的账号在另一台设备上登录，当前会话已被终止。请重新登录。",
	},
	"Token has been replaced": {
		title: "账号已在别处登录",
		content: "您的账号在另一台设备上登录，当前会话已被终止。请重新登录。",
	},
	"refresh token has been revoked": {
		title: "令牌失效提示",
		content: "令牌已过期，请重新登录",
	},
	"Invalid token": {
		title: "令牌失效提示",
		content: "令牌已过期，请重新登录",
	},
	"Token expired": {
		title: "令牌失效提示",
		content: "令牌已过期，请重新登录",
	},
} as const;

export type MessageKey = keyof typeof MessageType;
