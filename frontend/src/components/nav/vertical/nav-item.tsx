import Icon from "@/components/icon/icon";
import useLocale from "@/locales/use-locale";
import { TooltipContent } from "@/ui/tooltip";
import { Tooltip } from "@/ui/tooltip";
import { TooltipTrigger } from "@/ui/tooltip";
import { TooltipProvider } from "@/ui/tooltip";
import { cn } from "@/utils";
import { NavItemRenderer } from "../components";
import { useFavorites } from "../favorites-context";
import { navItemClasses, navItemStyles } from "../styles";
import type { NavItemProps } from "../types";

export function NavItem(item: NavItemProps) {
	const { title, icon, info, caption, open, active, disabled, depth, hasChild, hidden, path } = item;
	const { t } = useLocale();

	const { addFavorite, removeFavorite, isFavorite } = useFavorites();
	const favorite = isFavorite(path);

	const toggleFavorite = (e: React.MouseEvent) => {
		e.stopPropagation();
		e.preventDefault();

		if (favorite) {
			removeFavorite(path);
		} else {
			// 添加当前项到收藏夹，只保留必要的属性
			const favoriteItem = {
				title,
				path,
				icon,
			};
			addFavorite(favoriteItem);
		}
	};

	const content = (
		<>
			{/* Icon */}
			<span style={navItemStyles.icon} className="mr-3 items-center justify-center">
				{icon && typeof icon === "string" ? <Icon icon={icon} /> : icon}
			</span>

			{/* Texts */}
			<span style={navItemStyles.texts} className="min-h-[24px]">
				{/* Title */}
				<span style={navItemStyles.title}>{t(title)}</span>

				{/* Caption */}
				{caption && (
					<TooltipProvider>
						<Tooltip>
							<TooltipTrigger asChild>
								<span style={navItemStyles.caption}>{t(caption)}</span>
							</TooltipTrigger>
							<TooltipContent side="top" align="start">
								{t(caption)}
							</TooltipContent>
						</Tooltip>
					</TooltipProvider>
				)}
			</span>

			{/* Info */}
			{info && <span style={navItemStyles.info}>{info}</span>}

			{/* Arrow */}
			{hasChild && (
				<Icon
					icon="eva:arrow-ios-forward-fill"
					style={{
						...navItemStyles.arrow,
						transform: open ? "rotate(90deg)" : "rotate(0deg)",
					}}
				/>
			)}

			{/* Favorite Button */}
			{!hasChild && (
				<div
					onClick={toggleFavorite}
					className={cn(
						"h-6 w-6 p-0 ml-auto mr-1 flex items-center justify-center cursor-pointer rounded-full hover:bg-gray-200 dark:hover:bg-gray-700",
						favorite && "text-yellow-500 hover:text-yellow-600",
					)}
				>
					<Icon icon={favorite ? "eva:star-fill" : "eva:star-outline"} className="h-3.5 w-3.5" />
				</div>
			)}
		</>
	);

	const itemClassName = cn(
		navItemClasses.base,
		navItemClasses.hover,
		"min-h-[44px]",
		active && depth === 1 && navItemClasses.active,
		active && depth !== 1 && "bg-action-hover!",
		disabled && navItemClasses.disabled,
	);

	return (
		<NavItemRenderer hidden={Boolean(hidden)} item={item} className={itemClassName}>
			{content}
		</NavItemRenderer>
	);
}
