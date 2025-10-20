import { useRouter } from "@/routes/hooks";
import { useEffect, useRef } from "react";
import { Outlet } from "react-router";
import styled from "styled-components";
import SortableContainer from "./components/sortable-container";
import { SortableItem } from "./components/sortable-item";
import { TabItem } from "./components/tab-item";
import { useMultiTabsStyle } from "./hooks/use-tab-style";
import { useMultiTabsContext } from "./providers/multi-tabs-provider";
import type { KeepAliveTab } from "./types";

export default function MultiTabs() {
	const scrollContainer = useRef<HTMLUListElement>(null);

	const { tabs, activeTabRoutePath, setTabs } = useMultiTabsContext();
	const style = useMultiTabsStyle();
	const { push } = useRouter();

	const handleTabClick = ({ key }: KeepAliveTab) => {
		const tabKey = key;
		push(tabKey);
	};

	useEffect(() => {
		if (!scrollContainer.current) return;
		const tab = tabs.find((item) => item.key === activeTabRoutePath);
		const currentTabElement = scrollContainer.current.querySelector(`#tab${tab?.key.split("/").join("-")}`);
		if (currentTabElement) {
			currentTabElement.scrollIntoView({
				block: "nearest",
				behavior: "smooth",
			});
		}
	}, [tabs, activeTabRoutePath]);

	useEffect(() => {
		const container = scrollContainer.current;
		if (!container) return;

		const handleWheel = (e: WheelEvent) => {
			e.preventDefault();
			container.scrollLeft += e.deltaY;
		};

		container.addEventListener("mouseenter", () => {
			container.addEventListener("wheel", handleWheel);
		});

		container.addEventListener("mouseleave", () => {
			container.removeEventListener("wheel", handleWheel);
		});

		return () => {
			container.removeEventListener("wheel", handleWheel);
		};
	}, []);

	const handleDragEnd = (oldIndex: number, newIndex: number) => {
		const newTabs = Array.from(tabs);
		const [movedTab] = newTabs.splice(oldIndex, 1);
		newTabs.splice(newIndex, 0, movedTab);

		setTabs([...newTabs]);
	};

	const renderOverlay = (id: string | number) => {
		const tab = tabs.find((tab) => tab.key === id);
		if (!tab) return null;
		return <TabItem tab={tab} />;
	};

	return (
		<StyledMultiTabs>
			{/* 标签栏部分 */}
			<div style={style}>
				<SortableContainer items={tabs} onSortEnd={handleDragEnd} renderOverlay={renderOverlay}>
					<ul ref={scrollContainer} className="flex overflow-x-auto w-full px-2 h-full hide-scrollbar">
						{tabs.map((tab) => (
							<SortableItem tab={tab} key={tab.key} onClick={() => handleTabClick(tab)} />
						))}
					</ul>
				</SortableContainer>
			</div>

			{/* 内容部分 */}
			<div className="tab-content flex-1 overflow-hidden">
				{tabs.map((tab) => (
					<div
						key={tab.key}
						style={{
							display: tab.key === activeTabRoutePath ? "block" : "none",
							height: "100%",
						}}
					>
						{tab.children ? (
							<div key={tab.timeStamp}>{tab.children}</div>
						) : (
							<div
								key={tab.timeStamp}
								style={{
									height: "100%",
								}}
							>
								<Outlet />
							</div>
						)}
					</div>
				))}
			</div>
		</StyledMultiTabs>
	);
}

const StyledMultiTabs = styled.div`
  height: 100%;
  display: flex;
  flex-direction: column;

  .tab-content {
    flex: 1;
    overflow: hidden;
    margin-top: var(--layout-multi-tabs-height);
  }

  .ant-tabs {
    height: 100%;
    display: flex;
    flex-direction: column;

    .ant-tabs-content {
      flex: 1;
      height: 100%;
      overflow: auto;
    }

    .ant-tabs-tabpane {
      height: 100%;
      & > div {
        height: 100%;
      }
    }
  }

  .hide-scrollbar {
    overflow: scroll;
    scrollbar-width: none;
    -ms-overflow-style: none;
    will-change: transform;

    &::-webkit-scrollbar {
      display: none;
    }
  }
`;
