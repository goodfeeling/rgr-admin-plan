import Logo from "@/components/logo";
import { down, useMediaQuery } from "@/hooks";
import { useSettings } from "@/store/settingStore";
import { useMenu } from "@/store/useMenuStore";
import type { MenuTreeUserGroup } from "@/types/entity";
import { cn } from "@/utils";
import type { FC } from "react";
import { ThemeLayout } from "#/enum";
import Header from "./header";
import Main from "./main";
import { NavHorizontalLayout, NavMobileLayout, NavToggleButton, NavVerticalLayout } from "./nav";

interface LayoutProps {
	navData: MenuTreeUserGroup[];
}
// Dashboard Layout
const DashboardLayout: FC = () => {
	const isMobile = useMediaQuery(down("md"));
	const { themeLayout } = useSettings();
	const menuData = useMenu();

	return (
		<div
			data-slot="slash-layout-root"
			className={cn("w-full min-h-svh flex bg-background", {
				"flex-col": isMobile || themeLayout === ThemeLayout.Horizontal,
			})}
		>
			{isMobile ? <MobileLayout navData={menuData} /> : <PcLayout navData={menuData} />}
		</div>
	);
};
export default DashboardLayout;

// Pc Layout
function PcLayout({ navData }: LayoutProps) {
	const { themeLayout } = useSettings();
	if (themeLayout === ThemeLayout.Horizontal) return <PcHorizontalLayout navData={navData} />;
	return <PcVerticalLayout navData={navData} />;
}

function PcHorizontalLayout({ navData }: LayoutProps) {
	return (
		<div
			data-slot="slash-layout-content"
			className={cn("w-full h-screen flex flex-col transition-all duration-300 ease-in-out")}
		>
			<Header leftSlot={<Logo />} navData={navData} />
			<NavHorizontalLayout data={navData} />
			<Main />
		</div>
	);
}

function PcVerticalLayout({ navData }: LayoutProps) {
	const settings = useSettings();
	const { themeLayout } = settings;

	return (
		<>
			<NavVerticalLayout data={navData} />
			<div
				data-slot="slash-layout-content"
				className={cn("w-full h-screen flex flex-col transition-[padding] duration-300 ease-in-out", {
					"pl-[var(--layout-nav-width)]": themeLayout === ThemeLayout.Vertical,
					"pl-[var(--layout-nav-width-mini)]": themeLayout === ThemeLayout.Mini,
				})}
			>
				<Header leftSlot={<NavToggleButton />} navData={navData} />
				<Main />
			</div>
		</>
	);
}

// Mobile Layout
function MobileLayout({ navData }: LayoutProps) {
	return (
		<>
			<Header leftSlot={<NavMobileLayout data={navData} />} navData={navData} />
			<Main />
		</>
	);
}
