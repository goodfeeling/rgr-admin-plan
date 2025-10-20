import { Icon } from "@/components/icon";
import { useMenuBtn, useMenuBtnActions } from "@/store/menuBtnStore";
import { useMenuParameter, useMenuParameterActions } from "@/store/menuParameterStore";
import { Badge } from "@/ui/badge";
import { Button } from "@/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/ui/select";
import type { GetRef, InputRef, TableProps } from "antd";
import { Form, Input, InputNumber, Modal, Table, Tabs } from "antd";
import { createContext, useContext, useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import type { MenuBtn, MenuParameter } from "#/entity";
export type SettingType = {
	id: number;
};

export type SettingModalProps = {
	formValue: SettingType;
	title: string;
	show: boolean;
	onCancel: VoidFunction;
};

export default function SettingModal({ title, show, formValue, onCancel }: SettingModalProps) {
	const { id } = formValue;
	const { t } = useTranslation();

	const [open, setOpen] = useState(false);
	const handleCancel = () => {
		setOpen(false);
		onCancel();
	};
	useEffect(() => {
		setOpen(show);
	}, [show]);
	return (
		<Modal
			width={800}
			open={open}
			title={title}
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
			footer={false}
		>
			<Tabs
				defaultActiveKey="1"
				items={[
					{
						key: "1",
						label: t("sys.menu.system.menu_btn"),
						children: (
							<div className="max-h-[600px] overflow-y-auto">
								<BtnPage MenuId={id} />
							</div>
						),
					},
					{
						key: "2",
						label: t("sys.menu.system.menu_parameter"),
						children: (
							<div className="max-h-[600px] overflow-y-auto">
								<ParameterPage MenuId={id} />
							</div>
						),
					},
				]}
			/>
		</Modal>
	);
}

// 定义 cellType 类型
type CellType = "input" | "select" | "number";
type FormInstance<T> = GetRef<typeof Form<T>>;
const EditableContext = createContext<FormInstance<any> | null>(null);

interface EditableRowProps {
	index: number;
}

const EditableRow: React.FC<EditableRowProps> = ({ index, ...props }) => {
	const [form] = Form.useForm();
	return (
		<Form form={form} component={false}>
			<EditableContext.Provider value={form}>
				<tr {...props} />
			</EditableContext.Provider>
		</Form>
	);
};

interface EditableCellMenuBtnProps {
	title: React.ReactNode;
	editable: boolean;
	dataIndex: keyof MenuParameter;
	record: MenuParameter;
	handleSave: (record: MenuParameter) => void;
	cellType: CellType;
}

function ParameterPage({ MenuId }: { MenuId: number }) {
	const { t } = useTranslation();
	const { updateOrCreateMenu, deleteMenu, fetchMenu } = useMenuParameterActions();
	const menuParameterData = useMenuParameter();

	const EditableCell: React.FC<React.PropsWithChildren<EditableCellMenuBtnProps>> = ({
		cellType,
		title,
		editable,
		children,
		dataIndex,
		record,
		handleSave,
		...restProps
	}) => {
		const [editing, setEditing] = useState(false);
		const inputRef = useRef<InputRef>(null);
		// biome-ignore lint/style/noNonNullAssertion: <explanation>
		const form = useContext(EditableContext)!;

		useEffect(() => {
			if (editing) {
				inputRef.current?.focus();
			}
		}, [editing]);

		const toggleEdit = () => {
			setEditing(!editing);
			form.setFieldsValue({ [dataIndex]: record[dataIndex] });
		};

		const save = async () => {
			try {
				const values = await form.validateFields();
				toggleEdit();
				handleSave({ ...record, ...values });
			} catch (errInfo) {
				console.log("Save failed:", errInfo);
			}
		};
		// biome-ignore lint/suspicious/noImplicitAnyLet: <explanation>
		let inputNode;

		// 根据 type 字段的值来选择不同的输入组件
		switch (cellType) {
			case "select":
				inputNode = (
					<Select
						onValueChange={(value) => {
							form.setFieldsValue({ [dataIndex]: value });
							save();
						}}
						value={form.getFieldValue(dataIndex)}
					>
						<SelectTrigger>
							<SelectValue placeholder="Select Status" />
						</SelectTrigger>
						<SelectContent>
							<SelectItem value="params">
								<Badge variant="success">params</Badge>
							</SelectItem>
							<SelectItem value="query">
								<Badge variant="error">query</Badge>
							</SelectItem>
						</SelectContent>
					</Select>
				);
				break;
			case "number":
				inputNode = <InputNumber style={{ width: "100%" }} onPressEnter={save} onBlur={save} />;
				break;
			default:
				inputNode = <Input ref={inputRef} onPressEnter={save} onBlur={save} />;
		}

		let childNode = children;

		if (editable) {
			childNode = editing ? (
				<Form.Item
					style={{ margin: 0 }}
					name={dataIndex}
					rules={[{ required: true, message: `${title} is required.` }]}
				>
					{inputNode}
				</Form.Item>
			) : (
				<div className="editable-cell-value-wrap" style={{ paddingInlineEnd: 24 }} onClick={toggleEdit}>
					{children}
				</div>
			);
		}

		return <td {...restProps}>{childNode}</td>;
	};
	type ColumnTypes = Exclude<TableProps<MenuParameter>["columns"], undefined>;
	const [dataSource, setDataSource] = useState<MenuParameter[]>([]);

	useEffect(() => {
		fetchMenu(MenuId);
	}, [fetchMenu, MenuId]);

	useEffect(() => {
		setDataSource(menuParameterData);
	}, [menuParameterData]);
	const handleDelete = (id: number | undefined) => {
		if (!id) return;
		deleteMenu(id);
	};

	const defaultColumns: (ColumnTypes[number] & {
		editable?: boolean;
		dataIndex: string;
		cellType?: CellType;
	})[] = [
		{
			title: t("table.columns.menu_parameter.type"),
			dataIndex: "type",
			width: "30%",
			editable: true,
			cellType: "select",
		},
		{
			title: t("table.columns.menu_parameter.key"),
			dataIndex: "key",
			editable: true,
		},
		{
			title: t("table.columns.menu_parameter.value"),
			dataIndex: "value",
			editable: true,
		},
		{
			title: t("table.columns.common.operation"),
			dataIndex: "operation",
			render: (_, record) =>
				dataSource.length >= 1 ? (
					<Button variant="linkwarning" size="icon" onClick={() => handleDelete(record.id)}>
						<Icon icon="mingcute:delete-2-fill" size={18} className="text-error!" />
						<span>{t("table.button.delete")}</span>
					</Button>
				) : null,
		},
	];

	const handleAdd = () => {
		const newData: MenuParameter = {
			key: "new",
			type: "params",
			value: "string",
			sys_base_menu_id: MenuId,
		};
		updateOrCreateMenu(newData);
	};

	const handleSave = (row: MenuParameter) => {
		const newData = [...dataSource];
		const index = newData.findIndex((item) => row.id === item.id);
		const item = newData[index];
		newData.splice(index, 1, {
			...item,
			...row,
		});
		updateOrCreateMenu(row);
		setDataSource(newData);
	};

	const components = {
		body: {
			row: EditableRow,
			cell: EditableCell,
		},
	};

	const columns = defaultColumns.map((col) => {
		if (!col.editable) {
			return col;
		}
		return {
			...col,
			onCell: (record: MenuParameter) => ({
				record,
				editable: col.editable,
				dataIndex: col.dataIndex,
				title: col.title,
				handleSave,
				cellType: col.cellType || "input",
			}),
		};
	});
	return (
		<div className="flex flex-col gap-4">
			<div className="flex items-start justify-start">
				<Button onClick={handleAdd} type="button">
					<Icon icon="solar:add-circle-outline" size={18} />
					{t("table.button.add_a_row")}
				</Button>
			</div>
			<Table<MenuParameter>
				rowKey={(record) => record.id as number}
				components={components}
				rowClassName={() => "editable-row"}
				bordered
				dataSource={dataSource}
				columns={columns as ColumnTypes}
				pagination={false}
			/>
		</div>
	);
}
interface EditableCellMenuParameterProps {
	title: React.ReactNode;
	editable: boolean;
	dataIndex: keyof MenuParameter;
	record: MenuParameter;
	handleSave: (record: MenuParameter) => void;
}
function BtnPage({ MenuId }: { MenuId: number }) {
	const { t } = useTranslation();
	const { updateOrCreateMenu, deleteMenu, fetchMenu } = useMenuBtnActions();
	const menuBtnData = useMenuBtn();

	const EditableCell: React.FC<React.PropsWithChildren<EditableCellMenuParameterProps>> = ({
		title,
		editable,
		children,
		dataIndex,
		record,
		handleSave,
		...restProps
	}) => {
		const [editing, setEditing] = useState(false);
		const inputRef = useRef<InputRef>(null);
		// biome-ignore lint/style/noNonNullAssertion: <explanation>
		const form = useContext(EditableContext)!;

		useEffect(() => {
			if (editing) {
				inputRef.current?.focus();
			}
		}, [editing]);

		const toggleEdit = () => {
			setEditing(!editing);
			form.setFieldsValue({ [dataIndex]: record[dataIndex] });
		};

		const save = async () => {
			try {
				const values = await form.validateFields();

				toggleEdit();
				handleSave({ ...record, ...values });
			} catch (errInfo) {
				console.log("Save failed:", errInfo);
			}
		};
		let childNode = children;

		if (editable) {
			childNode = editing ? (
				<Form.Item
					style={{ margin: 0 }}
					name={dataIndex}
					rules={[{ required: true, message: `${title} is required.` }]}
				>
					<Input ref={inputRef} defaultValue={record?.value?.toString()} onPressEnter={save} onBlur={save} />
				</Form.Item>
			) : (
				<div className="editable-cell-value-wrap" style={{ paddingInlineEnd: 24 }} onClick={toggleEdit}>
					{children}
				</div>
			);
		}

		return <td {...restProps}>{childNode}</td>;
	};
	type ColumnTypes = Exclude<TableProps<MenuBtn>["columns"], undefined>;
	const [dataSource, setDataSource] = useState<MenuBtn[]>([]);

	useEffect(() => {
		fetchMenu(MenuId);
	}, [fetchMenu, MenuId]);

	useEffect(() => {
		setDataSource(menuBtnData);
	}, [menuBtnData]);
	const handleDelete = (id: number | undefined) => {
		if (!id) return;
		deleteMenu(id);
	};

	const defaultColumns: (ColumnTypes[number] & {
		editable?: boolean;
		dataIndex: string;
		cellType?: CellType;
	})[] = [
		{
			title: t("table.columns.menu_btn.name"),
			dataIndex: "name",
			width: "30%",
			editable: true,
		},
		{
			title: t("table.columns.menu_btn.desc"),
			dataIndex: "desc",
			editable: true,
		},
		{
			title: t("table.columns.common.operation"),
			dataIndex: "operation",
			render: (_, record) =>
				dataSource.length >= 1 ? (
					<Button variant="linkwarning" size="icon" onClick={() => handleDelete(record.id)}>
						<Icon icon="mingcute:delete-2-fill" size={18} className="text-error!" />
						<span>{t("table.button.delete")}</span>
					</Button>
				) : null,
		},
	];

	const handleAdd = () => {
		const newData: MenuBtn = {
			name: "New data",
			desc: "New data",
			sys_base_menu_id: MenuId,
		};
		updateOrCreateMenu(newData);
	};

	const handleSave = (row: MenuBtn) => {
		const newData = [...dataSource];
		const index = newData.findIndex((item) => row.id === item.id);
		const item = newData[index];
		newData.splice(index, 1, {
			...item,
			...row,
		});
		updateOrCreateMenu(row);
	};

	const components = {
		body: {
			row: EditableRow,
			cell: EditableCell,
		},
	};

	const columns = defaultColumns.map((col) => {
		if (!col.editable) {
			return col;
		}
		return {
			...col,
			onCell: (record: MenuBtn) => ({
				record,
				editable: col.editable,
				dataIndex: col.dataIndex,
				title: col.title,
				handleSave,
			}),
		};
	});
	return (
		<div className="flex flex-col gap-4">
			<div className="flex items-start justify-start">
				<Button onClick={handleAdd} variant="default">
					<Icon icon="solar:add-circle-outline" size={18} />
					{t("table.button.add_a_row")}
				</Button>
			</div>
			<Table<MenuBtn>
				rowKey={(record) => record.id as number}
				components={components}
				rowClassName={() => "editable-row"}
				bordered
				dataSource={dataSource}
				columns={columns as ColumnTypes}
				pagination={false}
			/>
		</div>
	);
}
