<?php
function fail($n) { return $n % 0; }
class C { function call($n) { return fail($n); } }
$o = new C; $o->call(7);
