import { buildTree } from "@/pages/management/system/role/modal";
import { useMenuActions } from "@/store/useMenuStore";
import { useUserActions, useUserInfo } from "@/store/userStore";
import { DownOutlined } from "@ant-design/icons";
import { Modal, Tree, type TreeProps } from "antd";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

export type SwitchModalProps = {
	title: string;
	show: boolean;
	onCancel: VoidFunction;
};
const SwitchModal = ({ show, onCancel }: SwitchModalProps) => {
	const { switchRole } = useUserActions();
	const { t } = useTranslation();
	const menuActions = useMenuActions();
	const { roles, current_role } = useUserInfo();
	const [selectedKeys, setSelectedKeys] = useState<string[]>([]);

	// 初始化默认选中 current_role
	useEffect(() => {
		if (current_role) {
			setSelectedKeys([String(current_role.id)]);
		}
	}, [current_role]);
	const handleCancel = () => {
		onCancel();
	};
	const onSelect: TreeProps["onSelect"] = async (keys) => {
		const selectedKey = keys[0]; // 单选，取第一个 key

		try {
			// 调用接口更新角色
			// 示例：await updateRole(selectedKey);
			console.log("调用接口更新角色:", selectedKey);

			// 接口调用成功后更新 store 中的 current_role
			setSelectedKeys([selectedKey as string]);
			// 等待角色切换完成
			await switchRole(selectedKey as number);
			// 确保菜单已更新
			await menuActions.fetchMenu();
		} catch (error) {
			console.error("更新角色失败:", error);
		}
	};

	return (
		<>
			<Modal
				open={show}
				title={`${t("sys.menu.current_role")}: ${current_role ? current_role.name : t("common.none")}`}
				onCancel={handleCancel}
				footer={null}
			>
				<Tree
					showLine
					switcherIcon={<DownOutlined />}
					defaultExpandAll={true}
					onSelect={onSelect}
					selectedKeys={selectedKeys} // 控制当前选中项
					treeData={buildTree(roles ? roles : [])}
				/>
			</Modal>
		</>
	);
};

export default SwitchModal;
