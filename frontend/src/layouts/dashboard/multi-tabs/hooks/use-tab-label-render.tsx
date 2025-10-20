import { USER_LIST } from "@/_mock/assets";
import { Icon } from "@/components/icon";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import type { KeepAliveTab } from "../types";

export function useTabLabelRender() {
	const { t } = useTranslation();

	const specialTabRenderMap = useMemo<Record<string, (tab: KeepAliveTab) => React.ReactNode>>(
		() => ({
			"sys.menu.system.user_detail": (tab: KeepAliveTab) => {
				const userId = tab.params?.id;
				const defaultLabel = t(tab.label);
				if (userId) {
					const user = USER_LIST.find((item) => item.id === userId);
					return `${user?.username}-${defaultLabel}`;
				}
				return defaultLabel;
			},
		}),
		[t],
	);

	const renderTabLabel = (tab: KeepAliveTab) => {
		const specialRender = specialTabRenderMap[tab.label];
		if (specialRender) {
			return specialRender(tab);
		}
		if (tab.icon) {
			return (
				<div className="flex items-center">
					<Icon icon={tab.icon as string} size={16} className="mr-2" />
					<span>{t(tab.label)}</span>
				</div>
			);
		}
		return t(tab.label);
	};

	return renderTabLabel;
}
