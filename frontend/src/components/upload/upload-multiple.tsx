import { UploadService } from "@/api/services/uploadService";
import { useOssUpload } from "@/hooks/ossUpload";
import { useRemoveFileInfoMutation } from "@/store/fileManageStore";
import { useSTSTokenLoading } from "@/store/stsTokenStore";
import useUserStore from "@/store/userStore";
import { LoadingOutlined } from "@ant-design/icons";
import type { UploadProps } from "antd";
import { App, Upload } from "antd";
import type { UploadListType } from "antd/es/upload/interface";
import type { UploadFile } from "antd/lib";
import type React from "react";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import "./upload-multiple.css";
import { Icon } from "@/components/icon";
import { Button } from "@/ui/button";
type ResultFile = {
	name?: string;
	url?: string;
};

interface UploadToolProps {
	uploadType?: string;
	onHandleSuccess?: (file: ResultFile) => void;
	listType: UploadListType | undefined;
	fileList?: UploadFile<any>[] | undefined;
	showUploadList?: boolean;
	renderType?: "button" | "image";
	renderImageUrl?: string;
	accept?: string;
	title?: string;
	uploadUri?: string;
}

const UploadTool: React.FC<Readonly<UploadToolProps>> = ({
	uploadType,
	onHandleSuccess,
	listType,
	showUploadList,
	fileList,
	renderType,
	renderImageUrl,
	title,
	accept,
	uploadUri,
}) => {
	const { t } = useTranslation();
	const { message } = App.useApp();
	const { userToken } = useUserStore.getState();
	const { uploadFile } = useOssUpload();
	const isOssLoading = useSTSTokenLoading();
	const removeMutation = useRemoveFileInfoMutation();
	const [imageUrl, setImageUrl] = useState<string>();
	const [loading, setLoading] = useState(false);

	useEffect(() => {
		setImageUrl(renderImageUrl);
	}, [renderImageUrl]);

	const props: UploadProps = {
		name: "file",
		listType: listType,
		showUploadList: showUploadList,
		defaultFileList: fileList,
		accept: accept ?? ".xlsx, .xls, image/*, .pdf, doc, .docx",
		action: `${import.meta.env.VITE_APP_BASE_API}${uploadUri ?? UploadService.Client.Multiple}`,
		headers: {
			Authorization: `Bearer ${userToken?.accessToken}`,
		},
		onChange(info) {
			const { status } = info.file;
			switch (status) {
				case "uploading":
					console.log(info.file, info.fileList);
					break;
				case "done":
					message.success(`${info.file.name} ${t("table.handle_message.upload_success")}`);
					if (onHandleSuccess) {
						onHandleSuccess({
							url: info.file.response.data[0].file_url,
							name: info.file.name,
						});
					}
					setImageUrl(info.file.response.data[0].file_url);
					break;
				case "error":
					message.error(`${info.file.name} ${t("table.handle_message.upload_error")}`);

					break;
				default:
					console.log(info);
			}
		},
		// 根据存储引擎类型决定是否使用自定义上传
		beforeUpload: async (file) => {
			if (uploadType === "aliyunoss") {
				setLoading(true);
				// 使用阿里云OSS上传
				const result = await uploadFile(file);
				setLoading(false);
				if (result.success) {
					// 上传成功后更新表单数据
					setImageUrl(result.url);
					if (onHandleSuccess) {
						onHandleSuccess({
							url: result.url,
							name: result.name || "",
						});
					}
					message.success(`${file.name} ${t("table.handle_message.upload_success")}`);
				} else {
					message.error(`${file.name} ${t("table.handle_message.upload_error")}`);
				}
				return false;
			}
			return true;
		},
		onRemove: async (file) => {
			removeMutation.mutate(file.response.data[0].id, {
				onSuccess: () => {
					message.success(t("table.handle_message.delete_success"));
				},
				onError: () => {
					message.error(t("table.handle_message.error"));
				},
			});
		},
		disabled: !isOssLoading,
	};

	const render = () => {
		if (renderType === "button") {
			return (
				<Button className="ml-2 text-white" variant="default">
					<Icon icon="solar:export-outline" size={18} />

					{title ?? t("sys.menu.upload")}
				</Button>
			);
		}
		return imageUrl ? (
			<div className="upload-image-container">
				<img
					src={imageUrl}
					alt="avatar"
					style={{
						width: "100%",
						height: "150px",
						objectFit: "cover",
						borderRadius: "6px",
						display: "block",
					}}
				/>
				<div className="upload-image-overlay">
					<Icon icon="solar:cloud-upload-broken" size={18} />
					{title ?? t("sys.menu.upload")}
				</div>
			</div>
		) : (
			<button style={{ border: 0, background: "none" }} type="button">
				{loading ? <LoadingOutlined /> : <Icon icon="solar:cloud-upload-broken" size={18} />}
				<div style={{ marginTop: 8 }}> {title ?? t("sys.menu.upload")}</div>
			</button>
		);
	};

	return <Upload {...props}>{render()}</Upload>;
};

export default UploadTool;
