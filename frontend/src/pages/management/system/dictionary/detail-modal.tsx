import { Card, CardContent } from "@/ui/card";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/ui/form";
import { Input } from "@/ui/input";

import UploadTool from "@/components/upload/upload-multiple";
import { useTranslationRule } from "@/hooks";
import { useDictionaryByTypeWithCache } from "@/hooks/dict";
import { Button, Modal, Radio, Select, Switch } from "antd";
import { type ReactNode, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import type { DictionaryDetail } from "#/entity";

export type DictionaryDetailModalProps = {
	formValue: DictionaryDetail;
	title: string;
	show: boolean;
	onOk: (values: DictionaryDetail) => Promise<boolean>;
	onCancel: VoidFunction;
};

export default function UserModal({ title, show, formValue, onOk, onCancel }: DictionaryDetailModalProps) {
	const form = useForm<DictionaryDetail>();
	const { data: status } = useDictionaryByTypeWithCache("status");
	const { t } = useTranslation();
	const [loading, setLoading] = useState(false);

	useEffect(() => {
		form.reset(formValue);
	}, [formValue, form]);

	const handleOk = async () => {
		form.handleSubmit(async (values) => {
			values.sort = Number(values.sort);
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

	return (
		<Modal
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
						name="label"
						rules={{
							required: useTranslationRule(t("table.columns.dictionary_detail.label")),
						}}
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.dictionary_detail.label")}</FormLabel>
								<FormControl>
									<Input {...field} />
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						rules={{
							required: useTranslationRule(t("table.columns.dictionary_detail.type")),
						}}
						name="type"
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.dictionary_detail.type")}</FormLabel>
								<FormControl>
									<Select
										style={{ width: 120 }}
										onChange={(value: string) => {
											field.onChange(value);
											form.setValue("value", "");
										}}
										value={field.value}
										options={[
											{ value: "string", label: "字符串" },
											{ value: "number", label: "数字" },
											{ value: "image", label: "图片" },
											{ value: "icon", label: "图标" },
										]}
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						rules={{
							required: useTranslationRule(t("table.columns.dictionary_detail.value")),
						}}
						name="value"
						render={({ field }) => {
							let result: ReactNode;

							switch (form.watch("type") || form.getValues().type) {
								case "string":
									result = <Input {...field} />;
									break;
								case "number":
									result = <Input type="number" {...field} />;
									break;
								case "boolean":
									result = (
										<div className="w-fit">
											<Switch checked={Boolean(field.value)} onChange={(value) => field.onChange(value)} />
										</div>
									);
									break;

								case "image":
									result = (
										<Card>
											<CardContent>
												<UploadTool
													onHandleSuccess={(result) => {
														if (result.url) {
															field.onChange(result.url);
														}
													}}
													listType="picture-card"
													renderType="image"
													showUploadList={false}
													renderImageUrl={field.value}
												/>
											</CardContent>
										</Card>
									);
									break;
								default:
									result = <Input {...field} />;
							}
							return (
								<FormItem>
									<FormLabel>{t("table.columns.dictionary_detail.value")}</FormLabel>
									<FormControl>{result}</FormControl>
									<FormMessage />
								</FormItem>
							);
						}}
					/>
					<FormField
						control={form.control}
						name="extend"
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.dictionary_detail.extend")}</FormLabel>
								<FormControl>
									<Input {...field} />
								</FormControl>
							</FormItem>
						)}
					/>

					<FormField
						control={form.control}
						name="status"
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.dictionary_detail.status")}</FormLabel>
								<FormControl>
									<Radio.Group
										onChange={(e) => {
											field.onChange(Number(e.target.value));
										}}
										value={String(field.value)}
									>
										{status?.map((item) => (
											<Radio.Button key={item.value} value={String(item.value)}>
												{item.label}
											</Radio.Button>
										))}
									</Radio.Group>
								</FormControl>
							</FormItem>
						)}
					/>

					<FormField
						control={form.control}
						name="sort"
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.dictionary_detail.sort")}</FormLabel>
								<FormControl>
									<Input type="number" {...field} />
								</FormControl>
							</FormItem>
						)}
					/>
				</form>
			</Form>
		</Modal>
	);
}
