import { appConfig } from "@/lib/config";
import { getDictionary } from "@/lib/i18n/dictionaries";
import { PermissionGuard } from "@/lib/rbac/permission-guard";

export default function Home() {
  const dict = getDictionary("en");

  // Simulated operational permissions for active back-office role
  const mockPermissions = ["fleets:list", "fleets:approve", "dashboard:read"];

  return (
    <div className="flex flex-col min-h-screen bg-zinc-50 font-sans dark:bg-black">
      {appConfig.enableDebugBanner && (
        <div className="w-full bg-amber-500 text-black py-1 px-4 text-xs font-mono font-bold uppercase tracking-wider text-center">
          Environment: {appConfig.environment} &bull; API: {appConfig.apiUrl}
        </div>
      )}
      <main className="flex flex-1 flex-col items-center justify-center p-8">
        <div className="max-w-xl text-center space-y-4">
          <h1 className="text-3xl font-bold tracking-tight text-zinc-900 dark:text-zinc-50">
            {dict.appTitle}
          </h1>
          <p className="text-zinc-600 dark:text-zinc-400">
            {dict.appDescription}
          </p>
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-mono bg-zinc-200 dark:bg-zinc-800 text-zinc-800 dark:text-zinc-200">
            <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
            {dict.connectedGateway}: {appConfig.environment}
          </div>

          {/* Dynamic RBAC Action Guard Demonstration (ADR-018) */}
          <div className="pt-6 border-t border-zinc-200 dark:border-zinc-800 flex justify-center gap-3">
            <PermissionGuard
              permission="fleets:approve"
              userPermissions={mockPermissions}
            >
              <button className="px-4 py-2 bg-emerald-600 text-white text-sm font-semibold rounded-lg shadow hover:bg-emerald-700 transition">
                {dict.actionApprove}
              </button>
            </PermissionGuard>

            <PermissionGuard
              permission="disputes:resolve"
              userPermissions={mockPermissions}
              fallback={
                <button
                  disabled
                  title={dict.permissionDenied}
                  className="px-4 py-2 bg-zinc-300 dark:bg-zinc-800 text-zinc-500 text-sm font-semibold rounded-lg cursor-not-allowed"
                >
                  {dict.navDisputes} (Locked)
                </button>
              }
            >
              <button className="px-4 py-2 bg-indigo-600 text-white text-sm font-semibold rounded-lg shadow">
                {dict.navDisputes}
              </button>
            </PermissionGuard>
          </div>
        </div>
      </main>
    </div>
  );
}
