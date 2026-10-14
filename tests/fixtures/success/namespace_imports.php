<?php
namespace X { class C {} function f() { return 4; } const K = 5; }
namespace Y { use X\C as Alias; use function X\f as call; use const X\K as NUM; echo Alias::class, "|", call(), "|", NUM, "\n"; }
