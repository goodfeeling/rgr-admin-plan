import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/ui/form";

import { useTranslationRule } from "@/hooks";
import { useDictionaryByTypeWithCache } from "@/hooks/dict";
import AdvancedCronField from "@/pages/components/cron";
import { Button, Input, Modal, Select } from "antd";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import type { ScheduledTask } from "#/entity";

export type ScheduledTaskModalProps = {
	formValue: ScheduledTask;
	title: string;
	show: boolean;
	onOk: (values: ScheduledTask) => Promise<boolean>;
	onCancel: VoidFunction;
};

export default function ScheduledTaskModal({ title, show, formValue, onOk, onCancel }: ScheduledTaskModalProps) {
	const { t } = useTranslation();
	const { data: taskTypes } = useDictionaryByTypeWithCache("task_type");
	const { data: apiMethod } = useDictionaryByTypeWithCache("api_method");
	const { data: taskExecType } = useDictionaryByTypeWithCache("task_exec");

	const form = useForm<ScheduledTask>();
	const [loading, setLoading] = useState(false);

	const taskType = form.watch("task_type");

	useEffect(() => {
		form.reset(formValue);
	}, [formValue, form]);

	const handleOk = async () => {
		form.handleSubmit(async (values) => {
			setLoading(true);
			const res = await onOk(values);
			if (res) {
				setLoading(false);
			}
		})();
	};

	const handleCancel = () => {
		onCancel();
	};

	// 根据task_type渲染不同的参数输入界面
	const renderTaskParams = () => {
		switch (taskType) {
			case "http_call":
				return (
					<div className="space-y-4">
						<FormField
							control={form.control}
							name="task_params.url"
							rules={{ required: "url is required" }}
							render={({ field }) => (
								<FormItem>
									<FormLabel>URL</FormLabel>
									<FormControl>
										<Input {...field} placeholder="http://example.com/api/health" />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="task_params.method"
							rules={{ required: "method is required" }}
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.schedule.method")}</FormLabel>
									<FormControl>
										<Select {...field} style={{ width: "100%" }} options={apiMethod} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="task_params.timeout"
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.schedule.timeout")}</FormLabel>
									<FormControl>
										<Input
											type="number"
											{...field}
											value={field.value || "30"}
											onChange={(e) => field.onChange(Number(e.target.value))}
											placeholder="30"
										/>
									</FormControl>
								</FormItem>
							)}
						/>
					</div>
				);

			case "function":
				return (
					<div className="space-y-4">
						<FormField
							control={form.control}
							name="task_params.function_name"
							rules={{ required: "function_name is required" }}
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.schedule.function_name")}</FormLabel>
									<FormControl>
										<Input {...field} placeholder="task function name" />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
					</div>
				);
			// 添加脚本任务类型
			case "script_exec":
				return (
					<div className="space-y-4">
						<FormField
							control={form.control}
							name="task_params.script_path"
							rules={{ required: "script_path is required" }}
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.schedule.script_path")}</FormLabel>
									<FormControl>
										<Input {...field} placeholder="/path/to/your/script.sh" />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="task_params.arguments"
							rules={{ required: "arguments is required" }}
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.schedule.arguments")}</FormLabel>
									<FormControl>
										<Input {...field} placeholder="arg1 arg2 arg3" />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="task_params.timeout"
							render={({ field }) => (
								<FormItem>
									<FormLabel> {t("table.columns.schedule.timeout")}</FormLabel>
									<FormControl>
										<Input
											type="number"
											{...field}
											value={field.value || ""}
											onChange={(e) => field.onChange(Number(e.target.value))}
											placeholder="60"
										/>
									</FormControl>
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="task_params.work_dir"
							rules={{ required: "work_dir is required" }}
							render={({ field }) => (
								<FormItem>
									<FormLabel> {t("table.columns.schedule.work_dir")}</FormLabel>
									<FormControl>
										<Input {...field} placeholder="/home/user" />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
					</div>
				);
			default:
				return (
					<FormField
						control={form.control}
						name="task_params"
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.schedule.task_params")}</FormLabel>
								<FormControl>
									<Input.TextArea
										rows={4}
										onChange={(e) => field.onChange(JSON.parse(e.target.value || "{}"))}
										value={JSON.stringify(field.value, null, 2)}
										placeholder='{"key": "value"}'
									/>
								</FormControl>
							</FormItem>
						)}
					/>
				);
		}
	};

	return (
		<Modal
			width={600}
			open={show}
			title={title}
			onOk={handleOk}
			onCancel={handleCancel}
			centered
			styles={{
				body: {
					maxHeight: "80vh",
					overflowY: "auto",
				},
			}}
			classNames={{
				body: "themed-scrollbar",
			}}
			footer={[
				<Button key="back" onClick={handleCancel}>
					{t("table.button.return")}
				</Button>,
				<Button key="submit" type="primary" loading={loading} onClick={handleOk}>
					{t("table.button.submit")}
				</Button>,
			]}
		>
			<Form {...form}>
				<form className="space-y-4">
					<FormField
						control={form.control}
						name="task_name"
						rules={{
							required: useTranslationRule(t("table.columns.schedule.task_name")),
						}}
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.schedule.task_name")}</FormLabel>
								<FormControl>
									<Input {...field} />
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						name="task_description"
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.schedule.task_description")}</FormLabel>
								<FormControl>
									<Input {...field} />
								</FormControl>
							</FormItem>
						)}
					/>

					<AdvancedCronField />

					<FormField
						control={form.control}
						name="exec_type"
						rules={{
							required: useTranslationRule(t("table.columns.schedule.exec_type")),
						}}
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.schedule.exec_type")}</FormLabel>
								<FormControl>
									<Select
										style={{ width: 150 }}
										onChange={(value: string) => {
											field.onChange(value);
										}}
										value={field.value}
										options={taskExecType}
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						name="task_type"
						rules={{
							required: useTranslationRule(t("table.columns.schedule.task_type")),
						}}
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.schedule.task_type")}</FormLabel>
								<FormControl>
									<Select
										style={{ width: 150 }}
										onChange={(value: string) => {
											field.onChange(value);
											form.setValue("task_params", {});
										}}
										value={String(field.value)}
										options={taskTypes}
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>

					{renderTaskParams()}
				</form>
			</Form>
		</Modal>
	);
}
