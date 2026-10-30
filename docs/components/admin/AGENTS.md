# Admin Console Agent Instructions & Engineering Protocols (`admin/`) — Dink Machinery

**Subsystem:** Next.js 16 Administrative Web Console & Operational Back-Office  
**Target:** AI Agents, Full-Stack Engineers & Front-End Architects  
**Scope:** Architecture, component composition, server-driven UI engines, state boundaries, accessibility, and CI/CD gates for `admin/`.  
**Governing Documents:** `AGENTS.md`, `TECH_STACK.md`, `CODE_STANDARDS.md`, `API_CONVENTIONS.md`, `ERROR_CATALOG.md`, `DEFINITION_OF_DONE.md`

---

## 1. Subsystem Mission & Operational Invariants

The **Next.js 16 Admin Web Console** (`admin/`) serves as the mission-critical operational cockpit for Dink Machinery platform administrators, dispute officers, and fleet compliance managers across Ethiopia and Kenya.

Any agent modifying or authoring code within `admin/` must adhere to these inviolable architectural constraints:

1. **Server-Driven Dynamic Metadata (`AD-06-BE`):** Operational tables, filter facets, sidebar menus, and actionable buttons must never be hardcoded. Layout configurations, visible columns, and RBAC actions are fetched dynamically from `/api/v1/admin/ui-config`.
2. **Contract Supremacy & Zero Manual Types:** Hand-written API DTOs or TypeScript response interfaces are strictly prohibited. All data contracts must be generated directly from `contracts/openapi.yaml` via `openapi-typescript` into `src/types/api.generated.ts`.
3. **React Server Components (RSC) by Default:** All layouts and pages under `src/app/` must be RSCs. The `'use client'` directive is confined strictly to interactive boundary leaf components (e.g., custom dropdowns, modal drawers, dynamic form controllers).
4. **RFC 9457 Problem Details Binding:** Form mutations that fail validation must map `error.invalidParams` arrays directly to field errors within `react-hook-form` using `@hookform/resolvers/zod`.
5. **Integer Minor-Unit Currency Math:** Floating-point representations for money are banned. Display amounts must be formatted from minor units (`amountMinor` int64 + `currency` code) via standard formatting utilities (e.g., Santims $\rightarrow$ ETB, Cents $\rightarrow$ KES).
6. **Zero Client-Side Token Storage (ASVS L3):** Authentication credentials and session tokens must never be stored in `localStorage`, `sessionStorage`, or window objects. Sessions must be managed exclusively via `HttpOnly`, `Secure`, `SameSite=Strict` cookies.
7. **Accessibility (WCAG 2.1 AA):** All interactive tables, drawers, and form elements must support 100% keyboard navigation, visible focus rings, ARIA live announcements, and pass automated `@axe-core/playwright` CI audits.

---

## 2. Technology Stack & Pinned Toolchain

| Layer / Concern | Technology & Version | Purpose & Package Scope |
|---|---|---|
| **Framework & Engine** | Next.js 16 (React 19) | App Router, Server Components (RSC), Server Actions (`'use server'`). |
| **Language & Types** | TypeScript 5.7+ | Strict mode (`strict: true`, `noImplicitAny: true`, `exactOptionalPropertyTypes: true`). |
| **Styling & Design Tokens** | Tailwind CSS v4 | CSS variables, utility tokens, dark mode, high-contrast outdoor theme. |
| **Data Grids & Virtualization** | `@tanstack/react-table` v8<br>`@tanstack/react-virtual` v3 | Dynamic headless grids, server-driven column models, virtual windowing. |
| **Data Fetching & Cache** | `@tanstack/react-query` v5 | Client-side cache synchronization, polling, optimistic updates. |
| **Form & Schema Validation** | `react-hook-form` v7<br>`@hookform/resolvers`<br>`zod` v3 | Type-safe form controllers, client boundary schemas, RFC 9457 binding. |
| **UI Primitives & Icons** | Radix UI / `@radix-ui/react-*`<br>`lucide-react` | Unstyled accessible primitives (dialogs, dropdowns, tooltips, popovers). |
| **Toast Notifications** | `sonner` | Accessible live region toasts displaying correlation `requestId`. |
| **Contract Synchronizer** | `openapi-typescript` v7 | Automated generation of TypeScript types from `contracts/openapi.yaml`. |

---

## 3. Directory Structure & Modular Layout

All code in `admin/` must conform to the following directory hierarchy:

```
admin/
├── src/
│   ├── app/                          # Next.js App Router (RSC by default)
│   │   ├── (auth)/                   # Authentication route group (Login, MFA challenge)
│   │   │   ├── login/
│   │   │   └── layout.tsx
│   │   ├── (dashboard)/              # Authenticated operator route group
│   │   │   ├── layout.tsx            # Main shell: Dynamic sidebar, header, breadcrumbs
│   │   │   ├── page.tsx              # Operations dashboard widgets
│   │   │   ├── fleets/               # Fleet approval queue & asset inspections
│   │   │   ├── bookings/             # Booking state machine & timeline tracker
│   │   │   ├── disputes/             # Escrow dispute resolution & evidence lightbox
│   │   │   ├── users/                # User management & KYC audit
│   │   │   └── settings/             # Country configs & role permission matrices
│   │   ├── api/                      # Edge proxy & Next.js internal route handlers
│   │   └── global-error.tsx          # Root RFC 9457 error boundary
│   ├── components/                   # Reusable UI primitives and operational blocks
│   │   ├── dynamic/                  # Server-driven UI renderers (DynamicTable, DynamicForm)
│   │   ├── layout/                   # Sidebar, Navbar, Breadcrumbs, UserNav
│   │   ├── feedback/                 # ProblemDetailsBanner, RequestIdBadge, Lightbox
│   │   └── ui/                       # Radix UI wrappers (Button, Modal, Input, Badge)
│   ├── hooks/                        # Custom client hooks (useTableConfig, useSession)
│   ├── lib/                          # Utility libraries & framework integrations
│   │   ├── api/                      # Typed HTTP client & Server Action dispatchers
│   │   ├── currency.ts               # Integer minor-unit formatters (Santims/Cents)
│   │   ├── errors.ts                 # RFC 9457 error parsers & field error binders
│   │   └── session.ts                # Server-side cookie session inspection
│   └── types/
│       ├── api.generated.ts          # Compiled OpenAPI TypeScript contracts (DO NOT EDIT)
│       └── dynamic-ui.ts             # Metadata schemas for Dynamic UI engine
├── tests/                            # Vitest component tests & Playwright E2E suites
├── tailwind.config.ts                # Tailwind CSS v4 theme configuration
└── package.json                      # Pinned dependencies & scripts
```

---

## 4. Server-Driven UI Engine (`AD-06-BE`)

To avoid recompilation and deployment cycles whenever administrative views, table columns, or permission rules change, the admin console renders data tables dynamically from `/api/v1/admin/ui-config`.

### 4.1. Metadata Response Schema
```typescript
// src/types/dynamic-ui.ts
export interface ColumnConfig {
  id: string;
  header: string;
  accessorKey: string;
  type: 'text' | 'money' | 'badge' | 'date' | 'actions';
  isSortable: boolean;
  isFilterable: boolean;
  align?: 'left' | 'center' | 'right';
  badgeVariants?: Record<string, 'default' | 'success' | 'warning' | 'destructive'>;
}

export interface TableConfig {
  tableId: string;
  title: string;
  defaultSortField: string;
  defaultSortDirection: 'asc' | 'desc';
  pageSize: number;
  columns: ColumnConfig[];
  actions: {
    id: string;
    label: string;
    requiredRole: string;
    variant: 'default' | 'destructive' | 'outline';
    endpoint: string;
    method: 'POST' | 'PUT' | 'PATCH';
  }[];
}
```

### 4.2. Headless Dynamic Grid Implementation
Dynamic tables must use `@tanstack/react-table` coupled with `@tanstack/react-virtual` to ensure zero performance degradation when displaying large datasets:

```tsx
// src/components/dynamic/DynamicTable.tsx
'use client';

import React from 'react';
import {
  useReactTable,
  getCoreRowModel,
  flexRender,
  ColumnDef,
} from '@tanstack/react-table';
import { TableConfig } from '@/types/dynamic-ui';
import { formatMinorUnits } from '@/lib/currency';

interface DynamicTableProps<TData> {
  config: TableConfig;
  data: TData[];
  onAction?: (actionId: string, row: TData) => Promise<void>;
}

export function DynamicTable<TData Record<string, any extends>>({
  config,
  data,
  onAction,
}: DynamicTableProps<TData>) {
  const columns = React.useMemo<ColumnDef<TData>[]>(() => {
    return config.columns.map((col) => ({
      id: col.id,
      header: col.header,
      accessorKey: col.accessorKey,
      cell: ({ getValue }) => {
        const val = getValue();
        if (col.type === 'money') {
          return formatMinorUnits(val.amountMinor, val.currency);
        }
        if (col.type === 'date') {
          return new Date(val).toLocaleDateString();
        }
        return String(val ?? '');
      },
    }));
  }, [config.columns]);

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="w-full overflow-x-auto rounded-md border border-border">
      <table className="w-full text-left text-sm">
        <thead className="bg-muted text-muted-foreground uppercase text-xs">
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <th key={header.id} className="px-4 py-3 font-medium">
                  {flexRender(header.column.columnDef.header, header.getContext())}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody className="divide-y divide-border">
          {table.getRowModel().rows.map((row) => (
            <tr key={row.id} className="hover:bg-muted/50 transition-colors">
              {row.getVisibleCells().map((cell) => (
                <td key={cell.id} className="px-4 py-3">
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
```

---

## 5. Server Actions, Form Validation & RFC 9457 Binding

All operational state changes (e.g., approving an equipment listing, resolving an escrow dispute, updating tax configurations) must execute via React Server Actions with strict Zod validation.

### 5.1. Server Action Dispatcher Pattern
```typescript
// src/app/(dashboard)/fleets/actions.ts
'use server';

import { revalidatePath } from 'next/cache';
import { z } from 'zod';
import { ProblemDetails } from '@/types/dynamic-ui';

const ApproveFleetSchema = z.object({
  fleetId: z.string().uuid(),
  version: z.number().int().positive(),
  notes: z.string().max(500).optional(),
});

export type ActionState = {
  success: boolean;
  problem?: ProblemDetails;
  data?: any;
};

export async function approveFleetAction(
  prevState: ActionState,
  formData: FormData
): Promise<ActionState> {
  const parsed = ApproveFleetSchema.safeParse({
    fleetId: formData.get('fleetId'),
    version: Number(formData.get('version')),
    notes: formData.get('notes'),
  });

  if (!parsed.success) {
    return {
      success: false,
      problem: {
        type: '[https://api.dinkmachinery.com/errors/request-validation-failed](https://api.dinkmachinery.com/errors/request-validation-failed)',
        title: 'Validation Failed',
        status: 400,
        detail: 'Invalid fleet approval submission parameters.',
        code: 'REQ_VALIDATION_FAILED',
        invalidParams: parsed.error.issues.map((issue) => ({
          name: issue.path.join('.'),
          reason: issue.message,
        })),
      },
    };
  }

  const response = await fetch(
    `${process.env.INTERNAL_API_URL}/api/v1/admin/fleets/${parsed.data.fleetId}:approve`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Idempotency-Key': crypto.randomUUID(),
      },
      body: JSON.stringify(parsed.data),
    }
  );

  const payload = await response.json();

  if (!response.ok) {
    return {
      success: false,
      problem: payload.error,
    };
  }

  revalidatePath('/fleets');
  return { success: true, data: payload.data };
}
```

### 5.2. RFC 9457 Form Error Binding
When a Server Action returns an RFC 9457 error payload containing `invalidParams`, bind the specific field errors into `react-hook-form`:

```tsx
// src/lib/errors.ts
import { FieldValues, UseFormSetError } from 'react-hook-form';
import { ProblemDetails } from '@/types/dynamic-ui';

export function bindProblemDetailsErrors<T FieldValues extends>(
  problem: ProblemDetails,
  setError: UseFormSetError<T>
) {
  if (!problem.invalidParams || problem.invalidParams.length === 0) {
    return;
  }

  problem.invalidParams.forEach((param) => {
    setError(param.name as any, {
      type: 'server',
      message: param.reason,
    });
  });
}
```

---

## 6. Security, Zero-PII & Access Control Guidelines

### 6.1. Client PII Sanitization
- Never render raw Tax Identification Numbers (TIN), Kebele cards, or bank accounts directly on unmasked tables.
- Use masking formatters in UI tables:
  ```typescript
  export function maskTIN(tin: string): string {
    if (!tin || tin.length < 4) return '****';
    return `TIN: *****${tin.slice(-4)}`;
  }
  ```

### 6.2. Toast Notification Governance (`sonner`)
- Global failures must display a non-blocking toast presenting the human-readable `detail` along with the correlation `requestId`:
  ```tsx
  toast.error(problem.title, {
    description: `${problem.detail} (Ref: ${requestId})`,
    duration: 6000,
  });
  ```

---

## 7. Automated Testing, Accessibility & Linting Gates

Every modification to `admin/` must clear the following local and CI verification checks before a pull request can be merged:

```bash
# 1. Typecheck & OpenAPI Contract Sync
cd admin
pnpm dlx openapi-typescript ../contracts/openapi.yaml -o src/types/api.generated.ts
pnpm run lint

# 2. Production Build Verification (Zero TypeScript/RSC warnings)
pnpm run build

# 3. Unit & Component Tests (Vitest)
pnpm run test:unit

# 4. Automated Accessibility Audit (Playwright + Axe-Core)
pnpm dlx playwright test tests/a11y/
```

### Pull Request Review Criteria for `admin/`
- [ ] No hardcoded tables or columns; layout derived from dynamic configuration (`AD-06-BE`).
- [ ] TypeScript types generated directly from `contracts/openapi.yaml`.
- [ ] Components use RSC by default; `'use client'` restricted to interactive leaves.
- [ ] Monetary displays format integer minor units (`amountMinor` int64); zero floating-point arithmetic.
- [ ] Forms bind RFC 9457 `invalidParams` into field error states.
- [ ] Complete keyboard accessibility with visible focus rings; passes Axe scans with 0 violations.
- [ ] Authentication handled exclusively via `HttpOnly` cookies; zero tokens in `localStorage`.
- [ ] Production build passes cleanly: `pnpm run build` succeeds with zero errors.