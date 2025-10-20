import LocalePicker from "@/components/locale-picker";
import Logo from "@/components/logo";
import { useMapBySystemConfig } from "@/hooks";
import SettingButton from "@/layouts/components/setting-button";
import { useUserInfo, useUserToken } from "@/store/userStore";
import { Navigate } from "react-router";
import LoginForm from "./login-form";
import MobileForm from "./mobile-form";
import { LoginProvider } from "./providers/login-provider";
import QrCodeFrom from "./qrcode-form";
import RegisterForm from "./register-form";
import ResetForm from "./reset-form";

function LoginPage() {
	const token = useUserToken();
	const { data: siteConfig } = useMapBySystemConfig();
	const { current_role: currentRole } = useUserInfo();
	if (token.accessToken) {
		return <Navigate to={currentRole?.default_router ? currentRole.default_router : "/dashboard/workbench"} replace />;
	}

	return (
		<div className="relative grid min-h-svh lg:grid-cols-2 bg-background">
			<div className="flex flex-col gap-4 p-6 md:p-10">
				<div className="flex justify-center gap-2 md:justify-start">
					<div className="flex items-center gap-2 font-medium cursor-pointer">
						<Logo size={28} />
						<span>{siteConfig?.name}</span>
					</div>
				</div>
				<div className="flex flex-1 items-center justify-center">
					<div className="w-full max-w-xs">
						<LoginProvider>
							<LoginForm />
							<MobileForm />
							<QrCodeFrom />
							<RegisterForm />
							<ResetForm />
						</LoginProvider>
					</div>
				</div>
			</div>

			<div className="relative hidden bg-background-paper lg:block">
				<img src={siteConfig?.login_img} alt="placeholder img" className="absolute inset-0 h-full w-full" />
			</div>

			<div className="absolute right-2 top-0 flex flex-row">
				<LocalePicker />
				<SettingButton />
			</div>
		</div>
	);
}
export default LoginPage;
