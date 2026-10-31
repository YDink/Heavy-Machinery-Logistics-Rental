import React, { ReactNode } from "react";

export interface PermissionGuardProps {
  permission: string;
  userPermissions?: string[];
  fallback?: ReactNode;
  children: ReactNode;
}

/**
 * PermissionGuard component conforming to ADR-018.
 * Renders children only if the specified permission key is present in userPermissions.
 */
export function PermissionGuard({
  permission,
  userPermissions = [],
  fallback = null,
  children,
}: PermissionGuardProps) {
  const hasPermission =
    userPermissions.includes(permission) ||
    userPermissions.includes("*") ||
    userPermissions.includes("super_admin");

  if (!hasPermission) {
    return <>{fallback}</>;
  }

  return <>{children}</>;
}
