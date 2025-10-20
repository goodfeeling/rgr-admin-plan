import { DndContext, KeyboardSensor, PointerSensor, closestCenter, useSensor, useSensors } from "@dnd-kit/core";
import {
	SortableContext,
	arrayMove,
	sortableKeyboardCoordinates,
	useSortable,
	verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { Button, Card, Checkbox, Col, Form, Input, InputNumber, Radio, Row, Select, Space, Typography } from "antd";
import { useState } from "react";
import { useTranslation } from "react-i18next";

const { Text } = Typography;
const { Option } = Select;

// 定义表单字段类型
interface FormField {
	id: string;
	type: "input" | "textarea" | "select" | "checkbox" | "radio" | "number";
	label: string;
	placeholder?: string;
	required?: boolean;
	options?: { id: string; label: string; value: string }[];
	defaultValue?: any;
}

// 定义表单配置
interface FormConfig {
	title: string;
	description: string;
	fields: FormField[];
}

// 可排序的表单项组件
const SortableItem = (props: {
	field: FormField;
	onClick: (field: FormField) => void;
}) => {
	const { field, onClick } = props;
	const [selectedField, _] = useState<FormField | null>(null);

	const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({ id: field.id });

	const style = {
		transform: CSS.Transform.toString(transform),
		transition,
	};

	const getFieldLabel = (type: FormField["type"]): string => {
		const labels: Record<FormField["type"], string> = {
			input: "单行文本",
			textarea: "多行文本",
			number: "数字输入",
			select: "下拉选择",
			checkbox: "复选框",
			radio: "单选框",
		};
		return labels[type];
	};

	return (
		<div
			ref={setNodeRef}
			style={style}
			{...attributes}
			onClick={() => onClick(field)}
			className={`p-3 mb-2 border rounded cursor-pointer ${
				selectedField?.id === field.id ? "border-blue-500 bg-blue-50" : "border-gray-300 hover:border-blue-300"
			} ${isDragging ? "shadow-lg" : ""}`}
		>
			<div className="flex justify-between items-center" {...listeners}>
				<Text strong>{field.label}</Text>
				<Text type="secondary">{getFieldLabel(field.type)}</Text>
			</div>
			<Text type="secondary" className="text-xs">
				{field.placeholder}
			</Text>
		</div>
	);
};

const FormBuilder: React.FC = () => {
	const { t } = useTranslation();
	const [formConfig, setFormConfig] = useState<FormConfig>({
		title: "表单标题",
		description: "表单描述",
		fields: [],
	});

	const [selectedField, setSelectedField] = useState<FormField | null>(null);
	const [form] = Form.useForm();

	const sensors = useSensors(
		useSensor(PointerSensor),
		useSensor(KeyboardSensor, {
			coordinateGetter: sortableKeyboardCoordinates,
		}),
	);

	// 可用的表单组件
	const componentPalette = [
		{ type: "input", label: "单行文本" },
		{ type: "textarea", label: "多行文本" },
		{ type: "number", label: "数字输入框" },
		{ type: "select", label: "下拉选择" },
		{ type: "checkbox", label: "复选框" },
		{ type: "radio", label: "单选框" },
	];

	// 添加字段到表单
	const addField = (type: FormField["type"]) => {
		const newField: FormField = {
			id: `field_${Date.now()}`,
			type,
			label: `${getFieldLabel(type)} ${formConfig.fields.length + 1}`,
			placeholder: `请输入${getFieldLabel(type)}`,
			required: false,
		};

		// 为选择类组件添加默认选项
		if (type === "select" || type === "checkbox" || type === "radio") {
			newField.options = [
				{ id: `option_${Date.now()}_1`, label: "选项1", value: "option1" },
				{ id: `option_${Date.now()}_2`, label: "选项2", value: "option2" },
			];
		}

		setFormConfig({
			...formConfig,
			fields: [...formConfig.fields, newField],
		});
	};

	// 获取字段显示名称
	const getFieldLabel = (type: FormField["type"]): string => {
		const labels: Record<FormField["type"], string> = {
			input: "单行文本",
			textarea: "多行文本",
			number: "数字输入",
			select: "下拉选择",
			checkbox: "复选框",
			radio: "单选框",
		};
		return labels[type];
	};

	// 处理拖拽结束事件
	const handleDragEnd = (event: any) => {
		const { active, over } = event;

		if (active.id !== over.id) {
			setFormConfig((config) => {
				const oldIndex = config.fields.findIndex((field) => field.id === active.id);
				const newIndex = config.fields.findIndex((field) => field.id === over.id);
				return {
					...config,
					fields: arrayMove(config.fields, oldIndex, newIndex),
				};
			});
		}
	};

	// 更新字段属性
	const updateField = (id: string, updates: Partial<FormField>) => {
		setFormConfig({
			...formConfig,
			fields: formConfig.fields.map((field) => (field.id === id ? { ...field, ...updates } : field)),
		});
	};

	// 删除字段
	const deleteField = (id: string) => {
		setFormConfig({
			...formConfig,
			fields: formConfig.fields.filter((field) => field.id !== id),
		});
		if (selectedField?.id === id) {
			setSelectedField(null);
		}
	};

	// 保存表单配置
	const saveForm = () => {
		console.log("保存表单配置:", formConfig);
		// 这里可以添加保存到数据库的逻辑
	};

	// 渲染表单字段编辑器
	const renderFieldEditor = () => {
		if (!selectedField) {
			return <div className="text-center p-4 text-gray-500">请选择一个字段进行编辑</div>;
		}

		return (
			<Form layout="vertical">
				<Form.Item label="字段标签">
					<Input
						value={selectedField.label}
						onChange={(e) => updateField(selectedField.id, { label: e.target.value })}
					/>
				</Form.Item>

				<Form.Item label="占位符">
					<Input
						value={selectedField.placeholder}
						onChange={(e) => updateField(selectedField.id, { placeholder: e.target.value })}
					/>
				</Form.Item>

				<Form.Item label="是否必填">
					<Checkbox
						checked={selectedField.required}
						onChange={(e) => updateField(selectedField.id, { required: e.target.checked })}
					>
						必填字段
					</Checkbox>
				</Form.Item>

				{(selectedField.type === "select" || selectedField.type === "checkbox" || selectedField.type === "radio") && (
					<Form.Item label="选项设置">
						<Space direction="vertical" style={{ width: "100%" }}>
							{selectedField.options?.map((option, index) => (
								<Space key={option.id}>
									<Input
										placeholder="选项标签"
										value={option.label}
										onChange={(e) => {
											const newOptions = [...(selectedField.options || [])];
											newOptions[index].label = e.target.value;
											updateField(selectedField.id, { options: newOptions });
										}}
									/>
									<Input
										placeholder="选项值"
										value={option.value}
										onChange={(e) => {
											const newOptions = [...(selectedField.options || [])];
											newOptions[index].value = e.target.value;
											updateField(selectedField.id, { options: newOptions });
										}}
									/>
									<Button
										danger
										onClick={() => {
											const newOptions = [...(selectedField.options || [])];
											newOptions.splice(index, 1);
											updateField(selectedField.id, { options: newOptions });
										}}
									>
										删除
									</Button>
								</Space>
							))}
							<Button
								onClick={() => {
									const newOptions = [
										...(selectedField.options || []),
										{
											id: `option_${Date.now()}_${selectedField.options?.length || 0}`,
											label: "新选项",
											value: `option${(selectedField.options?.length || 0) + 1}`,
										},
									];
									updateField(selectedField.id, { options: newOptions });
								}}
							>
								添加选项
							</Button>
						</Space>
					</Form.Item>
				)}

				<Button danger block onClick={() => deleteField(selectedField.id)}>
					删除字段
				</Button>
			</Form>
		);
	};

	// 渲染表单预览
	const renderFormPreview = () => {
		return (
			<Card title={formConfig.title} size="small">
				{formConfig.description && <p className="text-gray-600 mb-4">{formConfig.description}</p>}
				<Form layout="vertical" form={form}>
					{formConfig.fields.map((field) => (
						<Form.Item
							key={field.id}
							label={field.label}
							required={field.required}
							rules={[{ required: field.required, message: `请输入${field.label}` }]}
						>
							{field.type === "input" && <Input placeholder={field.placeholder} />}

							{field.type === "textarea" && <Input.TextArea placeholder={field.placeholder} rows={4} />}

							{field.type === "number" && <InputNumber placeholder={field.placeholder} style={{ width: "100%" }} />}

							{field.type === "select" && (
								<Select placeholder={field.placeholder} style={{ width: "100%" }}>
									{field.options?.map((option) => (
										<Option key={option.value} value={option.value}>
											{option.label}
										</Option>
									))}
								</Select>
							)}

							{field.type === "checkbox" && (
								<Checkbox.Group>
									<Space direction="vertical">
										{field.options?.map((option) => (
											<Checkbox key={option.value} value={option.value}>
												{option.label}
											</Checkbox>
										))}
									</Space>
								</Checkbox.Group>
							)}

							{field.type === "radio" && (
								<Radio.Group>
									<Space direction="vertical">
										{field.options?.map((option) => (
											<Radio key={option.value} value={option.value}>
												{option.label}
											</Radio>
										))}
									</Space>
								</Radio.Group>
							)}
						</Form.Item>
					))}
					<Form.Item>
						<Space>
							<Button type="primary">提交</Button>
							<Button>重置</Button>
						</Space>
					</Form.Item>
				</Form>
			</Card>
		);
	};

	return (
		<Card title={t("sys.menu.system.from_create")} size="small">
			<div className="p-4">
				<Row gutter={16}>
					{/* 左侧：组件库 */}
					<Col span={6}>
						<Card title="组件库" size="small">
							<Space direction="vertical" style={{ width: "100%" }}>
								{componentPalette.map((component) => (
									<Button key={component.type} block onClick={() => addField(component.type as FormField["type"])}>
										{component.label}
									</Button>
								))}
							</Space>
						</Card>

						<Card title="表单设置" size="small" className="mt-4">
							<Form layout="vertical">
								<Form.Item label="表单标题">
									<Input
										value={formConfig.title}
										onChange={(e) => setFormConfig({ ...formConfig, title: e.target.value })}
									/>
								</Form.Item>
								<Form.Item label="表单描述">
									<Input.TextArea
										value={formConfig.description}
										onChange={(e) =>
											setFormConfig({
												...formConfig,
												description: e.target.value,
											})
										}
										rows={3}
									/>
								</Form.Item>
							</Form>
						</Card>
					</Col>

					{/* 中间：表单设计区域 */}
					<Col span={12}>
						<Card
							title="表单设计"
							size="small"
							extra={
								<Space>
									<Button onClick={saveForm} type="primary">
										保存表单
									</Button>
								</Space>
							}
						>
							<DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
								<SortableContext
									items={formConfig.fields.map((field) => field.id)}
									strategy={verticalListSortingStrategy}
								>
									<div className="min-h-[400px]">
										{formConfig.fields.length === 0 ? (
											<div className="text-center text-gray-400 py-12">
												从左侧拖拽组件到此处，或点击组件添加到表单中
											</div>
										) : (
											formConfig.fields.map((field) => (
												<SortableItem key={field.id} field={field} onClick={setSelectedField} />
											))
										)}
									</div>
								</SortableContext>
							</DndContext>
						</Card>

						<Card title="表单预览" size="small" className="mt-4">
							{renderFormPreview()}
						</Card>
					</Col>

					{/* 右侧：字段属性编辑器 */}
					<Col span={6}>
						<Card title="字段属性" size="small">
							{renderFieldEditor()}
						</Card>
					</Col>
				</Row>
			</div>
		</Card>
	);
};

export default FormBuilder;
