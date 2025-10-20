import userService from "@/api/services/userService";
import { useRouter } from "@/routes/hooks";
import { useUserActions, useUserInfo } from "@/store/userStore";
import { Button } from "@/ui/button";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@/ui/dropdown-menu";
import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { NavLink } from "react-router";
import SwitchModal, { type SwitchModalProps } from "./switch-role";

/**
 * Account Dropdown
 */
export default function AccountDropdown() {
	const { replace } = useRouter();
	const { t } = useTranslation();
	const { user_name, email, header_img, current_role } = useUserInfo();
	const { clearUserInfoAndToken } = useUserActions();
	// const { backToLogin } = useLoginStateContext();

	const [switchDataModal, setSwitchDataModal] = useState<SwitchModalProps>({
		title: "Switch Role",
		show: false,
		onCancel: () => {
			setSwitchDataModal((prev) => ({ ...prev, show: false }));
		},
	});

	const logoutMutation = useMutation({
		mutationFn: userService.logout,
		onSuccess: () => {
			clearUserInfoAndToken();
			replace("/auth/login");
		},
		onError: (error) => {
			console.error("Logout failed:", error);
			clearUserInfoAndToken();
			replace("/auth/login");
		},
	});

	const handleLogout = () => {
		logoutMutation.mutate();
	};

	const handleSwitch = () => {
		setSwitchDataModal((prev) => ({
			...prev,
			show: true,
			title: "Switch Role",
		}));
	};

	return (
		<DropdownMenu>
			<DropdownMenuTrigger asChild>
				<Button variant="ghost" size="icon" className="rounded-full">
					<img className="h-6 w-6 rounded-full" src={header_img} alt="" />
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent className="w-56">
				<div className="flex items-center gap-2 p-2">
					<img className="h-10 w-10 rounded-full" src={header_img} alt="" />
					<div className="flex flex-col items-start">
						<div className="text-text-primary text-sm font-medium">{user_name}</div>
						<div className="text-text-secondary text-xs">{email}</div>
					</div>
				</div>
				<DropdownMenuSeparator />
				<DropdownMenuItem asChild>
					<span>{`${t("sys.menu.current_role")}: ${current_role ? current_role.name : t("common.none")}`}</span>
				</DropdownMenuItem>

				<DropdownMenuItem onClick={handleSwitch}>{t("sys.menu.user.switch_role")}</DropdownMenuItem>
				<DropdownMenuItem asChild>
					<NavLink to="/management/user/account">{t("sys.menu.user.account")}</NavLink>
				</DropdownMenuItem>
				<DropdownMenuSeparator />
				<DropdownMenuItem className="font-bold text-warning" onClick={handleLogout}>
					{t("sys.login.logout")}
				</DropdownMenuItem>
			</DropdownMenuContent>
			<SwitchModal {...switchDataModal} />
		</DropdownMenu>
	);
}
