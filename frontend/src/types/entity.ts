import type { GetProp, TableProps } from "antd";
import type { SorterResult } from "antd/es/table/interface";

export type ColumnsType<T extends object = object> = TableProps<T>["columns"];
export type TablePaginationConfig = Exclude<GetProp<TableProps, "pagination">, boolean>;
export interface TableParams {
	pagination?: TablePaginationConfig;
	sortField?: SorterResult<any>["field"];
	sortOrder?: SorterResult<any>["order"];
	filters?: Parameters<GetProp<TableProps, "onChange">>[1];
	searchParams?: { [key: string]: any };
}
export interface PageList<T> {
	list: T[];
	total: number;
	page: number;
	page_size: number;
	total_page: number;
	filters?: DataFilters;
}

type StringMapOfStringArray = { [key: string]: string[] };
interface DataRangeFilter {
	field: string;
	start: string;
	end: string;
}
interface DataFilters {
	likeFilters: StringMapOfStringArray;
	matches: StringMapOfStringArray;
	dataRanges: DataRangeFilter[] | null;
	sortBy: string[];
	sortDirection: string;
	page: number;
	pageSize: number;
}

export interface UserToken {
	accessToken?: string;
	refreshToken?: string;
	expirationAccessDateTime?: string;
	expirationRefreshDateTime?: string;
}

export interface UserInfo {
	id: number;
	email: string;
	user_name: string;
	nick_name: string;
	header_img: string;
	phone: string;
	status?: number;
	created_at: string;
	updated_at: string;
	current_role?: Role;
	roles?: Role[];
}

export interface UpdateUser {
	email: string;
	user_name: string;
	nick_name: string;
	header_img: string;
	phone: string;
	status?: number;
}

export interface Role {
	id: number;
	parent_id: number;
	name: string;
	label: string;
	status: number;
	order?: number;
	description?: string;
	created_at: string;
	updated_at: string;
	default_router: string;
	children?: Role[];
	path: number[];
}

export interface RoleTree {
	value: string;
	title: string;
	key: string;
	children: RoleTree[];
	path: number[];
}

export interface FileInfo {
	id: number;
	file_name: string;
	file_path: string;
	file_md5: string;
	file_url: string;
	storage_engine: string;
	file_origin_name: string;
	created_at: string;
	updated_at: string;
}

export interface Operation {
	id: number;
	name: string;
	path: string;
	ip: string;
	method: string;
	latency: number;
	agent: string;
	error_message: string;
	body: string;
	status: number;
	created_at: string;
	updated_at: string;
}

export interface Api {
	id: number;
	path: string;
	api_group: string;
	method: string;
	description: string;
	created_at: string;
	updated_at: string;
}

export interface ApiGroup {
	api_group: { [key: string]: any };
	groups: string[];
}

export interface ApiGroupItem {
	title: string;
	key: string;
	disableCheckbox: boolean;
	children: ApiGroupItem[];
}

export interface Dictionary {
	id: number;
	name: string;
	type: string;
	status: number;
	desc: string;
	is_generate_file: number;
	created_at: string;
	updated_at: string;
	details: DictionaryDetail[];
}

export interface DictionaryDetail {
	id: number;
	label: string;
	value: string;
	extend: string;
	status: number;
	sort: number;
	type: string;
	sys_dictionary_Id: number | null;
	created_at: string;
	updated_at: string;
}

export interface Menu {
	id: number;
	menu_level: number;
	parent_id: number;
	name: string;
	path: string;
	hidden: boolean;
	component: string;
	sort: number;
	keep_alive: number;
	title: string;
	icon: string;
	menu_group_id: number;
	created_at: string;
	updated_at: string;
	level: number[];
	children?: Menu[];
	menu_btns?: MenuBtn[];
	menu_parameters?: MenuParameter[];
	btn_slice?: string[];
}

export interface MenuTree {
	value: string;
	title: string;
	key: string;
	origin?: Menu;
	children: MenuTree[];
	path: number[];
	disabled?: boolean;
}

export interface roleSetting {
	role_menus: { [key: string]: number[] };
	role_apis: string[];
	role_btns: { [key: string]: number[] };
}

export interface MenuGroup {
	id: number;
	name: string;
	path: string;
	status: number;
	sort: number;
	created_at: string;
	updated_at: string;
}

export interface MenuTreeUserGroup {
	id?: number;
	name?: string;
	path?: string;
	items: Menu[];
}

export interface MenuBtn {
	id?: number;
	name: string;
	desc: string;
	sys_base_menu_id: number;
}

export interface MenuParameter {
	id?: number;
	type: string;
	key: string;
	value: string;
	sys_base_menu_id: number;
}

export type PasswordEditReq = {
	oldPassword: string;
	newPassword: string;
	confirmPassword: string;
};

export type STSToken = {
	access_key_id: string;
	access_key_secret: string;
	security_token: string;
	expiration: string;
	bucket_name: string;
	region: string;
	refresh_token: string;
};

export type ScheduledTask = {
	id: number;
	task_name: string;
	task_description: string;
	cron_expression: string;
	task_type: string;
	task_params: { [key: string]: string };
	status: number;
	exec_type: string;
	last_execute_time: string;
	next_execute_time: string;
	created_at: string;
	updated_at: string;
};

export type TaskExecutionLog = {
	id: number;
	task_id: number;
	execute_time: string;
	execute_result: number;
	execute_duration: number;
	error_message: string;
	created_at: string;
	updated_at: string;
};

export type WebSocketMessage<T> = {
	type: string;
	data: T;
};
export type ConfigResponse = {
	data: { [key: string]: { [key: string]: any } };
};
