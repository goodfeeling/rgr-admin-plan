import { UploadService } from "@/api/services/uploadService";
import userService from "@/api/services/userService";
import { UploadAvatar } from "@/components/upload";
import { useTranslationRule } from "@/hooks";
import { useDictionaryByTypeWithCache } from "@/hooks/dict";
import RoleSelect from "@/pages/components/role-select/RoleSelect";
import useUserStore from "@/store/userStore";
import type { UserInfo } from "@/types/entity";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/ui/form";
import { Input } from "@/ui/input";
import { Button, Modal, Radio } from "antd";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";

export type UserModalProps = {
	formValue: UserInfo;
	treeData: any[];
	title: string;
	show: boolean;
	onOk: (values: UserInfo) => Promise<boolean>;
	onCancel: VoidFunction;
};

const UserNewModal = ({ title, show, formValue, treeData, onOk, onCancel }: UserModalProps) => {
	const { t } = useTranslation();
	const { data: status } = useDictionaryByTypeWithCache("status");
	const [loading, setLoading] = useState(false);
	const { userToken } = useUserStore.getState();
	const form = useForm<UserInfo>({
		defaultValues: formValue,
	});

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
		<>
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
							name="header_img"
							rules={{
								required: useTranslationRule(t("table.columns.user.avatar")),
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.user.avatar")}</FormLabel>
									<FormControl>
										<UploadAvatar
											defaultAvatar={field.value}
											onHeaderImgChange={(fileUrl: string) => {
												form.setValue("header_img", fileUrl);
											}}
											action={`${import.meta.env.VITE_APP_BASE_API}${UploadService.Client.Single}`}
											headers={{
												Authorization: `Bearer ${userToken?.accessToken}`,
											}}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="user_name"
							rules={{
								required: useTranslationRule(t("table.columns.user.user_name")),
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.user.user_name")}</FormLabel>
									<FormControl>
										<Input {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="email"
							rules={{
								required: useTranslationRule(t("table.columns.user.email")),
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.user.email")}</FormLabel>
									<FormControl>
										<Input {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="nick_name"
							rules={{
								required: useTranslationRule(t("table.columns.user.nick_name")),
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.user.nick_name")}</FormLabel>
									<FormControl>
										<Input {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="roles"
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.user.roles")}</FormLabel>
									<FormControl>
										<RoleSelect
											roles={field.value ?? []}
											treeData={treeData}
											recordKey={`modal_${form.getValues().id}`}
											onChange={async (values) => {
												try {
													await userService.bindRole(form.getValues().id, values);
													console.log("update success");
												} catch (error) {
													console.error("update error:", error);
												}
											}}
										/>
									</FormControl>
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="phone"
							rules={{
								required: useTranslationRule(t("table.columns.user.phone")),
							}}
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.user.phone")}</FormLabel>
									<FormControl>
										<Input {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="status"
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.user.status")}</FormLabel>
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

export default UserNewModal;
