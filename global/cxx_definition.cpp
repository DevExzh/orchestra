#include "cxx_definition.h"
#ifdef __cplusplus
extern "C"{
#endif
#include <stdint.h>
#include <stdarg.h>
#define ARGS_EXPAND_0
#define ARGS_EXPAND_1 args[0]
#define ARGS_EXPAND_2 ARGS_EXPAND_1, args[1]
#define ARGS_EXPAND_3 ARGS_EXPAND_2, args[2]
#define ARGS_EXPAND_4 ARGS_EXPAND_3, args[3]
#define ARGS_EXPAND_5 ARGS_EXPAND_4, args[4]
#define ARGS_EXPAND_6 ARGS_EXPAND_5, args[5]
#define ARGS_EXPAND_7 ARGS_EXPAND_6, args[6]
#define ARGS_EXPAND_8 ARGS_EXPAND_7, args[7]
#define ARGS_EXPAND_9 ARGS_EXPAND_8, args[8]
#define ARGS_EXPAND_10 ARGS_EXPAND_9, args[9]
#define ARGS_EXPAND_11 ARGS_EXPAND_10 , args[10]
#define ARGS_EXPAND_12 ARGS_EXPAND_11 , args[11]
#define ARGS_EXPAND_13 ARGS_EXPAND_12 , args[12]
#define ARGS_EXPAND_14 ARGS_EXPAND_13 , args[13]
#define ARGS_EXPAND_15 ARGS_EXPAND_14 , args[14]
#define ARGS_EXPAND_16 ARGS_EXPAND_15 , args[15]
#define ARGS_EXPAND_17 ARGS_EXPAND_16 , args[16]
#define ARGS_EXPAND_18 ARGS_EXPAND_17 , args[17]
#define ARGS_EXPAND_19 ARGS_EXPAND_18 , args[18]
#define ARGS_EXPAND_20 ARGS_EXPAND_19 , args[19]
#define ARGS_EXPAND_21 ARGS_EXPAND_20 , args[20]
#define ARGS_EXPAND_22 ARGS_EXPAND_21 , args[21]
#define ARGS_EXPAND_23 ARGS_EXPAND_22 , args[22]
#define ARGS_EXPAND_24 ARGS_EXPAND_23 , args[23]
#define ARGS_EXPAND_25 ARGS_EXPAND_24 , args[24]
#define ARGS_EXPAND_26 ARGS_EXPAND_25 , args[25]
#define ARGS_EXPAND_27 ARGS_EXPAND_26 , args[26]
#define ARGS_EXPAND_28 ARGS_EXPAND_27 , args[27]
#define ARGS_EXPAND_29 ARGS_EXPAND_28 , args[28]
#define ARGS_EXPAND_30 ARGS_EXPAND_29 , args[29]
#define ARGS_EXPAND_31 ARGS_EXPAND_30 , args[30]
#define ARGS_EXPAND_32 ARGS_EXPAND_31 , args[31]
#define ARGS_EXPAND_33 ARGS_EXPAND_32 , args[32]
#define ARGS_EXPAND_34 ARGS_EXPAND_33 , args[33]
#define ARGS_EXPAND_35 ARGS_EXPAND_34 , args[34]
#define ARGS_EXPAND_36 ARGS_EXPAND_35 , args[35]
#define ARGS_EXPAND_37 ARGS_EXPAND_36 , args[36]
#define ARGS_EXPAND_38 ARGS_EXPAND_37 , args[37]
#define ARGS_EXPAND_39 ARGS_EXPAND_38 , args[38]
#define ARGS_EXPAND_40 ARGS_EXPAND_39 , args[39]
#define ARGS_EXPAND_41 ARGS_EXPAND_40 , args[40]
#define ARGS_EXPAND_42 ARGS_EXPAND_41 , args[41]
#define ARGS_EXPAND_43 ARGS_EXPAND_42 , args[42]
#define ARGS_EXPAND_44 ARGS_EXPAND_43 , args[43]
#define ARGS_EXPAND_45 ARGS_EXPAND_44 , args[44]
#define ARGS_EXPAND_46 ARGS_EXPAND_45 , args[45]
#define ARGS_EXPAND_47 ARGS_EXPAND_46 , args[46]
#define ARGS_EXPAND_48 ARGS_EXPAND_47 , args[47]
#define ARGS_EXPAND_49 ARGS_EXPAND_48 , args[48]
#define ARGS_EXPAND_50 ARGS_EXPAND_49 , args[49]
#define ARGS_EXPAND_51 ARGS_EXPAND_50 , args[50]
#define ARGS_EXPAND_52 ARGS_EXPAND_51 , args[51]
#define ARGS_EXPAND_53 ARGS_EXPAND_52 , args[52]
#define ARGS_EXPAND_54 ARGS_EXPAND_53 , args[53]
#define ARGS_EXPAND_55 ARGS_EXPAND_54 , args[54]
#define ARGS_EXPAND_56 ARGS_EXPAND_55 , args[55]
#define ARGS_EXPAND_57 ARGS_EXPAND_56 , args[56]
#define ARGS_EXPAND_58 ARGS_EXPAND_57 , args[57]
#define ARGS_EXPAND_59 ARGS_EXPAND_58 , args[58]
#define ARGS_EXPAND_60 ARGS_EXPAND_59 , args[59]
#define ARGS_EXPAND_61 ARGS_EXPAND_60 , args[60]
#define ARGS_EXPAND_62 ARGS_EXPAND_61 , args[61]
#define ARGS_EXPAND_63 ARGS_EXPAND_62 , args[62]
#define ARGS_EXPAND_64 ARGS_EXPAND_63 , args[63]
#define FUNCTION_CAST(number, ...) case number: ret = reinterpret_cast<uint64_t(*)(REPEAT(number, void*))>(address)(ARGS(number)); break;

uint64_t SystemCallN(void* address, uint64_t count, ...) {
    va_list param;
    va_start(param, count);
    uint64_t ret;
    void** args = new void*[count];
    for(uint64_t i = 0; i < count; ++i) {
        args[i] = va_arg(param, void*);
    }
    switch (count) {
        FUNCTION_CAST(0)
        FUNCTION_CAST(1)
        FUNCTION_CAST(2)
        FUNCTION_CAST(3)
        FUNCTION_CAST(4)
        FUNCTION_CAST(5)
        FUNCTION_CAST(6)
        FUNCTION_CAST(7)
        FUNCTION_CAST(8)
        FUNCTION_CAST(9)
        FUNCTION_CAST(10)
        FUNCTION_CAST(11)
        FUNCTION_CAST(12)
        FUNCTION_CAST(13)
        FUNCTION_CAST(14)
        FUNCTION_CAST(15)
        FUNCTION_CAST(16)
        FUNCTION_CAST(17)
        FUNCTION_CAST(18)
        FUNCTION_CAST(19)
        FUNCTION_CAST(20)
        FUNCTION_CAST(21)
        FUNCTION_CAST(22)
        FUNCTION_CAST(23)
        FUNCTION_CAST(24)
        FUNCTION_CAST(25)
        FUNCTION_CAST(26)
        FUNCTION_CAST(27)
        FUNCTION_CAST(28)
        FUNCTION_CAST(29)
        FUNCTION_CAST(30)
        FUNCTION_CAST(31)
        FUNCTION_CAST(32)
        FUNCTION_CAST(33)
        FUNCTION_CAST(34)
        FUNCTION_CAST(35)
        FUNCTION_CAST(36)
        FUNCTION_CAST(37)
        FUNCTION_CAST(38)
        FUNCTION_CAST(39)
        FUNCTION_CAST(40)
        FUNCTION_CAST(41)
        FUNCTION_CAST(42)
        FUNCTION_CAST(43)
        FUNCTION_CAST(44)
        FUNCTION_CAST(45)
        FUNCTION_CAST(46)
        FUNCTION_CAST(47)
        FUNCTION_CAST(48)
        FUNCTION_CAST(49)
        FUNCTION_CAST(50)
        FUNCTION_CAST(51)
        FUNCTION_CAST(52)
        FUNCTION_CAST(53)
        FUNCTION_CAST(54)
        FUNCTION_CAST(55)
        FUNCTION_CAST(56)
        FUNCTION_CAST(57)
        FUNCTION_CAST(58)
        FUNCTION_CAST(59)
        FUNCTION_CAST(60)
        FUNCTION_CAST(61)
        FUNCTION_CAST(62)
        FUNCTION_CAST(63)
        FUNCTION_CAST(64)
        default:
            break;
    }
    va_end(param);
    return ret;
}
uint64_t SystemCall(void* addr) { return reinterpret_cast<uint64_t(*)()>(addr)(); }
uint64_t SystemCall_1(void* addr, void* p1) { return SystemCallN(addr, 1, p1); }
uint64_t SystemCall_2(void* addr, void* p1, void* p2) { return SystemCallN(addr, 2, p1, p2); }
uint64_t SystemCall_3(void* addr, void* p1, void* p2, void* p3) { return SystemCallN(addr, 3, p1, p2, p3); }
uint64_t SystemCall_4(void* addr, void* p1, void* p2, void* p3, void* p4) { return SystemCallN(addr, 4, p1, p2, p3, p4); }

#if defined(_WIN32)
#include <Windows.h>
#define Handle HINSTANCE

uintptr_t OpenLibrary(const char* path) {
    ::SetErrorMode(SEM_FAILCRITICALERRORS|SEM_NOOPENFILEERRORBOX);
    Handle hnd = nullptr;
    hnd = ::LoadLibrary(path);
    if(!hnd) return 0;
    return reinterpret_cast<uintptr_t>(hnd);
}

const char* GetLibraryError() {
    LPSTR messageBuffer = nullptr;
    FormatMessageA(FORMAT_MESSAGE_ALLOCATE_BUFFER | FORMAT_MESSAGE_FROM_SYSTEM | FORMAT_MESSAGE_IGNORE_INSERTS,
                   nullptr,
                   GetLastError(),
                   MAKELANGID(LANG_ENGLISH, SUBLANG_ENGLISH_US),
                   (LPSTR)&messageBuffer,
                   0,
                   nullptr);
    return messageBuffer;
}

bool CloseLibrary(uintptr_t handle) {
    return FreeLibrary(reinterpret_cast<HMODULE>(handle));
}

void* FindSymbol(uintptr_t handle, const char* symbol) {
    void* sym = reinterpret_cast<void*>(GetProcAddress(reinterpret_cast<HMODULE>(handle), symbol));
    if(sym == NULL) return 0;
    return sym;
}
# elif defined(__linux__) || defined(__APPLE__)
#include <dlfcn.h>
#include <limits.h>
#include <stdlib.h>
#include <stdint.h>
#define Handle void*

uintptr_t OpenLibrary(const char* path) {
    Handle h = dlopen(path, RTLD_LAZY|RTLD_GLOBAL);
    if (h == nullptr) {
        return 0;
    }
    return reinterpret_cast<uintptr_t>(h);
}

bool CloseLibrary(uintptr_t handle) {
    dlclose(reinterpret_cast<void*>(handle));
    return dlerror() == nullptr;
}

const char* GetLibraryError() {
    return dlerror();
}

void* FindSymbol(uintptr_t handle, const char* symbol) {
    void* r = dlsym(reinterpret_cast<void*>(h), name);
    if (r == nullptr) {
        return 0;
    }
	return r;
}
#endif

#ifdef __cplusplus
}
#endif