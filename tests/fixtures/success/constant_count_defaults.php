<?php
function size($a, $mode = COUNT_RECURSIVE) { return count($a, $mode); }
class C { public $mode = COUNT_NORMAL; }
echo COUNT_NORMAL, ':', COUNT_RECURSIVE, ':', size([[1]]), ':', (new C())->mode, "\n";
