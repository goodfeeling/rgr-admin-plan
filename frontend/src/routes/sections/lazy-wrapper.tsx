import { Suspense } from "react";
import type { ComponentType, LazyExoticComponent, ReactNode } from "react";
import KeepAlive from "react-activation";

interface LazyLoadOptions {
	keepAlive?: boolean;
	keepAliveName?: string;
}

export function lazyLoad(Component: LazyExoticComponent<ComponentType<any>>, options: LazyLoadOptions = {}): ReactNode {
	const { keepAlive = false, keepAliveName } = options;

	const element = (
		<Suspense fallback={<div>Loading...</div>}>
			<Component />
		</Suspense>
	);

	// 如果需要 keepAlive，则包装 KeepAlive 组件
	if (keepAlive) {
		return (
			<KeepAlive name={keepAliveName} saveScrollPosition="screen">
				{element}
			</KeepAlive>
		);
	}

	return element;
}
