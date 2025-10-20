import type { TreeNode } from "@/ui/tree-select-input";
import { chain } from "ramda";

/**
 * Flatten an array containing a tree structure
 * @param {T[]} trees - An array containing a tree structure
 * @returns {T[]} - Flattened array
 */
export function flattenTrees<T extends { children?: T[] }>(trees: T[] = []): T[] {
	return chain((node) => {
		const children = node.children || [];
		return [node, ...flattenTrees(children)];
	}, trees);
}

export function buildFileTree(paths: string[]): TreeNode | null {
	if (!paths.length) return null;

	const root: TreeNode = {
		title: "",
		value: "",
		key: "",
		path: [],
		children: [],
		isLast: false,
	};

	let i = 1;
	for (const path of paths) {
		const segments = path.split(/[/\\]/);
		let currentNode = root;

		let pathArr: number[] = [];
		let pathStr = "";
		for (const [index, segment] of segments.entries()) {
			const isLast = index === segments.length - 1;
			if (!segment) continue;
			pathArr = [...pathArr, i];

			if (segment !== "pages") {
				if (pathStr === "") {
					pathStr = `${segment}`;
				} else {
					pathStr = `${pathStr}/${segment}`;
				}
			}

			let existingChild = currentNode.children?.find((child) => child.title === segment);

			if (!existingChild) {
				existingChild = {
					title: segment,
					value: pathStr,
					key: pathStr,
					path: pathArr,
					children: [],
					isLast: isLast,
				};

				if (!currentNode.children) currentNode.children = [];
				currentNode.children.push(existingChild);
			}

			currentNode = existingChild;
			i++;
		}
	}

	return root.children?.[0] || null; // 返回第一个根节点作为实际的根目录
}
