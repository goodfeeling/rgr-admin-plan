import apiClient from "@/api/apiClient";

export class CaptchaService {
	/**
	 * 生成验证码
	 * @param id
	 */
	generate(id: string) {
		return apiClient.get<CaptchaService.CaptchaGenerateResult>({
			url: `${CaptchaService.Client.Generate}?captcha_id=${id}`,
		});
	}
}

export namespace CaptchaService {
	export enum Client {
		Generate = "/captcha/generate",
	}

	export interface CaptchaGenerateResult {
		id: string;
		b64s: string; // base64 image string
		config: {
			width: number;
			height: number;
			length: number;
		};
	}
}

export default new CaptchaService();
