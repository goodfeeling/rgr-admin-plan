import type { UserService } from "@/api/services/userService";
import userService from "@/api/services/userService";
import LocalePicker from "@/components/locale-picker";
import Logo from "@/components/logo";
import { useMapBySystemConfig } from "@/hooks";
import SettingButton from "@/layouts/components/setting-button";
import { Button } from "@/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/ui/form";
import { Input } from "@/ui/input";
import { Loader2 } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { useSearchParams } from "react-router";
import { LoginProvider } from "./providers/login-provider";

function ResetPasswordPage() {
	const [searchParams] = useSearchParams();
	const resetToken = searchParams.get("token");
	const { data: siteConfig } = useMapBySystemConfig();
	const { t } = useTranslation();
	const form = useForm<UserService.PasswordResetReq>();
	const [loading, setLoading] = useState(false);

	if (!resetToken) {
		return (
			<div className="relative grid min-h-svh lg:grid-cols-2 bg-background">
				<div className="flex flex-col gap-4 p-6 md:p-10">
					<div className="flex justify-center gap-2 md:justify-start">
						<div className="flex items-center gap-2 font-medium cursor-pointer">
							<Logo size={28} />
						</div>
					</div>
				</div>
			</div>
		);
	}

	const handleFinish = async (values: UserService.PasswordResetReq) => {
		setLoading(true);
		try {
			await userService.changePassword(values, resetToken || "");
		} finally {
			setLoading(false);
		}
	};
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
							<div className="flex flex-col gap-6">
								<Form {...form}>
									<form onSubmit={form.handleSubmit(handleFinish)} className="space-y-4">
										<div className="flex flex-col items-center gap-2 text-center">
											<h1 className="text-2xl font-bold">{t("sys.login.forgetFormTitle")}</h1>
											<p className="text-balance text-sm text-muted-foreground">
												{t("sys.login.resetPasswordDescription")}
											</p>
										</div>

										<FormField
											control={form.control}
											name="new_password"
											rules={{ required: t("sys.login.passwordPlaceholder") }}
											render={({ field }) => (
												<FormItem>
													<FormLabel>{t("sys.login.newPassword")}</FormLabel>
													<FormControl>
														<Input
															type="password"
															placeholder={t("sys.login.passwordPlaceholder")}
															{...field}
															suppressHydrationWarning
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>

										<FormField
											control={form.control}
											name="confirm_password"
											rules={{
												required: t("sys.login.confirmPasswordPlaceholder"),
											}}
											render={({ field }) => (
												<FormItem>
													<FormLabel>{t("sys.login.confirmPassword")}</FormLabel>
													<FormControl>
														<Input type="password" placeholder={t("sys.login.confirmPasswordPlaceholder")} {...field} />
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>

										{/* 提交按钮 */}
										<Button type="submit" className="w-full">
											{loading && <Loader2 className="animate-spin mr-2" />}
											{t("common.okText")}
										</Button>
									</form>
								</Form>
							</div>
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

export default ResetPasswordPage;
