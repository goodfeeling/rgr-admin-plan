import schedule_taskClient from "../apiClient";

import type { PageList, ScheduledTask } from "#/entity";

export class ScheduledTaskService {
	/**
	 * 获取定时任务列表
	 * @returns
	 */
	getScheduledTasks() {
		return schedule_taskClient.get<ScheduledTask[]>({
			url: `${ScheduledTaskService.Client.ScheduledTask}`,
		});
	}

	/**
	 * 更新定时任务
	 * @param id
	 * @param schedule_taskInfo
	 * @returns
	 */
	updateScheduledTask(id: number, schedule_taskInfo: ScheduledTask) {
		return schedule_taskClient.put<ScheduledTask>({
			url: `${ScheduledTaskService.Client.ScheduledTask}/${id}`,
			data: schedule_taskInfo,
		});
	}

	/**
	 * 创建定时任务
	 * @param schedule_taskInfo
	 * @returns
	 */
	createScheduledTask(schedule_taskInfo: ScheduledTask) {
		return schedule_taskClient.post<ScheduledTask>({
			url: `${ScheduledTaskService.Client.ScheduledTask}`,
			data: schedule_taskInfo,
		});
	}

	/**
	 * 搜索定时任务
	 * @param searchStr
	 * @returns
	 */
	searchPageList(searchStr: string) {
		return schedule_taskClient.get<PageList<ScheduledTask>>({
			url: `${ScheduledTaskService.Client.SearchScheduledTask}?${searchStr}`,
		});
	}

	/**
	 * 删除定时任务
	 * @param id
	 * @returns
	 */
	deleteScheduledTask(id: number) {
		return schedule_taskClient.delete<string>({
			url: `${ScheduledTaskService.Client.ScheduledTask}/${id}`,
		});
	}

	/**
	 * 批量删除定时任务
	 * @param ids
	 * @returns
	 */
	deleteBatch(ids: number[]) {
		return schedule_taskClient.post<number>({
			url: `${ScheduledTaskService.Client.DeleteBatch}`,
			data: { ids },
		});
	}

	/**
	 * 启用定时任务
	 * @param id
	 * @returns
	 */
	enableTask(id: number) {
		return schedule_taskClient.post<string>({
			url: `${ScheduledTaskService.Client.ScheduledTask}/enable/${id}`,
		});
	}

	/**
	 * 禁用定时任务
	 * @param id
	 * @returns
	 */
	disableTask(id: number) {
		return schedule_taskClient.post<string>({
			url: `${ScheduledTaskService.Client.ScheduledTask}/disable/${id}`,
		});
	}

	/**
	 * 重新加载定时任务
	 * @returns
	 */
	reloadTask() {
		return schedule_taskClient.post<string>({
			url: `${ScheduledTaskService.Client.ScheduledTask}/reload`,
		});
	}
}

export namespace ScheduledTaskService {
	export enum Client {
		ScheduledTask = "/scheduled_task",
		SearchScheduledTask = "/scheduled_task/search",
		DeleteBatch = "/scheduled_task/batch",
	}
}

export default new ScheduledTaskService();
