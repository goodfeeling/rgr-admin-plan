declare module "ali-oss" {
	export default class OSS {
		constructor(options: {
			region: string;
			accessKeyId: string;
			accessKeySecret: string;
			stsToken?: string;
			bucket?: string;
		});

		put(
			name: string,
			file: File | string | Buffer,
			options?: object,
		): Promise<{
			url: string;
			name: string;
			[key: string]: any;
		}>;
	}
}
