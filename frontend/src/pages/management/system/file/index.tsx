import { Icon } from "@/components/icon";
import {
	useBatchRemoveFileInfoMutation,
	useFileInfoActions,
	useFileInfoManageCondition,
	useFileInfoQuery,
	useRemoveFileInfoMutation,
	useUpdateOrCreateFileInfoMutation,
} from "@/store/fileManageStore";
import { Button } from "@/ui/button";
import { CardContent, CardHeader } from "@/ui/card";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/ui/form";

import { useDictionaryByTypeWithCache } from "@/hooks/dict";
import { useQueryClient } from "@tanstack/react-query";
import type { TableProps } from "antd";
import { Card, Input, Popconfirm, Select, Table } from "antd";
import type { TableRowSelection } from "antd/es/table/interface";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import type { ColumnsType, FileInfo } from "#/entity";
import FileModal, { type FileModalProps } from "./modal";

const defaultFileValue: FileInfo = {
	id: 0,
	file_name: "",
	file_md5: "",
	file_path: "",
	file_url: "",
	file_origin_name: "",
	storage_engine: "",
	created_at: "",
	updated_at: "",
};

type SearchFormFieldType = {
	file_origin_name?: string;
	storage_engine?: string;
};

const searchDefaultValue = {
	file_origin_name: undefined,
	storage_engine: undefined,
};

const App: React.FC = () => {
	const queryClient = useQueryClient();
	const { t } = useTranslation();
	const searchForm = useForm<SearchFormFieldType>({
		defaultValues: searchDefaultValue,
	});

	const updateOrCreateMutation = useUpdateOrCreateFileInfoMutation();
	const removeMutation = useRemoveFileInfoMutation();
	const batchRemoveMutation = useBatchRemoveFileInfoMutation();
	const { data, isLoading } = useFileInfoQuery();
	const condition = useFileInfoManageCondition();
	const { setCondition } = useFileInfoActions();
	const { data: storageEngine } = useDictionaryByTypeWithCache("file_storage_engine");

	const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
	const [apiModalProps, setFileModalProps] = useState<FileModalProps>({
		formValue: { ...defaultFileValue },
		title: t("table.button.upload"),
		storageEngine: undefined,
		show: false,
		onOk: async (values: FileInfo | null) => {
			if (values) {
				updateOrCreateMutation.mutate(values, {
					onSuccess: () => {
						console.log("created file info success");
					},
				});
			} else {
				// manual refresh
				queryClient.invalidateQueries({ queryKey: ["fileManageList"] });
			}

			return true;
		},
		onCancel: () => {
			setFileModalProps((prev) => ({ ...prev, show: false }));
		},
	});

	useEffect(() => {
		setFileModalProps((prev) => ({ ...prev, storageEngine }));
	}, [storageEngine]);

	const handleTableChange: TableProps<FileInfo>["onChange"] = (pagination, filters, sorter) => {
		setCondition({
			...condition,
			pagination,
			filters,
			sortOrder: Array.isArray(sorter) ? undefined : condition.sortOrder,
			sortField: Array.isArray(sorter) ? undefined : condition.sortField,
		});
	};

	const onCreate = () => {
		setFileModalProps((prev) => ({
			...prev,
			show: true,
			...defaultFileValue,
			title: t("table.button.upload"),
			formValue: { ...defaultFileValue },
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

	const columns: ColumnsType<FileInfo> = [
		{
			title: "ID",
			dataIndex: "id",
			key: "id",
		},
		{
			title: t("table.columns.file.file_origin_name"),
			dataIndex: "file_origin_name",
			key: "file_origin_name",
			ellipsis: true,
			render: (_, record) => {
				return (
					<div className="flex">
						<a target="_blank" href={record.file_url} rel="noopener noreferrer">
							{record.file_origin_name}
						</a>
					</div>
				);
			},
		},

		{
			title: t("table.columns.file.storage_engine"),
			dataIndex: "storage_engine",
			key: "storage_engine",
		},
		{
			title: t("table.columns.common.created_at"),
			dataIndex: "created_at",
			key: "created_at",
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
		console.log("selectedRowKeys changed: ", newSelectedRowKeys);
		setSelectedRowKeys(newSelectedRowKeys);
	};

	const rowSelection: TableRowSelection<FileInfo> = {
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

	return (
		<div className="flex flex-col gap-4">
			<Card>
				<CardContent>
					<Form {...searchForm}>
						<div className="flex items-center gap-4">
							<FormField
								control={searchForm.control}
								name="file_origin_name"
								render={({ field }) => (
									<FormItem>
										<FormLabel>{t("table.columns.file.file_origin_name")}</FormLabel>
										<FormControl>
											<Input {...field} />
										</FormControl>
									</FormItem>
								)}
							/>

							<FormField
								control={searchForm.control}
								name="storage_engine"
								render={({ field }) => (
									<FormItem>
										<FormLabel>{t("table.columns.file.storage_engine")}</FormLabel>
										<FormControl>
											<Select
												style={{ width: 150 }}
												onChange={(value: string) => {
													field.onChange(value);
												}}
												value={field.value}
												options={storageEngine}
												placeholder={`${t("table.handle_message.select")}${t("table.columns.file.storage_engine")}`}
											/>
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
			<Card title={t("sys.menu.system.file")} size="small">
				<CardHeader>
					<div className="flex items-start justify-start">
						<Button onClick={() => onCreate()} className="text-white" variant="default">
							<Icon icon="solar:add-circle-outline" size={18} />
							{t("table.button.upload")}
						</Button>
						<Button
							onClick={() => handleDeleteSelection()}
							variant="destructive"
							className="ml-2 text-white"
							disabled={!(selectedRowKeys.length > 0)}
						>
							<Icon icon="solar:trash-bin-minimalistic-outline" size={18} />
							{t("table.button.delete")}
						</Button>
					</div>
				</CardHeader>

				<CardContent>
					<Table<FileInfo>
						rowKey={(record) => record.id}
						rowSelection={rowSelection}
						scroll={{ x: "max-content" }}
						columns={columns}
						pagination={{
							current: data?.page || condition.pagination?.current || 1,
							pageSize: data?.page_size || condition.pagination?.pageSize || 10,
							total: data?.total || condition?.pagination?.total || 0,
							showTotal: (total) => `${t("table.page.total")} ${total} ${t("table.page.items")}`,
							showSizeChanger: true,
							pageSizeOptions: ["10", "20", "50", "100"],
						}}
						dataSource={data?.list}
						loading={isLoading}
						onChange={handleTableChange}
					/>
				</CardContent>
				<FileModal {...apiModalProps} />
			</Card>
		</div>
	);
};

export default App;
