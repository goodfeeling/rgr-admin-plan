import { Icon } from "@/components/icon";
import useLocale from "@/locales/use-locale";
import { themeVars } from "@/theme/theme.css";
import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/ui/hover-card";
import { cn } from "@/utils";
import { useLocation } from "react-router";
import { useFavorites } from "../favorites-context";
import { navItemStyles } from "../styles";

export function FavoritesBar() {
	const { favorites } = useFavorites();
	const { t } = useLocale();
	const location = useLocation();

	if (favorites.length === 0) {
		return null;
	}

	return (
		<>
			{favorites.length > 0 && (
				<HoverCard openDelay={100}>
					<HoverCardTrigger asChild>
						<div
							className={cn(
								"inline-flex items-center justify-center rounded-md px-2 py-1.5 text-sm transition-all duration-300 ease-in-out cursor-pointer",
								"hover:bg-action-hover!",
							)}
						>
							<div style={{ color: themeVars.colors.text.secondary }}>
								<span style={navItemStyles.icon} className="items-center justify-center">
									<Icon icon="eva:star-fill" className="h-5 w-5" />
								</span>
								<span className="ml-2 flex-auto!">{t("sys.menu.favorites")}</span>
							</div>
						</div>
					</HoverCardTrigger>
					<HoverCardContent side="bottom" sideOffset={10} className="p-1">
						<div className="min-w-48 rounded-md border bg-popover p-2 text-popover-foreground shadow-md">
							{favorites.map((item, index) => {
								const isActive = location.pathname === item.path;
								return (
									<div
										key={`favorite_${String(index)}`}
										className={cn(
											"flex w-full cursor-pointer items-center rounded-md px-2 py-1.5 text-sm transition-all duration-300 ease-in-out",
											isActive ? "bg-primary/10 text-primary hover:bg-primary/20" : "hover:bg-action-hover!",
										)}
										onClick={() => {
											window.location.hash = `#${item.path}`;
										}}
										style={{ color: themeVars.colors.text.secondary }}
									>
										{item.icon && typeof item.icon === "string" ? (
											<Icon icon={item.icon} className="mr-2 h-5 w-5" />
										) : (
											item.icon && <span className="mr-2">{item.icon}</span>
										)}
										<span>{t(item.title)}</span>
									</div>
								);
							})}
						</div>
					</HoverCardContent>
				</HoverCard>
			)}
		</>
	);
}
