import taskExecutionLogService, { TaskExecutionLogService } from "@/api/services/taskExecutionLogService";
import { useSharedWebSocket } from "@/hooks/sharedWebSocket";
import type { TaskExecutionLog } from "@/types/entity";
import { Badge } from "@/ui/badge";
import { Button } from "@/ui/button";
import { Card, CardTitle } from "@/ui/card";
import { ScrollArea } from "@/ui/scroll-area";
import { Modal } from "antd";
import { format } from "date-fns";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

export type LogModalProps = {
	title: string;
	show: boolean;
	id: number;
	onCancel?: () => void;
};

const LogPage = ({ title, show, onCancel, id }: LogModalProps) => {
	const { t } = useTranslation();
	const [expandedLogs, setExpandedLogs] = useState<Record<number, boolean>>({});
	// 存储实时日志
	const [realTimeLogs, setRealTimeLogs] = useState<TaskExecutionLog[]>([]);
	// 存储历史日志（用于分页）
	const [historicalLogs, setHistoricalLogs] = useState<TaskExecutionLog[]>([]);
	const [currentPage, setCurrentPage] = useState(1);
	const [totalLogs, setTotalLogs] = useState(0);
	const [loading, setLoading] = useState(false);
	// 控制显示模式：实时模式或历史模式
	const [viewMode, setViewMode] = useState<"realtime" | "historical">("realtime");
	const { connected, message, sendMessage } = useSharedWebSocket(TaskExecutionLogService.Client.WsUri);

	// 发送获取数据请求
	useEffect(() => {
		if (show && id) {
			// 延迟发送消息，确保WebSocket连接已建立
			const timer = setTimeout(() => {
				sendMessage(JSON.stringify({ taskId: id, limit: 100 }));
			}, 100);

			return () => clearTimeout(timer);
		}
	}, [show, sendMessage, id]);

	// 处理接收到的消息
	useEffect(() => {
		if (!show) return;

		// 实时接收新日志
		if (message?.data) {
			if (message?.type === "log_update") {
				setRealTimeLogs((prevLogs) => [message.data, ...prevLogs].slice(0, 100));
			} else {
				setRealTimeLogs(message.data);
			}
		}

		// 加载第一页历史日志
		loadHistoricalLogs(6);
	}, [message, show]);

	// 加载历史日志
	const loadHistoricalLogs = async (page: number) => {
		setLoading(true);
		try {
			const response = await taskExecutionLogService.searchPageList(`taskId_match=${id}&page=${page}&pageSize=20`);
			setHistoricalLogs(response.list);
			setTotalLogs(response.total);
			setLoading(false);
			console.log("load data with", page, "page history logs");
		} catch (error) {
			console.error("load history log error:", error);
		} finally {
			setLoading(false);
		}
	};

	const getStatusVariant = (result: number) => {
		if (result === 1) return "default";
		if (result === 0) return "destructive";
		return "secondary";
	};

	const handleCancel = () => {
		setRealTimeLogs([]);
		setHistoricalLogs([]);
		setExpandedLogs({});
		onCancel?.(); // 调用父组件的 onCancel 回调
		setCurrentPage(6);
	};

	// 判断文本是否需要折叠显示
	const shouldTruncate = (text: string, maxLength = 100) => {
		return text && text.length > maxLength;
	};

	// 截取文本
	const truncateText = (text: string, maxLength = 100) => {
		return text ? `${text.substring(0, maxLength)}...` : "";
	};

	const toggleExpand = (logId: number) => {
		setExpandedLogs((prev) => ({
			...prev,
			[logId]: !prev[logId],
		}));
	};

	// 获取当前显示的日志
	const currentLogs = viewMode === "realtime" ? realTimeLogs : historicalLogs;

	// 如果模态框不显示，则不渲染任何内容
	if (!show) {
		return null;
	}

	return (
		<Modal
			width={800}
			open={show}
			title={
				<div className="flex items-center gap-2">
					<CardTitle>{title}</CardTitle>
					<Badge variant={connected ? "default" : "destructive"}>
						{connected ? t("table.handle_message.connected") : t("table.handle_message.disconnected")}
					</Badge>
				</div>
			}
			onCancel={handleCancel}
			centered
			footer={false}
		>
			{/* 添加视图切换 */}
			<div className="flex justify-between items-center mb-2">
				<div className="flex gap-2">
					<Button
						key="realtime-button"
						variant={viewMode === "realtime" ? "default" : "outline"}
						size="sm"
						onClick={() => setViewMode("realtime")}
					>
						{t("table.handle_message.realtime")}
					</Button>
					<Button
						key="historical-button"
						variant={viewMode === "historical" ? "default" : "outline"}
						size="sm"
						onClick={() => setViewMode("historical")}
					>
						{t("table.handle_message.historical")}
					</Button>
				</div>
				{viewMode === "historical" && (
					<div className="text-sm text-muted-foreground">
						{`${t("table.page.total")} ${totalLogs} ${t("table.page.items")}`}
					</div>
				)}
			</div>

			<ScrollArea className="h-[calc(100vh-220px)]">
				<div className="space-y-4">
					{currentLogs.length === 0 ? (
						<div className="text-center py-8 text-muted-foreground" key="no-logs">
							{viewMode === "realtime"
								? connected
									? t("table.handle_message.waiting_load_data")
									: t("table.handle_message.connecting")
								: loading
									? t("table.handle_message.loading")
									: t("table.handle_message.no_history_log")}
						</div>
					) : (
						<>
							{currentLogs.map((log) => {
								const isErrorMessageLong = shouldTruncate(log.error_message, 100);
								return (
									<Card key={`log-card-${log.id}`} className="p-4">
										<div className="flex justify-between items-start">
											<div>
												<div className="text-sm text-muted-foreground mt-1">
													{t("table.columns.task_log.execute_time")}:{" "}
													{format(new Date(log.execute_time), "yyyy-MM-dd HH:mm:ss")}
												</div>
											</div>
											<Badge variant={getStatusVariant(log.execute_result)}>
												{log.execute_result === 1
													? t("table.handle_message.execution_success")
													: t("table.handle_message.execution_failed")}
											</Badge>
										</div>
										<div className="mt-3 text-sm">
											<div className="font-medium">{t("table.columns.task_log.execute_result")}:</div>
											<div className="mt-1 whitespace-pre-wrap">
												{log.execute_result === 1
													? t("table.handle_message.execution_success")
													: t("table.handle_message.execution_failed")}
											</div>
										</div>
										{log.error_message && (
											<div className="mt-2 text-sm text-destructive">
												<div className="font-medium">{t("table.columns.task_log.error_message")}:</div>
												<div className="mt-1 whitespace-pre-wrap">
													{isErrorMessageLong && !expandedLogs[log.id] ? (
														<>
															{truncateText(log.error_message, 100)}
															<Button
																key={`expand-button-${log.id}`}
																variant="link"
																size="sm"
																className="ml-2 h-auto p-0 text-gray-500"
																onClick={() => toggleExpand(log.id)}
															>
																{t("table.handle_message.show_more")}
															</Button>
														</>
													) : (
														<>
															{log.error_message}
															{isErrorMessageLong && (
																<Button
																	key={`collapse-button-${log.id}`}
																	variant="link"
																	size="sm"
																	className="ml-2 h-auto p-0 text-gray-500"
																	onClick={() => toggleExpand(log.id)}
																>
																	{t("table.handle_message.show_less")}
																</Button>
															)}
														</>
													)}
												</div>
											</div>
										)}
										<div className="mt-3 text-xs text-muted-foreground flex justify-between">
											<span>
												{t("table.columns.task_log.execute_duration")}: {log.execute_duration}ms
											</span>
											<span>
												{t("table.columns.task_log.created_at")}:{" "}
												{format(new Date(log.created_at), "yyyy-MM-dd HH:mm:ss")}
											</span>
										</div>
									</Card>
								);
							})}

							{/* 历史日志分页 */}
							{viewMode === "historical" && totalLogs > 0 && (
								<div className="flex justify-center mt-4" key="pagination">
									{/* 这里添加分页组件 */}
									<div className="flex gap-2">
										{/* 示例分页按钮 */}
										<Button
											key="prev-button"
											variant="outline"
											size="sm"
											disabled={currentPage === 1}
											onClick={() => {
												const newPage = currentPage - 1;
												setCurrentPage(newPage);
												loadHistoricalLogs(newPage);
											}}
										>
											{t("table.handle_message.previous")}
										</Button>
										<span className="self-center text-sm" key="page-info">
											{`${t("table.handle_message.page")}  ${currentPage}`}
										</span>
										<Button
											key="next-button"
											variant="outline"
											size="sm"
											// 这里需要根据实际总页数判断是否禁用
											onClick={() => {
												const newPage = currentPage + 1;
												setCurrentPage(newPage);
												loadHistoricalLogs(newPage);
											}}
										>
											{t("table.handle_message.next")}
										</Button>
									</div>
								</div>
							)}
						</>
					)}
				</div>
			</ScrollArea>
		</Modal>
	);
};

export default LogPage;
