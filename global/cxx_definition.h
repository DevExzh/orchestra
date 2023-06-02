#ifndef CXX_DEFINITION
#define CXX_DEFINITION

#ifdef __cplusplus
extern "C"{
#endif
#include <stdint.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>

#define REPEAT(number, ...) REPEAT_##number(__VA_ARGS__)
#define ARGS(number) ARGS_EXPAND_##number
#define REPEAT_0(macro)
#define REPEAT_1(macro) macro
#define REPEAT_2(m, ...) REPEAT_1(m), m
#define REPEAT_3(m, ...) REPEAT_2(m), m
#define REPEAT_4(m, ...) REPEAT_3(m), m
#define REPEAT_5(m, ...) REPEAT_4(m), m
#define REPEAT_6(m, ...) REPEAT_5(m), m
#define REPEAT_7(m, ...) REPEAT_6(m), m
#define REPEAT_8(m, ...) REPEAT_7(m), m
#define REPEAT_9(m, ...) REPEAT_8(m), m
#define REPEAT_10(m, ...) REPEAT_9(m), m
#define REPEAT_11(m, ...) REPEAT_10(m), m
#define REPEAT_12(m, ...) REPEAT_11(m), m
#define REPEAT_13(m, ...) REPEAT_12(m), m
#define REPEAT_14(m, ...) REPEAT_13(m), m
#define REPEAT_15(m, ...) REPEAT_14(m), m
#define REPEAT_16(m, ...) REPEAT_15(m), m
#define REPEAT_17(m, ...) REPEAT_16(m), m
#define REPEAT_18(m, ...) REPEAT_17(m), m
#define REPEAT_19(m, ...) REPEAT_18(m), m
#define REPEAT_20(m, ...) REPEAT_19(m), m
#define REPEAT_21(m, ...) REPEAT_20(m), m
#define REPEAT_22(m, ...) REPEAT_21(m), m
#define REPEAT_23(m, ...) REPEAT_22(m), m
#define REPEAT_24(m, ...) REPEAT_23(m), m
#define REPEAT_25(m, ...) REPEAT_24(m), m
#define REPEAT_26(m, ...) REPEAT_25(m), m
#define REPEAT_27(m, ...) REPEAT_26(m), m
#define REPEAT_28(m, ...) REPEAT_27(m), m
#define REPEAT_29(m, ...) REPEAT_28(m), m
#define REPEAT_30(m, ...) REPEAT_29(m), m
#define REPEAT_31(m, ...) REPEAT_30(m), m
#define REPEAT_32(m, ...) REPEAT_31(m), m
#define REPEAT_33(m, ...) REPEAT_32(m), m
#define REPEAT_34(m, ...) REPEAT_33(m), m
#define REPEAT_35(m, ...) REPEAT_34(m), m
#define REPEAT_36(m, ...) REPEAT_35(m), m
#define REPEAT_37(m, ...) REPEAT_36(m), m
#define REPEAT_38(m, ...) REPEAT_37(m), m
#define REPEAT_39(m, ...) REPEAT_38(m), m
#define REPEAT_40(m, ...) REPEAT_39(m), m
#define REPEAT_41(m, ...) REPEAT_40(m), m
#define REPEAT_42(m, ...) REPEAT_41(m), m
#define REPEAT_43(m, ...) REPEAT_42(m), m
#define REPEAT_44(m, ...) REPEAT_43(m), m
#define REPEAT_45(m, ...) REPEAT_44(m), m
#define REPEAT_46(m, ...) REPEAT_45(m), m
#define REPEAT_47(m, ...) REPEAT_46(m), m
#define REPEAT_48(m, ...) REPEAT_47(m), m
#define REPEAT_49(m, ...) REPEAT_48(m), m
#define REPEAT_50(m, ...) REPEAT_49(m), m
#define REPEAT_51(m, ...) REPEAT_50(m), m
#define REPEAT_52(m, ...) REPEAT_51(m), m
#define REPEAT_53(m, ...) REPEAT_52(m), m
#define REPEAT_54(m, ...) REPEAT_53(m), m
#define REPEAT_55(m, ...) REPEAT_54(m), m
#define REPEAT_56(m, ...) REPEAT_55(m), m
#define REPEAT_57(m, ...) REPEAT_56(m), m
#define REPEAT_58(m, ...) REPEAT_57(m), m
#define REPEAT_59(m, ...) REPEAT_58(m), m
#define REPEAT_60(m, ...) REPEAT_59(m), m
#define REPEAT_61(m, ...) REPEAT_60(m), m
#define REPEAT_62(m, ...) REPEAT_61(m), m
#define REPEAT_63(m, ...) REPEAT_62(m), m
#define REPEAT_64(m, ...) REPEAT_63(m), m
#define SYSCALL_N(num) uint64_t SystemCall_##num(void*, REPEAT(num, void*))

uint64_t SystemCall(void* address);
SYSCALL_N(1);
SYSCALL_N(2);
SYSCALL_N(3);
SYSCALL_N(4);
uintptr_t OpenLibrary(const char* path);
const char* GetLibraryError();
bool CloseLibrary(uintptr_t handle);
void* FindSymbol(uintptr_t handle, const char* symbol);

#ifdef __cplusplus
}
#endif

#endif