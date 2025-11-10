import type { VariantProps } from "class-variance-authority"
import { cva } from "class-variance-authority"

export { default as Badge } from "./Badge.vue"

export const badgeVariants = cva(
    "inline-flex items-center justify-center rounded-md border px-2 py-0.5 text-xs font-medium w-fit whitespace-nowrap shrink-0 [&>svg]:size-3 gap-1 [&>svg]:pointer-events-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive transition-[color,box-shadow] overflow-hidden",
    {
        variants: {
            variant: {
                default:
                    "border-transparent bg-primary text-primary-foreground [a&]:hover:bg-primary/90",
                secondary:
                    "border-transparent bg-secondary text-secondary-foreground [a&]:hover:bg-secondary/90",
                destructive:
                    "border-transparent bg-destructive text-white [a&]:hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/60",
                outline:
                    "text-foreground [a&]:hover:bg-accent [a&]:hover:text-accent-foreground",

                // ✅ 新增柔和风格
                success:
                    "border-transparent bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300 [a&]:hover:bg-emerald-200",
                warning:
                    "border-transparent bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300 [a&]:hover:bg-amber-200",
                danger:
                    "border-transparent bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-300 [a&]:hover:bg-rose-200",
                info:
                    "border-transparent bg-sky-100 text-sky-700 dark:bg-sky-900/30 dark:text-sky-300 [a&]:hover:bg-sky-200",
                purple:
                    "border-transparent bg-violet-100 text-violet-700 dark:bg-violet-900/30 dark:text-violet-300 [a&]:hover:bg-violet-200",
            },
        },
        defaultVariants: {
            variant: "default",
        },
    },
)

export type BadgeVariants = VariantProps<typeof badgeVariants>
