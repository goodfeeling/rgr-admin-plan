import { cn } from "@/utils";
import type { NavProps } from "../types";
import { FavoritesBar } from "./favorites-bar";
import { NavGroup } from "./nav-group";

export function NavHorizontal({ data, className, ...props }: NavProps) {
	return (
		<nav className={cn("flex items-center gap-1 min-h-[56px] border-b border-dashed", className)} {...props}>
			{/* favorites bar */}
			<FavoritesBar />
			{/* other bar */}
			{data.map((group, index) => (
				<NavGroup key={group.name || index} name={group.name} items={group.items} />
			))}
		</nav>
	);
}
