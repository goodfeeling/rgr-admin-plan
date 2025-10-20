import { useTranslationRule } from "@/hooks";
import { useDictionaryByTypeWithCache } from "@/hooks/dict";
import type { Role, RoleTree } from "@/types/entity";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/ui/form";
import { Input } from "@/ui/input";
import { Button, Modal, Radio, TreeSelect } from "antd";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";

export type RoleModalProps = {
	formValue: Role;
	treeRawData: Role[];
	title: string;
	show: boolean;
	onOk: (values: Role) => Promise<boolean>;
	onCancel: VoidFunction;
};
export function buildTree(tree: Role[]): RoleTree[] {
	return tree.map((item: Role): RoleTree => {
		return {
			value: item.id.toString(),
			title: item.name,
			key: item.id.toString(),
			path: item.path,
			children: item.children ? buildTree(item.children) : [],
		};
	});
}
const RoleNewModal = ({ title, show, treeRawData, formValue, onOk, onCancel }: RoleModalProps) => {
	const { t } = useTranslation();
	const { data: status } = useDictionaryByTypeWithCache("status");
	const [loading, setLoading] = useState(false);
	const [treeData, setTreeData] = useState<RoleTree[]>([]);
	const form = useForm<Role>({
		defaultValues: formValue,
	});

	useEffect(() => {
		form.reset(formValue);
		setTreeData([
			{
				value: "0",
				title: t("table.columns.common.root_node"),
				key: "0",
				path: [0],
				children: buildTree(treeRawData),
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

	return (
		<>
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
							name="name"
							rules={{
								required: useTranslationRule(t("table.columns.role.name")),
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.role.name")}</FormLabel>
									<FormControl>
										<Input {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="parent_id"
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.role.parent_id")}</FormLabel>
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
											treeDefaultExpandAll
											onChange={(value) => {
												field.onChange(Number(value));
											}}
											treeData={treeData}
										/>
									</FormControl>
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="label"
							rules={{
								required: useTranslationRule(t("table.columns.role.label")),
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.role.label")}</FormLabel>
									<FormControl>
										<Input {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="description"
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.role.description")}</FormLabel>
									<FormControl>
										<Input {...field} />
									</FormControl>
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="order"
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.role.order")}</FormLabel>
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
									<FormLabel>{t("table.columns.role.status")}</FormLabel>
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
		</>
	);
};

export default RoleNewModal;
