import { Icon } from "@/components/icon";
import { Button } from "@/ui/button";
import { Popconfirm, Table } from "antd";
import { useEffect, useState } from "react";
import type { ColumnsType, Menu } from "#/entity";

import { useMenuQuery, useRemoveMenuMutation, useUpdateOrCreateMenuMutation } from "@/store/menuManageStore";
import { Badge } from "@/ui/badge";
import { CardContent, CardHeader } from "@/ui/card";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import MenuModal, { type MenuModalProps } from "./menu-modal";
import SettingModal, { type SettingModalProps } from "./setting-modal";

const MenuList = ({ selectedId }: { selectedId: number | null }) => {
	const { t } = useTranslation();
	const defaultValue: Menu = {
		id: 0,
		menu_level: 0,
		parent_id: 0,
		name: "",
		path: "",
		hidden: false,
		component: "",
		sort: 0,
		keep_alive: 0,
		title: "",
		icon: "",
		menu_group_id: selectedId ? selectedId : 0,
		created_at: "",
		updated_at: "",
		level: [],
		children: [],
	};

	const updateOrCreateMutation = useUpdateOrCreateMenuMutation();
	const removeMutation = useRemoveMenuMutation();
	const { data, isLoading } = useMenuQuery(selectedId ?? 0);
	const [expandedKeys, setExpandedKeys] = useState<number[]>([]);

	const [menuModalProps, setUserModalProps] = useState<MenuModalProps>({
		formValue: { ...defaultValue },
		title: t("table.button.add"),
		show: false,
		treeRawData: [],
		onOk: async (values: Menu): Promise<boolean> => {
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

	const [settingModalProps, setSettingModalProps] = useState<SettingModalProps>({
		formValue: { id: 0 },
		title: t("table.button.add"),
		show: false,
		onCancel: () => {
			setSettingModalProps((prev) => ({ ...prev, show: false }));
		},
	});

	useEffect(() => {
		if (data) {
			setUserModalProps((prev) => ({
				...prev,
				treeRawData: data,
			}));
		}
	}, [data]);

	// create menu
	const onCreate = (formValue: Menu | undefined) => {
		const setValue = defaultValue;
		if (formValue !== undefined) {
			setValue.parent_id = formValue.id;
		}

		setUserModalProps((prev) => ({
			...prev,
			show: true,
			...setValue,
			title: t("table.button.add"),
			formValue: { ...setValue },
		}));
	};

	const onEdit = (formValue: Menu) => {
		setUserModalProps((prev) => ({
			...prev,
			show: true,
			title: t("table.button.edit"),
			formValue,
		}));
	};

	const onSetting = (formValue: Menu) => {
		setSettingModalProps((prev) => ({
			...prev,
			show: true,
			title: t("table.button.button_and_parameter"),
			formValue: { id: formValue.id },
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
	const handleExpand = (expanded: boolean, record: Menu) => {
		const keys = expanded ? [...expandedKeys, record.id] : expandedKeys.filter((key) => key !== record.id);
		setExpandedKeys(keys);
	};
	const columns: ColumnsType<Menu> = [
		{
			title: t("table.columns.menu.menu_id"),
			dataIndex: "expand",
			render: (_, record) => {
				const level = record.level.length;
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
			title: t("table.columns.menu.title"),
			dataIndex: "title",
			render: (_, record) => {
				return <span>{t(record.title)}</span>;
			},
		},
		{
			title: t("table.columns.menu.icon"),
			dataIndex: "icon",
			render: (_, record) => {
				return (
					<div>
						<Icon icon={record.icon} size={18} />
						<span> {record.icon}</span>
					</div>
				);
			},
		},
		{
			title: t("table.columns.menu.name"),
			dataIndex: "name",
		},
		{
			title: t("table.columns.menu.path"),
			dataIndex: "path",
		},
		{
			title: t("table.columns.menu.hidden"),
			dataIndex: "hidden",
			render: (_, record) => {
				return (
					<Badge variant={record.hidden ? "success" : "error"}>
						{record.hidden ? t("table.button.yes") : t("table.button.no")}
					</Badge>
				);
			},
		},
		{
			title: t("table.columns.menu.parent_id"),
			dataIndex: "parent_id",
		},
		{
			title: t("table.columns.menu.sort"),
			dataIndex: "sort",
		},
		{
			title: t("table.columns.menu.component"),
			dataIndex: "component",
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
			width: 300,
			fixed: "right",
			render: (_, record) => (
				<div className="grid grid-cols-2 gap-2 text-gray-500">
					<Button
						variant="link"
						size="sm"
						onClick={() => onCreate(record)}
						className="whitespace-nowrap justify-start "
					>
						<div className="flex items-center">
							<Icon icon="solar:add-square-bold" size={18} />
							<span className="ml-1">{t("table.button.add_sub_route")}</span>
						</div>
					</Button>
					<Button variant="link" size="sm" onClick={() => onEdit(record)} className="whitespace-nowrap justify-start ">
						<div className="flex items-center">
							<Icon icon="solar:pen-bold-duotone" size={18} />
							<span className="ml-1">{t("table.button.edit")}</span>
						</div>
					</Button>
					<Button
						variant="link"
						size="sm"
						onClick={() => onSetting(record)}
						className="whitespace-nowrap justify-start"
					>
						<div className="flex items-center">
							<Icon icon="solar:pen-new-square-outline" size={18} />
							<span className="ml-1">{t("table.button.button_and_parameter")}</span>
						</div>
					</Button>

					<Popconfirm
						title={t("table.handle_message.delete_prompt")}
						description={t("table.handle_message.confirm_delete")}
						onConfirm={() => handleDelete(record.id)}
						okText={t("table.button.yes")}
						cancelText={t("table.button.no")}
					>
						<Button variant="link" size="sm" className="whitespace-nowrap justify-start text-white">
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
		<>
			<CardHeader className="p-0">
				<div className="flex items-start justify-start">
					<Button onClick={() => onCreate(undefined)} className="text-white">
						<Icon icon="solar:add-circle-outline" size={18} />
						{t("table.button.add")}
					</Button>
				</div>
			</CardHeader>

			<CardContent className="p-0">
				<Table<Menu>
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
				<MenuModal {...menuModalProps} />
				<SettingModal {...settingModalProps} />
			</CardContent>
		</>
	);
};

export default MenuList;
