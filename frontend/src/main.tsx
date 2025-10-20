import "./global.css";
import "./theme/theme.css";
import "./locales/i18n";

import { Suspense, useEffect, useState } from "react";
import ReactDOM from "react-dom/client";
import { ErrorBoundary } from "react-error-boundary";
import { Navigate, Outlet, RouterProvider, createHashRouter } from "react-router";
import App from "./App";
import { registerLocalIcons } from "./components/icon";
import { LineLoading } from "./components/loading";
import DashboardLayout from "./layouts/dashboard";
import PageError from "./pages/sys/error/PageError";
import AuthGuard from "./routes/components/auth-guard";
import useAppMenu from "./routes/hooks/use-menu";
import { authRoutes } from "./routes/sections/auth";
import { buildRoutes, convertMenuTreeUserGroupToMenus } from "./routes/sections/build-routes";
import { mainRoutes } from "./routes/sections/main";
import { useUserInfo } from "./store/userStore";
import type { Menu, Role } from "./types/entity";

// create function router
function createAppRouter(menuData: Menu[], role: Role | undefined) {
	const routesSection = buildRoutes(menuData);
	const defaultRouter = role?.default_router;

	return createHashRouter([
		{
			Component: () => (
				<App>
					<Outlet />
				</App>
			),
			errorElement: <ErrorBoundary fallbackRender={PageError} />,
			children: [
				{
					path: "/",
					element: (
						<AuthGuard>
							<Suspense fallback={<LineLoading />}>
								<DashboardLayout />
							</Suspense>
						</AuthGuard>
					),
					children: [
						{
							index: true,
							element: <Navigate to={!defaultRouter ? "/dashboard/workbench" : defaultRouter} replace />,
						},
						...routesSection,
					],
				},
				...authRoutes,
				...mainRoutes,
				{ path: "*", element: <Navigate to="/404" replace /> },
			],
		},
	]);
}

// top app component
function AppWrapper() {
	const { menuData, loading, error } = useAppMenu();
	const { current_role: currentRole } = useUserInfo();
	const [router, setRouter] = useState<any>(null);
	const [initialized, setInitialized] = useState(false);

	useEffect(() => {
		let isMounted = true;

		const initializeRouter = async () => {
			try {
				let routesData: Menu[] = [];
				if (Array.isArray(menuData) && menuData.length > 0) {
					routesData = convertMenuTreeUserGroupToMenus(menuData);
				}

				if (isMounted) {
					const newRouter = createAppRouter(routesData, currentRole);
					setRouter(newRouter);
					setInitialized(true);
				}
			} catch (err) {
				console.error("Failed to initialize router:", err);
				if (isMounted) {
					// 即使出错也要确保有基本路由
					const fallbackRouter = createAppRouter([], currentRole);
					setRouter(fallbackRouter);
					setInitialized(true);
				}
			}
		};

		if (!loading) {
			initializeRouter();
		}

		return () => {
			isMounted = false;
		};
	}, [menuData, loading, currentRole]);

	if (loading || !initialized) {
		return <LineLoading />;
	}

	if (error) {
		return (
			<PageError
				error={error}
				resetErrorBoundary={(): void => {
					window.location.reload();
				}}
			/>
		);
	}

	if (!router) {
		return <LineLoading />;
	}

	return <RouterProvider router={router} />;
}
// 入口函数
async function initApp() {
	await registerLocalIcons();

	ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(<AppWrapper />);
}

initApp();
