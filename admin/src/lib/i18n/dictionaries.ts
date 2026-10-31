export type Locale = "en" | "am";

export interface TranslationDictionary {
  appTitle: string;
  appDescription: string;
  connectedGateway: string;
  navDashboard: string;
  navFleets: string;
  navBookings: string;
  navApprovals: string;
  navDisputes: string;
  navFinance: string;
  navSettings: string;
  actionApprove: string;
  actionReject: string;
  actionSuspend: string;
  actionExport: string;
  permissionDenied: string;
}

export const dictionaries: Record<Locale, TranslationDictionary> = {
  en: {
    appTitle: "Dink Machinery Admin",
    appDescription: "Operations & Governance Console for Dink Machinery Marketplace.",
    connectedGateway: "Connected to gateway",
    navDashboard: "Dashboard",
    navFleets: "Fleet Inventory",
    navBookings: "Bookings",
    navApprovals: "Compliance Approvals",
    navDisputes: "Disputes & Arbitrations",
    navFinance: "Commissions & Payouts",
    navSettings: "System Configuration",
    actionApprove: "Approve Listing",
    actionReject: "Reject Listing",
    actionSuspend: "Suspend Asset",
    actionExport: "Export CSV",
    permissionDenied: "You do not hold permissions to access this control.",
  },
  am: {
    appTitle: "የድንቅ ማሽነሪ አስተዳዳሪ ፖርታል",
    appDescription: "የድንቅ ማሽነሪ የገበያ ቦታ ስራዎች እና አስተዳደር መቆጣጠሪያ።",
    connectedGateway: "ከኔትወርክ መስመር ጋር ተገናኝቷል",
    navDashboard: "ዳሽቦርድ",
    navFleets: "የማሽነሪዎች ዝርዝር",
    navBookings: "የኪራይ ቦታ ማስያዣዎች",
    navApprovals: "የማረጋገጫ ጥያቄዎች",
    navDisputes: "የክርክር አፈታት",
    navFinance: "ኮሚሽንና ክፍያዎች",
    navSettings: "የሲስተም ማስተካከያ",
    actionApprove: "ፈቃድ ስጥ",
    actionReject: "ውድቅ አድርግ",
    actionSuspend: "አግድ",
    actionExport: "መረጃ አውርድ (CSV)",
    permissionDenied: "ይህን መቆጣጠሪያ ለመጠቀም የሚያስችል ፈቃድ የለዎትም።",
  },
};

export function getDictionary(locale: Locale = "en"): TranslationDictionary {
  return dictionaries[locale] || dictionaries.en;
}
