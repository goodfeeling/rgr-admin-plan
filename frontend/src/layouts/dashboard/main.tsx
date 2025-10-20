import { LineLoading } from "@/components/loading";
import { useSettings } from "@/store/settingStore";
import { ScrollArea } from "@/ui/scroll-area";
import { cn } from "@/utils";
import { Suspense } from "react";
import { Outlet } from "react-router";
import MultiTabs from "./multi-tabs";
import { MultiTabsProvider } from "./multi-tabs/providers/multi-tabs-provider";

const Main = () => {
	const { themeStretch, multiTab } = useSettings();
	return (
		<main
			data-slot="slash-layout-main"
			className={cn("flex w-full grow bg-background", {
				"md:pt-[var(--layout-multi-tabs-height)]": multiTab,
			})}
		>
			<ScrollArea
				className={cn(
					"h-full w-full p-2 mx-auto transition-all duration-300 ease-in-out overscroll-none",
					themeStretch ? "" : "xl:max-w-screen-xl",
					{
						// 当启用多标签页且在PC端时，需要额外的顶部padding来为固定定位的标签栏留出空间
						"md:pt-0": multiTab,
					},
				)}
			>
				{multiTab ? (
					<MultiTabsProvider>
						<MultiTabs />
					</MultiTabsProvider>
				) : (
					<Suspense fallback={<LineLoading />}>
						<Outlet />
					</Suspense>
				)}
			</ScrollArea>
		</main>
	);
};

export default Main;
