import emailService from "@/api/services/emailService";
import { Icon } from "@/components/icon";
import { Button } from "@/ui/button";
import { Form, FormControl, FormField, FormItem, FormMessage } from "@/ui/form";
import { Input } from "@/ui/input";
import { useMutation } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { ReturnButton } from "./components/ReturnButton";
import { LoginStateEnum, useLoginStateContext } from "./providers/login-provider";

type ResetFromEmail = { email: string };

function ResetForm() {
	const { t } = useTranslation();
	const { loginState, backToLogin } = useLoginStateContext();
	const form = useForm<ResetFromEmail>();
	const sendEmailMutation = useMutation({
		mutationFn: async (email: string) => {
			return await emailService.sendForgetPassword(email);
		},
		onSuccess: () => {
			toast.success("success!");
		},
	});
	const onFinish = (values: ResetFromEmail) => {
		sendEmailMutation.mutate(values.email);
	};

	// 修改返回按钮的处理函数
	const handleReturn = () => {
		// 重置mutation状态
		sendEmailMutation.reset();
		// 调用原有的返回逻辑
		backToLogin();
	};
	if (loginState !== LoginStateEnum.RESET_PASSWORD) return null;

	// 请求成功后显示成功消息
	if (sendEmailMutation.isSuccess) {
		return (
			<>
				<div className="mb-8 text-center">
					<Icon icon="local:ic-reset-password" size="100" className="text-primary!" />
				</div>
				<div className="space-y-4">
					<div className="flex flex-col items-center gap-2 text-center">
						<h1 className="text-2xl font-bold">{t("sys.login.forgetFormTitle")}</h1>
						<p className="text-balance text-sm text-muted-foreground">
							{t("sys.login.resetPasswordRequestToYourEmail")}
						</p>
					</div>
					<ReturnButton onClick={backToLogin} />
				</div>
			</>
		);
	}

	return (
		<>
			<div className="mb-8 text-center">
				<Icon icon="local:ic-reset-password" size="100" className="text-primary!" />
			</div>
			<Form {...form}>
				<form onSubmit={form.handleSubmit(onFinish)} className="space-y-4">
					<div className="flex flex-col items-center gap-2 text-center">
						<h1 className="text-2xl font-bold">{t("sys.login.forgetFormTitle")}</h1>
						<p className="text-balance text-sm text-muted-foreground">{t("sys.login.forgetFormSecondTitle")}</p>
					</div>

					<FormField
						control={form.control}
						name="email"
						render={({ field }) => (
							<FormItem>
								<FormControl>
									<Input placeholder={t("sys.login.email")} {...field} />
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<Button type="submit" className="w-full" disabled={sendEmailMutation.isPending}>
						{sendEmailMutation.isPending ? t("sys.login.sendingEmailButton") : t("sys.login.sendEmailButton")}
					</Button>
					<ReturnButton onClick={handleReturn} />
				</form>
			</Form>
		</>
	);
}

export default ResetForm;
