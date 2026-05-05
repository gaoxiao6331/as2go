// easy_syntax.ts
// AssemblyScript 易转换语法测试文件
// 涵盖：变量声明、基础类型、for/while 循环、switch 分支、函数、数组、数学运算、位运算

// ─────────────────────────────────────────────
// 1. 基础类型与变量声明
// ─────────────────────────────────────────────
const MAX_SIZE: i32 = 100;
const PI: f64 = 3.14159265358979;

let counter: i32 = 0;
let flag: bool = true;
let ratio: f64 = 0.0;
let message: string = "hello";

// ─────────────────────────────────────────────
// 2. 算术 & 位运算
// ─────────────────────────────────────────────
export function bitwiseOps(a: i32, b: i32): i32 {
  let andResult: i32 = a & b;
  let orResult: i32 = a | b;
  let xorResult: i32 = a ^ b;
  let shiftLeft: i32 = a << 2;
  let shiftRight: i32 = b >> 1;
  let notA: i32 = ~a;
  return andResult + orResult + xorResult + shiftLeft + shiftRight + notA;
}

export function mathOps(x: f64, y: f64): f64 {
  let sum: f64 = x + y;
  let diff: f64 = x - y;
  let product: f64 = x * y;
  let quotient: f64 = x / y;
  let remainder: f64 = x % y;
  return sum + diff + product + quotient + remainder;
}

// ─────────────────────────────────────────────
// 3. if / else if / else
// ─────────────────────────────────────────────
export function classify(n: i32): string {
  if (n < 0) {
    return "negative";
  } else if (n === 0) {
    return "zero";
  } else if (n < 10) {
    return "small";
  } else if (n < 100) {
    return "medium";
  } else {
    return "large";
  }
}

// ─────────────────────────────────────────────
// 4. switch / case / default
// ─────────────────────────────────────────────
export function dayName(day: i32): string {
  switch (day) {
    case 0:
      return "Sunday";
    case 1:
      return "Monday";
    case 2:
      return "Tuesday";
    case 3:
      return "Wednesday";
    case 4:
      return "Thursday";
    case 5:
      return "Friday";
    case 6:
      return "Saturday";
    default:
      return "Unknown";
  }
}

export function switchWithFallthrough(code: i32): i32 {
  let result: i32 = 0;
  switch (code) {
    case 1:
    case 2:
      result = 10;
      break;
    case 3:
      result = 30;
      break;
    case 4:
    case 5:
    case 6:
      result = 60;
      break;
    default:
      result = -1;
  }
  return result;
}

// ─────────────────────────────────────────────
// 5. for 循环（标准 / 倒序 / 步进）
// ─────────────────────────────────────────────
export function sumUpTo(n: i32): i32 {
  let total: i32 = 0;
  for (let i: i32 = 1; i <= n; i++) {
    total += i;
  }
  return total;
}

export function countDown(from: i32): i32 {
  let steps: i32 = 0;
  for (let i: i32 = from; i > 0; i--) {
    steps++;
  }
  return steps;
}

export function sumEvens(limit: i32): i32 {
  let s: i32 = 0;
  for (let i: i32 = 0; i <= limit; i += 2) {
    s += i;
  }
  return s;
}

// ─────────────────────────────────────────────
// 6. while / do-while
// ─────────────────────────────────────────────
export function collatz(n: i32): i32 {
  let steps: i32 = 0;
  while (n !== 1) {
    if (n % 2 === 0) {
      n = n / 2;
    } else {
      n = 3 * n + 1;
    }
    steps++;
  }
  return steps;
}

export function doWhileExample(start: i32): i32 {
  let val: i32 = start;
  do {
    val = val * 2;
  } while (val < 1000);
  return val;
}

// ─────────────────────────────────────────────
// 7. break / continue
// ─────────────────────────────────────────────
export function firstDivisibleBy7(limit: i32): i32 {
  for (let i: i32 = 1; i <= limit; i++) {
    if (i % 7 === 0) {
      return i;
    }
  }
  return -1;
}

export function sumSkipMultiplesOf3(n: i32): i32 {
  let total: i32 = 0;
  for (let i: i32 = 1; i <= n; i++) {
    if (i % 3 === 0) {
      continue;
    }
    total += i;
  }
  return total;
}

// ─────────────────────────────────────────────
// 8. 嵌套循环
// ─────────────────────────────────────────────
export function multiplicationTable(size: i32): i32 {
  let checksum: i32 = 0;
  for (let i: i32 = 1; i <= size; i++) {
    for (let j: i32 = 1; j <= size; j++) {
      checksum += i * j;
    }
  }
  return checksum;
}

export function bubbleSort(arr: Int32Array): void {
  let n: i32 = arr.length;
  for (let i: i32 = 0; i < n - 1; i++) {
    for (let j: i32 = 0; j < n - i - 1; j++) {
      if (arr[j] > arr[j + 1]) {
        let tmp: i32 = arr[j];
        arr[j] = arr[j + 1];
        arr[j + 1] = tmp;
      }
    }
  }
}

// ─────────────────────────────────────────────
// 9. 函数：递归、默认参数、多返回值模拟
// ─────────────────────────────────────────────
export function factorial(n: i32): i64 {
  if (n <= 1) return 1;
  return n * factorial(n - 1);
}

export function fibonacci(n: i32): i32 {
  if (n <= 1) return n;
  return fibonacci(n - 1) + fibonacci(n - 2);
}

export function power(base: f64, exp: i32 = 2): f64 {
  let result: f64 = 1.0;
  for (let i: i32 = 0; i < exp; i++) {
    result *= base;
  }
  return result;
}

export function gcd(a: i32, b: i32): i32 {
  while (b !== 0) {
    let t: i32 = b;
    b = a % b;
    a = t;
  }
  return a;
}

// ─────────────────────────────────────────────
// 10. 静态数组操作
// ─────────────────────────────────────────────
export function arraySum(arr: Int32Array): i32 {
  let total: i32 = 0;
  for (let i: i32 = 0; i < arr.length; i++) {
    total += arr[i];
  }
  return total;
}

export function arrayMax(arr: Int32Array): i32 {
  let max: i32 = arr[0];
  for (let i: i32 = 1; i < arr.length; i++) {
    if (arr[i] > max) {
      max = arr[i];
    }
  }
  return max;
}

export function reverseArray(arr: Int32Array): void {
  let left: i32 = 0;
  let right: i32 = arr.length - 1;
  while (left < right) {
    let tmp: i32 = arr[left];
    arr[left] = arr[right];
    arr[right] = tmp;
    left++;
    right--;
  }
}

// ─────────────────────────────────────────────
// 11. 内存操作（load / store）
// ─────────────────────────────────────────────
export function rawMemoryExample(ptr: usize): i32 {
  store<i32>(ptr, 42);
  let val: i32 = load<i32>(ptr);
  store<i32>(ptr + 4, val * 2);
  return load<i32>(ptr + 4);
}

// ─────────────────────────────────────────────
// 12. 三元运算符 & 逻辑运算符
// ─────────────────────────────────────────────
export function ternaryExample(x: i32): i32 {
  return x > 0 ? x : -x; // abs
}

export function logicalOps(a: bool, b: bool, c: bool): bool {
  return (a && b) || (!a && c) || (a && !b && !c);
}

// ─────────────────────────────────────────────
// 13. 类型转换
// ─────────────────────────────────────────────
export function typeCasts(n: i32): f64 {
  let asF64: f64 = <f64>n;
  let asI64: i64 = <i64>n;
  let backToI32: i32 = <i32>asF64;
  let u: u32 = <u32>backToI32;
  return asF64 + <f64>asI64 + <f64>u;
}

// ─────────────────────────────────────────────
// 14. 内置数学函数
// ─────────────────────────────────────────────
export function mathBuiltins(x: f64): f64 {
  let s: f64 = Math.sqrt(x);
  let a: f64 = Math.abs(x - 10.0);
  let fl: f64 = Math.floor(x);
  let ce: f64 = Math.ceil(x);
  let mn: f64 = Math.min(x, 5.0);
  let mx: f64 = Math.max(x, 5.0);
  let pw: f64 = Math.pow(x, 2.0);
  return s + a + fl + ce + mn + mx + pw;
}

// ─────────────────────────────────────────────
// 15. 简单 class（无继承）
// ─────────────────────────────────────────────
class Point {
  x: f64;
  y: f64;

  constructor(x: f64, y: f64) {
    this.x = x;
    this.y = y;
  }

  distanceTo(other: Point): f64 {
    let dx: f64 = this.x - other.x;
    let dy: f64 = this.y - other.y;
    return Math.sqrt(dx * dx + dy * dy);
  }

  scale(factor: f64): void {
    this.x *= factor;
    this.y *= factor;
  }
}

export function pointDistance(x1: f64, y1: f64, x2: f64, y2: f64): f64 {
  let p1 = new Point(x1, y1);
  let p2 = new Point(x2, y2);
  return p1.distanceTo(p2);
}

// ─────────────────────────────────────────────
// 16. 导出入口
// ─────────────────────────────────────────────
export function run(): i32 {
  let a = sumUpTo(10);
  let b = collatz(27);
  let c = fibonacci(10);
  let d = multiplicationTable(5);
  let pt = pointDistance(0.0, 0.0, 3.0, 4.0); // should be 5.0
  return a + b + c + d + <i32>pt;
}