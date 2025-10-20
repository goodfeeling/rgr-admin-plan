import roleService from "@/api/services/roleService";
import userService from "@/api/services/userService";
import { Icon } from "@/components/icon";
import { useDictionaryByTypeWithCache } from "@/hooks/dict";
import RoleSelect from "@/pages/components/role-select/RoleSelect";
import {
	usePasswordResetMutation,
	useRemoveUserMutation,
	useUpdateOrCreateUserMutation,
	useUserManageActions,
	useUserManageCondition,
	useUserQuery,
} from "@/store/userManageStore";
import { Badge } from "@/ui/badge";
import { Button } from "@/ui/button";
import { CardContent, CardHeader } from "@/ui/card";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/ui/form";
import type { TableProps } from "antd";
import { Card, Input, Popconfirm, Select, Table } from "antd";
import { useCallback, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import type { ColumnsType, RoleTree, UserInfo } from "#/entity";
import { buildTree } from "../role/modal";
import PermissionModal, { type UserModalProps } from "./modal";

const defaultUserValue: UserInfo = {
	id: 0,
	email: "",
	user_name: "",
	nick_name: "",
	header_img: "",
	phone: "",
	status: 2,
	created_at: "",
	updated_at: "",
	roles: [],
	current_role: undefined,
};

type SearchFormFieldType = {
	user_name: string;
	status: string;
};

const searchDefaultValue = {
	user_name: "",
	status: undefined,
};

const App: React.FC = () => {
	const { t } = useTranslation();

	const updateOrCreateMutation = useUpdateOrCreateUserMutation();
	const removeMutation = useRemoveUserMutation();
	const passwordResetMutation = usePasswordResetMutation();

	const { data, isLoading } = useUserQuery();
	const condition = useUserManageCondition();
	const { setCondition } = useUserManageActions();

	const { data: statusType } = useDictionaryByTypeWithCache("status");
	const [treeData, setTreeData] = useState<RoleTree[]>([]);

	const searchForm = useForm<SearchFormFieldType>({
		defaultValues: searchDefaultValue,
	});

	const [userModalProps, setUserModalProps] = useState<UserModalProps>({
		formValue: { ...defaultUserValue },
		title: t("table.button.add"),
		treeData: [],
		show: false,
		onOk: async (values: UserInfo): Promise<boolean> => {
			updateOrCreateMutation.mutate(values, {
				onSuccess: () => {
					toast.success(t("table.handle_message.success"));
					setUserModalProps((prev) => ({ ...prev, show: false }));
				},
			});
			return true;
		},
		onCancel: () => {
			setUserModalProps((prev) => ({ ...prev, show: false }));
		},
	});

	const handleTableChange: TableProps<UserInfo>["onChange"] = (pagination, filters, sorter) => {
		setCondition({
			...condition,
			pagination,
			filters,
			sortOrder: Array.isArray(sorter) ? undefined : condition.sortOrder,
			sortField: Array.isArray(sorter) ? undefined : condition.sortField,
		});
	};

	const handleDelete = async (id: number) => {
		removeMutation.mutate(id, {
			onSuccess: () => {
				toast.success(t("table.handle_message.success"));
			},
			onError: () => {
				toast.error(t("table.handle_message.error"));
			},
		});
	};

	const onResetPassword = (id: number) => {
		passwordResetMutation.mutate(id, {
			onSuccess: () => {
				toast.success(t("table.handle_message.success"));
			},
			onError: () => {
				toast.error(t("table.handle_message.error"));
			},
		});
	};

	const getTreeData = useCallback(async () => {
		const response = await roleService.getRoles();
		const treeData = buildTree(response);
		setTreeData(treeData);
		setUserModalProps((prev) => ({
			...prev,
			treeData,
		}));
	}, []);

	useEffect(() => {
		getTreeData();
	}, [getTreeData]);

	const onReset = () => {
		setCondition({
			...condition,
			searchParams: searchDefaultValue,
			pagination: {
				...condition.pagination,
				current: 1,
			},
		});
		searchForm.reset();
	};

	const onSearch = () => {
		const values = searchForm.getValues();
		setCondition({
			...condition,
			searchParams: {
				...values,
			},
			pagination: {
				...condition.pagination,
				current: 1,
			},
		});
	};

	const onCreate = () => {
		setUserModalProps((prev) => ({
			...prev,
			show: true,
			...defaultUserValue,
			title: t("table.button.add"),
			formValue: { ...defaultUserValue },
		}));
	};

	const onEdit = (formValue: UserInfo) => {
		setUserModalProps((prev) => ({
			...prev,
			show: true,
			title: t("table.button.edit"),
			formValue,
		}));
	};

	const columns: ColumnsType<UserInfo> = [
		{
			title: "ID",
			dataIndex: "id",
			sorter: true,
			width: "5%",
		},
		{
			title: t("table.columns.user.user_name"),
			dataIndex: "user_name",
			width: 300,
			render: (_, record) => {
				return (
					<div className="flex">
						<img alt="" src={record.header_img} className="h-10 w-10 rounded-full" />
						<div className="ml-2 flex flex-col">
							<span className="text-sm">{record.user_name}</span>
							<span className="text-xs text-text-secondary">{record.email}</span>
						</div>
					</div>
				);
			},
		},

		{
			title: t("table.columns.user.nick_name"),
			dataIndex: "nick_name",
		},
		{
			title: t("table.columns.user.phone"),
			dataIndex: "phone",
		},
		{
			title: t("table.columns.user.status"),
			dataIndex: "status",
			align: "center",
			width: 120,
			render: (status) => {
				const statusItem = statusType?.find((item) => Number(item.value) === status);

				return <Badge variant={status === 1 ? "success" : "error"}>{statusItem?.label}</Badge>;
			},
		},
		{
			title: t("table.columns.user.roles"),
			dataIndex: "roles",
			width: 350,
			render: (roles, record) => {
				return (
					<RoleSelect
						roles={roles}
						treeData={treeData}
						recordKey={`list_${record.id}`}
						onChange={async (values) => {
							try {
								await userService.bindRole(record.id, values);
								toast.success(t("table.handle_message.success"));
							} catch (error) {
								console.error("更新失败:", error);
							}
						}}
					/>
				);
			},
		},
		{
			title: t("table.columns.common.created_at"),
			dataIndex: "created_at",
			key: "created_at",
		},
		{
			title: t("table.columns.common.updated_at"),
			dataIndex: "updated_at",
			key: "updated_at",
		},
		{
			title: t("table.columns.common.operation"),
			key: "operation",
			align: "center",
			width: 220,
			fixed: "right",
			render: (_, record) => (
				<div className="grid grid-cols-3 gap-1 text-gray-500">
					<Button variant="link" size="icon" onClick={() => onEdit(record)} className="whitespace-nowrap justify-start">
						<div className="flex items-center">
							<Icon icon="solar:pen-bold-duotone" size={18} />
							<span className="ml-1"> {t("table.button.edit")}</span>
						</div>
					</Button>
					<Popconfirm
						title={t("table.handle_message.reset_prompt")}
						description={t("table.handle_message.confirm_reset_password")}
						onConfirm={() => onResetPassword(record.id)}
						okText={t("table.button.yes")}
						cancelText={t("table.button.no")}
					>
						<Button variant="link" size="icon" className="whitespace-nowrap justify-start">
							<div className="flex items-center">
								<Icon icon="solar:restart-line-duotone" size={18} color="orange" />
								<span className="ml-1 text-orange-500">{t("table.button.reset_password")}</span>
							</div>
						</Button>
					</Popconfirm>

					<Popconfirm
						title={t("table.handle_message.delete_prompt")}
						description={t("table.handle_message.confirm_delete")}
						onConfirm={() => handleDelete(record.id)}
						okText={t("table.button.yes")}
						cancelText={t("table.button.no")}
					>
						<Button variant="link" size="icon" className="whitespace-nowrap justify-start">
							<div className="flex items-center">
								<Icon icon="mingcute:delete-2-fill" size={18} color="red" />
								<span className="ml-1 text-red-500">{t("table.button.delete")}</span>
							</div>
						</Button>
					</Popconfirm>
				</div>
			),
		},
	];

	return (
		<div className="flex flex-col gap-4">
			<Card>
				<CardContent>
					<Form {...searchForm}>
						<div className="flex items-center gap-4">
							<FormField
								control={searchForm.control}
								name="user_name"
								render={({ field }) => (
									<FormItem>
										<FormLabel>{t("table.columns.user.user_name")}</FormLabel>
										<FormControl>
											<Input {...field} />
										</FormControl>
									</FormItem>
								)}
							/>
							<FormField
								control={searchForm.control}
								name="status"
								render={({ field }) => (
									<FormItem>
										<FormLabel>{t("table.columns.user.status")}</FormLabel>
										<Select
											onChange={(value: string) => {
												field.onChange(value);
											}}
											value={field.value}
											options={statusType}
											placeholder={`${t("table.handle_message.select")}${t("table.columns.user.status")}`}
										/>
									</FormItem>
								)}
							/>
							<div className="flex ml-auto">
								<Button variant="outline" onClick={() => onReset()}>
									<Icon icon="solar:restart-line-duotone" size={18} />
									{t("table.button.reset")}
								</Button>
								<Button variant="default" className="ml-4 text-white" onClick={() => onSearch()}>
									<Icon icon="solar:rounded-magnifer-linear" size={18} />
									{t("table.button.search")}
								</Button>
							</div>
						</div>
					</Form>
				</CardContent>
			</Card>
			<Card title={t("sys.menu.system.user")} size="small">
				<CardHeader>
					<div className="flex items-center justify-between">
						<Button onClick={() => onCreate()} className="text-white">
							<Icon icon="solar:add-circle-outline" size={18} />
							{t("table.button.add")}
						</Button>
					</div>
				</CardHeader>

				<CardContent>
					<Table<UserInfo>
						rowKey={(record) => record.id}
						scroll={{ x: "max-content" }}
						columns={columns}
						pagination={{
							current: data?.page || 1,
							pageSize: data?.page_size || 10,
							total: data?.total || 0,
							showTotal: (total) => `${t("table.page.total")} ${total} ${t("table.page.items")}`,
							showSizeChanger: true,
							pageSizeOptions: ["10", "20", "50", "100"],
						}}
						dataSource={data?.list}
						loading={isLoading}
						onChange={handleTableChange}
					/>
				</CardContent>
				<PermissionModal {...userModalProps} />
			</Card>
		</div>
	);
};

export default App;
