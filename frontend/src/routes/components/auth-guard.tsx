import PermissionGuard from "@/components/premission/guard";
import { useUserToken } from "@/store/userStore";
import { useCallback, useEffect } from "react";
import { useRouter } from "../hooks";

type Props = {
	children: React.ReactNode;
};
export default function AuthGuard({ children }: Props) {
	const router = useRouter();
	const { accessToken } = useUserToken();

	const check = useCallback(() => {
		if (!accessToken) {
			router.replace("/auth/login");
		}
	}, [router, accessToken]);

	// biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
	useEffect(() => {
		check();
	}, [check, accessToken]);

	return <PermissionGuard fallback={<div> No Permission!</div>}>{children}</PermissionGuard>;
}
