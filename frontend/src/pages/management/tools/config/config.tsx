import configService from "@/api/services/configService";
import UploadTool from "@/components/upload/upload-multiple";
import { useDictionaryByTypeWithCache } from "@/hooks";
import { themeVars } from "@/theme/theme.css";
import { BasicStatus } from "@/types/enum";
import { Button, Form, Input, InputNumber, Radio, Select, Tag } from "antd";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

const ConfigBox = ({
	configData,
	module,
}: {
	configData: {
		[key: string]: any;
	};
	module: string;
}) => {
	const { t } = useTranslation();
	const [formData, setFormData] = useState<{ [key: string]: any }>(configData);
	const [loading, setLoading] = useState(false);
	const [form] = Form.useForm();
	const { data: status } = useDictionaryByTypeWithCache("status");

	const { data: storageEngine } = useDictionaryByTypeWithCache("file_storage_engine");
	const { data: database } = useDictionaryByTypeWithCache("database_type");
	const { data: eventBus } = useDictionaryByTypeWithCache("event_bus");
	const { data: logLevel } = useDictionaryByTypeWithCache("log_level");
	const { data: encodeLevel } = useDictionaryByTypeWithCache("encode_level");
	const { data: outputType } = useDictionaryByTypeWithCache("output_type");
	const { data: apiMethod } = useDictionaryByTypeWithCache("api_method");
	// define select options
	const configSelectMap: Record<string, Record<string, any[] | undefined> | undefined> = {
		server: {
			storage_engine: storageEngine,
			database: database,
			event_bus: eventBus,
		},
		zap: {
			level: logLevel,
			encode_level: encodeLevel,
			encoding: outputType,
		},
	};

	const handleSave = (key: string) => {
		form
			.validateFields()
			.then(async (values) => {
				try {
					setLoading(true);
					await configService.updateConfig(values, key);
					toast.success(`${key} 配置已保存，下次重启服务后生效`);
					setFormData({
						...formData,
						data: {
							...formData.data,
							[key]: values,
						},
					});
				} catch (error) {
					console.error("Save failed:", error);
					toast.error("保存配置失败");
				} finally {
					setLoading(false);
				}
			})
			.catch((error) => {
				console.error("Validate Failed:", error);
				toast.error("表单验证失败");
			});
	};

	const render = (sectionKey: string, configKey: string, configValue: string) => {
		if (sectionKey === "site" && ["logo", "favicon", "login_img"].includes(configKey)) {
			return (
				<UploadTool
					onHandleSuccess={(result) => {
						if (result.url) {
							form.setFieldValue(configKey, result.url);
							handleSave(sectionKey);
						}
					}}
					listType="text"
					renderType="image"
					showUploadList={false}
					renderImageUrl={configValue}
				/>
			);
		}

		// 检查是否需要渲染Select组件
		const selectOptions = configSelectMap[sectionKey]?.[configKey];
		if (selectOptions) {
			return <Select options={selectOptions} />;
		}

		// allowed_methods
		if (configKey === "allowed_methods") {
			return (
				<Select
					mode="multiple"
					style={{ width: "100%" }}
					tagRender={(props) => {
						const { label, value, closable, onClose } = props;
						const onPreventMouseDown = (event: React.MouseEvent<HTMLSpanElement>) => {
							event.preventDefault();
							event.stopPropagation();
						};
						return (
							<Tag
								color={value}
								onMouseDown={onPreventMouseDown}
								closable={closable}
								onClose={onClose}
								style={{
									marginInlineEnd: 4,
									color: themeVars.colors.text.primary,
								}}
							>
								{label}
							</Tag>
						);
					}}
					placeholder="Please select"
					options={apiMethod}
				/>
			);
		}

		return typeof configValue === "number" ? (
			<InputNumber style={{ width: "100%" }} />
		) : typeof configValue === "boolean" ? (
			<Radio.Group
				onChange={(e) => {
					form.setFieldValue(configKey, e.target.value);
				}}
				value={configValue}
			>
				{status?.map((item) => (
					<Radio.Button key={item.value} value={Number(item.value) === BasicStatus.ENABLE}>
						{item.label}
					</Radio.Button>
				))}
			</Radio.Group>
		) : typeof configValue === "object" ? (
			<Input.TextArea rows={4} defaultValue={JSON.stringify(configValue, null, 2)} />
		) : (
			<Input defaultValue={String(configValue)} />
		);
	};

	return (
		<Form form={form} layout="vertical" initialValues={configData} onFinish={() => handleSave(module)}>
			<div className="p-4">
				<div className="grid grid-cols-1 md:grid-cols-2 gap-4">
					{Object.entries(configData || {}).map(([configKey, configValue]) => (
						<Form.Item
							key={configKey}
							label={t(`sys.config.columns.${module}.${configKey}`, configKey)}
							name={configKey}
							rules={[{ required: true, message: `请输入 ${configKey}` }]}
						>
							{render(module, configKey, configValue)}
						</Form.Item>
					))}
				</div>
				<div className="flex justify-end mt-6 gap-5">
					<Button type="primary" onClick={() => handleSave(module)} loading={loading}>
						{`${t("common.saveText")}${t("common.configuration")}`}
					</Button>
				</div>
			</div>
		</Form>
	);
};

export default ConfigBox;
