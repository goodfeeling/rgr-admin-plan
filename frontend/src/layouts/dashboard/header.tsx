import LocalePicker from "@/components/locale-picker";
import type { NavItemDataProps } from "@/components/nav/types";
import { useSettings } from "@/store/settingStore";
import { cn } from "@/utils";
import type { ReactNode } from "react";
import AccountDropdown from "../components/account-dropdown";
import BreadCrumb from "../components/bread-crumb";
import SearchBar from "../components/search-bar";
import SettingButton from "../components/setting-button";

interface HeaderProps {
	leftSlot?: ReactNode;
	navData: {
		name?: string;
		items: NavItemDataProps[];
	}[];
}

export default function Header({ leftSlot, navData }: HeaderProps) {
	const { breadCrumb } = useSettings();
	return (
		<header
			data-slot="slash-layout-header"
			className={cn(
				"sticky z-app-bar top-0 right-0 left-auto flex items-center bg-background justify-between px-2 ml-[1px]",
				"h-[var(--layout-header-height)] grow-0 shrink-0",
			)}
		>
			<div className="flex items-center">
				{leftSlot}

				<div className="hidden md:block ml-4">{breadCrumb && <BreadCrumb navData={navData} />}</div>
			</div>

			<div className="flex items-center gap-1">
				<SearchBar />
				<LocalePicker />
				{/* <NoticeButton /> */}
				<SettingButton />
				<AccountDropdown />
			</div>
		</header>
	);
}
