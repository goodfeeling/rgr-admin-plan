import { Icon } from "@/components/icon";
import { Button } from "@/ui/button";
import { Card } from "antd";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import MenuList from "./base";
import MenuGroupList from "./group";

const App: React.FC = () => {
	const [selectedId, setSelectedId] = useState<number | null>(null);
	const [isCollapsed, setIsCollapsed] = useState(false);
	const { t } = useTranslation();

	const toggleCollapse = () => {
		setIsCollapsed(!isCollapsed);
	};

	return (
		<div className="flex w-full gap-4">
			<div className={`${isCollapsed ? "w-12" : "w-1/4"} pr-2 transition-all duration-300`}>
				<Card
					title={isCollapsed ? "" : t("sys.menu.system.menu_group")}
					size="small"
					extra={
						<Button variant="ghost" size="icon" onClick={toggleCollapse} className="h-6 w-6">
							<Icon icon={isCollapsed ? "ic:round-arrow-forward" : "ic:round-arrow-back"} className="h-4 w-4" />
						</Button>
					}
				>
					{isCollapsed ? (
						<div className="flex justify-center pt-4">
							<Button variant="ghost" size="icon" onClick={toggleCollapse} className="h-6 w-6">
								<Icon icon="ic:round-arrow-forward" className="h-4 w-4" />
							</Button>
						</div>
					) : (
						<MenuGroupList onSelect={setSelectedId} />
					)}
				</Card>
			</div>
			<div className={`${isCollapsed ? "w-[calc(100%-3rem-1.5rem)]" : "w-3/4"} pl-2 transition-all duration-300`}>
				<Card title={t("sys.menu.system.menu")} size="small">
					<MenuList selectedId={selectedId} />
				</Card>
			</div>
		</div>
	);
};

export default App;
