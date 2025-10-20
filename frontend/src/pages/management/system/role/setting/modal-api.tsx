import apiService from "@/api/services/apisService";
import { useRoleSettingActions } from "@/store/roleSettingStore";
import type { ApiGroupItem } from "@/types/entity";
import { Input, List, Tag, Tree } from "antd";
import type { TreeProps } from "antd";
import { useCallback, useEffect, useRef, useState } from "react";
const { Search } = Input;

type ApiSettingProps = {
	id: number;
	apiIds: string[];
};
const methodColorMap: Record<string, string> = {
	GET: "var(--colors-palette-primary-default)", // 蓝色
	POST: "var(--colors-palette-success-default)", // 绿色
	PUT: "var(--colors-palette-warning-default)", // 橙色
	DELETE: "var(--colors-palette-error-default)", // 红色
	PATCH: "var(--colors-palette-info-default)", // 青色
	OPTIONS: "var(--colors-palette-gray-500)", // 灰色
	HEAD: "var(--colors-palette-gray-700)", // 深灰色
};
const ApiSetting = ({ id, apiIds }: ApiSettingProps) => {
	const { updateApis } = useRoleSettingActions();
	const [checkedKeys, setCheckedKeys] = useState<React.Key[]>([]);
	const [treeData, setTreeData] = useState<ApiGroupItem[]>([]);
	const searchTimeoutRef = useRef<NodeJS.Timeout | null>(null);
	useEffect(() => {
		setCheckedKeys(
			apiIds.filter((item) => {
				return item.startsWith("0---") === false;
			}),
		);
	}, [apiIds]);

	// 加载菜单树
	const onLoadMenuTree = useCallback(async (path?: string) => {
		const response = await apiService.getApiGroupList(path);
		setTreeData(response);
	}, []);

	useEffect(() => {
		onLoadMenuTree("");
	}, [onLoadMenuTree]);

	const onCheck: TreeProps["onCheck"] = (checkedKeysValue) => {
		console.log("onCheck", checkedKeysValue);
		setCheckedKeys(checkedKeysValue as React.Key[]);
		updateApis(
			id,
			(checkedKeysValue as string[]).filter((item) => {
				return item.startsWith("0---") === false;
			}),
		);
	};

	const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
		const searchValue = e.target.value;
		if (searchTimeoutRef.current) {
			clearTimeout(searchTimeoutRef.current);
		}
		searchTimeoutRef.current = setTimeout(() => {
			onLoadMenuTree(searchValue);
		}, 300); // 300ms 防抖延迟
	};

	return (
		<div>
			<Search style={{ marginBottom: 8 }} placeholder="Search" onChange={handleSearch} />
			<Tree
				checkable
				selectable={false}
				onCheck={onCheck}
				checkedKeys={checkedKeys}
				treeData={treeData}
				multiple={true}
				className="no-hover-tree"
				titleRender={(node) => {
					const keys = node.key.split("---");
					return (
						<List.Item>
							<List.Item.Meta
								title={
									<span style={{ color: "var(--primary)" }}>
										{node.children === null ? <Tag color={methodColorMap[keys[1]]}>{keys[1]}</Tag> : ""}
										{node.title}
									</span>
								}
								description={node.children === null ? keys[0] : ""}
							/>
						</List.Item>
					);
				}}
			/>
		</div>
	);
};

export default ApiSetting;
