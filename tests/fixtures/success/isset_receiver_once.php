<?php
class C { public $a = [null, 5]; public $calls = 0; function target() { ++$this->calls; echo 'T'; return $this; } function key() { ++$this->calls; echo 'K'; return 1; } }
$o = new C; echo isset($o->target()->a[$o->key()]) + 0, ':', $o->calls, "\n";
echo empty($o->target()->a[$o->key()]) + 0, ':', $o->calls, "\n";
