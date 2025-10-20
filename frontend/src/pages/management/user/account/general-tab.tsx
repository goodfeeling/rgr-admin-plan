import { UploadService } from "@/api/services/uploadService";
import userService from "@/api/services/userService";
import { UploadAvatar } from "@/components/upload";
import useUserStore, { useUserInfo } from "@/store/userStore";
import userStore from "@/store/userStore";
import { Button } from "@/ui/button";
import { Card, CardContent, CardFooter } from "@/ui/card";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/ui/form";
import { Input } from "@/ui/input";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import type { UpdateUser } from "#/entity";
type FieldType = {
	user_name?: string;
	email?: string;
	nick_name: string;
	header_img: string;
	phone?: string;
};

export default function GeneralTab() {
	const { t } = useTranslation();
	const { userToken } = useUserStore.getState();
	const { header_img, user_name, email, phone, nick_name, id = 0 } = useUserInfo();
	const form = useForm<FieldType>({
		defaultValues: {
			user_name,
			email,
			phone,
			nick_name,
		},
	});

	const handleClick = async () => {
		try {
			const { user_name = "", email = "", phone = "", header_img, nick_name = "" } = await form.getValues();
			const updateUser: UpdateUser = {
				user_name,
				email,
				phone,
				header_img,
				nick_name,
			};
			const { actions } = userStore.getState();
			const userInfo = await userService.updateUser(id, updateUser);
			actions.setUserInfo(userInfo);
			toast.success("Update success!");
		} catch (error) {
			console.error("error:", error);
		}
	};

	const onHeaderImgChange = (fileUrl: string) => {
		form.setValue("header_img", fileUrl);
	};

	return (
		<div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
			<div className="flex-1">
				<Card className="flex-col px-6! pb-10! pt-20!">
					<UploadAvatar
						defaultAvatar={header_img}
						onHeaderImgChange={onHeaderImgChange}
						action={`${import.meta.env.VITE_APP_BASE_API}${UploadService.Client.Single}`}
						headers={{
							Authorization: `Bearer ${userToken?.accessToken}`,
						}}
					/>
					<Button variant="destructive">{t("common.delAccountText")}</Button>
				</Card>
			</div>
			<div className="flex-2">
				<Card>
					<CardContent>
						<Form {...form}>
							<div className="grid grid-cols-1 gap-4 md:grid-cols-2">
								<FormField
									control={form.control}
									name="user_name"
									disabled
									render={({ field }) => (
										<FormItem>
											<FormLabel>{t("sys.account.username")}</FormLabel>
											<FormControl>
												<Input {...field} />
											</FormControl>
										</FormItem>
									)}
								/>
								<FormField
									control={form.control}
									name="email"
									render={({ field }) => (
										<FormItem>
											<FormLabel>{t("sys.account.email")}</FormLabel>
											<FormControl>
												<Input {...field} />
											</FormControl>
										</FormItem>
									)}
								/>
								<FormField
									control={form.control}
									name="phone"
									render={({ field }) => (
										<FormItem>
											<FormLabel>{t("sys.account.phone")}</FormLabel>
											<FormControl>
												<Input {...field} />
											</FormControl>
										</FormItem>
									)}
								/>
								<FormField
									control={form.control}
									name="nick_name"
									render={({ field }) => (
										<FormItem>
											<FormLabel>{t("sys.account.nickname")}</FormLabel>
											<FormControl>
												<Input {...field} />
											</FormControl>
										</FormItem>
									)}
								/>
							</div>
						</Form>
					</CardContent>
					<CardFooter>
						<Button onClick={handleClick}>{t("common.saveText")}</Button>
					</CardFooter>
				</Card>
			</div>
		</div>
	);
}
