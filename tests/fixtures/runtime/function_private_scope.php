<?php
class C { private $x = 7; function call() { return readPrivate($this); } }
function readPrivate($o) { return $o->x; }
$o = new C; $o->call();
