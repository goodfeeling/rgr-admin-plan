import configService from "@/api/services/configService";
import type { ConfigResponse } from "@/types/entity";
import { Card, Tabs, message } from "antd";
import { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import ConfigBox from "./config";
function handleData(cr: ConfigResponse) {
	if (cr.data.cors) {
		cr.data.cors.allowed_methods = cr.data.cors.allowed_methods.split(",");
	}

	return cr;
}
const ConfigTabs = () => {
	const [formData, setFormData] = useState<ConfigResponse>({ data: {} });
	const configKeys = Object.keys(formData.data);
	const { t } = useTranslation();
	const fetchConfigData = useCallback(async () => {
		try {
			const response = await configService.getConfigs();
			setFormData(handleData(response));
		} catch (error) {
			console.error("Failed to fetch config data:", error);
			message.error("获取配置数据失败");
		}
	}, []);

	useEffect(() => {
		fetchConfigData();
	}, [fetchConfigData]);

	const items = configKeys.map((key) => {
		const configData = formData.data[key];
		return {
			label: t(`sys.config.tabs.${key}`, { defaultValue: key }),
			key: key,
			children: <ConfigBox configData={configData} module={key} />,
		};
	});

	return (
		<div className="p-6">
			<Card>
				<h1 className="text-2xl font-bold mb-6">{t("sys.menu.system.setting")}</h1>
				<Tabs defaultActiveKey={configKeys[0]} items={items} tabPosition="top" />
			</Card>
		</div>
	);
};

export default ConfigTabs;
