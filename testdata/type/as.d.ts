// ======================================================
// AssemblyScript Type Mapping (TypeScript Lightweight Layer)
// 目标：仅补充 TS 无法表达的 WASM / AS 语义
// 不重复 lib.es / lib.dom / 标准 TS 类型
// ======================================================

/**
 * =========================
 * 1. Integer Types (WASM)
 * =========================
 * TS 没有 bit-width int，因此全部语义化为 number / bigint
 */

// signed integers
 type i8 = number;
 type i16 = number;
 type i32 = number;

// unsigned integers
 type u8 = number;
 type u16 = number;
 type u32 = number;

// 64-bit integers (TS only supports bigint)
 type i64 = number;
 type u64 = number;

/**
 * =========================
 * 2. Floating Point Types
 * =========================
 */

 type f32 = number;
 type f64 = number;

/**
 * =========================
 * 3. Memory / Pointer Types
 * =========================
 * AssemblyScript uses linear memory addressing
 * TS has no pointer concept
 */

 type ptr = number;     // memory address / offset
 type usize = number;   // pointer-sized unsigned int
 type isize = number;   // pointer-sized signed int

/**
 * =========================
 * 4. Boolean (semantic only)
 * =========================
 */

 type bool = boolean;

/**
 * =========================
 * 5. AssemblyScript Semantics
 * =========================
 */

/**
 * Nullable type (AS-style)
 */
 type nullable<T> = T | null;

/**
 * Reference marker (no runtime effect)
 */
 type ref<T> = T;

/**
 * Unreachable type (used in compiler semantics)
 */
 type unreachable = never;

/**
 * =========================
 * 6. Optional Utility Types (AS-like patterns)
 * =========================
 */

/**
 * Static array (fixed-length semantic)
 * TS cannot enforce fixed length at type level reliably
 */
 type StaticArray<T> = readonly T[];

/**
 * AS-style dynamic array alias (no override of global Array!)
 */
 type ASArray<T> = T[];

/**
 * Function type helper
 */
 type Func<Args extends any[] = any[], R = any> = (...args: Args) => R;

/**
 * =========================
 * 7. Optional Memory Views (re- safe)
 * =========================
 * We DO NOT redeclare built-in TS types, only alias if needed
 */

 type MemoryBuffer = ArrayBuffer;

/**
 * =========================
 * 8. Notes (important design rule)
 * =========================
 *
 * ❌ DO NOT redefine:
 * - string
 * - Array
 * - Uint8Array / TypedArray
 * - DataView
 * - Function
 *
 * Reason: already defined in lib.es202x + lib.dom
 *
 * ✔ Only define:
 * - bit-width integers
 * - pointer / memory semantics
 * - AS compiler semantics
 */