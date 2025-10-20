import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/ui/form";
import { Input } from "@/ui/input";

import { useTranslationRule } from "@/hooks";
import { Button, Modal, Select } from "antd";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import type { Api, DictionaryDetail } from "#/entity";

export type ApiModalProps = {
	formValue: Api;
	apiGroup: DictionaryDetail[] | undefined;
	apiMethod: DictionaryDetail[] | undefined;
	title: string;
	show: boolean;
	onOk: (values: Api) => Promise<boolean>;
	onCancel: VoidFunction;
};

export default function ApiModal({ title, show, formValue, apiGroup, apiMethod, onOk, onCancel }: ApiModalProps) {
	const { t } = useTranslation();
	const form = useForm<Api>({
		defaultValues: formValue,
	});
	const [loading, setLoading] = useState(false);

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

	return (
		<Modal
			open={show}
			title={title}
			onOk={handleOk}
			onCancel={handleCancel}
			centered
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
						name="path"
						rules={{
							required: useTranslationRule(t("table.columns.api.path")),
						}}
						render={({ field }) => (
							<FormItem>
								<FormLabel> {t("table.columns.api.path")}</FormLabel>
								<FormControl>
									<Input {...field} />
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>

					<FormField
						control={form.control}
						name="method"
						rules={{
							required: useTranslationRule(t("table.columns.api.method")),
						}}
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.api.method")}</FormLabel>

								<Select
									style={{ width: 150 }}
									onChange={(value: string) => {
										field.onChange(value);
									}}
									value={field.value}
									options={apiMethod}
								/>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						name="api_group"
						rules={{
							required: useTranslationRule(t("table.columns.api.api_group")),
						}}
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.api.api_group")}</FormLabel>

								<Select
									style={{ width: 150 }}
									onChange={(value: string) => {
										field.onChange(value);
									}}
									value={field.value}
									options={apiGroup}
								/>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
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
				</form>
			</Form>
		</Modal>
	);
}
