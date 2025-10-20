import task_execution_logClient from "../apiClient";

import type { PageList, TaskExecutionLog } from "#/entity";

export class TaskExecutionLogService {
	/**
	 * 搜索任务执行日志分页列表
	 * @param searchStr 搜索字符串
	 */
	searchPageList(searchStr: string) {
		return task_execution_logClient.get<PageList<TaskExecutionLog>>({
			url: `${TaskExecutionLogService.Client.SearchTaskExecutionLog}?${searchStr}`,
		});
	}
}

export namespace TaskExecutionLogService {
	export enum Client {
		TaskExecutionLog = "/task_execution_log",
		SearchTaskExecutionLog = "/task_execution_log/search",
		WsUri = "/ws/scheduleLog",
	}
}

export default new TaskExecutionLogService();
