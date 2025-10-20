import { useMenu, useMenuActions, useMenuError, useMenuLoading } from "@/store/useMenuStore";
import type { MenuTreeUserGroup } from "@/types/entity";
import { useEffect, useState } from "react";

type UseAppMenuResult = {
	menuData: MenuTreeUserGroup[];
	loading: boolean;
	error: Error | null;
};
function useAppMenu(): UseAppMenuResult {
	const menuData = useMenu();
	const fetchMenu = useMenuActions().fetchMenu;
	const loading = useMenuLoading();
	const errorString = useMenuError();
	const [error, setError] = useState<Error | null>(null);

	useEffect(() => {
		fetchMenu().catch((err) => {
			setError(err instanceof Error ? err : new Error("Unknown error"));
		});
	}, [fetchMenu]);

	return {
		menuData,
		loading,
		error: error ?? (errorString ? new Error(errorString) : null),
	};
}

export default useAppMenu;
