package codegen

const runtimeFuncs = `
// --- as2go runtime ---
// abort：对应 AS 的 abort() 内置
func __as_abort(msg string) {
    panic("abort: " + msg)
}
`
