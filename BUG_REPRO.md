# Bug 复现说明：核验不存在的参保人 panic

## Bug 是什么
`InsuredPersonRepository` 在记录不存在时返回 `nil, nil`（typed-nil 旁路），service 判 `err==nil` 误以为查询成功，handler 解引用 nil 触发 panic，接口回 500 而非 404。

## 如何触发
POST /api/v1/clients/verify 核验不存在的身份证号+医保卡号。

## 真实错误信息
```
--- FAIL: TestInsuredVerifyMissingNoPanicP401
panic: runtime error: invalid memory address or nil pointer dereference [recovered, repanicked]
panic({0x1034c5440?, 0x1038fe720?})
--- FAIL: TestInsuredHandlerVerifyMissing404P403
time=... level=ERROR msg="panic recovered" panic="runtime error: invalid memory address or nil pointer dereference" path=/api/v1/clients/verify
```
