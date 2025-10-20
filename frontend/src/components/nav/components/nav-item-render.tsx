import { RouterLink } from "@/routes/components/router-link";
import type { NavItemProps } from "../types";

type NavItemRendererProps = {
	item: NavItemProps;
	className: string;
	children: React.ReactNode;
	hidden: boolean;
};

/**
 * Renderer for Navigation Items.
 * Handles disabled, external link, clickable child container, and internal link logic.
 */
export const NavItemRenderer: React.FC<NavItemRendererProps> = ({ item, className, children, hidden }) => {
	const { disabled, externalLink, hasChild, path, onClick } = item;

	const style = hidden ? { display: "none" } : {};

	if (disabled) {
		return (
			<div className={className} style={style}>
				{children}
			</div>
		);
	}

	if (externalLink) {
		return (
			<a href={path} target="_blank" rel="noopener noreferrer" className={className} style={style}>
				{children}
			</a>
		);
	}

	if (hasChild) {
		// Vertical nav items with children are clickable containers
		return (
			<div className={className} onClick={onClick} style={style}>
				{children}
			</div>
		);
	}

	// Default: internal link
	return (
		<RouterLink href={path} className={className} style={style}>
			{children}
		</RouterLink>
	);
};
