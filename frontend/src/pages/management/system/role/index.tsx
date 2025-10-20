import { Icon } from "@/components/icon";
import { useDictionaryByTypeWithCache } from "@/hooks/dict";
import { useRemoveRoleMutation, useRoleQuery, useUpdateOrCreateRoleMutation } from "@/store/roleManageStore";
import { Badge } from "@/ui/badge";
import { Button } from "@/ui/button";
import { CardContent, CardHeader } from "@/ui/card";
import type { TableProps } from "antd";
import { Card, Popconfirm, Table } from "antd";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import type { Role } from "#/entity";
import RoleModal, { type RoleModalProps } from "./modal";
import SettingModal, { type SettingModalProps } from "./setting/index";

type ColumnsType<T extends object = object> = TableProps<T>["columns"];

const defaultValue: Role = {
	id: 0,
	parent_id: 0,
	name: "",
	label: "",
	order: 0,
	description: "",
	status: 2,
	created_at: "",
	updated_at: "",
	default_router: "",
	children: [],
	path: [],
};

const App: React.FC = () => {
	const { t } = useTranslation();

	const { data: statusType } = useDictionaryByTypeWithCache("status");
	const updateOrCreateMutation = useUpdateOrCreateRoleMutation();
	const removeMutation = useRemoveRoleMutation();
	const { data, isLoading } = useRoleQuery();
	const [expandedKeys, setExpandedKeys] = useState<number[]>([]);

	const [settingModalPros, setSettingModalProps] = useState<SettingModalProps>({
		id: 0,
		roleData: { ...defaultValue },
		title: t("table.button.add"),
		show: false,
		onCancel: () => {
			setSettingModalProps((prev) => ({ ...prev, show: false }));
		},
	});

	const [roleModalProps, setRoleModalProps] = useState<RoleModalProps>({
		formValue: { ...defaultValue },
		title: t("table.button.add"),
		show: false,
		treeRawData: [],
		onOk: async (values: Role): Promise<boolean> => {
			updateOrCreateMutation.mutate(values, {
				onSuccess: () => {
					toast.success("success!");
					setRoleModalProps((prev) => ({ ...prev, show: false }));
				},
			});
			return true;
		},
		onCancel: () => {
			setRoleModalProps((prev) => ({ ...prev, show: false }));
		},
	});

	useEffect(() => {
		if (data) {
			setRoleModalProps((prev) => ({
				...prev,
				treeRawData: data,
			}));
		}
	}, [data]);

	const onCreate = (formValue: Role | undefined) => {
		const setValue = defaultValue;
		if (formValue !== undefined) {
			setValue.parent_id = formValue.id;
		}
		setRoleModalProps((prev) => ({
			...prev,
			show: true,
			...setValue,
			title: t("table.button.add"),
			formValue: { ...setValue },
		}));
	};

	const onEdit = (formValue: Role) => {
		setRoleModalProps((prev) => ({
			...prev,
			show: true,
			title: t("table.button.edit"),
			formValue,
		}));
	};

	const onSetting = (value: Role) => {
		setSettingModalProps((prev) => ({
			...prev,
			show: true,
			title: t("table.button.role_setting"),
			id: value.id,
			roleData: value,
		}));
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

	const handleExpand = (expanded: boolean, record: Role) => {
		const keys = expanded ? [...expandedKeys, record.id] : expandedKeys.filter((key) => key !== record.id);
		setExpandedKeys(keys);
	};

	const columns: ColumnsType<Role> = [
		{
			title: t("table.columns.role.role_id"),
			dataIndex: "expand",
			render: (_, record) => {
				const level = record.path.length;
				return record.children?.length ? (
					<Button
						onClick={() => handleExpand(!expandedKeys.includes(record.id), record)}
						variant="ghost"
						size="icon"
						style={{
							marginLeft: record.parent_id !== 0 ? `${level * 20}px` : "",
						}}
					>
						{expandedKeys.includes(record.id) ? "▼" : "▶"}
						<span>{record.id}</span>
					</Button>
				) : (
					<span style={{ marginLeft: `${level * 20}px` }}>{record.id}</span>
				);
			},

			width: 90,
		},
		{
			title: t("table.columns.role.name"),
			dataIndex: "name",
		},
		{
			title: t("table.columns.role.label"),
			dataIndex: "label",
		},
		{
			title: t("table.columns.role.order"),
			dataIndex: "order",
		},
		{
			title: t("table.columns.role.description"),
			dataIndex: "description",
		},
		{
			title: t("table.columns.role.status"),
			dataIndex: "status",
			align: "center",
			width: 120,
			render: (status) => {
				const statusItem = statusType?.find((item) => Number(item.value) === status);

				return <Badge variant={status === 1 ? "success" : "error"}>{statusItem?.label}</Badge>;
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
			width: 250,
			fixed: "right",
			render: (_, record) => (
				<div className="grid grid-cols-2 gap-2 text-gray-500">
					<Button
						variant="link"
						size="icon"
						onClick={() => onSetting(record)}
						className="whitespace-nowrap justify-start"
					>
						<Icon icon="solar:settings-bold" size={18} />
						<span className="ml-1">{t("table.button.role_setting")}</span>
					</Button>
					<Button
						variant="link"
						size="icon"
						onClick={() => onCreate(record)}
						className="whitespace-nowrap justify-start"
					>
						<div className="flex items-center">
							<Icon icon="solar:add-square-bold" size={18} />
							<span className="ml-1">{t("table.button.add_sub_role")}</span>
						</div>
					</Button>
					<Button variant="link" size="icon" onClick={() => onEdit(record)} className="whitespace-nowrap justify-start">
						<div className="flex items-center">
							<Icon icon="solar:pen-bold-duotone" size={18} />
							<span className="ml-1">{t("table.button.edit")}</span>
						</div>
					</Button>
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
		<Card title={t("sys.menu.system.role")} size="small">
			<CardHeader>
				<div className="flex items-center justify-between">
					<Button onClick={() => onCreate(undefined)} className="text-white">
						<Icon icon="solar:add-circle-outline" size={18} />
						{t("table.button.add")}
					</Button>
				</div>
			</CardHeader>

			<CardContent>
				<Table
					rowKey={(record) => record.id}
					scroll={{ x: "max-content" }}
					columns={columns}
					dataSource={data}
					loading={isLoading}
					pagination={false}
					expandable={{
						showExpandColumn: false,
						expandedRowKeys: expandedKeys,
						onExpand: (expanded, record) => handleExpand(expanded, record),
					}}
				/>
			</CardContent>
			<RoleModal {...roleModalProps} />
			<SettingModal {...settingModalPros} />
		</Card>
	);
};

export default App;
