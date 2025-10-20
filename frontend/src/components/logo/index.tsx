import { useMapBySystemConfig } from "@/hooks";
import { useTheme } from "@/theme/hooks";
import { cn } from "@/utils";
import { NavLink } from "react-router";
import { Icon } from "../icon";

interface Props {
	size?: number | string;
	className?: string;
}
function Logo({ size = 50, className }: Props) {
	const { themeTokens } = useTheme();
	const { data: siteConfig } = useMapBySystemConfig();
	if (siteConfig?.logo?.endsWith(".svg") || siteConfig?.logo?.endsWith(".png") || siteConfig?.logo?.endsWith(".jpg")) {
		// 显示网络SVG图片
		return (
			<NavLink to="/" className={cn(className)}>
				<img src={siteConfig?.logo} alt="Logo" style={{ width: size, height: size }} />
			</NavLink>
		);
	}
	return (
		<NavLink to="/" className={cn(className)}>
			<Icon icon="solar:code-square-bold" color={themeTokens.color.palette.primary.default} size={size} />
		</NavLink>
	);
}

export default Logo;
