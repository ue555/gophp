<?php
echo isset($this) + 0, empty($this) + 0, ':';
function f() { return isset($this); }
class C { function check() { return isset($this) && !f(); } }
echo (int)(new C())->check(), "\n";
