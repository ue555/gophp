<?php
class Box {
    public $a = ['x' => [4]]; public $calls = 0;
    function target() { ++$this->calls; echo 'T'; return $this; }
    function key() { ++$this->calls; echo 'K'; return 'x'; }
}
$o = new Box; unset($o->target()->a[$o->key()][0]);
echo ':', count($o->a['x']), ':', $o->calls, "\n";
