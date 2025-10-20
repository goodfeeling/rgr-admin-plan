import { Icon } from "@/components/icon";
import useLocale from "@/locales/use-locale";
import { themeVars } from "@/theme/theme.css";
import { Button } from "@/ui/button";
import { useLocation } from "react-router";
import { useFavorites } from "../favorites-context";

export function FavoritesBar() {
	const { favorites, removeFavorite } = useFavorites();
	const { t } = useLocale();
	const location = useLocation();

	if (favorites.length === 0) {
		return null;
	}

	return (
		<div className="p-2 border-b border-border">
			<div className="flex items-center justify-between mb-2">
				<div className="flex items-center" style={{ color: themeVars.colors.text.secondary }}>
					<Icon icon={"eva:star-fill"} className="h-3.5 w-3.5 mr-1" />

					<span className="text-sm font-medium">{t("sys.menu.favorites")}</span>
				</div>
			</div>

			<div className="grid grid-cols-2 gap-2">
				{favorites.map((item, index) => (
					<div
						key={`favorite_${String(index)}`}
						className={`flex flex-col items-center justify-center p-2 rounded-lg cursor-pointer transition-colors relative group ${
							location.pathname === item.path
								? "bg-primary/10 border border-primary/20"
								: "hover:bg-gray-100 dark:hover:bg-gray-800"
						}`}
						onClick={() => {
							window.location.hash = `#${item.path}`;
						}}
						style={{ color: themeVars.colors.text.secondary }}
					>
						<Button
							variant="ghost"
							size="sm"
							className="absolute top-0 right-0 h-5 w-5 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
							onClick={(e) => {
								e.stopPropagation();
								removeFavorite(item.path);
							}}
						>
							<Icon icon="eva:close-circle-fill" className="h-4 w-4 text-red-500" />
						</Button>
						<div className="mb-1">
							{item.icon && typeof item.icon === "string" ? <Icon icon={item.icon} className="h-6 w-6" /> : item.icon}
						</div>
						<span className="text-xs text-center truncate w-full px-1">{t(item.title)}</span>
					</div>
				))}
			</div>
		</div>
	);
}
