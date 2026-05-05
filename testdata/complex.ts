// hard_syntax.ts
// AssemblyScript 难转换语法测试文件
// 涵盖：try/catch/finally、异常类型、泛型、闭包/函数类型、
//        动态 Array<T>、Map/Set、接口、继承多态、操作符重载、
//        装饰器、nullability、类型谓词

// ─────────────────────────────────────────────
// 1. try / catch / finally（AssemblyScript 仅支持抛出整数/引用）
// ─────────────────────────────────────────────

// AS 中 throw 只能抛出继承自 Error 的对象（或 i32）
class DivisionError extends Error {
  dividend: i32;
  divisor: i32;

  constructor(dividend: i32, divisor: i32) {
    super("division by zero");
    this.dividend = dividend;
    this.divisor = divisor;
  }
}

class RangeError extends Error {
  value: i32;
  min: i32;
  max: i32;

  constructor(value: i32, min: i32, max: i32) {
    super("value out of range");
    this.value = value;
    this.min = min;
    this.max = max;
  }
}

export function safeDivide(a: i32, b: i32): i32 {
  try {
    if (b === 0) {
      throw new DivisionError(a, b);
    }
    return a / b;
  } catch (e: DivisionError) {
    // 捕获特定类型——Go 需要用 errors.As 或自定义 panic/recover 模拟
    return -1;
  } finally {
    // finally 块在 Go 中需转换为 defer
    // (side-effect only, no return here to keep semantics simple)
  }
}

export function rangeCheck(value: i32, min: i32, max: i32): i32 {
  try {
    if (value < min || value > max) {
      throw new RangeError(value, min, max);
    }
    return value;
  } catch (e: RangeError) {
    return min; // clamp to min on error
  } finally {
    // cleanup: 对应 Go 的 defer cleanup()
  }
}

// 嵌套 try/catch
export function nestedTryCatch(x: i32): i32 {
  try {
    try {
      if (x < 0) throw new DivisionError(x, 0);
      return safeDivide(100, x);
    } catch (inner: DivisionError) {
      // 内层捕获后再次抛出——Go 需要 panic(recover()) 组合
      throw new RangeError(x, 0, 1000);
    }
  } catch (outer: RangeError) {
    return -999;
  }
}

// ─────────────────────────────────────────────
// 2. 泛型函数 & 泛型类
// ─────────────────────────────────────────────

// 泛型函数——Go 1.18+ generics 可以对应，但约束语义不同
function identity<T>(value: T): T {
  return value;
}

function swap<T>(a: T, b: T): void {
  let tmp: T = a;
  a = b;
  b = tmp;
}

function clamp<T extends number>(value: T, lo: T, hi: T): T {
  if (value < lo) return lo;
  if (value > hi) return hi;
  return value;
}

// 泛型类
class Stack<T> {
  private data: Array<T> = new Array<T>();

  push(item: T): void {
    this.data.push(item);
  }

  pop(): T {
    if (this.data.length === 0) {
      throw new Error("stack underflow");
    }
    return this.data.pop();
  }

  peek(): T {
    return this.data[this.data.length - 1];
  }

  get size(): i32 {
    return this.data.length;
  }

  isEmpty(): bool {
    return this.data.length === 0;
  }
}

export function genericStackDemo(): i32 {
  let s = new Stack<i32>();
  s.push(1);
  s.push(2);
  s.push(3);
  let top = s.pop();   // 3
  return top + s.peek(); // 3 + 2 = 5
}

// ─────────────────────────────────────────────
// 3. 动态 Array<T>（不同于静态 TypedArray）
// ─────────────────────────────────────────────

export function dynamicArrayOps(): i32 {
  let arr = new Array<i32>();

  // push / pop
  for (let i: i32 = 0; i < 5; i++) {
    arr.push(i * i);
  }
  let popped: i32 = arr.pop(); // 16

  // splice——Go 没有直接对应，需手动重建切片
  arr.splice(1, 2); // remove 2 elements at index 1

  // unshift / shift
  arr.unshift(100);
  let shifted: i32 = arr.shift(); // 100

  // indexOf / includes
  let idx: i32 = arr.indexOf(4);
  let has: bool = arr.includes(0);

  // slice (returns new array)
  let sliced = arr.slice(0, 2);

  return popped + idx + (has ? 1 : 0) + sliced.length;
}

// ─────────────────────────────────────────────
// 4. Map<K, V> & Set<T>
// ─────────────────────────────────────────────

export function mapOperations(): i32 {
  let m = new Map<string, i32>();

  m.set("alpha", 1);
  m.set("beta", 2);
  m.set("gamma", 3);

  let hasAlpha: bool = m.has("alpha");
  let val: i32 = m.get("beta");   // AS: returns value or 0 (no optional)
  m.delete("gamma");

  // keys() / values() iteration — 在 Go 中 map iteration 顺序不定
  let keys = m.keys();
  let total: i32 = 0;
  for (let i: i32 = 0; i < keys.length; i++) {
    total += m.get(keys[i]);
  }
  return total + val + (hasAlpha ? 10 : 0);
}

export function setOperations(): i32 {
  let s = new Set<i32>();

  s.add(1);
  s.add(2);
  s.add(3);
  s.add(2); // duplicate, ignored

  let size: i32 = s.size;         // 3
  let has2: bool = s.has(2);
  s.delete(2);
  let afterDel: i32 = s.size;     // 2

  let vals = s.values();
  let sum: i32 = 0;
  for (let i: i32 = 0; i < vals.length; i++) {
    sum += vals[i];
  }
  return size + afterDel + sum + (has2 ? 1 : 0);
}

// ─────────────────────────────────────────────
// 5. 接口（Interface）
// ─────────────────────────────────────────────

interface Shape {
  area(): f64;
  perimeter(): f64;
  describe(): string;
}

interface Scalable {
  scale(factor: f64): void;
}

class Circle implements Shape, Scalable {
  radius: f64;

  constructor(r: f64) {
    this.radius = r;
  }

  area(): f64 {
    return Math.PI * this.radius * this.radius;
  }

  perimeter(): f64 {
    return 2.0 * Math.PI * this.radius;
  }

  describe(): string {
    return "circle";
  }

  scale(factor: f64): void {
    this.radius *= factor;
  }
}

class Rectangle implements Shape, Scalable {
  width: f64;
  height: f64;

  constructor(w: f64, h: f64) {
    this.width = w;
    this.height = h;
  }

  area(): f64 {
    return this.width * this.height;
  }

  perimeter(): f64 {
    return 2.0 * (this.width + this.height);
  }

  describe(): string {
    return "rectangle";
  }

  scale(factor: f64): void {
    this.width *= factor;
    this.height *= factor;
  }
}

// 多态调用——Go 需要用 interface{} 或 type switch 来模拟
export function totalArea(shapes: Array<Shape>): f64 {
  let total: f64 = 0.0;
  for (let i: i32 = 0; i < shapes.length; i++) {
    total += shapes[i].area();
  }
  return total;
}

// ─────────────────────────────────────────────
// 6. 继承 & 方法重写（super 调用）
// ─────────────────────────────────────────────

class Animal {
  name: string;

  constructor(name: string) {
    this.name = name;
  }

  speak(): string {
    return this.name + " makes a sound";
  }

  toString(): string {
    return "Animal(" + this.name + ")";
  }
}

class Dog extends Animal {
  breed: string;

  constructor(name: string, breed: string) {
    super(name); // super 调用
    this.breed = breed;
  }

  speak(): string {
    return super.speak() + ": Woof!"; // super 方法调用
  }

  fetch(): string {
    return this.name + " fetches the ball";
  }
}

class GuideDog extends Dog {
  owner: string;

  constructor(name: string, breed: string, owner: string) {
    super(name, breed);
    this.owner = owner;
  }

  speak(): string {
    return super.speak() + " (guide dog)";
  }
}

export function polymorphismDemo(): string {
  let animals: Array<Animal> = new Array<Animal>();
  animals.push(new Animal("Cat"));
  animals.push(new Dog("Rex", "Labrador"));
  animals.push(new GuideDog("Buddy", "Golden", "Alice"));

  let result: string = "";
  for (let i: i32 = 0; i < animals.length; i++) {
    result = result + animals[i].speak() + ";";
  }
  return result;
}

// ─────────────────────────────────────────────
// 7. 操作符重载
// ─────────────────────────────────────────────

class Vec2 {
  x: f64;
  y: f64;

  constructor(x: f64, y: f64) {
    this.x = x;
    this.y = y;
  }

  // 操作符重载——Go 完全不支持，需转换为方法调用
  @operator("+")
  add(other: Vec2): Vec2 {
    return new Vec2(this.x + other.x, this.y + other.y);
  }

  @operator("-")
  sub(other: Vec2): Vec2 {
    return new Vec2(this.x - other.x, this.y - other.y);
  }

  @operator("*")
  scale(scalar: f64): Vec2 {
    return new Vec2(this.x * scalar, this.y * scalar);
  }

  @operator("==")
  equals(other: Vec2): bool {
    return this.x === other.x && this.y === other.y;
  }

  magnitude(): f64 {
    return Math.sqrt(this.x * this.x + this.y * this.y);
  }
}

export function vec2Demo(): f64 {
  let v1 = new Vec2(1.0, 2.0);
  let v2 = new Vec2(3.0, 4.0);
  let v3 = v1 + v2;           // 操作符重载语法
  let v4 = v3 - v1;
  let v5 = v4 * 2.0;
  return v5.magnitude();
}

// ─────────────────────────────────────────────
// 8. 闭包 / 一等函数 / 函数类型
// ─────────────────────────────────────────────

// AS 中函数类型用 (a: T) => R 表示，但捕获外部变量支持有限
type Predicate<T> = (item: T) => bool;
type Transformer<T, R> = (item: T) => R;
type Reducer<T, R> = (acc: R, item: T) => R;

function filter<T>(arr: Array<T>, pred: Predicate<T>): Array<T> {
  let result = new Array<T>();
  for (let i: i32 = 0; i < arr.length; i++) {
    if (pred(arr[i])) result.push(arr[i]);
  }
  return result;
}

function map<T, R>(arr: Array<T>, fn: Transformer<T, R>): Array<R> {
  let result = new Array<R>();
  for (let i: i32 = 0; i < arr.length; i++) {
    result.push(fn(arr[i]));
  }
  return result;
}

function reduce<T, R>(arr: Array<T>, fn: Reducer<T, R>, initial: R): R {
  let acc: R = initial;
  for (let i: i32 = 0; i < arr.length; i++) {
    acc = fn(acc, arr[i]);
  }
  return acc;
}

export function higherOrderDemo(): i32 {
  let nums = new Array<i32>();
  for (let i: i32 = 1; i <= 10; i++) nums.push(i);

  // 内联 lambda（捕获外部变量在 AS 中受限）
  let evens = filter<i32>(nums, (n: i32): bool => n % 2 === 0);
  let doubled = map<i32, i32>(evens, (n: i32): i32 => n * 2);
  let sum = reduce<i32, i32>(doubled, (acc: i32, n: i32): i32 => acc + n, 0);
  return sum; // (2+4+6+8+10)*2 = 60
}

// ─────────────────────────────────────────────
// 9. Nullable 类型 / 空检查（AS 的 | null）
// ─────────────────────────────────────────────

class Node {
  value: i32;
  next: Node | null;

  constructor(value: i32) {
    this.value = value;
    this.next = null;
  }
}

class LinkedList {
  head: Node | null = null;
  private _size: i32 = 0;

  append(value: i32): void {
    let node = new Node(value);
    if (this.head === null) {
      this.head = node;
    } else {
      let current: Node = this.head!; // non-null assertion
      while (current.next !== null) {
        current = current.next!;
      }
      current.next = node;
    }
    this._size++;
  }

  find(value: i32): Node | null {
    let current: Node | null = this.head;
    while (current !== null) {
      if (current.value === value) return current;
      current = current.next;
    }
    return null;
  }

  get size(): i32 {
    return this._size;
  }
}

export function linkedListDemo(): i32 {
  let list = new LinkedList();
  list.append(1);
  list.append(2);
  list.append(3);

  let found: Node | null = list.find(2);
  // null 检查——Go 对应 if found != nil
  if (found !== null) {
    return found.value + list.size; // 2 + 3 = 5
  }
  return -1;
}

// ─────────────────────────────────────────────
// 10. 装饰器（@inline、@lazy、自定义）
// ─────────────────────────────────────────────

// @inline 提示内联——Go 中无直接对应（编译器自行决定）
@inline
function fastAbs(x: i32): i32 {
  return x < 0 ? -x : x;
}

// @lazy 延迟初始化——Go 中需用 sync.Once 模拟
@lazy
let expensiveValue: i32 = computeExpensive();

function computeExpensive(): i32 {
  let result: i32 = 0;
  for (let i: i32 = 0; i < 1000; i++) {
    result += i;
  }
  return result;
}

// 自定义装饰器（元编程）
function deprecated(target: Function): Function {
  // In AS, decorators are limited; this is illustrative
  return target;
}

@deprecated
function oldApi(x: i32): i32 {
  return x * 2;
}

export function decoratorDemo(): i32 {
  return fastAbs(-42) + expensiveValue;
}

// ─────────────────────────────────────────────
// 11. 异步模拟（AS 本身不支持 async/await，
//     但转换工具可能遇到用户手写的类似模式）
// ─────────────────────────────────────────────

// AS 无原生 Promise，这里用回调函数模式模拟
type Callback<T> = (err: Error | null, result: T) => void;

function asyncAdd(a: i32, b: i32, cb: Callback<i32>): void {
  // 同步执行但接口语义是异步的
  try {
    if (a < 0 || b < 0) {
      throw new Error("negative input");
    }
    cb(null, a + b);
  } catch (e: Error) {
    cb(e, 0);
  }
}

export function callbackDemo(): i32 {
  let result: i32 = 0;
  asyncAdd(3, 4, (err: Error | null, val: i32): void => {
    if (err === null) {
      result = val;
    }
  });
  return result; // 7
}

// ─────────────────────────────────────────────
// 12. 复合：泛型 + 异常 + 接口组合场景
// ─────────────────────────────────────────────

interface Repository<T> {
  findById(id: i32): T | null;
  save(item: T): void;
  delete(id: i32): bool;
}

class InMemoryRepo<T> implements Repository<T> {
  private store: Map<i32, T> = new Map<i32, T>();
  private nextId: i32 = 1;

  findById(id: i32): T | null {
    if (this.store.has(id)) {
      return this.store.get(id);
    }
    return null;
  }

  save(item: T): void {
    try {
      this.store.set(this.nextId++, item);
    } catch (e: Error) {
      throw new Error("save failed: " + e.message);
    }
  }

  delete(id: i32): bool {
    if (!this.store.has(id)) return false;
    this.store.delete(id);
    return true;
  }

  get count(): i32 {
    return this.store.size;
  }
}

export function repositoryDemo(): i32 {
  let repo = new InMemoryRepo<i32>();
  repo.save(100);
  repo.save(200);
  repo.save(300);

  let found: i32 | null = repo.findById(2);
  let deleted: bool = repo.delete(1);

  let val: i32 = found !== null ? found! : 0;
  return val + (deleted ? 1 : 0) + repo.count; // 200 + 1 + 2 = 203
}