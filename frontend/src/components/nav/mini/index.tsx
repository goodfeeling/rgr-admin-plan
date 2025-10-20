import { cn } from "@/utils";
import type { NavProps } from "../types";
import { FavoritesBar } from "./favorites-bar";
import { NavGroup } from "./nav-group";

export const NavMini = ({ data, className, ...props }: NavProps) => {
	return (
		<nav className={cn("flex flex-col", className)} {...props}>
			{/* favorites bar */}
			<FavoritesBar />
			{/* other bar */}
			<ul className="flex flex-col gap-1">
				{data.map((item, index) => (
					<NavGroup key={item.name || index} items={item.items} />
				))}
			</ul>
		</nav>
	);
};
