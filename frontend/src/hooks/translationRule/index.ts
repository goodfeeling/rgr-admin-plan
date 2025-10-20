import { useTranslation } from "react-i18next";

export function useTranslationRule(field: string) {
	const { t } = useTranslation();
	return `${field} ${t("table.handle_message.is_required")}`;
}
