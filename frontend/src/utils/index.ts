import type { TableParams } from "@/types/entity";
import type { AnyObject } from "antd/es/_util/type";
import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function toURLSearchParams<T extends AnyObject>(record: T) {
	const params = new URLSearchParams();
	for (const [key, value] of Object.entries(record)) {
		params.append(key, value);
	}
	return params;
}

export function getRandomUserParams(
	params: TableParams,
	callback?: (params: Record<string, any>, searchParams: { [key: string]: any }) => void,
) {
	const { pagination, filters, sortField, sortOrder, searchParams, ...restParams } = params;
	const result: Record<string, any> = {};

	// https://github.com/mockapi-io/docs/wiki/Code-examples#pagination
	result.pageSize = pagination?.pageSize;
	result.page = pagination?.current;

	// https://github.com/mockapi-io/docs/wiki/Code-examples#filtering
	if (filters) {
		for (const [key, value] of Object.entries(filters)) {
			if (value !== undefined && value !== null) {
				result[`${key}_match`] = value;
			}
		}
	}

	// https://github.com/mockapi-io/docs/wiki/Code-examples#sorting
	if (sortField) {
		result.sortBy = sortField;
		result.sortDirection = sortOrder === "ascend" ? "asc" : "desc";
	}

	// 处理其他参数
	for (const [key, value] of Object.entries(restParams)) {
		if (value !== undefined && value !== null) {
			result[key] = value;
		}
	}

	if (callback && searchParams) callback(result, searchParams);

	return result;
}
export const parseUriFromUrl = (fileUrl: string): string => {
	try {
		const url = new URL(fileUrl);
		return url.pathname.startsWith("/") ? url.pathname.substring(1) : url.pathname;
	} catch (error) {
		// 如果不是有效的URL，返回原始字符串或处理错误
		console.error("Invalid URL:", error);
		return fileUrl;
	}
};
