import { Icon } from "@/components/icon";
import { Button } from "@/ui/button";
import type { TableProps } from "antd";
import { Popconfirm, Table } from "antd";
import { useEffect, useState } from "react";
import type { ColumnsType, Dictionary } from "#/entity";

import { CardContent, CardHeader } from "@/ui/card";
import { toast } from "sonner";

import {
	useDictionaryActions,
	useDictionaryManageCondition,
	useDictionaryQuery,
	useRemoveDictionaryMutation,
	useUpdateOrCreateDictionaryMutation,
} from "@/store/dictionaryManageStore";
import { useTranslation } from "react-i18next";
import DictionaryModal, { type DictionaryModalProps } from "./modal";

const DictionaryList = ({
	onSelect,
}: {
	onSelect?: (id: number | null) => void;
}) => {
	const defaultDictionaryValue: Dictionary = {
		id: 0,
		name: "",
		type: "",
		status: 2,
		desc: "",
		is_generate_file: 0,
		created_at: "",
		updated_at: "",
		details: [],
	};
	const { t } = useTranslation();
	const updateOrCreateMutation = useUpdateOrCreateDictionaryMutation();
	const removeMutation = useRemoveDictionaryMutation();
	const { data, isLoading } = useDictionaryQuery();
	const condition = useDictionaryManageCondition();
	const { setCondition } = useDictionaryActions();

	const [selectedId, setSelectedId] = useState<number | null>(null);
	const [apiModalProps, setDictionaryModalProps] = useState<DictionaryModalProps>({
		formValue: { ...defaultDictionaryValue },
		title: t("table.button.add"),
		show: false,
		onOk: async (values: Dictionary) => {
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
		if (data && data.list.length > 0 && onSelect) {
			setSelectedId(data.list[0].id);
			onSelect(data.list[0].id);
		}
	}, [data]);

	const handleTableChange: TableProps<Dictionary>["onChange"] = (pagination, filters, sorter) => {
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

	const onEdit = (formValue: Dictionary) => {
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

	const columns: ColumnsType<Dictionary> = [
		{
			title: t("table.columns.dictionary.name"),
			dataIndex: "name",
			key: "name",
			ellipsis: true,
		},
		{
			title: t("table.columns.common.operation"),
			key: "operation",
			align: "center",
			width: 120,
			render: (_, record) => (
				<div className="flex w-full justify-center text-gray-500">
					<Button
						variant="link"
						size="icon"
						onClick={() => onEdit(record)}
						className="flex flex-row  items-center justify-center gap-1 px-2 py-1"
					>
						<Icon icon="solar:pen-bold-duotone" size={18} />
					</Button>
					<Popconfirm
						title="Delete the task"
						description="Are you sure to delete this task?"
						onConfirm={() => handleDelete(record.id)}
						okText="Yes"
						cancelText="No"
					>
						<Button variant="link" size="icon">
							<Icon icon="mingcute:delete-2-fill" size={18} />
						</Button>
					</Popconfirm>
				</div>
			),
		},
	];

	const handleRowClick = (record: Dictionary) => {
		setSelectedId(record.id);
		if (onSelect) {
			onSelect(record.id);
		}
	};

	return (
		<>
			<CardHeader className="p-0">
				<div className="flex items-start justify-start">
					<Button onClick={() => onCreate()} className="text-white" variant="default">
						<Icon icon="solar:add-circle-outline" size={18} />
						{t("table.button.add")}
					</Button>
				</div>
			</CardHeader>

			<CardContent className="p-0">
				<Table<Dictionary>
					rowKey={(record) => record.id}
					scroll={{ x: "100%" }}
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
					onRow={(record) => ({
						onClick: () => handleRowClick(record),
					})}
					rowClassName={(record: Dictionary) => {
						return record.id === selectedId
							? "bg-primary  shadow hover:bg-primary/90"
							: "text-gray-700 dark:text-gray-300";
					}}
				/>
				<DictionaryModal {...apiModalProps} />
			</CardContent>
		</>
	);
};

export default DictionaryList;
