@echo off

del oppai.dll
del oppai.lib
del oppai.exp

cl  -D_CRT_SECURE_NO_WARNINGS ^
	-DOPPAI_IMPLEMENTATION ^
    -DNOMINMAX ^
	-DBUILD_LIBRARY^
    -O2 ^
    -nologo -MD -GR -EHsc -W4 ^
    /LD /fp:fast rtpp.c -Feoppai.dll