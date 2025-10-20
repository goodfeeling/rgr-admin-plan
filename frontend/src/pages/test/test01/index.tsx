import PermissionButton from "@/components/premission/button";
import { useState } from "react";
const SelectDemo = () => {
	const [count, setCount] = useState(0);
	const [inputValue, setInputValue] = useState("");
	return (
		<>
			<PermissionButton
				permissionString={"test"}
				onClick={() => {
					console.log("hello world");
				}}
				className="btn btn-primary"
				variant="link"
			>
				编辑
			</PermissionButton>

			<div className="p-4">
				<h2>Keep Alive 测试页面</h2>
				<p>这个页面用于测试 keep alive 功能</p>

				<div className="mt-4">
					<p>计数器: {count}</p>
					{/* biome-ignore lint/a11y/useButtonType: <explanation> */}
					<button onClick={() => setCount((c) => c + 1)} className="bg-blue-500 text-white px-4 py-2 rounded mr-2">
						增加计数
					</button>
					{/* biome-ignore lint/a11y/useButtonType: <explanation> */}
					<button onClick={() => setCount(0)} className="bg-gray-500 text-white px-4 py-2 rounded">
						重置计数
					</button>
				</div>

				<div className="mt-4">
					<label htmlFor="inputTest" className="block mb-2">
						输入测试:
					</label>
					<input
						id="inputTest"
						type="text"
						value={inputValue}
						onChange={(e) => setInputValue(e.target.value)}
						className="border p-2 rounded w-full"
						placeholder="输入一些文字，然后切换页面再回来"
					/>
				</div>

				<div className="mt-4">
					<p>当前时间: {new Date().toLocaleTimeString()}</p>
				</div>
			</div>
		</>
	);
};

export default SelectDemo;
