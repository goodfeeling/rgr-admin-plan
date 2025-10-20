import type React from "react";
import { useEffect, useRef, useState } from "react";
import { Tree } from "antd";
import { useSettings } from "@/store/settingStore";
import { paletteColors } from "@/theme/tokens/color";
export interface TreeNode {
  value: string;
  title: string;
  key: string;
  children?: TreeNode[];
  path: number[];
  isLast?: boolean;
}

interface TreeSelectInputProps {
  treeData: TreeNode[];
  value?: string;
  disabled?: boolean;
  onChange?: (value: string) => void;
  placeholder?: string;
}

const TreeSelectInput: React.FC<TreeSelectInputProps> = ({
  treeData,
  value,
  disabled,
  onChange,
  placeholder = "请选择",
}) => {
  const { themeMode } = useSettings();
  const [selectedKey, setSelectedKey] = useState<string>(value || "");
  const [selectedLabel, setSelectedLabel] = useState<string>("");
  const [isTreeVisible, setIsTreeVisible] = useState<boolean>(false);
  const [selectedPath, setSelectedPath] = useState<number[]>([]);

  // 用 ref 获取外层容器用于判断是否点击外部
  const containerRef = useRef<HTMLDivElement>(null);

  const filterTreeNode = (node: any): boolean => {
    if (selectedPath === null) {
      return false;
    }
    return node.value === "0"
      ? true
      : selectedPath.includes(Number(node.value));
  };

  // 查找默认选中的节点 label
  useEffect(() => {
    if (value && treeData.length > 0) {
      const findNode = (
        nodes: TreeNode[],
        targetValue: string
      ): TreeNode | null => {
        for (const node of nodes) {
          if (node.value === targetValue) return node;
          if (node.children) {
            const found = findNode(node.children, targetValue);
            if (found) return found;
          }
        }
        return null;
      };

      const matchedNode = findNode(treeData, value);
      if (matchedNode) {
        setSelectedLabel(matchedNode.title);
        setSelectedKey(matchedNode.value);
        setSelectedPath(matchedNode.path);
      }
    }
  }, [value, treeData]);

  const handleSelect = (keys: React.Key[], info: any) => {
    const key = keys[0];
    if (key) {
      const label = info.node.title;
      const path = info.node.path;
      setSelectedKey(key as string);
      setSelectedLabel(label);
      setSelectedPath(path);
      if (onChange) {
        onChange(key as string);
      }
    }
    setIsTreeVisible(false);
  };

  // 点击外部区域关闭下拉面板
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        setIsTreeVisible(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  return (
    <div
      ref={containerRef} // 👈 设置 ref 用于检测外部点击
      style={{ position: "relative", width: "100%" }}
    >
      {/* 模拟 Input 显示选中项 */}
      <input
        type="text"
        readOnly
        value={selectedLabel}
        onClick={() => setIsTreeVisible(true)}
        placeholder={placeholder}
        style={{
          width: "100%",
          backgroundColor: disabled
            ? themeMode === "light"
              ? paletteColors.gray["300"]
              : paletteColors.gray["800"]
            : "color-mix(in oklab, var(--input) var(--opacity-30), transparent)",
          padding: "8px 12px",
          border: "1px solid rgba(145 158 171 / 20%)",
          borderRadius: 4,
          cursor: disabled ? "not-allowed" : "pointer",
        }}
      />

      {/* 树形选择面板 */}
      {isTreeVisible && !disabled && (
        <div
          style={{
            position: "absolute",
            top: "100%",
            left: 0,
            right: 0,
            maxHeight: 200,
            overflowY: "auto",
            border: "1px solid #eaeaea",
            backgroundColor: "#ffffff",
            zIndex: 1000,
            boxShadow: "0 2px 8 rgba(0, 0, 0, 0.1)",
          }}
          tabIndex={-1} // 可聚焦但不影响 tab 导航
        >
          <Tree
            treeData={treeData}
            selectedKeys={[selectedKey]}
            filterTreeNode={filterTreeNode}
            onSelect={handleSelect}
          />
        </div>
      )}
    </div>
  );
};

export default TreeSelectInput;
