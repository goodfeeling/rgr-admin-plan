import { up, useMediaQuery } from "@/hooks";
import { useSettings } from "@/store/settingStore";
import { themeVars } from "@/theme/theme.css";
import { rgbAlpha } from "@/utils/theme";
import { type CSSProperties, useMemo } from "react";
import { ThemeLayout } from "#/enum";

export function useMultiTabsStyle() {
	const { themeLayout } = useSettings();
	const isPc = useMediaQuery(up("md"));

	return useMemo(() => {
		const style: CSSProperties = {
			position: "fixed",
			top: "var(--layout-header-height)", // 紧接在头部下方
			height: "36px", // 设置标签页高度
			backgroundColor: rgbAlpha(themeVars.colors.background.defaultChannel, 0.9),
			transition: "all 200ms cubic-bezier(0.4, 0, 0.2, 1) 0ms",
			width: "100%",
		};

		// 根据布局类型调整宽度
		if (themeLayout === ThemeLayout.Horizontal) {
			// 水平布局时保持全宽
		} else if (isPc) {
			// 垂直布局时需要考虑侧边栏宽度
			style.width = "calc(100% - var(--layout-nav-width))";
		} else {
			// 移动端保持全宽
			style.width = "100%";
			style.left = 0;
		}

		return style;
	}, [themeLayout, isPc]);
}
