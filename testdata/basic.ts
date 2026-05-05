// 基础函数
export function add(a: i32, b: i32): i32 {
    return a + b;
}

// 枚举
export const enum Color {
    Red,
    Green,
    Blue,
}

// 类
export class Vec2 {
    x: f64;
    y: f64;

    length(): f64 {
        return Math.sqrt(this.x * this.x + this.y * this.y);
    }
}

// 数组
export function sum(arr: Array<i32>): i32 {
    let total: i32 = 0;
    for (let i: i32 = 0; i < arr.length; i++) {
        total += arr[i];
    }
    return total;
}
