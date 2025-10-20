import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/ui/form";
import { Input } from "@/ui/input";

import { useTranslationRule } from "@/hooks";
import { useDictionaryByTypeWithCache } from "@/hooks/dict";
import { Button, Modal, Radio, Switch } from "antd";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import type { Dictionary } from "#/entity";

export type DictionaryModalProps = {
	formValue: Dictionary;
	title: string;
	show: boolean;
	onOk: (values: Dictionary) => Promise<boolean>;
	onCancel: VoidFunction;
};

export default function UserModal({ title, show, formValue, onOk, onCancel }: DictionaryModalProps) {
	const form = useForm<Dictionary>({
		defaultValues: formValue,
	});
	const { data: status } = useDictionaryByTypeWithCache("status");
	const { t } = useTranslation();
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
			styles={{
				body: {
					maxHeight: "80vh",
					overflowY: "auto",
				},
			}}
			classNames={{
				body: "themed-scrollbar",
			}}
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
						name="name"
						rules={{
							required: useTranslationRule(t("table.columns.dictionary.name")),
						}}
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.dictionary.name")}</FormLabel>
								<FormControl>
									<Input {...field} />
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						name="type"
						rules={{
							required: useTranslationRule(t("table.columns.dictionary.type")),
						}}
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.dictionary.type")}</FormLabel>
								<FormControl>
									<Input {...field} />
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						name="is_generate_file"
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.dictionary.is_generate_file")}</FormLabel>
								<FormControl>
									<div className="w-fit">
										<Switch checked={Boolean(field.value)} onChange={(value) => field.onChange(Number(value))} />
									</div>
								</FormControl>
							</FormItem>
						)}
					/>

					<FormField
						control={form.control}
						name="status"
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.dictionary.status")}</FormLabel>
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
						name="desc"
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.dictionary.desc")}</FormLabel>
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
