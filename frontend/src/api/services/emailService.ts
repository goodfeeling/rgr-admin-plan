import apiClient from "../apiClient";

class EmailService {
	/**
	 * 发送忘记密码邮件
	 * @param email
	 */
	sendForgetPassword(email: string) {
		return apiClient.post<boolean>({
			url: `${EmailService.Client.Email}/forget-password`,
			data: { email },
		});
	}
}

namespace EmailService {
	export enum Client {
		Email = "/email",
	}
}

export default new EmailService();
