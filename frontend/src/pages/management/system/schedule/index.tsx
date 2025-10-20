import { Icon } from "@/components/icon";

import { useDictionaryByTypeWithCache } from "@/hooks/dict";
import {
	useBatchRemoveScheduledTaskMutation,
	useDisableTaskMutation,
	useEnableTaskMutation,
	useReloadTaskMutation,
	useRemoveScheduledTaskMutation,
	useScheduledTaskManageCondition,
	useScheduledTaskManegeActions,
	useScheduledTaskQuery,
	useUpdateOrCreateScheduledTaskMutation,
} from "@/store/scheduleManageStore";
import { Button } from "@/ui/button";
import { CardContent, CardHeader } from "@/ui/card";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/ui/form";

import { Badge } from "@/ui/badge";
import type { TableProps } from "antd";
import { Card, Input, Popconfirm, Select, Table } from "antd";
import type { TableRowSelection } from "antd/es/table/interface";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import type { ColumnsType, ScheduledTask } from "#/entity";
import ExecLogModal, { type LogModalProps } from "./log";
import ScheduledTaskModal, { type ScheduledTaskModalProps } from "./modal";

const defaultScheduledTaskValue: ScheduledTask = {
	id: 0,
	task_name: "",
	task_description: "",
	cron_expression: "",
	task_type: "",
	task_params: {},
	status: 1,
	exec_type: "",
	last_execute_time: "",
	next_execute_time: "",
	created_at: "",
	updated_at: "",
};

type SearchFormFieldType = {
	task_name?: string;
	status?: string;
	task_type?: string;
};

const searchDefaultValue = {
	task_name: undefined,
	status: undefined,
	task_type: undefined,
};

const App: React.FC = () => {
	const { t } = useTranslation();
	const searchForm = useForm<SearchFormFieldType>({
		defaultValues: searchDefaultValue,
	});

	const updateOrCreateMutation = useUpdateOrCreateScheduledTaskMutation();
	const removeMutation = useRemoveScheduledTaskMutation();
	const batchRemoveMutation = useBatchRemoveScheduledTaskMutation();
	const enableTaskMutation = useEnableTaskMutation();
	const disableTaskMutation = useDisableTaskMutation();
	const reloadTaskMutation = useReloadTaskMutation();

	// load data
	const { data, isLoading } = useScheduledTaskQuery({ enablePolling: true });
	const condition = useScheduledTaskManageCondition();
	const { setCondition } = useScheduledTaskManegeActions();

	// enum type
	const { data: taskTypes } = useDictionaryByTypeWithCache("task_type");
	const { data: statusType } = useDictionaryByTypeWithCache("task_status");
	const statusTypeMap = new Map<string, string>(statusType?.map((item) => [item.value, item.label]));

	const [processingTaskIds, setProcessingTaskIds] = useState<Set<number>>(new Set());
	const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
	const [scheduleModalProps, setScheduledTaskModalProps] = useState<ScheduledTaskModalProps>({
		formValue: { ...defaultScheduledTaskValue },
		title: t("table.button.add"),
		show: false,
		onOk: async (values: ScheduledTask) => {
			updateOrCreateMutation.mutate(values, {
				onSuccess: () => {
					toast.success(t("table.handle_message.success"));
					setScheduledTaskModalProps((prev) => ({ ...prev, show: false }));
				},
			});
			return true;
		},
		onCancel: () => {
			setScheduledTaskModalProps((prev) => ({ ...prev, show: false }));
		},
	});

	const [logModalProps, setLogModalProps] = useState<LogModalProps>({
		title: t("table.button.log"),
		show: false,
		id: 0,
		onCancel: () => {
			setLogModalProps((prev) => ({ ...prev, show: false }));
		},
	});
	const handleTableChange: TableProps<ScheduledTask>["onChange"] = (pagination, filters, sorter) => {
		setCondition({
			...condition,
			pagination,
			filters,
			sortOrder: Array.isArray(sorter) ? undefined : condition.sortOrder,
			sortField: Array.isArray(sorter) ? undefined : condition.sortField,
		});
	};

	const onCreate = () => {
		setScheduledTaskModalProps((prev) => ({
			...prev,
			show: true,
			...defaultScheduledTaskValue,
			title: t("table.button.add"),

			formValue: { ...defaultScheduledTaskValue },
		}));
	};

	const onEdit = (formValue: ScheduledTask) => {
		setScheduledTaskModalProps((prev) => ({
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
	// show log
	const onShowLog = (id: number) => {
		console.log(id);

		setLogModalProps((prev) => ({
			...prev,
			show: true,
			id,
		}));
	};

	const columns: ColumnsType<ScheduledTask> = [
		{
			title: "ID",
			dataIndex: "id",
			key: "id",
		},
		{
			title: t("table.columns.schedule.task_name"),
			dataIndex: "task_name",
			key: "task_name",
			ellipsis: true,
		},
		{
			title: t("table.columns.schedule.cron_expression"),
			dataIndex: "cron_expression",
			key: "cron_expression",
		},
		{
			title: t("table.columns.schedule.task_type"),
			dataIndex: "task_type",
			key: "task_type",
			render: (task_type) => {
				const taskType = taskTypes?.find((item) => item.value === task_type);
				return taskType?.label || task_type;
			},
		},
		{
			title: t("table.columns.schedule.status"),
			dataIndex: "status",
			key: "status",
			render: (status: number) => {
				const statusResult = statusTypeMap.get(status.toString())?.toLowerCase();
				return <Badge variant={(statusResult as any) ?? "default"}>{statusResult}</Badge>;
			},
		},

		{
			title: t("table.columns.schedule.last_execute_time"),
			dataIndex: "last_execute_time",
			key: "last_execute_time",
		},
		{
			title: t("table.columns.schedule.next_execute_time"),
			dataIndex: "next_execute_time",
			key: "next_execute_time",
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
			width: 300,
			render: (_, record) => (
				<div className="grid grid-cols-3 gap-2 text-gray-500">
					<Button
						variant="link"
						size="icon"
						onClick={() => onEnableTask(record)}
						disabled={record.status === 0 || record.status === 2 || processingTaskIds.has(record.id)}
						className="whitespace-nowrap justify-start"
					>
						{processingTaskIds.has(record.id) && record.status !== 2 ? (
							<>
								<Icon icon="svg-spinners:bars-rotate-fade" size={18} />
								<span className="ml-1">{t("table.button.handling")}</span>
							</>
						) : (
							<>
								<Icon icon="solar:rewind-back-line-duotone" size={18} />
								<span className="ml-1">{t("table.button.start")}</span>
							</>
						)}
					</Button>
					<Button
						variant="link"
						size="icon"
						onClick={() => onDisableTask(record)}
						disabled={record.status === 1 || processingTaskIds.has(record.id)}
						className="whitespace-nowrap justify-start"
					>
						<div className="flex items-center">
							{processingTaskIds.has(record.id) && record.status !== 1 ? (
								<>
									<Icon icon="svg-spinners:bars-rotate-fade" size={18} />
									<span className="ml-1">{t("table.button.handling")}</span>
								</>
							) : (
								<>
									<Icon icon="solar:stop-circle-outline" size={18} />
									<span className="ml-1">{t("table.button.stop")}</span>
								</>
							)}
						</div>
					</Button>
					<Button
						variant="link"
						size="icon"
						onClick={() => onShowLog(record.id)}
						className="whitespace-nowrap justify-start"
					>
						<div className="flex items-center">
							<Icon icon="solar:menu-dots-circle-linear" size={18} />
							<span className="ml-1">
								<span className="ml-1"> {t("table.button.log")}</span>
							</span>
						</div>
					</Button>
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

	// running task
	const onEnableTask = async (formValue: ScheduledTask) => {
		setProcessingTaskIds((prev) => new Set(prev).add(formValue.id));
		enableTaskMutation.mutate(formValue.id, {
			onSuccess: () => {
				toast.success(t("table.handle_message.success"));
				setProcessingTaskIds((prev) => {
					const newSet = new Set(prev);
					newSet.delete(formValue.id);
					return newSet;
				});
			},
			onError: () => {
				toast.error(t("table.handle_message.error"));
				setProcessingTaskIds((prev) => {
					const newSet = new Set(prev);
					newSet.delete(formValue.id);
					return newSet;
				});
			},
		});
	};

	// stop task
	const onDisableTask = (formValue: ScheduledTask) => {
		// 添加任务到处理中集合
		setProcessingTaskIds((prev) => new Set(prev).add(formValue.id));

		disableTaskMutation.mutate(formValue.id, {
			onSuccess: () => {
				toast.success(t("table.handle_message.success"));
				// 从处理中集合移除任务
				setProcessingTaskIds((prev) => {
					const newSet = new Set(prev);
					newSet.delete(formValue.id);
					return newSet;
				});
			},
			onError: () => {
				toast.error(t("table.handle_message.error"));
				// 从处理中集合移除任务
				setProcessingTaskIds((prev) => {
					const newSet = new Set(prev);
					newSet.delete(formValue.id);
					return newSet;
				});
			},
		});
	};

	// reload task
	const onReloadTask = () => {
		reloadTaskMutation.mutate(undefined, {
			onSuccess: () => {
				toast.success(t("table.handle_message.success"));
			},
			onError: () => {
				toast.error(t("table.handle_message.error"));
			},
		});
	};

	// 选择改变
	const onSelectChange = (newSelectedRowKeys: React.Key[]) => {
		console.log("selectedRowKeys changed: ", newSelectedRowKeys);
		setSelectedRowKeys(newSelectedRowKeys);
	};

	const rowSelection: TableRowSelection<ScheduledTask> = {
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
								name="task_name"
								render={({ field }) => (
									<FormItem>
										<FormLabel>{t("table.columns.schedule.task_name")}</FormLabel>
										<FormControl>
											<Input {...field} />
										</FormControl>
									</FormItem>
								)}
							/>
							<FormField
								control={searchForm.control}
								name="task_type"
								render={({ field }) => (
									<FormItem>
										<FormLabel>{t("table.columns.schedule.task_type")}</FormLabel>
										<Select
											onChange={(value: string) => {
												field.onChange(value);
											}}
											value={field.value}
											options={taskTypes}
											placeholder={`${t("table.handle_message.select")}${t("table.columns.schedule.task_type")}`}
										/>
									</FormItem>
								)}
							/>

							<FormField
								control={searchForm.control}
								name="status"
								render={({ field }) => (
									<FormItem>
										<FormLabel>{t("table.columns.schedule.status")}</FormLabel>
										<Select
											onChange={(value: string) => {
												field.onChange(value);
											}}
											value={field.value}
											options={statusType}
											placeholder={`${t("table.handle_message.select")}${t("table.columns.schedule.status")}`}
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
			<Card title={t("sys.menu.system.schedule")} size="small">
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
						<Button onClick={() => onReloadTask()} variant="default" className="ml-2 text-white">
							<Icon icon="solar:refresh-bold" size={18} />
							{t("table.button.reload_all_task")}
						</Button>
					</div>
				</CardHeader>

				<CardContent>
					<Table<ScheduledTask>
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
						onRow={(record) => ({
							className: processingTaskIds.has(record.id) ? "opacity-75" : "",
						})}
					/>
				</CardContent>
				<ScheduledTaskModal {...scheduleModalProps} />
				<ExecLogModal {...logModalProps} />
			</Card>
		</div>
	);
};

export default App;
