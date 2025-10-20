import UploadTool from "@/components/upload/upload-multiple";
import type { DictionaryDetail, FileInfo } from "@/types/entity";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/ui/form";
import { parseUriFromUrl } from "@/utils";
import { Modal, Select } from "antd";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";

export type FileModalProps = {
	formValue: FileInfo;
	title: string;
	storageEngine: DictionaryDetail[] | undefined;
	show: boolean;
	onOk: (values: FileInfo | null) => Promise<boolean>;
	onCancel: VoidFunction;
};
const FileNewModal = ({ title, show, formValue, storageEngine, onOk, onCancel }: FileModalProps) => {
	const { t } = useTranslation();
	const [selectedStorageEngine, setSelectedStorageEngine] = useState<string>("local");
	const form = useForm<FileInfo>({
		defaultValues: formValue,
	});
	useEffect(() => {
		form.reset(formValue);
	}, [formValue, form]);

	const handleCancel = () => {
		onCancel();
	};

	return (
		<>
			<Modal open={show} title={title} onCancel={handleCancel} centered footer={false}>
				<Form {...form}>
					<form className="space-y-4">
						<FormField
							control={form.control}
							name="storage_engine"
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t("table.columns.file.file_origin_name")}</FormLabel>
									<FormControl>
										<Select
											defaultValue="local"
											style={{ width: 120 }}
											onChange={(value: string) => {
												field.onChange(value);
												setSelectedStorageEngine(value);
											}}
											options={storageEngine}
										/>
									</FormControl>
								</FormItem>
							)}
						/>
						<UploadTool
							onHandleSuccess={(result) => {
								if (result.url) {
									form.setValue("file_url", result.url ?? "");
									form.setValue("file_name", result.name ?? "");
									form.setValue("file_origin_name", result.name || "");
									if (result.url) {
										const uri = parseUriFromUrl(result.url);
										form.setValue("file_path", uri);
									}
									onOk(form.getValues());
								}
							}}
							listType="text"
							renderType="button"
							uploadType={selectedStorageEngine}
						/>
					</form>
				</Form>
			</Modal>
		</>
	);
};

export default FileNewModal;
