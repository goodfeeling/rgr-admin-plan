// import { useUserInfo } from "@/store/userStore";
// import { useRoleSettingBtnIds } from "@/store/roleSettingStore";

// interface UsePermissionProps {
//   menuId?: number;
//   buttonId?: number;
// }

// export const usePermission = ({
//   menuId,
//   buttonId,
// }: UsePermissionProps = {}) => {
//   const userInfo = useUserInfo();
//   const roleBtnIds = useRoleSettingBtnIds();

//   const hasButtonPermission = (): boolean => {
//     // 如果没有指定 menuId 或 buttonId，则默认有权限
//     if (menuId === undefined || buttonId === undefined) {
//       return true;
//     }

//     // 超级管理员拥有所有权限
//     if (userInfo.role?.name === "Admin") {
//       return true;
//     }

//     // 检查当前角色是否拥有该按钮权限
//     if (roleBtnIds && roleBtnIds[menuId]) {
//       return roleBtnIds[menuId].includes(buttonId);
//     }

//     return false;
//   };

//   const hasPermission = (permissions: string[] = []): boolean => {
//     // 如果没有指定特定权限，则只检查按钮权限
//     if (!permissions.length) {
//       return hasButtonPermission();
//     }

//     // 超级管理员拥有所有权限
//     if (userInfo.role?.name === "Admin") {
//       return true;
//     }

//     // 检查用户是否拥有指定权限
//     const userPermissions = userInfo.permissions || [];
//     return permissions.some((permission) =>
//       userPermissions.some((userPerm) => userPerm.label === permission)
//     );
//   };

//   return {
//     hasPermission,
//     hasButtonPermission,
//   };
// };
