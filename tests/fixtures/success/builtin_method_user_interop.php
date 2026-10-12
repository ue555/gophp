<?php
class C { function count($a) { return count($a); } function strlen() { return 'method'; } }
function size($a) { return (new C())->count($a); }
echo size([1, 2]), ':', (new C())->strlen(), ':', count([1]), "\n";
