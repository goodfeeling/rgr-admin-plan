import { DEFAULT_USER } from "@/_mock/assets";
import captchaService, { type CaptchaService } from "@/api/services/captchaService";
import type { UserService } from "@/api/services/userService";
import { Icon } from "@/components/icon";
import { useSignIn } from "@/store/userStore";
import { Button } from "@/ui/button";
import { Checkbox } from "@/ui/checkbox";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/ui/form";
import { Input } from "@/ui/input";
import { cn } from "@/utils";
import { Loader2 } from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { LoginStateEnum, useLoginStateContext } from "./providers/login-provider";

export function LoginForm({ className, ...props }: React.ComponentPropsWithoutRef<"form">) {
	const { t } = useTranslation();
	const [loading, setLoading] = useState(false);
	const [remember, setRemember] = useState(true);
	const [captcha, setCaptcha] = useState<CaptchaService.CaptchaGenerateResult | null>(null);
	const [captchaLoading, setCaptchaLoading] = useState(false);

	const { loginState, setLoginState } = useLoginStateContext();
	const signIn = useSignIn();

	const form = useForm<UserService.SignInReq & { captcha_code: string; captcha_id: string }>({
		defaultValues: {
			user_name: DEFAULT_USER.username,
			password: DEFAULT_USER.password,
			captcha_code: "",
			captcha_id: "",
		},
	});

	// 获取验证码
	const fetchCaptcha = useCallback(
		async (currentCaptchaId?: string) => {
			setCaptchaLoading(true);
			try {
				const res = await captchaService.generate(currentCaptchaId || "");
				setCaptcha(res);
				// 设置验证码ID到表单中
				form.setValue("captcha_id", res.id);
			} catch (error) {
				console.error("Failed to fetch captcha", error);
			} finally {
				setCaptchaLoading(false);
			}
		},
		[form],
	);

	// 初次加载时获取验证码
	useEffect(() => {
		if (loginState === LoginStateEnum.LOGIN) {
			fetchCaptcha();
		}
	}, [loginState, fetchCaptcha]);

	const handleFinish = async (values: UserService.SignInReq & { captcha_code: string; captcha_id: string }) => {
		setLoading(true);
		try {
			// 先验证验证码
			if (!values.captcha_id || !values.captcha_code) {
				toast.error(t("sys.login.captchaPlaceholder"));
				return;
			}

			await signIn({
				user_name: values.user_name,
				password: values.password,
				captcha_answer: values.captcha_code,
				captcha_id: values.captcha_id,
			});
		} finally {
			// 登录后刷新验证码
			fetchCaptcha();
			setLoading(false);
		}
	};
	if (loginState !== LoginStateEnum.LOGIN) return null;

	return (
		<div className={cn("flex flex-col gap-6", className)}>
			<Form {...form} {...props}>
				<form onSubmit={form.handleSubmit(handleFinish)} className="space-y-4">
					<div className="flex flex-col items-center gap-2 text-center">
						<h1 className="text-2xl font-bold">{t("sys.login.signInFormTitle")}</h1>
						<p className="text-balance text-sm text-muted-foreground">{t("sys.login.signInFormDescription")}</p>
					</div>

					<FormField
						control={form.control}
						name="user_name"
						rules={{ required: t("sys.login.accountPlaceholder") }}
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("sys.login.userName")}</FormLabel>
								<FormControl>
									<Input placeholder={t("sys.login.accountPlaceholder")} {...field} />
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>

					<FormField
						control={form.control}
						name="password"
						rules={{ required: t("sys.login.passwordPlaceholder") }}
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("sys.login.password")}</FormLabel>
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
						name="captcha_code"
						rules={{ required: t("sys.login.captchaPlaceholder") }}
						render={({ field }) => (
							<FormItem>
								<FormLabel>{t("sys.login.captcha")}</FormLabel>
								<div className="flex flex-col gap-2">
									<FormControl>
										<Input
											placeholder={t("sys.login.captchaPlaceholder")}
											{...field}
											maxLength={captcha?.config.length}
										/>
									</FormControl>
									<div
										className="flex cursor-pointer items-center justify-center rounded border bg-slate-50"
										onClick={() => fetchCaptcha(captcha?.id)}
										style={{
											minHeight: `${captcha?.config.height}px`,
											minWidth: `${captcha?.config.width}px}`,
										}}
									>
										{captchaLoading ? (
											<Loader2 className="animate-spin" />
										) : captcha ? (
											<img
												src={captcha.b64s}
												alt="captcha"
												className="object-contain"
												style={{ maxWidth: "100%", maxHeight: "100%" }}
											/>
										) : (
											<span className="text-xs text-muted-foreground">验证码</span>
										)}
									</div>
								</div>
								<FormMessage />
							</FormItem>
						)}
					/>
					<input type="hidden" {...form.register("captcha_id")} />

					{/* 记住我/忘记密码 */}
					<div className="flex flex-row justify-between">
						<div className="flex items-center space-x-2">
							<Checkbox
								id="remember"
								checked={remember}
								onCheckedChange={(checked) => setRemember(checked === "indeterminate" ? false : checked)}
							/>
							<label
								htmlFor="remember"
								className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
							>
								{t("sys.login.rememberMe")}
							</label>
						</div>
						<Button variant="link" onClick={() => setLoginState(LoginStateEnum.RESET_PASSWORD)} size="sm">
							{t("sys.login.forgetPassword")}
						</Button>
					</div>

					{/* 登录按钮 */}
					<Button type="submit" className="w-full" disabled={captchaLoading}>
						{loading && <Loader2 className="animate-spin mr-2" />}
						{t("sys.login.loginButton")}
					</Button>

					{/* 手机登录/二维码登录 */}
					<div className="grid gap-4 sm:grid-cols-2">
						<Button variant="outline" className="w-full" onClick={() => setLoginState(LoginStateEnum.MOBILE)}>
							<Icon icon="uil:mobile-android" size={20} />
							{t("sys.login.mobileSignInFormTitle")}
						</Button>
						<Button variant="outline" className="w-full" onClick={() => setLoginState(LoginStateEnum.QR_CODE)}>
							<Icon icon="uil:qrcode-scan" size={20} />
							{t("sys.login.qrSignInFormTitle")}
						</Button>
					</div>

					{/* 其他登录方式 */}
					{/* <div className="relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t after:border-border">
						<span className="relative z-10 bg-background px-2 text-muted-foreground">{t("sys.login.otherSignIn")}</span>
					</div>
					<div className="flex cursor-pointer justify-around text-2xl">
						<Button variant="ghost" size="icon">
							<Icon icon="mdi:github" size={24} />
						</Button>
						<Button variant="ghost" size="icon">
							<Icon icon="mdi:wechat" size={24} />
						</Button>
						<Button variant="ghost" size="icon">
							<Icon icon="ant-design:google-circle-filled" size={24} />
						</Button>
					</div> */}

					{/* 注册 */}
					<div className="text-center text-sm">
						{t("sys.login.noAccount")}
						<Button variant="link" className="px-1" onClick={() => setLoginState(LoginStateEnum.REGISTER)}>
							{t("sys.login.signUpFormTitle")}
						</Button>
					</div>
				</form>
			</Form>
		</div>
	);
}

export default LoginForm;
