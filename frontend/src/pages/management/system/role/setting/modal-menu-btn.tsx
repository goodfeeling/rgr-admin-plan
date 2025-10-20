import { useRoleSettingActions, useRoleSettingBtnIds } from "@/store/roleSettingStore";
import type { MenuBtn } from "@/types/entity";
import { Table } from "antd";
import type { TableColumnsType, TableProps } from "antd";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import type { selectMenuData } from "./index";

type MenuBtnSettingProps = {
	selectMenuBtn: selectMenuData;
	roleId: number;
};
const MenuBtnSetting = ({ selectMenuBtn, roleId }: MenuBtnSettingProps) => {
	const { t } = useTranslation();
	const menuBtnIds = useRoleSettingBtnIds();
	const { updateRoleBtns } = useRoleSettingActions();
	const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
	const [menuBtnData, setMenuBtnData] = useState<MenuBtn[]>([]);
	const columns: TableColumnsType<MenuBtn> = [
		{
			title: t("table.columns.menu_btn.name"),
			dataIndex: "name",
		},
		{
			title: t("table.columns.menu_btn.desc"),
			dataIndex: "desc",
		},
	];

	useEffect(() => {
		setMenuBtnData(selectMenuBtn.menuBtns);
		setSelectedRowKeys(menuBtnIds[selectMenuBtn.menuId]);
	}, [menuBtnIds, selectMenuBtn]);

	const rowSelection: TableProps<MenuBtn>["rowSelection"] = {
		selectedRowKeys,
		type: "checkbox",
		onChange: (selectedRowKeys) => {
			setSelectedRowKeys(selectedRowKeys);
			updateRoleBtns(roleId, selectMenuBtn.menuId, selectedRowKeys as number[]);
		},
	};
	return (
		<>
			<Table<MenuBtn>
				rowKey={"id"}
				rowSelection={rowSelection}
				columns={columns}
				dataSource={menuBtnData}
				pagination={false}
				onChange={(value) => {
					console.log("selectedRowKeys changed: ", value);
				}}
			/>
		</>
	);
};

export default MenuBtnSetting;
