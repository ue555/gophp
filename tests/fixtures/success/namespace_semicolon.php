<?php
namespace Acme\One;
const ANSWER = 42; function f() { return ANSWER; } class Thing {}
echo f(), "|", Thing::class, "|", __NAMESPACE__, "\n";
