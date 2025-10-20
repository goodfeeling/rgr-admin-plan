import { useTranslationRule } from "@/hooks";
import { useDictionaryByTypeWithCache } from "@/hooks/dict";
import useLangTree from "@/hooks/langTree";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/ui/form";
import { Input } from "@/ui/input";
import { Button, Cascader, Modal, Radio } from "antd";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import type { MenuGroup } from "#/entity";

export type MenuGroupModalProps = {
	formValue: MenuGroup;
	title: string;
	show: boolean;
	onOk: (values: MenuGroup) => Promise<boolean>;
	onCancel: VoidFunction;
};

export default function UserModal({ title, show, formValue, onOk, onCancel }: MenuGroupModalProps) {
	const { t, i18n } = useTranslation();
	const { data: status } = useDictionaryByTypeWithCache("status");

	const form = useForm<MenuGroup>({
		defaultValues: formValue,
	});
	const [loading, setLoading] = useState(false);
	const [isManualTitleInput, setIsManualTitleInput] = useState(true);
	const langTree = useLangTree(i18n.store.data[i18n.language].translation);

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
	const handleValue = (value: any) => {
		if (Array.isArray(value)) {
			return value;
		}
		if (typeof value === "string") {
			return value.split("/");
		}
		return undefined;
	};

	return (
		<Modal
			width={400}
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
						name="name"
						rules={{
							required: useTranslationRule(t("table.columns.menu_group.name")),
						}}
						render={({ field }) => (
							<FormItem>
								<FormLabel className="flex items-center justify-between">
									<span>{t("table.columns.menu_group.name")}</span>
									<Button
										type="link"
										size="small"
										onClick={() => {
											setIsManualTitleInput(!isManualTitleInput);
											field.onChange(undefined);
										}}
									>
										{isManualTitleInput
											? t("table.button.switch_convenient_selection")
											: t("table.button.switch_manual_input")}
									</Button>
								</FormLabel>
								<FormControl>
									{isManualTitleInput ? (
										<Input
											{...field}
											placeholder={t("table.handle_message.title_placeholder")}
											value={field.value || ""}
										/>
									) : (
										<div className="flex gap-2">
											<Cascader
												style={{ flex: 1 }}
												value={handleValue(field.value)}
												options={langTree}
												onChange={(value) => {
													// 将数组形式的路径值转换为字符串
													if (Array.isArray(value)) {
														field.onChange(value[value.length - 1]);
													} else {
														field.onChange(value);
													}
												}}
												placeholder="please select name"
												popupMenuColumnStyle={{
													width: "200px",
													whiteSpace: "normal",
												}}
												showSearch
											/>
										</div>
									)}
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						name="path"
						rules={{
							required: useTranslationRule(t("table.columns.menu_group.path")),
						}}
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.menu_group.path")}</FormLabel>
								<FormControl>
									<Input {...field} />
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						name="sort"
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("table.columns.menu_group.sort")}</FormLabel>
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
								<FormLabel>{t("table.columns.menu_group.status")}</FormLabel>
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
				</form>
			</Form>
		</Modal>
	);
}
