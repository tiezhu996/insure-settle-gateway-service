# Bug 复现说明：过期 token 触发 nil 解引用 panic

## Bug 是什么
`ParseToken` 对过期 token 吞掉错误并返回 `(nil, nil)`（typed-nil），`JwtAuthRequired` 判 `err==nil` 后解引用 nil claims 触发 panic；结算详情接口对过期 token 回 500 而非 401。

## 如何触发
用过期 token 请求 /api/v1/settlements/{no} 等 JWT 保护接口。

## 真实错误信息
```
--- FAIL: TestParseTokenInvalidTypedNilP801
--- FAIL: TestJwtAuthInvalidTokenNoPanicP802
[Recovery] ... panic recovered:
runtime error: invalid memory address or nil pointer dereference
--- FAIL: TestSettleDetailBadToken401P803
time=... msg="panic recovered" panic="runtime error: invalid memory address or nil pointer dereference" path=/api/v1/settlements/S1
--- FAIL: TestSettleSubmitMissingClientIDP805
msg="panic recovered" panic="interface conversion: interface {} is nil, not uint" path=/api/v1/settlements/submit
```
