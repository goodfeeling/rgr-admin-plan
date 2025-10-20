import { TreeSelect } from "antd";
import { useEffect, useState } from "react";
import type { Role } from "#/entity";

interface RoleSelectProps {
	roles: Role[];
	treeData: any[];
	recordKey: string;
	onChange?: (values: string[]) => void;
}

const RoleSelect: React.FC<RoleSelectProps> = ({ roles, treeData, recordKey, onChange }) => {
	const [selectedRoleIds, setSelectedRoleIds] = useState<string[]>([]);

	useEffect(() => {
		const ids = roles.map((x: Role) => String(x.id));
		setSelectedRoleIds(ids);
	}, [roles]);

	const handleSelectChange = (values: string[]) => {
		setSelectedRoleIds(values);
		if (onChange) {
			onChange(values);
		}
	};

	return (
		<TreeSelect
			key={recordKey}
			treeData={treeData}
			value={selectedRoleIds}
			onChange={handleSelectChange}
			treeCheckable={true}
			showCheckedStrategy={TreeSelect.SHOW_ALL}
			placeholder="place select role"
			style={{ width: "100%" }}
		/>
	);
};

export default RoleSelect;
