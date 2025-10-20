import menuService from "@/api/services/menuService";
import { useRoleSettingActions } from "@/store/roleSettingStore";
import type { MenuBtn, MenuTree, MenuTreeUserGroup } from "@/types/entity";
import { Button } from "@/ui/button";
import { List, Tag, Tree } from "antd";
import type { TreeProps } from "antd";
import { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { buildTree } from "../../menu/base/menu-modal";
import { getAllKeys, type selectMenuData } from "./index";

type MenuSettingProps = {
	id: number;
	defaultRoleRouter: string;
	menuGroupIds: { [key: string]: number[] };
	setSelectMenuBtn: React.Dispatch<React.SetStateAction<selectMenuData>>;
};

const MenuSetting = ({ id, defaultRoleRouter, menuGroupIds, setSelectMenuBtn }: MenuSettingProps) => {
	const { t } = useTranslation();
	const [roleMenuData, setRoleMenuData] = useState<{ [key: string]: any }>([]);
	const [groupData, setGroupData] = useState<MenuTreeUserGroup[]>([]);
	const [defaultRouter, setDefaultRouter] = useState<string>("");
	// 加载菜单树
	const onLoadMenuTree = useCallback(async () => {
		const response = await menuService.getUserMenu(true);
		setGroupData(response);
	}, []);

	useEffect(() => {
		setRoleMenuData(menuGroupIds);
	}, [menuGroupIds]);

	useEffect(() => {
		onLoadMenuTree();
	}, [onLoadMenuTree]);

	// 默认路由
	useEffect(() => {
		setDefaultRouter(defaultRoleRouter);
	}, [defaultRoleRouter]);

	return (
		<div>
			<List
				itemLayout="vertical"
				dataSource={groupData}
				renderItem={(item) => {
					const treeData = buildTree(item.items, t);
					return (
						<List.Item>
							<List.Item.Meta title={t(item.name || "default")} />
							<div>
								<TreeList
									id={id}
									groupId={item.id || 0}
									defaultRouter={defaultRouter}
									treeData={treeData}
									checkKeys={roleMenuData}
									setSelectMenuBtn={setSelectMenuBtn}
									setDefaultRouter={setDefaultRouter}
								/>
							</div>
						</List.Item>
					);
				}}
			/>
		</div>
	);
};

type checkedKeys = {
	checked: number[];
	halfChecked: number[];
};

type TreeListProps = {
	id: number;
	defaultRouter: string;
	treeData: MenuTree[];
	checkKeys: { [key: string]: any };
	groupId: number;
	setSelectMenuBtn: (selectMenuData: selectMenuData) => void;
	setDefaultRouter: React.Dispatch<React.SetStateAction<string>>;
};
const TreeList = ({
	id,
	defaultRouter,
	treeData,
	checkKeys,
	groupId,
	setSelectMenuBtn,
	setDefaultRouter,
}: TreeListProps) => {
	const { t } = useTranslation();
	const { updateMenus, updateRouterPath } = useRoleSettingActions();
	const [checkedKeys, setCheckedKeys] = useState<React.Key[]>([]);
	const [expandedKeys, setExpandedKeys] = useState<React.Key[]>([]);

	useEffect(() => {
		const tempData = checkKeys[groupId];
		if (tempData && tempData.length > 0) {
			setCheckedKeys(tempData.map((item: number) => String(item)));
		}
	}, [checkKeys, groupId]);

	// 初始化展开所有节点
	useEffect(() => {
		setExpandedKeys(getAllKeys(treeData));
	}, [treeData]);

	// 组件卸载时清理状态
	useEffect(() => {
		return () => {
			setCheckedKeys([]);
			setExpandedKeys([]);
		};
	}, []);

	const onCheck: TreeProps["onCheck"] = (checkedKeysValue) => {
		console.log("onCheck", checkedKeysValue);
		const temp: checkedKeys = checkedKeysValue as checkedKeys;
		setCheckedKeys(temp.checked);
		updateMenus(id, String(groupId), temp.checked);
	};

	// 处理节点展开/收起
	const onExpand: TreeProps["onExpand"] = (expandedKeysValue) => {
		setExpandedKeys(expandedKeysValue);
	};

	// 更新默认路由
	const updateDefaultRouter = (data: MenuTree) => {
		const routerPath = data.origin ? data.origin.path : "";
		setDefaultRouter(routerPath);
		updateRouterPath(id, routerPath);
	};
	// 配置可控按钮
	const settingRoleBtn = (data: MenuTree, menuBtns: MenuBtn[] | undefined) => {
		setSelectMenuBtn({
			menuId: data.origin ? data.origin.id : 0,
			menuBtns: menuBtns ? menuBtns : [],
			isSet: true,
		});
	};
	return (
		<Tree
			checkable
			selectable={false}
			checkStrictly={true}
			expandedKeys={expandedKeys}
			onExpand={onExpand}
			onCheck={onCheck}
			checkedKeys={checkedKeys}
			treeData={treeData}
			multiple={true}
			className="no-hover-tree"
			titleRender={(node) => {
				const hasBtn = node.origin?.menu_btns && node.origin.menu_btns.length > 0;
				return (
					<span className="flex justify-between items-center w-full">
						<span>{node.title}</span>
						{defaultRouter === node.origin?.path ? (
							<Tag color="blue">{t("table.button.home_page")}</Tag>
						) : (
							<Button
								variant="link"
								onClick={(e) => {
									e.stopPropagation();
									updateDefaultRouter(node);
								}}
							>
								{t("table.button.set_home_page")}
							</Button>
						)}
						<Button
							variant="ghost"
							onClick={(e) => {
								e.stopPropagation();
								settingRoleBtn(node, node.origin?.menu_btns);
							}}
							hidden={!hasBtn}
						>
							{t("table.button.set_role_button")}
						</Button>
					</span>
				);
			}}
		/>
	);
};

export default MenuSetting;
