import { Icon } from "@/components/icon";
import { Button } from "@/ui/button";
import type { TableProps } from "antd";
import { Image, Popconfirm, Table } from "antd";
import type { TableRowSelection } from "antd/es/table/interface";
import { useEffect, useState } from "react";
import type { ColumnsType, DictionaryDetail } from "#/entity";

import { useDictionaryByTypeWithCache } from "@/hooks/dict";
import {
	useBatchRemoveDictionaryDetailMutation,
	useDictionaryDetailActions,
	useDictionaryDetailManageCondition,
	useDictionaryDetailQuery,
	useRemoveDictionaryDetailMutation,
	useUpdateOrCreateDictionaryDetailMutation,
} from "@/store/dictionaryDetailManageStore";
import { Badge } from "@/ui/badge";
import { CardContent, CardHeader } from "@/ui/card";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import DictionaryDetailModal, { type DictionaryDetailModalProps } from "./detail-modal";

const DictionaryDetailList = ({
	selectedDictId,
}: {
	selectedDictId: number | null;
}) => {
	const { t } = useTranslation();
	const defaultDictionaryValue: DictionaryDetail = {
		id: 0,
		label: "",
		value: "",
		extend: "",
		status: 2,
		sort: 0,
		type: "string",
		sys_dictionary_Id: selectedDictId,
		created_at: "",
		updated_at: "",
	};
	const { data: statusType } = useDictionaryByTypeWithCache("status");
	const updateOrCreateMutation = useUpdateOrCreateDictionaryDetailMutation();
	const removeMutation = useRemoveDictionaryDetailMutation();
	const batchRemoveMutation = useBatchRemoveDictionaryDetailMutation();
	const { data, isLoading } = useDictionaryDetailQuery();
	const condition = useDictionaryDetailManageCondition();
	const { setCondition } = useDictionaryDetailActions();

	const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
	const [apiModalProps, setDictionaryModalProps] = useState<DictionaryDetailModalProps>({
		formValue: { ...defaultDictionaryValue },
		title: t("table.button.add"),
		show: false,
		onOk: async (values: DictionaryDetail) => {
			updateOrCreateMutation.mutate(values, {
				onSuccess: () => {
					toast.success(t("table.handle_message.success"));
					setDictionaryModalProps((prev) => ({ ...prev, show: false }));
				},
			});
			return true;
		},
		onCancel: () => {
			setDictionaryModalProps((prev) => ({ ...prev, show: false }));
		},
	});

	// biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
	useEffect(() => {
		if (!selectedDictId) {
			return;
		}
		setCondition({
			...condition,
			searchParams: {
				selectedDictId: selectedDictId || 0,
			},
			pagination: {
				...condition.pagination,
				current: 1,
			},
		});
	}, [selectedDictId]);

	const handleTableChange: TableProps<DictionaryDetail>["onChange"] = (pagination, filters, sorter) => {
		setCondition({
			...condition,
			pagination,
			filters,
			sortOrder: Array.isArray(sorter) ? undefined : condition.sortOrder,
			sortField: Array.isArray(sorter) ? undefined : condition.sortField,
		});
	};

	const onCreate = () => {
		setDictionaryModalProps((prev) => ({
			...prev,
			show: true,
			...defaultDictionaryValue,
			title: t("table.button.add"),
			formValue: { ...defaultDictionaryValue },
		}));
	};

	const onEdit = (formValue: DictionaryDetail) => {
		setDictionaryModalProps((prev) => ({
			...prev,
			show: true,
			title: t("table.button.edit"),
			formValue,
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

	const handleDeleteSelection = async () => {
		batchRemoveMutation.mutate(selectedRowKeys as number[], {
			onSuccess: () => {
				toast.success(t("table.handle_message.success"));
			},
			onError: () => {
				toast.error(t("table.handle_message.error"));
			},
		});
	};

	const columns: ColumnsType<DictionaryDetail> = [
		{
			title: "ID",
			dataIndex: "id",
			key: "id",
		},
		{
			title: t("table.columns.dictionary_detail.label"),
			dataIndex: "label",
			key: "label",
			ellipsis: true,
		},
		{
			title: t("table.columns.dictionary_detail.value"),
			dataIndex: "value",
			key: "value",
			render: (_, record) => {
				switch (record.type) {
					case "image":
						return <Image src={record.value} width={50} height={50} />;
					case "icon":
						return <Icon icon={record.value} size={18} />;
					default:
						return record.value;
				}
			},
		},
		{
			title: t("table.columns.dictionary_detail.extend"),
			dataIndex: "extend",
			key: "extend",
		},
		{
			title: t("table.columns.dictionary_detail.status"),
			dataIndex: "status",
			key: "status",
			ellipsis: true,
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
			dataIndex: "operation",
			key: "operation",
			align: "center",
			width: 150,
			fixed: "right",
			render: (_, record) => (
				<div className="grid grid-cols-2 gap-2 text-gray-500">
					<Button variant="link" size="icon" onClick={() => onEdit(record)} className="whitespace-nowrap justify-start">
						<div className="flex items-center">
							<Icon icon="solar:pen-bold-duotone" size={18} />
							<span className="ml-1"> {t("table.button.edit")}</span>
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

	const onSelectChange = (newSelectedRowKeys: React.Key[]) => {
		console.log("selectedRowKeys changed: ", newSelectedRowKeys);
		setSelectedRowKeys(newSelectedRowKeys);
	};

	const rowSelection: TableRowSelection<DictionaryDetail> = {
		selectedRowKeys,
		onChange: onSelectChange,
		selections: [
			Table.SELECTION_ALL,
			Table.SELECTION_INVERT,
			Table.SELECTION_NONE,
			{
				key: "odd",
				text: t("table.columns.common.select_odd_row"),
				onSelect: (changeableRowKeys) => {
					let newSelectedRowKeys = [];
					newSelectedRowKeys = changeableRowKeys.filter((_, index) => {
						if (index % 2 !== 0) {
							return false;
						}
						return true;
					});
					setSelectedRowKeys(newSelectedRowKeys);
				},
			},
			{
				key: "even",
				text: t("table.columns.common.select_even_row"),
				onSelect: (changeableRowKeys) => {
					let newSelectedRowKeys = [];
					newSelectedRowKeys = changeableRowKeys.filter((_, index) => {
						if (index % 2 !== 0) {
							return true;
						}
						return false;
					});
					setSelectedRowKeys(newSelectedRowKeys);
				},
			},
		],
	};

	const hasSelected = selectedRowKeys.length > 0;

	return (
		<>
			<CardHeader className="p-0">
				<div className="flex items-start justify-start">
					<Button onClick={() => onCreate()} className="text-white" variant="default">
						<Icon icon="solar:add-circle-outline" size={18} />
						{t("table.button.add")}
					</Button>
					<Button
						onClick={() => handleDeleteSelection()}
						variant="destructive"
						className="ml-2 text-white"
						disabled={!hasSelected}
					>
						<Icon icon="solar:trash-bin-minimalistic-outline" size={18} />
						{t("table.button.delete")}
					</Button>
				</div>
			</CardHeader>

			<CardContent className="p-0">
				<Table<DictionaryDetail>
					rowKey={(record) => record.id}
					rowSelection={rowSelection}
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
				<DictionaryDetailModal {...apiModalProps} />
			</CardContent>
		</>
	);
};

export default DictionaryDetailList;
