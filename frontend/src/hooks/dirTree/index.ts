import type { TreeNode } from "@/ui/tree-select-input";
import { buildFileTree } from "@/utils/tree";

export default function useDirTree(): TreeNode[] {
	// file dir tree
	const modules = import.meta.glob("/src/pages/**/*.tsx");
	const filePaths = Object.keys(modules).map((path) => path.replace("/src", "").replace(".tsx", ""));
	const tree = buildFileTree(filePaths);

	return tree ? (tree.children ? tree.children : []) : [];
}
