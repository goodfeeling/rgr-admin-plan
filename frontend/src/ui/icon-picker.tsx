// @/components/icon-picker-custom.tsx
import { Icon } from "@/components/icon";
import { useDictionaryByTypeWithCache } from "@/hooks/dict";
import { useTheme } from "@/theme/hooks";
import { useEffect, useRef, useState } from "react";

export type IconPickerProps = {
  value?: string;
  onChange: (icon: string) => void;
};

export const IconPicker = ({ value, onChange }: IconPickerProps) => {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const containerRef = useRef<HTMLDivElement>(null);
  const { themeTokens } = useTheme();
  const { data: icons } = useDictionaryByTypeWithCache("icons");
  // 设置进度条颜色，优先使用传入的颜色，否则使用主题色
  const backgroundColor = themeTokens.color.background.default;
  const filteredIcons = search
    ? icons
        ?.filter((icon) =>
          icon.value.toLowerCase().includes(search.toLowerCase())
        )
        .map((icon) => icon.value)
    : icons?.map((icon) => icon.value);

  // 点击外部关闭弹窗
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(e.target as Node)
      ) {
        setOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div ref={containerRef} className="relative inline-block w-full">
      <button
        type="button"
        className="flex w-full items-center justify-between gap-2 rounded-md border border-input bg-background px-3 py-2 text-sm"
        onClick={() => setOpen(!open)}
      >
        <div className="flex items-center gap-2 truncate">
          {value && <Icon icon={value} size={18} />}
          <span>{value || "Select Icon"}</span>
        </div>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          className={`transition-transform duration-200 ${
            open ? "rotate-180" : ""
          }`}
        >
          <polyline points="6 9 12 15 18 9" />
        </svg>
      </button>

      {/* 弹窗内容 */}
      {open && (
        <div
          className="absolute z-50 mt-1 w-full rounded-md  border p-2 shadow-md"
          style={{ backgroundColor: backgroundColor }}
        >
          <input
            type="text"
            placeholder="Search icon..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full rounded-md border px-3 py-1 text-sm"
          />

          <div className="mt-2 max-h-60 overflow-y-auto pr-1">
            {filteredIcons?.length === 0 ? (
              <div className="py-2 text-center text-sm text-muted-foreground">
                No icon found.
              </div>
            ) : (
              <div className="grid grid-cols-4 gap-2">
                {filteredIcons?.map((icon) => (
                  <div
                    key={icon}
                    className="flex cursor-pointer flex-col items-center justify-center rounded p-2 hover:bg-gray-200 dark:hover:bg-gray-800"
                    onClick={() => {
                      onChange(icon);
                      setOpen(false);
                    }}
                  >
                    <Icon icon={icon} size={16} />
                    <span className="text-xs" style={{ textAlign: "center" }}>
                      {icon}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};
