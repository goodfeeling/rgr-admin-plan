import apisService, { ApisService } from "@/api/services/apisService";
import { Icon } from "@/components/icon";
import UploadTool from "@/components/upload/upload-multiple";
import { useDictionaryByTypeWithCache } from "@/hooks/dict";
import {
	useApiActions,
	useApiManageCondition,
	useApiQuery,
	useBatchRemoveApiMutation,
	useRemoveApiMutation,
	useSynchronizeApiMutation,
	useUpdateOrCreateApiMutation,
} from "@/store/apiManageStore";
import { Button } from "@/ui/button";
import { CardContent, CardHeader } from "@/ui/card";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/ui/form";
import type { TableProps } from "antd";
import { Card, Input, Popconfirm, Select, Table } from "antd";
import type { TableRowSelection } from "antd/es/table/interface";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import type { Api, ColumnsType } from "#/entity";
import ApiModal, { type ApiModalProps } from "./api-modal";

const defaultApiValue: Api = {
	id: 0,
	path: "",
	api_group: "",
	method: "",
	description: "",
	created_at: "",
	updated_at: "",
};

type SearchFormFieldType = {
	path?: string;
	description?: string;
	api_group?: string;
	method?: string;
};

const searchDefaultValue = {
	path: "",
	description: "",
	api_group: undefined,
	method: undefined,
};

const App: React.FC = () => {
	const { t } = useTranslation();
	const searchForm = useForm<SearchFormFieldType>({
		defaultValues: searchDefaultValue,
	});

	const updateOrCreateMutation = useUpdateOrCreateApiMutation();
	const removeMutation = useRemoveApiMutation();
	const batchRemoveMutation = useBatchRemoveApiMutation();
	const synchronizeMutation = useSynchronizeApiMutation();
	const { data, isLoading } = useApiQuery();
	const condition = useApiManageCondition();
	const { setCondition } = useApiActions();
	const { data: apiGroup } = useDictionaryByTypeWithCache("api_group");
	const { data: apiMethod } = useDictionaryByTypeWithCache("api_method");

	const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
	const [apiModalProps, setApiModalProps] = useState<ApiModalProps>({
		formValue: { ...defaultApiValue },
		apiGroup: undefined,
		apiMethod: undefined,
		title: t("table.button.add"),
		show: false,
		onOk: async (values: Api) => {
			updateOrCreateMutation.mutate(values, {
				onSuccess: () => {
					toast.success(t("table.handle_message.success"));
					setApiModalProps((prev) => ({ ...prev, show: false }));
				},
			});
			return true;
		},
		onCancel: () => {
			setApiModalProps((prev) => ({ ...prev, show: false }));
		},
	});

	useEffect(() => {
		setApiModalProps((prev) => ({
			...prev,
			apiGroup: apiGroup,
			apiMethod: apiMethod,
		}));
	}, [apiGroup, apiMethod]);

	const handleTableChange: TableProps<Api>["onChange"] = (pagination, filters, sorter) => {
		setCondition({
			...condition,
			pagination,
			filters,
			sortOrder: Array.isArray(sorter) ? undefined : condition.sortOrder,
			sortField: Array.isArray(sorter) ? undefined : condition.sortField,
		});
	};

	const onCreate = () => {
		setApiModalProps((prev) => ({
			...prev,
			show: true,
			...defaultApiValue,
			title: t("table.button.add"),
			formValue: { ...defaultApiValue },
		}));
	};

	const onEdit = (formValue: Api) => {
		setApiModalProps((prev) => ({
			...prev,
			show: true,
			title: t("table.button.edit"),
			formValue,
		}));
	};

	// single delete
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

	// batch delete
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

	const columns: ColumnsType<Api> = [
		{
			title: "ID",
			dataIndex: "id",
			key: "id",
		},
		{
			title: t("table.columns.api.path"),
			dataIndex: "path",
			key: "path",
			ellipsis: true,
		},
		{
			title: t("table.columns.api.api_group"),
			dataIndex: "api_group",
			key: "api_group",
		},
		{
			title: t("table.columns.api.method"),
			dataIndex: "method",
			key: "method",
		},
		{
			title: t("table.columns.api.description"),
			dataIndex: "description",
			key: "description",
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
			width: 150,
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

	// sync data
	const onSynchronize = () => {
		synchronizeMutation.mutate();
	};

	// selector change
	const onSelectChange = (newSelectedRowKeys: React.Key[]) => {
		console.log("selectedRowKeys changed: ", newSelectedRowKeys);
		setSelectedRowKeys(newSelectedRowKeys);
	};

	const rowSelection: TableRowSelection<Api> = {
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

	const handleExportMenuClick = async () => {
		// 处理导出全部数据逻辑
		try {
			// 调用 downloadTemplate 方法获取 Blob 数据
			const response = await apisService.exportApi();

			// 创建一个隐藏的 a 标签用于触发下载
			const url = window.URL.createObjectURL(new Blob([response]));
			const link = document.createElement("a");
			link.href = url;
			link.setAttribute("download", "apis_export.xlsx"); // 设置下载文件名
			document.body.appendChild(link);
			link.click(); // 触发点击下载

			// 清理创建的 URL 对象和 a 标签
			window.URL.revokeObjectURL(url);
			document.body.removeChild(link);
		} catch (error) {
			console.error("Download failed:", error);
			// 处理错误情况，例如显示通知给用户
		}
	};

	const downloadTemplate = async () => {
		try {
			// 调用 downloadTemplate 方法获取 Blob 数据
			const response = await apisService.downloadTemplate();

			// 创建一个隐藏的 a 标签用于触发下载
			const url = window.URL.createObjectURL(new Blob([response]));
			const link = document.createElement("a");
			link.href = url;
			link.setAttribute("download", "apis_export.xlsx"); // 设置下载文件名
			document.body.appendChild(link);
			link.click(); // 触发点击下载

			// 清理创建的 URL 对象和 a 标签
			window.URL.revokeObjectURL(url);
			document.body.removeChild(link);
		} catch (error) {
			console.error("Download failed:", error);
			// 处理错误情况，例如显示通知给用户
		}
	};

	return (
		<div className="flex flex-col gap-4">
			<Card>
				<CardContent>
					<Form {...searchForm}>
						<div className="flex items-center gap-4">
							<FormField
								control={searchForm.control}
								name="path"
								render={({ field }) => (
									<FormItem>
										<FormLabel>{t("table.columns.api.path")}</FormLabel>
										<FormControl>
											<Input {...field} />
										</FormControl>
									</FormItem>
								)}
							/>

							<FormField
								control={searchForm.control}
								name="description"
								render={({ field }) => (
									<FormItem>
										<FormLabel>{t("table.columns.api.description")}</FormLabel>
										<FormControl>
											<Input {...field} />
										</FormControl>
									</FormItem>
								)}
							/>
							<FormField
								control={searchForm.control}
								name="method"
								render={({ field }) => (
									<FormItem>
										<FormLabel>{t("table.columns.api.method")}</FormLabel>
										<Select
											onChange={(value: string) => {
												field.onChange(value);
											}}
											value={field.value}
											options={apiMethod}
											placeholder={`${t("table.handle_message.select")}${t("table.columns.api.method")}`}
										/>
									</FormItem>
								)}
							/>
							<FormField
								control={searchForm.control}
								name="api_group"
								render={({ field }) => (
									<FormItem>
										<FormLabel>{t("table.columns.api.api_group")}</FormLabel>
										<Select
											onChange={(value: string) => {
												field.onChange(value);
											}}
											value={field.value}
											options={apiGroup}
											placeholder={`${t("table.handle_message.select")}${t("table.columns.api.api_group")}`}
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
			<Card title={t("sys.menu.system.api")} size="small">
				<CardHeader>
					<div className="flex items-start justify-start">
						<Button onClick={() => onCreate()} variant="default" className="text-white">
							<Icon icon="solar:add-circle-outline" size={18} />
							{t("table.button.add")}
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
						<Button onClick={() => onSynchronize()} variant="outline" className="ml-2 text-white">
							<Icon icon="solar:refresh-outline" size={18} />
							{t("table.button.synchronize")}
						</Button>
						<Button onClick={() => downloadTemplate()} className="ml-2 text-white" variant="default">
							<Icon icon="solar:cloud-download-outline" size={18} />
							{t("table.button.download_template")}
						</Button>
						<UploadTool
							title={t("table.button.import")}
							listType="text"
							renderType="button"
							accept=".xlsx"
							showUploadList={false}
							uploadUri={ApisService.Client.Import}
						/>
						<Button className="ml-2 text-white" variant="default" onClick={() => handleExportMenuClick()}>
							<Icon icon="solar:import-outline" size={18} />
							{t("table.button.export")}
						</Button>
					</div>
				</CardHeader>

				<CardContent>
					<Table<Api>
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
				<ApiModal {...apiModalProps} />
			</Card>
		</div>
	);
};

export default App;
