import { appConfig } from "@/lib/config";

export default function Home() {
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
            {appConfig.appTitle}
          </h1>
          <p className="text-zinc-600 dark:text-zinc-400">
            Operations &amp; Governance Console for Dink Machinery Marketplace.
          </p>
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-mono bg-zinc-200 dark:bg-zinc-800 text-zinc-800 dark:text-zinc-200">
            <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
            Connected to {appConfig.environment} gateway
          </div>
        </div>
      </main>
    </div>
  );
}
