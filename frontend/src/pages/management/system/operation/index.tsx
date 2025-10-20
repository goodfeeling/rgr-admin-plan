import { Icon } from "@/components/icon";
import { useDictionaryByTypeWithCache } from "@/hooks/dict";
import {
	useBatchRemoveOperationMutation,
	useOperationActions,
	useOperationManageCondition,
	useOperationQuery,
	useRemoveOperationMutation,
} from "@/store/operationManageStore";
import { Button } from "@/ui/button";
import { CardContent, CardHeader } from "@/ui/card";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/ui/form";
import type { TableProps } from "antd";
import { Card, Input, Popconfirm, Select, Table } from "antd";
import type { TableRowSelection } from "antd/es/table/interface";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import type { ColumnsType, Operation } from "#/entity";

type SearchFormFieldType = {
	method: string;
	path: string;
	status: number;
};

const searchDefaultValue = { path: undefined, method: undefined, status: 0 };

const App: React.FC = () => {
	const { t } = useTranslation();
	const searchForm = useForm<SearchFormFieldType>({
		defaultValues: searchDefaultValue,
	});
	const removeMutation = useRemoveOperationMutation();
	const batchRemoveMutation = useBatchRemoveOperationMutation();
	const { data, isLoading } = useOperationQuery();
	const condition = useOperationManageCondition();
	const { setCondition } = useOperationActions();
	const { data: apiMethod } = useDictionaryByTypeWithCache("api_method");
	const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

	const handleTableChange: TableProps<Operation>["onChange"] = (pagination, filters, sorter) => {
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

	const columns: ColumnsType<Operation> = [
		{
			title: "ID",
			dataIndex: "id",
			key: "id",
		},
		{
			title: t("table.columns.operation.ip"),
			dataIndex: "ip",
			key: "ip",
		},
		{
			title: t("table.columns.operation.path"),
			dataIndex: "path",
			key: "path",
		},
		{
			title: t("table.columns.operation.method"),
			dataIndex: "method",
			key: "method",
		},
		{
			title: t("table.columns.operation.status"),
			dataIndex: "status",
			key: "status",
		},
		{
			title: t("table.columns.operation.latency"),
			dataIndex: "latency",
			key: "latency",
		},
		{
			title: t("table.columns.operation.agent"),
			dataIndex: "agent",
			key: "agent",
		},
		{
			title: t("table.columns.operation.error"),
			dataIndex: "error_message",
			key: "error_message",
		},
		{
			title: t("table.columns.operation.body"),
			dataIndex: "body",
			key: "body",
			ellipsis: true,
		},
		{
			title: t("table.columns.operation.resp"),
			dataIndex: "resp",
			key: "resp",
			ellipsis: true,
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
			fixed: "right",
			width: 100,
			render: (_, record) => (
				<div className="grid grid-cols-1 gap-2 text-gray-500">
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

	const onSelectChange = (newSelectedRowKeys: React.Key[]) => {
		setSelectedRowKeys(newSelectedRowKeys);
	};

	const rowSelection: TableRowSelection<Operation> = {
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
		<div className="flex flex-col gap-4">
			<Card>
				<CardContent>
					<Form {...searchForm}>
						<div className="flex items-center gap-4">
							<FormField
								control={searchForm.control}
								name="method"
								render={({ field }) => (
									<FormItem>
										<FormLabel>{t("table.columns.operation.method")}</FormLabel>
										<Select
											onChange={(value: string) => {
												field.onChange(value);
											}}
											value={field.value}
											options={apiMethod}
											placeholder={`${t("table.handle_message.select")}${t("table.columns.operation.method")}`}
										/>
									</FormItem>
								)}
							/>
							<FormField
								control={searchForm.control}
								name="path"
								render={({ field }) => (
									<FormItem>
										<FormLabel>{t("table.columns.operation.path")}</FormLabel>
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
										<FormLabel>{t("table.columns.operation.status")}</FormLabel>
										<FormControl>
											<Input {...field} />
										</FormControl>
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
			<Card title={t("sys.menu.system.operation")} size="small">
				<CardHeader>
					<div className="flex items-center justify-between">
						<Button
							variant="destructive"
							className="text-white"
							onClick={() => handleDeleteSelection()}
							disabled={!hasSelected}
						>
							<Icon icon="solar:trash-bin-minimalistic-outline" size={18} />
							{t("table.button.delete")}
						</Button>
					</div>
				</CardHeader>
				<CardContent>
					<Table<Operation>
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
				</CardContent>
			</Card>
		</div>
	);
};

export default App;
