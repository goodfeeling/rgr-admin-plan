import { useRouter } from "@/routes/hooks";
import { Button } from "@/ui/button";
import { ArrowLeft } from "lucide-react";
import { FC } from "react";

interface BackButtonProps {
  /**
   * 可选的自定义点击处理函数
   * 如果提供，将替代默认的后退行为
   */
  onClick?: () => void;

  /**
   * 按钮文本，默认为"返回"
   */
  text?: string;

  /**
   * 是否只显示图标，不显示文本
   */
  iconOnly?: boolean;
}

/**
 * 返回上一页按钮组件
 * 使用浏览器历史记录API或自定义回调函数实现返回功能
 */
const BackButton: FC<BackButtonProps> = ({
  onClick,
  text = "返回",
  iconOnly = false,
}) => {
  const router = useRouter();

  const handleClick = () => {
    if (onClick) {
      onClick();
    } else {
      router.back();
    }
  };

  return (
    <Button
      variant="outline"
      onClick={handleClick}
      className="flex items-center gap-2"
    >
      <ArrowLeft size={16} />
      {!iconOnly && text}
    </Button>
  );
};

export default BackButton;
