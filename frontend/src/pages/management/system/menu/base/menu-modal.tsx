import type { Menu, MenuTree } from "@/types/entity";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/ui/form";
import { IconPicker } from "@/ui/icon-picker";
import { Input } from "@/ui/input";

import { useTranslationRule } from "@/hooks";
import useDirTree from "@/hooks/dirTree";
import useLangTree from "@/hooks/langTree";
import { Button, Cascader, Modal, Switch, TreeSelect } from "antd";
import type { TFunction } from "i18next";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";

export type MenuModalProps = {
	formValue: Menu;
	treeRawData: Menu[];
	title: string;
	show: boolean;
	onOk: (values: Menu) => Promise<boolean>;
	onCancel: VoidFunction;
};
// 构建树形结构
export const buildTree = (tree: Menu[], t: TFunction<"translation", undefined>, disabledId?: string): MenuTree[] => {
	return tree.map((item: Menu): MenuTree => {
		return {
			value: item.id.toString(),
			title: t(item.title),
			key: item.id.toString(),
			path: item.level,
			origin: item,
			disabled: item.id.toString() === disabledId,
			children: item.children ? buildTree(item.children, t) : [],
		};
	});
};

const MenuNewModal = ({ title, show, treeRawData, formValue, onOk, onCancel }: MenuModalProps) => {
	const { t, i18n } = useTranslation();

	const [loading, setLoading] = useState(false);
	const [treeData, setTreeData] = useState<MenuTree[]>([]);
	const dirTree = useDirTree();
	const [isManualInput, setIsManualInput] = useState(true);
	const [isManualTitleInput, setIsManualTitleInput] = useState(true);
	const langTree = useLangTree(i18n.store.data[i18n.language].translation);
	const form = useForm<Menu>({
		defaultValues: formValue,
	});

	useEffect(() => {
		form.reset(formValue);
		const currentId = formValue.id.toString();
		setTreeData([
			{
				value: "0",
				title: t("table.columns.common.root_node"),
				key: "0",
				path: [0],
				children: buildTree(treeRawData, t, currentId),
			},
		]);
	}, [formValue, treeRawData, form, t]);

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
		<>
			<Modal
				width={600}
				open={show}
				title={title}
				onOk={handleOk}
				styles={{
					body: {
						maxHeight: "80vh",
						overflowY: "auto",
					},
				}}
				classNames={{
					body: "themed-scrollbar",
				}}
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
							name="component"
							rules={{
								required: useTranslationRule(t("table.columns.menu.component")),
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel className="flex items-center justify-between">
										<span>{t("table.columns.menu.component")}</span>
										<Button
											type="link"
											size="small"
											onClick={() => {
												setIsManualInput(!isManualInput);
												field.onChange(undefined);
											}}
										>
											{isManualInput
												? t("table.button.switch_convenient_selection")
												: t("table.button.switch_manual_input")}
										</Button>
									</FormLabel>
									<FormControl>
										{isManualInput ? (
											<Input
												{...field}
												placeholder={t("table.handle_message.file_path_placeholder")}
												value={field.value || ""}
											/>
										) : (
											<div className="flex gap-2">
												<Cascader
													style={{ flex: 1 }}
													fieldNames={{
														label: "title",
														children: "children",
													}}
													value={handleValue(field.value)}
													options={dirTree}
													onChange={(value) => {
														if (Array.isArray(value)) {
															field.onChange(value[value.length - 1]);
														} else {
															field.onChange(value);
														}
													}}
													placeholder="please select file path"
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
							name="title"
							rules={{
								required: useTranslationRule(t("table.columns.menu.title")),
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel className="flex items-center justify-between">
										<span>{t("table.columns.menu.title")}</span>
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
														if (Array.isArray(value)) {
															field.onChange(value[value.length - 1]);
														} else {
															field.onChange(value);
														}
													}}
													placeholder="please select title"
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
							name="name"
							rules={{
								required: useTranslationRule(t("table.columns.menu.name")),
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.menu.name")}</FormLabel>
									<FormControl>
										<Input {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="path"
							rules={{
								required: useTranslationRule(t("table.columns.menu.path")),
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.menu.path")}</FormLabel>
									<FormControl>
										<Input {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="hidden"
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.menu.hidden")}</FormLabel>
									<FormControl>
										<div className="w-fit">
											<Switch checked={field.value} onChange={(value) => field.onChange(value)} />
										</div>
									</FormControl>
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="keep_alive"
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.menu.keep_alive")}</FormLabel>
									<FormControl>
										<div className="w-fit">
											<Switch
												checked={field.value === 1}
												onChange={(value) => field.onChange(value === true ? 1 : 0)}
											/>
										</div>
									</FormControl>
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="parent_id"
							rules={{
								required: useTranslationRule(t("table.columns.menu.parent_id")),
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.menu.parent_id")}</FormLabel>
									<FormControl>
										<TreeSelect
											showSearch
											style={{ width: "100%" }}
											value={String(field.value)}
											styles={{
												popup: { root: { maxHeight: 400, overflow: "auto" } },
											}}
											placeholder="Please select"
											allowClear
											onChange={(value) => {
												field.onChange(Number(value));
											}}
											treeData={treeData}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="icon"
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.menu.icon")}</FormLabel>
									<FormControl>
										<IconPicker value={field.value} onChange={field.onChange} />
									</FormControl>
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="sort"
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.menu.sort")}</FormLabel>
									<FormControl>
										<Input
											{...field}
											type="number"
											onChange={(e) => {
												field.onChange(Number(e.target.value));
											}}
										/>
									</FormControl>
								</FormItem>
							)}
						/>
					</form>
				</Form>
			</Modal>
		</>
	);
};

export default MenuNewModal;
