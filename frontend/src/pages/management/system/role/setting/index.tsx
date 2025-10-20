import { Icon } from "@/components/icon";
import { useRoleSettingActions, useRoleSettingApiIds, useRoleSettingMenuIds } from "@/store/roleSettingStore";
import type { MenuBtn, Role } from "@/types/entity";
import { Button } from "@/ui/button";
import { Card, Modal, Tabs } from "antd";
import type { TreeDataNode } from "antd";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import ApiSetting from "./modal-api";
import MenuSetting from "./modal-menu";
import MenuBtnSetting from "./modal-menu-btn";

export const getAllKeys = (data: TreeDataNode[]): React.Key[] => {
	return data.reduce((acc, node) => {
		acc.push(node.key);
		if (node.children) {
			acc.push(...getAllKeys(node.children));
		}
		return acc;
	}, [] as React.Key[]);
};

export type SettingModalProps = {
	roleData: Role;
	id: number;
	title: string;
	show: boolean;
	onCancel: VoidFunction;
};

export type selectMenuData = {
	menuId: number;
	menuBtns: MenuBtn[];
	isSet: boolean;
};

export default function SettingModal({ roleData, title, show, onCancel }: SettingModalProps) {
	const { t } = useTranslation();
	const menuGroupIds = useRoleSettingMenuIds();
	const ApiIds = useRoleSettingApiIds();
	const { fetch } = useRoleSettingActions();
	const [selectMenuBtn, setSelectMenuBtn] = useState<selectMenuData>({
		menuId: 0,
		menuBtns: [],
		isSet: false,
	});
	useEffect(() => {
		if (show && roleData.id) {
			fetch(roleData.id);
		}
	}, [roleData.id, fetch, show]);
	return (
		<Modal
			title={title}
			closable={{ "aria-label": "Custom Close Button" }}
			open={show}
			onCancel={() => {
				onCancel();
				// clear data
				setSelectMenuBtn({
					menuId: 0,
					menuBtns: [],
					isSet: false,
				});
			}}
			styles={{
				body: {
					maxHeight: "80vh",
					overflowY: "auto",
				},
			}}
			classNames={{
				body: "themed-scrollbar",
			}}
			footer={null}
			width={800}
			centered
		>
			<Tabs
				hidden={selectMenuBtn.isSet}
				defaultActiveKey="1"
				items={[
					{
						key: "1",
						label: t("table.button.role_menu"),
						children: (
							<div className="max-h-[600px] overflow-y-auto">
								<MenuSetting
									id={roleData.id}
									defaultRoleRouter={roleData.default_router}
									menuGroupIds={menuGroupIds}
									setSelectMenuBtn={setSelectMenuBtn}
								/>
							</div>
						),
					},
					{
						key: "2",
						label: t("table.button.role_api"),
						children: (
							<div className="max-h-[600px] overflow-y-auto">
								<ApiSetting id={roleData.id} apiIds={ApiIds} />
							</div>
						),
					},
				]}
				onChange={(key: string) => {
					console.log(key);
				}}
			/>
			<Card
				hidden={!selectMenuBtn.isSet}
				title={
					<Button
						variant="ghost"
						size="sm"
						className="flex items-center gap-2 px-3 py-2 hover:bg-accent rounded-md transition-colors"
						onClick={() =>
							setSelectMenuBtn((prev) => ({
								...prev,
								isSet: false,
							}))
						}
					>
						<Icon icon="solar:alt-arrow-left-outline" className="text-3xl" />
						<span className="text-sm">{t("table.button.go_back")}</span>
					</Button>
				}
			>
				<MenuBtnSetting selectMenuBtn={selectMenuBtn} roleId={roleData.id} />
			</Card>
		</Modal>
	);
}
