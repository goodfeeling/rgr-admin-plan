import { QueryClientProvider } from "@tanstack/react-query";
import { QueryClient } from "@tanstack/react-query";
import { Analytics as VercelAnalytics } from "@vercel/analytics/react";
import { Helmet, HelmetProvider } from "react-helmet-async";
import { MotionLazy } from "./components/animate/motion-lazy";
import { RouteLoadingProgress } from "./components/loading";
import Toast from "./components/toast";
import { AntdAdapter } from "./theme/adapter/antd.adapter";
import { ThemeProvider } from "./theme/theme-provider";
import "@ant-design/v5-patch-for-react-19";
import { AliveScope } from "react-activation";
import { FavoritesProvider } from "./components/nav/favorites-context";
import { useMapBySystemConfig } from "./hooks";
import useUserStatusNotification from "./hooks/userStatusWatch/useUserStatusNotification";
import { getOrCreateDeviceId } from "./utils/deviceId";

// initial device ID
getOrCreateDeviceId();
const queryClient = new QueryClient();

function App({ children }: { children: React.ReactNode }) {
	return (
		<HelmetProvider>
			<QueryClientProvider client={queryClient}>
				<FavoritesProvider>
					<AppContent>{children}</AppContent>
				</FavoritesProvider>
			</QueryClientProvider>
		</HelmetProvider>
	);
}

function AppContent({ children }: { children: React.ReactNode }) {
	const { data: systemConfig } = useMapBySystemConfig();
	useUserStatusNotification();
	return (
		<ThemeProvider adapters={[AntdAdapter]}>
			<AliveScope>
				<VercelAnalytics />
				<Helmet>
					<title>{systemConfig?.name || "My Admin"}</title>
					<link rel="icon" href={systemConfig?.favicon || "/favicon.ico"} />
				</Helmet>
				<Toast />
				<RouteLoadingProgress />
				<MotionLazy>{children}</MotionLazy>
			</AliveScope>
		</ThemeProvider>
	);
}

export default App;
