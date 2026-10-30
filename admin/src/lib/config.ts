export type Environment = "development" | "staging" | "production";

export interface AppConfig {
  environment: Environment;
  apiUrl: string;
  appTitle: string;
  enableDebugBanner: boolean;
  enableAuthMock: boolean;
  isDevelopment: boolean;
  isStaging: boolean;
  isProduction: boolean;
}

const env = (process.env.NEXT_PUBLIC_ENVIRONMENT || "development") as Environment;

export const appConfig: AppConfig = {
  environment: env,
  apiUrl: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1",
  appTitle: process.env.NEXT_PUBLIC_APP_TITLE || "Dink Machinery Admin",
  enableDebugBanner: process.env.NEXT_PUBLIC_ENABLE_DEBUG_BANNER === "true" || env !== "production",
  enableAuthMock: process.env.NEXT_PUBLIC_ENABLE_AUTH_MOCK === "true" && env === "development",
  isDevelopment: env === "development",
  isStaging: env === "staging",
  isProduction: env === "production",
};
