export type LangTree = {
	label: string;
	value: string;
	children: LangTree[];
};

export default function useLangTree(lang: any): LangTree[] {
	// 构建树结构
	const buildLangTree = (subChildren: any, parentKey: string): LangTree[] => {
		const result: LangTree[] = [];
		for (const k in subChildren) {
			const currentKey = parentKey ? `${parentKey}.${k}` : k;
			let temp: LangTree = {
				label: k,
				value: currentKey,
				children: [],
			};

			const subTemp = (subChildren as any)[k];
			if (subTemp instanceof Object) {
				temp.children = buildLangTree((subChildren as any)[k], currentKey);
			} else {
				temp = {
					label: `${currentKey}-${subTemp}`,
					value: currentKey,
					children: [],
				};
			}
			result.push(temp);
		}

		return result;
	};

	return buildLangTree(lang, "");
}
